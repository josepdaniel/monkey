package lexer

import (
	"monkey/token"
	"testing"
)

func TestInNextPosition(t *testing.T) {
	// Check that the position pointer will not exceed the length of the input
	input := "foo"
	lexer := New(&input)

	// Try to advance the pointer beyond the input length
	for i := 0; i < len(input)*2; i++ {
		lexer = lexer.inNextPosition()
	}
	if lexer.position > len(input) {
		t.Fatalf("Expected position to be l.t.e \"%d\", got \"%d\"", len(input), lexer.position)
	}
}

func TestNexToken(t *testing.T) {

	testCase := func(input *string, expected *[]token.Token) {
		lexer := New(input)

		for i, expected := range *expected {
			var tok token.Token
			lexer, tok = lexer.Next()

			if tok.Type != expected.Type {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, expected.Type, tok.Type)
			}

			if tok.Literal != expected.Literal {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
					i, expected.Literal, tok.Literal)
			}
		}
	}

	input1 := `=+(){},;`
	testCase(&input1, &[]token.Token{
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	})

	input2 := ""
	testCase(&input2, &[]token.Token{{Type: token.EOF, Literal: ""}})

	input3 := ".let five = 5;"
	testCase(&input3, &[]token.Token{
		{Type: token.LET, Literal: ".let"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	})

	input4 := ".return x = 3 == 5;"
	testCase(&input4, &[]token.Token{
		{Type: token.RETURN, Literal: ".return"},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "3"},
		{Type: token.EQ, Literal: "=="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	})

	input5 := "`"
	testCase(&input5, &[]token.Token{
		{Type: token.ILLEGAL, Literal: "`"},
		{Type: token.EOF, Literal: ""},
	})

	input6 := ".let letter = 5;"
	testCase(&input6, &[]token.Token{
		{Type: token.LET, Literal: ".let"},
		{Type: token.IDENT, Literal: "letter"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	})
}

func TestReadWord(t *testing.T) {

	input := "foo bar baz"
	expected := []string{"foo", "bar", "baz"}

	lexer := New(&input)

	for _, exp := range expected {
		var word *string
		lexer, word = readWord(lexer)
		if *word != exp {
			t.Fatalf("Expected \"%s\", got \"%s\"", exp, *word)
		}
		// Skip over the whitespace character to get to the start of the next word
		lexer = lexer.inNextPosition()
	}

	// test the invalid input case
	input2 := "2"
	lexer2 := New(&input2)
	_, word := readWord(lexer2)

	if word != nil {
		t.Fatalf("Expected \"nil\", got \"%s\"", *word)
	}

}

func TestReadLiteral(t *testing.T) {

	input := "foo"
	lexer := New(&input)
	lexer, result := readLiteral("foo")(lexer)
	if *result != "foo" {
		t.Fatalf("Expected \"foo\", got \"%s\"", *result)
	}

	input = "foo"
	lexer = New(&input)
	lexer, result = readLiteral("fob")(lexer)
	if result != nil {
		t.Fatalf("Expected \"nil\", got \"%s\"", *result)
	}
	if lexer.position != 2 {
		t.Fatalf("Expected position to be \"2\", got \"%d\"", lexer.position)
	}
}

func TestBacktrack(t *testing.T) {

	testCase := func(input string, lit string, expectedPos int) {
		lexer := New(&input)
		reader := withBacktrack(readLiteral(lit))
		lexer, _ = reader(lexer)
		if lexer.position != expectedPos {
			t.Fatalf("Expected position to be \"%d\", got \"%d\"", expectedPos, lexer.position)
		}
	}

	testCase("foo", "foo", 3)
	testCase("foo", "fob", 0)

}
