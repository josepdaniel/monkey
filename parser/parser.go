package parser

import (
	"errors"
	"fmt"
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

// An expression is a start followed by an end
// expr := start, end
func parseExpression(l lexer.Lexer) (lexer.Lexer, Expression) {

	l, start := parseExpressionStart(l)
	if start == nil {
		return l, nil
	}

	l, end := parseExpressionEnd(l, start)
	if end != nil {
		return l, end
	}

	return l, start
}

// A start is a simple, non-recursive expression
// start := ident | int | Nothing
func parseExpressionStart(l lexer.Lexer) (lexer.Lexer, Expression) {
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

// An end is an expression that is recursive.
// end := add | sub | Nothing
func parseExpressionEnd(l lexer.Lexer, start Expression) (lexer.Lexer, Expression) {
	l, add := parseAddition(l, start)
	if add != nil {
		return l, add
	}

	l, sub := parseSubtraction(l, start)
	if sub != nil {
		return l, sub
	}

	return l, nil
}

// add := "+", expr
func parseAddition(l lexer.Lexer, lhs Expression) (lexer.Lexer, *AddExpr) {
	new, plus := l.Next()
	if plus.Type != lexer.PLUS {
		return l, nil
	}

	new, rhs := parseExpression(new)
	if rhs == nil {
		return l, nil
	}
	return new, &AddExpr{lhs, rhs}
}

// sub := "-", expr
func parseSubtraction(l lexer.Lexer, lhs Expression) (lexer.Lexer, *SubExpr) {
	new, minus := l.Next()
	if minus.Type != lexer.MINUS {
		return l, nil
	}

	new, rhs := parseExpression(new)
	if rhs == nil {
		return l, nil
	}
	return new, &SubExpr{lhs, rhs}
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

func ParseStatement(l lexer.Lexer) (lexer.Lexer, Statement, error) {
	l, assign := parseAssignment(l)
	if assign != nil {
		return l, assign, nil
	}

	errorMsg := fmt.Sprintln("Expected statement at: ", l.Position)
	return l, nil, errors.New(errorMsg)
}
