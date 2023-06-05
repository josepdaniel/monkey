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

	var env = &Env{
		locals: make([]string, 0),
	}
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
	compiledExpression, err := compileExpression(statement.Rhs, env)
	if err != nil {
		return []Instruction{}, env, err
	}
	output := append(compiledExpression, PUSH("rax"))
	env = env.addLocal(statement.Lhs)

	return output, env, nil
}

func compileExpression(expression parser.Expression, env *Env) ([]Instruction, error) {
	switch expression := expression.(type) {
	case *parser.IntExpr:
		return compileIntegerExpression(*expression)
	case *parser.IdentExpr:
		return compileIdentExpression(*expression, env)
	case *parser.AddExpr:
		return compileAddExpression(*expression, env)
	case *parser.SubExpr:
		return compileSubExpression(*expression, env)
	}

	return []Instruction{}, errors.New("[compiler err]: unexpected expression type")
}

func compileIntegerExpression(expression parser.IntExpr) ([]Instruction, error) {
	return []Instruction{
		MOV("rax", fmt.Sprint(expression.Value)),
	}, nil
}

func compileIdentExpression(expression parser.IdentExpr, env *Env) ([]Instruction, error) {
	address, err := env.lexicalAddress(expression.Name)
	if err != nil {
		return []Instruction{}, err
	}
	return []Instruction{
		// For now we are assuming all data types are 8 bytes wide
		// TODO: implement a types table that can be used to look up the size of a type
		MOV("rax", fmt.Sprintf("[rsp+%d]", address*8)),
	}, nil
}

func compileAddExpression(expression parser.AddExpr, env *Env) ([]Instruction, error) {
	// Compile the left hand side, and push the result onto the stack
	left, err := compileExpression(expression.Lhs, env)
	if err != nil {
		return []Instruction{}, err
	}
	output := append(left, PUSH("rax"))

	// Compile the right hand side
	// Note - the number of 'locals' in the environment will need to be bumped up by 1
	// because we pushed a value onto the stack in the LHS, without binding it to a local.
	// "_" is not an allowed identifier, so we can be sure it won't exist in the environment.
	tmpEnv := env.addLocal("_")
	right, err := compileExpression(expression.Rhs, tmpEnv)
	if err != nil {
		return []Instruction{}, err
	}

	output = append(output, right...)
	// The result of the RHS expression will be in RAX. Add it to the LHS expression.
	// Also pop the top element of the stack.
	output = append(output, []Instruction{
		ADD("rax", "[rsp]"),
		ADD("rsp", "8"),
	}...)

	return output, nil
}

func compileSubExpression(expression parser.SubExpr, env *Env) ([]Instruction, error) {
	left, err := compileExpression(expression.Rhs, env)
	if err != nil {
		return []Instruction{}, err
	}
	output := append(left, PUSH("rax"))

	tmpEnv := env.addLocal("_")
	right, err := compileExpression(expression.Lhs, tmpEnv)
	if err != nil {
		return []Instruction{}, err
	}

	output = append(output, right...)
	// The result of the RHS expression will be in RAX. Add it to the LHS expression.
	// Also pop the top element of the stack.
	output = append(output, []Instruction{
		SUB("rax", "[rsp]"),
		ADD("rsp", "8"),
	}...)

	return output, nil
}
