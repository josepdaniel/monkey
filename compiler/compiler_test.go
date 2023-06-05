package compiler

import (
	"fmt"
	"monkey/lexer"
	"monkey/parser"
	"testing"
)

func parseHelper(t *testing.T, s string) parser.Program {
	lexer := lexer.New(&s)
	program, error := parser.ParseProgram(lexer)
	if error != nil {
		t.Fatal("Could not parse program: ", error)
	}
	return *program
}

func TestCompile(t *testing.T) {
	program := parseHelper(t, "let x: int = 3 let y: int = x")

	var compiled, err = Compile(program)
	if err != nil {
		t.Fatalf("Failed to compile: %s", err.Error())
	}
	fmt.Println(Render(compiled))

}
