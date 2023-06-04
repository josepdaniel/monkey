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
	return []Instruction{}, env, errors.New("expected statement")
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
		return compileIntegerExpression(expression)
	}
	return []Instruction{}, errors.New("expected expression")
}

func compileIntegerExpression(expression *parser.IntExpr) ([]Instruction, error) {
	return []Instruction{
		MOV("rax", fmt.Sprint(expression.Value)),
	}, nil
}
