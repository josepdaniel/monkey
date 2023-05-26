package parser

import (
	"monkey/lexer"
	"strconv"
)

func parseIdentifier(l lexer.Lexer) (lexer.Lexer, *IdentExpr) {

	if new, tok := l.Next(); tok.Type == lexer.IDENT {
		ident := &IdentExpr{tok.Lexeme}
		return new, ident
	} else {
		return l, nil
	}

}

func parseInt(l lexer.Lexer) (lexer.Lexer, *IntExpr) {

	if new, tok := l.Next(); tok.Type == lexer.INT {
		value, err := strconv.Atoi(tok.Lexeme)
		if err != nil {
			return l, nil
		}
		return new, &IntExpr{value}
	} else {
		return l, nil
	}
}

func parseExpression(l lexer.Lexer) (lexer.Lexer, Expression) {

	l, ident := parseIdentifier(l)
	if ident != nil {
		return l, ident
	}

	l, integer := parseInt(l)
	if integer != nil {
		return l, integer
	}

	return l, nil
}

func allOf(l lexer.Lexer, types ...lexer.TokenType) (lexer.Lexer, []lexer.Token) {
	var tokens []lexer.Token
	for _, t := range types {
		if new, tok := l.Next(); tok.Type == t {
			tokens = append(tokens, tok)
			l = new
		} else {
			return l, nil
		}
	}
	return l, tokens
}

func parseAssignment(l lexer.Lexer) (lexer.Lexer, Statement) {

	if new, toks := allOf(l, lexer.LET, lexer.IDENT, lexer.ASSIGN_T, lexer.IDENT, lexer.ASSIGN); toks != nil {
		if new, rhs := parseExpression(new); rhs != nil {
			return new, &AssignStmt{Lhs: toks[1].Lexeme, Tipe: toks[3].Lexeme, Rhs: rhs}
		} else {
			return l, nil
		}
	} else {
		return l, nil
	}
}
