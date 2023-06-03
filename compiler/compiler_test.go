package compiler

import (
	"monkey/parser"
	"testing"
)

func TestCompile(t *testing.T) {
	program := parser.Program{
		Statements: []parser.Statement{
			&parser.AssignStmt{
				Lhs:  "foo",
				Tipe: "int",
				Rhs:  &parser.IntExpr{Value: 3},
			},
		},
	}

	var compiled, err = compile(program)
	if err != nil {
		t.Fatalf("Failed to compile: %s", err.Error())
	}
	println(*compiled)

}
