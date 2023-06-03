package compiler

import (
	"errors"
	"fmt"
	"monkey/parser"
)

type Env struct {
	locals []*string
}

func (env *Env) addLocal(s *string) *Env {
	return &Env{
		locals: append(env.locals, s),
	}
}

func compile(program parser.Program) (*string, error) {
	var output string = `
	BITS 64
	CPU X64
	section .text
	global _start
_start:
`

	compiledStatements, err := compileStatements(program.Statements)
	output += *compiledStatements

	if err != nil {
		return nil, err
	} else {
		return &output, nil
	}

}

func compileStatements(statements []parser.Statement) (*string, error) {

	var env = &Env{
		locals: make([]*string, 0),
	}
	var output string

	for _, statement := range statements {
		var res *string
		var err error

		res, env, err = compileStatement(statement, env)
		if err != nil {
			return nil, err
		}
		output += *res
	}
	return &output, nil
}

func compileStatement(statement parser.Statement, env *Env) (*string, *Env, error) {
	switch statement := statement.(type) {
	case *parser.AssignStmt:
		return compileAssignStmt(statement, env)
	}
	return nil, env, errors.New("unknown statement type")
}

func compileAssignStmt(statement *parser.AssignStmt, env *Env) (*string, *Env, error) {
	var output string
	compiledExpression, err := compileExpression(statement.Rhs, env)
	if err != nil {
		return nil, env, err
	}
	output += *compiledExpression
	output += fmt.Sprintln("\tpush rax")

	env = env.addLocal(&statement.Lhs)

	return &output, env, nil
}

func compileExpression(expression parser.Expression, env *Env) (*string, error) {
	switch expression := expression.(type) {
	case *parser.IntExpr:
		return compileIntegerExpression(expression)
	}
	return nil, errors.New("unknown expression type")
}

func compileIntegerExpression(expression *parser.IntExpr) (*string, error) {
	var result = fmt.Sprintf("\tmov rax %d\n", expression.Value)
	return &result, nil
}
