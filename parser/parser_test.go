package parser

import (
	"encoding/json"
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
			exp, _ := json.Marshal(expected)
			res, _ := json.Marshal(result)
			t.Fatalf("Expected \"%s\", got \"%s\"", exp, res)
		}
	}

	test("foo", &IdentExpr{"foo", 0})
	test("123", &IntExpr{123, 0})
	test("[1", nil)
	test("true", &BoolExpr{true, 0})
	test("3 < 5", &LessThanExpr{&IntExpr{3, 0}, &IntExpr{5, 3}, 1})

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
		exp, _ := json.Marshal(expected)
		res, _ := json.Marshal(result)
		if result == nil {
			t.Fatalf("Expected \"%s\", got \"nil\"", exp)
		}
		if !(reflect.DeepEqual(result, expected)) {
			t.Fatalf("Expected \"%s\", got \"%s\"", exp, res)
		}
	}

	test("let foo: int = 123", &AssignStmt{"foo", "int", &IntExpr{123, 14}, 0})
	test("let bar: long = -1928", &AssignStmt{"bar", "long", &IntExpr{-1928, 15}, 0})

	test("let bar: long = 3 - 4 + foo", &AssignStmt{"bar", "long",
		&AddExpr{
			&SubExpr{&IntExpr{3, 15}, &IntExpr{4, 19}, 17},
			&IdentExpr{"foo", 23},
			21,
		}, 0})
	test("let x: bool = false", &AssignStmt{"x", "bool", &BoolExpr{false, 13}, 0})
	test("let x: bool = 5 < 4", &AssignStmt{"x", "bool", &LessThanExpr{&IntExpr{5, 13}, &IntExpr{4, 17}, 15}, 0})
}

func TestParseProgram(t *testing.T) {
	input := `
		let foo: int = 123
		let bar: int = foo - 4
	`
	l := lexer.New(&input)
	program, err := ParseProgram(l)
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 2 {
		t.Fatal("Expected two statements")
	}

	difference, _ := diff.Diff(
		&AssignStmt{
			Lhs:  "foo",
			Tipe: "int",
			Rhs:  &IntExpr{123, 17},
			Pos:  0,
		},
		program.Statements[0],
	)
	if err != nil {
		t.Error(err)
	}
	if len(difference) != 0 {
		t.Errorf("Failed to parse first statement: %+v", difference)
	}

	difference, _ = diff.Diff(
		&AssignStmt{
			Lhs:  "bar",
			Tipe: "int",
			Rhs: &SubExpr{
				Lhs: &IdentExpr{"foo", 38},
				Rhs: &IntExpr{4, 44},
				Pos: 42,
			},
			Pos: 21,
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

func TestParseBinaryExprAssociaticity(t *testing.T) {
	input := "1 + 2 - 3"
	lexer := lexer.New(&input)
	_, node := parseExpression(lexer)
	expected := &SubExpr{
		Lhs: &AddExpr{
			Lhs: &IntExpr{1, 0},
			Rhs: &IntExpr{2, 3},
			Pos: 1,
		},
		Rhs: &IntExpr{3, 7},
		Pos: 5,
	}

	difference, err := diff.Diff(expected, node)
	if err != nil {
		t.Error(err)
	}
	if len(difference) != 0 {
		nodeRepr, _ := json.MarshalIndent(node, "", "  ")
		expectedRepr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Failed to parse expression. Got %s, expected %s", nodeRepr, expectedRepr)

	}
}
