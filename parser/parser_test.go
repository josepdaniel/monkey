package parser

import (
	"monkey/lexer"
	"reflect"
	"testing"

	"github.com/r3labs/diff/v3"
)

func TestParseIdentifier(t *testing.T) {
	input := "foo bar"
	l := lexer.New(&input)
	l, result := parseIdentifier(l)

	if result == nil {
		t.Fatalf("Expected \"foo\", got \"nil\"")
	}

	if result.Name != "foo" {
		t.Fatalf("Expected \"foo\", got \"%s\"", result.Name)
	}

	// The lexer should have advanced to the next token
	_, result = parseIdentifier(l)
	if result == nil {
		t.Fatalf("Expected \"bar\", got \"nil\"")
	}
	if result.Name != "bar" {
		t.Fatalf("Expected \"bar\", got \"%s\"", result.Name)
	}

	// Try to parse an invalid identifier
	input = "* foo"
	l = lexer.New(&input)
	_, result = parseIdentifier(l)
	if result != nil {
		t.Fatalf("Expected \"nil\", got \"%s\"", result.Name)
	}
}

func TestParseInt(t *testing.T) {
	input := "123 -12"
	l := lexer.New(&input)

	l, firstInt := parseInt(l)
	_, secondInt := parseInt(l)

	if firstInt.Value != 123 {
		t.Fatalf("Expected \"123\", got \"%d\"", firstInt.Value)
	}
	if secondInt.Value != -12 {
		t.Fatalf("Expected \"-12\", got \"%d\"", secondInt.Value)
	}

}

func TestParseExpression(t *testing.T) {

	var parseHelper = func(input string) Expression {
		l := lexer.New(&input)
		_, node := parseExpression(l)
		return node
	}

	var test = func(input string, expected Expression) {
		result := parseHelper(input)
		if !(reflect.DeepEqual(result, expected)) {
			t.Fatalf("Expected \"%s\", got \"%s\"", expected, result)
		}
	}

	test("foo", &IdentExpr{"foo"})
	test("123", &IntExpr{123})
	test("[1", nil)

}

func TestAllOf(t *testing.T) {
	input := "let x: int = -145"
	l := lexer.New(&input)

	_, tokens := allOf(l, lexer.LET, lexer.IDENT, lexer.ASSIGN_T, lexer.IDENT, lexer.ASSIGN)
	if len(tokens) != 5 {
		t.Fatalf("Expected 5 tokens, got %d", len(tokens))
	}
}

func TestParseAssignment(t *testing.T) {

	var parseHelper = func(input string) Statement {
		l := lexer.New(&input)
		_, node := parseAssignment(l)
		return node
	}

	var test = func(input string, expected Statement) {
		result := parseHelper(input)
		if result == nil {
			t.Fatalf("Expected \"%s\", got \"nil\"", expected)
		}
		if !(reflect.DeepEqual(result, expected)) {
			t.Fatalf("Expected \"%s\", got \"%s\"", expected, result)
		}
	}

	test("let foo: int = 123", &AssignStmt{"foo", "int", &IntExpr{123}})
	test("let bar: long = -1928", &AssignStmt{"bar", "long", &IntExpr{-1928}})

	test("let bar: long = 3 - 4 + foo", &AssignStmt{"bar", "long",
		&SubExpr{&IntExpr{3}, &AddExpr{&IntExpr{4}, &IdentExpr{"foo"}}}})

}

func TestParseProgram(t *testing.T) {
	input := `
		let foo: int = 123
		let bar: int = foo - 4
	`
	lexer := lexer.New(&input)
	program, err := ParseProgram(lexer)
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 2 {
		t.Fatal("Expected two statements")
	}

	difference, err := diff.Diff(
		&AssignStmt{
			Lhs:  "foo",
			Tipe: "int",
			Rhs:  &IntExpr{Value: 123},
		},
		program.Statements[0],
	)
	if err != nil {
		t.Error(err)
	}
	if len(difference) != 0 {
		t.Errorf("Failed to parse first statement: %+v", difference)
	}

	difference, err = diff.Diff(
		&AssignStmt{
			Lhs:  "bar",
			Tipe: "int",
			Rhs: &SubExpr{
				Lhs: &IdentExpr{Name: "foo"},
				Rhs: &IntExpr{4},
			},
		},
		program.Statements[1],
		diff.AllowTypeMismatch(false),
	)
	if err != nil {
		t.Error(err)
	}
	if len(difference) != 0 {
		t.Errorf("Failed to parse second statement: %+v", difference)
	}

}
