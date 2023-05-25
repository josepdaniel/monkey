package lexer

type TokenType string

type Token struct {
	Type   TokenType
	Lexeme string
}

var Keywords = map[string]TokenType{
	"def":    FUNCTION,
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
	"&&": AND,
	"||": OR,
}

var Operators = map[string]TokenType{
	":": ASSIGN_T,
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
	FLOAT     TokenType = "FLOAT"
	EQ        TokenType = "=="
	NEQ       TokenType = "!="
	AND       TokenType = "&&"
	OR        TokenType = "||"
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
	ASSIGN_T  TokenType = ":"
	FUNCTION  TokenType = "FUNCTION"
	LET       TokenType = "LET"
	TRUE      TokenType = "TRUE"
	FALSE     TokenType = "FALSE"
	IF        TokenType = "IF"
	ELSE      TokenType = "ELSE"
	RETURN    TokenType = "RETURN"
)
