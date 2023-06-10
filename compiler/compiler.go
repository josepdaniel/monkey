package compiler

import (
	"errors"
	"fmt"
	"monkey/parser"
)

func Compile(program parser.Program) ([]Instruction, error) {

	var prelude = []Instruction{
		SECTION(".text"),
		GLOBAL("_start"),
		LABEL("_start"),
	}

	var epilogue = []Instruction{
		MOV("rdi", "rax"),       // exit code - for now make it whatever was in rax
		MOV("rax", "0x2000001"), // exit syscall
		SYSCALL(),
	}

	compiledStatements, err := compileStatements(program.Statements)
	output := append(append(prelude, compiledStatements...), epilogue...)

	if err != nil {
		return []Instruction{}, err
	} else {
		return output, nil
	}

}

func compileStatements(statements []parser.Statement) ([]Instruction, error) {

	var env = NewEnv()
	var output []Instruction

	for _, statement := range statements {
		var res []Instruction
		var err error

		res, env, err = compileStatement(statement, env)
		if err != nil {
			return []Instruction{}, err
		}
		output = append(output, res...)
	}
	return output, nil
}

func compileStatement(statement parser.Statement, env *Env) ([]Instruction, *Env, error) {
	switch statement := statement.(type) {
	case *parser.AssignStmt:
		return compileAssignStmt(statement, env)
	}
	return []Instruction{}, env, errors.New("[compiler err]: unexpected statement type")
}

func compileAssignStmt(statement *parser.AssignStmt, env *Env) ([]Instruction, *Env, error) {
	compiledExpression, exprTipe, err := compileExpression(statement.Rhs, env)
	if err != nil {
		return []Instruction{}, env, err
	}
	assignmentTipe, ok := env.lookupTipe(statement.Tipe)

	if !ok {
		err := fmt.Sprint("[type err]: type not found ", statement.Tipe)
		return []Instruction{}, env, errors.New(err)
	}

	if exprTipe != assignmentTipe {
		err := fmt.Sprint("[type err]: cannot cannot assign type: ", exprTipe.Name, " to ", statement.Tipe)
		return []Instruction{}, env, errors.New(err)
	}
	output := append(compiledExpression, PUSH("rax"))
	env, err = env.addBinding(statement.Lhs, statement.Tipe)

	if err != nil {
		return []Instruction{}, env, err
	}

	return output, env, nil
}

func compileExpression(expression parser.Expression, env *Env) ([]Instruction, Tipe, error) {
	switch expression := expression.(type) {
	case *parser.IntExpr:
		return compileIntegerExpression(*expression)
	case *parser.IdentExpr:
		return compileIdentExpression(*expression, env)
	case *parser.AddExpr:
		return compileAddExpression(*expression, env)
	case *parser.SubExpr:
		return compileSubExpression(*expression, env)
	case *parser.BoolExpr:
		return compileBoolExpression(*expression, env)
	case *parser.LessThanExpr:
		return compileLessThanExpression(*expression, env)
	}

	return []Instruction{}, T_NEVER(0), errors.New("[compiler err]: unexpected expression type")
}

func compileIntegerExpression(expression parser.IntExpr) ([]Instruction, Tipe, error) {
	return []Instruction{
		MOV("rax", fmt.Sprint(expression.Value)),
	}, T_INT, nil
}

func compileBoolExpression(expression parser.BoolExpr, env *Env) ([]Instruction, Tipe, error) {
	val := 0
	if expression.Value {
		val = 1
	}
	return []Instruction{
		MOV("rax", fmt.Sprint(val)),
	}, T_BOOL, nil
}

func compileLessThanExpression(expression parser.LessThanExpr, env *Env) ([]Instruction, Tipe, error) {
	left, leftTipe, err := compileExpression(expression.Lhs, env)
	if err != nil {
		return []Instruction{}, T_NEVER(0), err
	}
	output := append(left, PUSH("rax"))
	tmpEnv := env.addNever(leftTipe.Size)

	right, rightTipe, err := compileExpression(expression.Rhs, tmpEnv)
	if err != nil {
		return []Instruction{}, T_NEVER(0), err
	}

	if (leftTipe != T_INT) || (rightTipe != T_INT) {
		err := fmt.Sprint("[type err]: cannot cannot compare types: ", leftTipe.Name, " and ", rightTipe.Name)
		return []Instruction{}, T_NEVER(0), errors.New(err)
	}

	output = append(output, right...)

	ifTrue := genLabel()
	done := genLabel()

	output = append(output, []Instruction{
		CMP("[rsp]", "rax"),
		JL(ifTrue),
		MOV("rax", "0"),
		JMP(done),
		LABEL(ifTrue),
		MOV("rax", "1"),
		LABEL(done),
	}...)

	return output, T_BOOL, nil
}

func compileIdentExpression(expression parser.IdentExpr, env *Env) ([]Instruction, Tipe, error) {
	address, tipe, err := env.lexicalAddress(expression.Name)
	if err != nil {
		return []Instruction{}, T_NEVER(0), err
	}
	return []Instruction{
		MOV("rax", fmt.Sprintf("[rsp+%d]", address)),
	}, tipe, nil
}

func compileAddExpression(expression parser.AddExpr, env *Env) ([]Instruction, Tipe, error) {
	left, leftTipe, err := compileExpression(expression.Lhs, env)
	if err != nil {
		return []Instruction{}, T_NEVER(0), err
	}
	output := append(left, PUSH("rax"))
	tmpEnv := env.addNever(leftTipe.Size)

	right, rightTipe, err := compileExpression(expression.Rhs, tmpEnv)
	if err != nil {
		return []Instruction{}, T_NEVER(0), err
	}

	if (leftTipe != T_INT) || (rightTipe != T_INT) {
		err := fmt.Sprint("[type err]: cannot cannot add types: ", leftTipe.Name, " and ", rightTipe.Name)
		return []Instruction{}, T_NEVER(0), errors.New(err)
	}

	output = append(output, right...)
	output = append(output, []Instruction{
		ADD("rax", "[rsp]"),
		ADD("rsp", "8"),
	}...)

	return output, leftTipe, nil
}

func compileSubExpression(expression parser.SubExpr, env *Env) ([]Instruction, Tipe, error) {
	left, leftTipe, err := compileExpression(expression.Rhs, env)
	if err != nil {
		return []Instruction{}, T_NEVER(0), err
	}
	output := append(left, PUSH("rax"))
	tmpEnv := env.addNever(leftTipe.Size)

	right, rightTipe, err := compileExpression(expression.Lhs, tmpEnv)
	if err != nil {
		return []Instruction{}, T_NEVER(0), err
	}

	if (leftTipe != T_INT) || (rightTipe != T_INT) {
		err := fmt.Sprint("[type err]: cannot cannot subtract types: ", leftTipe.Name, " and ", rightTipe.Name)
		return []Instruction{}, T_NEVER(0), errors.New(err)
	}

	output = append(output, right...)
	output = append(output, []Instruction{
		SUB("rax", "[rsp]"),
		ADD("rsp", "8"),
	}...)

	return output, leftTipe, nil
}
