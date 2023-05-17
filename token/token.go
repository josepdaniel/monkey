package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var Keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
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
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	EQ        = "=="
	NEQ       = "!="
	ASSIGN    = "="
	PLUS      = "+"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	MINUS     = "-"
	BANG      = "!"
	SLASH     = "/"
	ASTERISK  = "*"
	LT        = "<"
	GT        = ">"
	FUNCTION  = "FUNCTION"
	LET       = "LET"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
)
