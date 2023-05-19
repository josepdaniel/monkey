package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// We define keywords with a leading period to distinguish them from regular identifiers
var Keywords = map[string]TokenType{
	".fn":     FUNCTION,
	".let":    LET,
	".true":   TRUE,
	".false":  FALSE,
	".if":     IF,
	".else":   ELSE,
	".return": RETURN,
}

var DoubleCharOperators = map[string]TokenType{
	"==": EQ,
	"!=": NEQ,
}

var Operators = map[string]TokenType{
	"=": ASSIGN,
	"+": PLUS,
	"-": MINUS,
	"!": BANG,
	"/": SLASH,
	"*": ASTERISK,
	"<": LT,
	">": GT,
}

var Delimiters = map[string]TokenType{
	",": COMMA,
	";": SEMICOLON,
	"(": LPAREN,
	")": RPAREN,
	"{": LBRACE,
	"}": RBRACE,
}

const (
	ILLEGAL   TokenType = "ILLEGAL"
	EOF       TokenType = "EOF"
	IDENT     TokenType = "IDENT"
	INT       TokenType = "INT"
	EQ        TokenType = "=="
	NEQ       TokenType = "!="
	ASSIGN    TokenType = "="
	PLUS      TokenType = "+"
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	MINUS     TokenType = "-"
	BANG      TokenType = "!"
	SLASH     TokenType = "/"
	ASTERISK  TokenType = "*"
	LT        TokenType = "<"
	GT        TokenType = ">"
	FUNCTION  TokenType = "FUNCTION"
	LET       TokenType = "LET"
	TRUE      TokenType = "TRUE"
	FALSE     TokenType = "FALSE"
	IF        TokenType = "IF"
	ELSE      TokenType = "ELSE"
	RETURN    TokenType = "RETURN"
)
