package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input    *string
	position int
}

type reader func(lexer Lexer) (Lexer, *string)

func New(input *string) Lexer {
	return Lexer{input: input, position: 0}
}

// Returns a new lexer with the new state
func (lexer Lexer) inNextPosition() Lexer {
	nextPosition := minInt(len(*lexer.input), lexer.position+1)
	return Lexer{input: lexer.input, position: nextPosition}
}

// Return the current character
func (lexer Lexer) currentChar() byte {
	if lexer.position >= len(*lexer.input) {
		return 0
	} else {
		return (*lexer.input)[lexer.position]
	}
}

// Reads until the stop function matches the current letter
func (lexer Lexer) until(stop func(ch byte) bool) (Lexer, *string) {

	// If the first character is already a stop character, return nil
	if stop(lexer.currentChar()) {
		return lexer, nil
	}

	initialPosition := lexer.position
	for !stop(lexer.currentChar()) {
		lexer = lexer.inNextPosition()
	}
	result := (*lexer.input)[initialPosition:lexer.position]
	return lexer, &result
}

// Reads a word (letters only)
var readWord reader = func(lexer Lexer) (Lexer, *string) {
	return lexer.until(func(ch byte) bool {
		return !isLetter(ch)
	})
}

// Reads a number (no floats)
var readNumber reader = func(lexer Lexer) (Lexer, *string) {
	return lexer.until(func(ch byte) bool {
		return !isDigit(ch)
	})
}

// Reads a literal string
func readLiteral(literal string) reader {
	return func(lexer Lexer) (Lexer, *string) {
		for i := 0; i < len(literal); i++ {
			if lexer.currentChar() != literal[i] {
				return lexer, nil
			}
			lexer = lexer.inNextPosition()
		}
		return lexer, &literal
	}
}

// Backtracks the reader to its original position if it fails
func withBacktrack(reader reader) reader {
	return func(lexer Lexer) (Lexer, *string) {
		newLexer, result := reader(lexer)
		if result == nil {
			return lexer, nil
		} else {
			return newLexer, result
		}
	}
}

// Advance until the lexer is not in a whitespace position
func (lexer Lexer) skipWhitespace() Lexer {
	lexer, _ = lexer.until(func(ch byte) bool {
		return ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r'
	})
	return lexer
}

func readToken(lexer Lexer, litToTokenMapping map[string]token.TokenType) (Lexer, *token.Token) {

	for key, value := range litToTokenMapping {
		if lexer, lit := withBacktrack(readLiteral(key))(lexer); lit != nil {
			return lexer, &token.Token{Type: value, Literal: *lit}
		}
	}
	return lexer, nil
}

func (lexer Lexer) Next() (Lexer, token.Token) {
	lexer = lexer.skipWhitespace()

	if lexer.currentChar() == 0 {
		return lexer, token.Token{Type: token.EOF, Literal: ""}
	} else if lexer, tok := readToken(lexer, token.DoubleCharOperators); tok != nil {
		return lexer, *tok
	} else if lexer, tok := readToken(lexer, token.Operators); tok != nil {
		return lexer, *tok
	} else if lexer, tok := readToken(lexer, token.Delimiters); tok != nil {
		return lexer, *tok
	} else if lexer, tok := readToken(lexer, token.Keywords); tok != nil {
		return lexer, *tok
	} else if lexer, lit := withBacktrack(readWord)(lexer); lit != nil {
		return lexer, token.Token{Type: token.IDENT, Literal: *lit}
	} else if lexer, lit := withBacktrack(readNumber)(lexer); lit != nil {
		return lexer, token.Token{Type: token.INT, Literal: *lit}
	} else {
		return lexer.inNextPosition(), token.Token{Type: token.ILLEGAL, Literal: string(lexer.currentChar())}
	}

}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }
