package lexer

import (
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
	if lexer.Position > len(input) {
		t.Fatalf("Expected position to be l.t.e \"%d\", got \"%d\"", len(input), lexer.Position)
	}
}

func TestNexToken(t *testing.T) {

	testCase := func(input *string, expected *[]Token) {
		lexer := New(input)

		for i, expected := range *expected {
			var tok Token
			lexer, tok = lexer.Next()

			if tok.Type != expected.Type {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, expected.Type, tok.Type)
			}

			if tok.Lexeme != expected.Lexeme {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
					i, expected.Lexeme, tok.Lexeme)
			}
		}
	}

	input1 := `=+(){},;`
	testCase(&input1, &[]Token{
		{Type: ASSIGN, Lexeme: "="},
		{Type: PLUS, Lexeme: "+"},
		{Type: LPAREN, Lexeme: "("},
		{Type: RPAREN, Lexeme: ")"},
		{Type: LBRACE, Lexeme: "{"},
		{Type: RBRACE, Lexeme: "}"},
		{Type: COMMA, Lexeme: ","},
		{Type: SEMICOLON, Lexeme: ";"},
		{Type: EOF, Lexeme: ""},
	})

	input2 := ""
	testCase(&input2, &[]Token{{Type: EOF, Lexeme: ""}})

	input3 := "let five = 5;"
	testCase(&input3, &[]Token{
		{Type: LET, Lexeme: "let"},
		{Type: IDENT, Lexeme: "five"},
		{Type: ASSIGN, Lexeme: "="},
		{Type: INT, Lexeme: "5"},
		{Type: SEMICOLON, Lexeme: ";"},
		{Type: EOF, Lexeme: ""},
	})

	input4 := "return x = 3 == 5;"
	testCase(&input4, &[]Token{
		{Type: RETURN, Lexeme: "return"},
		{Type: IDENT, Lexeme: "x"},
		{Type: ASSIGN, Lexeme: "="},
		{Type: INT, Lexeme: "3"},
		{Type: EQ, Lexeme: "=="},
		{Type: INT, Lexeme: "5"},
		{Type: SEMICOLON, Lexeme: ";"},
		{Type: EOF, Lexeme: ""},
	})

	input5 := "`"
	testCase(&input5, &[]Token{
		{Type: ILLEGAL, Lexeme: "`"},
		{Type: EOF, Lexeme: ""},
	})

	input6 := "let letter = 5;"
	testCase(&input6, &[]Token{
		{Type: LET, Lexeme: "let"},
		{Type: IDENT, Lexeme: "letter"},
		{Type: ASSIGN, Lexeme: "="},
		{Type: INT, Lexeme: "5"},
		{Type: SEMICOLON, Lexeme: ";"},
		{Type: EOF, Lexeme: ""},
	})

	input7 := "def myFunc(arg): int = return 5;"
	testCase(&input7, &[]Token{
		{Type: FUNCTION, Lexeme: "def"},
		{Type: IDENT, Lexeme: "myFunc"},
		{Type: LPAREN, Lexeme: "("},
		{Type: IDENT, Lexeme: "arg"},
		{Type: RPAREN, Lexeme: ")"},
		{Type: ASSIGN_T, Lexeme: ":"},
		{Type: IDENT, Lexeme: "int"},
		{Type: ASSIGN, Lexeme: "="},
		{Type: RETURN, Lexeme: "return"},
		{Type: INT, Lexeme: "5"},
	})

	input8 := "8 -9 100 8.2"
	testCase(&input8, &[]Token{
		{Type: INT, Lexeme: "8"},
		{Type: INT, Lexeme: "-9"},
		{Type: INT, Lexeme: "100"},
		{Type: FLOAT, Lexeme: "8.2"},
	})

	// input9 := "let x: int = 5; def isMultipleof5And2(n: int) = {}"
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
	if lexer.Position != 2 {
		t.Fatalf("Expected position to be \"2\", got \"%d\"", lexer.Position)
	}
}

func TestBacktrack(t *testing.T) {

	testCase := func(input string, lit string, expectedPos int) {
		lexer := New(&input)
		reader := withBacktrack(readLiteral(lit))
		lexer, _ = reader(lexer)
		if lexer.Position != expectedPos {
			t.Fatalf("Expected position to be \"%d\", got \"%d\"", expectedPos, lexer.Position)
		}
	}

	testCase("foo", "foo", 3)
	testCase("foo", "fob", 0)

}
