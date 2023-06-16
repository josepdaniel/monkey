package lexer

import (
	"strconv"
	"strings"
)

type Lexer struct {
	Input    *string
	Position int
}

// Convention for a reader:
//   - The returned lexer should have its pointer advanced to the point
//     where reading stopped, even if the reader failed
//   - The returned string pointer should be nil if the reader failed,
//     otherwise it should contain the raw string read from the input
type reader func(lexer Lexer) (Lexer, *string)

func New(input *string) Lexer {
	return Lexer{Input: input, Position: 0}
}

// Returns a new lexer with the new state
func (lexer Lexer) inNextPosition() Lexer {
	nextPosition := minInt(len(*lexer.Input), lexer.Position+1)
	return Lexer{Input: lexer.Input, Position: nextPosition}
}

func (lexer Lexer) CurrentLine() string {
	var start, end int = lexer.Position, lexer.Position + 1
	for {
		if start == 0 || (*lexer.Input)[start-1] == '\n' || (*lexer.Input)[start] == '\n' {
			break
		}
		start--
	}

	for (*lexer.Input)[end] == '\n' {
		end++
	}

	for {
		if end == len(*lexer.Input) || (*lexer.Input)[end] == '\n' {
			break
		}
		end++
	}
	return strings.TrimSpace((*lexer.Input)[start:end])
}

// Return the current character
func (lexer Lexer) currentChar() byte {
	if lexer.Position >= len(*lexer.Input) {
		return 0
	} else {
		return (*lexer.Input)[lexer.Position]
	}
}

// Reads until the stop function matches the current letter
func (lexer Lexer) until(stop func(ch byte) bool) (Lexer, *string) {

	// If the first character is already a stop character, return nil
	if stop(lexer.currentChar()) {
		return lexer, nil
	}

	initialPosition := lexer.Position
	for !stop(lexer.currentChar()) {
		lexer = lexer.inNextPosition()
	}
	result := (*lexer.Input)[initialPosition:lexer.Position]
	return lexer, &result
}

// Reads a word (letters only)
var readWord reader = func(lexer Lexer) (Lexer, *string) {
	return lexer.until(func(ch byte) bool {
		return !isLetter(ch)
	})
}

// Reads a floating point number
var readInt reader = func(lexer Lexer) (Lexer, *string) {
	lexer, lexeme := lexer.until(func(ch byte) bool {
		return !isDigit(ch) && ch != '.' && ch != '-'
	})

	if lexeme == nil {
		return lexer, lexeme
	}
	_, err := strconv.ParseInt(*lexeme, 10, 16)

	if err != nil {
		return lexer, nil
	}

	return lexer, lexeme
}

var readFloat reader = func(lexer Lexer) (Lexer, *string) {
	lexer, lexeme := lexer.until(func(ch byte) bool {
		return !isDigit(ch) && ch != '.' && ch != '-'
	})

	if lexeme == nil {
		return lexer, lexeme
	}
	_, err := strconv.ParseFloat(*lexeme, 32)

	if err != nil {
		return lexer, nil
	}

	return lexer, lexeme
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

func readToken(lexer Lexer, litToTokenMapping map[string]TokenType) (Lexer, *Token) {

	for key, value := range litToTokenMapping {
		if lexer, lit := withBacktrack(readLiteral(key))(lexer); lit != nil {
			return lexer, &Token{Type: value, Lexeme: *lit}
		}
	}
	return lexer, nil
}

func (lexer Lexer) Next() (Lexer, Token) {
	lexer = lexer.skipWhitespace()

	if lexer.currentChar() == 0 {
		return lexer, Token{Type: EOF, Lexeme: ""}
	} else if lexer, lit := withBacktrack(readInt)(lexer); lit != nil {
		return lexer, Token{Type: INT, Lexeme: *lit}
	} else if lexer, lit := withBacktrack(readFloat)(lexer); lit != nil {
		return lexer, Token{Type: FLOAT, Lexeme: *lit}
	} else if lexer, tok := readToken(lexer, DoubleCharOperators); tok != nil {
		return lexer, *tok
	} else if lexer, tok := readToken(lexer, Operators); tok != nil {
		return lexer, *tok
	} else if lexer, tok := readToken(lexer, Delimiters); tok != nil {
		return lexer, *tok
	} else if lexer, lit := withBacktrack(readWord)(lexer); lit != nil {

		// If the word is a keyword, treat is as such
		if tok, ok := Keywords[*lit]; ok {
			return lexer, Token{Type: tok, Lexeme: *lit}
		}

		// Otherwise treat the word as an identifier (e.g. variable name)
		return lexer, Token{Type: IDENT, Lexeme: *lit}

	} else {
		return lexer.inNextPosition(), Token{Type: ILLEGAL, Lexeme: string(lexer.currentChar())}
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
