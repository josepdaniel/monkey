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

func parseBool(l lexer.Lexer) (lexer.Lexer, *BoolExpr) {

	if new, tok := l.Next(); tok.Type == lexer.TRUE || tok.Type == lexer.FALSE {
		value := tok.Type == lexer.TRUE
		return new, &BoolExpr{value}
	} else {
		return l, nil
	}
}

// An expression is a start followed by an end
// expr := start, end
func parseExpression(l lexer.Lexer) (lexer.Lexer, Expression) {

	new, tree := parseExpressionStart(l)
	if tree == nil {
		return l, nil
	}

	new, tree = parseExpressionEnd(new, tree)

	return new, tree
}

// A start is a simple, non-recursive expression
// start := enclosedExpression | ident | int | bool | Nothing
func parseExpressionStart(l lexer.Lexer) (lexer.Lexer, Expression) {

	new, tree := parseEnclosedExpression(l)
	if tree != nil {
		return new, tree
	}

	new, ident := parseIdentifier(l)
	if ident != nil {
		return new, ident
	}

	new, integer := parseInt(l)
	if integer != nil {
		return new, integer
	}

	new, bool := parseBool(l)
	if bool != nil {
		return new, bool
	}

	return l, nil
}

// An end is zero or more recursive expressions
// end := {add | sub | lt | gt}
func parseExpressionEnd(l lexer.Lexer, start Expression) (lexer.Lexer, Expression) {
	new, add := parseAddition(l, start)
	if add != nil {
		return parseExpressionEnd(new, add)
	}

	new, sub := parseSubtraction(l, start)
	if sub != nil {
		return parseExpressionEnd(new, sub)
	}

	new, lt := parseLessThan(l, start)
	if lt != nil {
		return parseExpressionEnd(new, lt)
	}

	new, gt := parseGreaterThan(l, start)
	if gt != nil {
		return parseExpressionEnd(new, gt)
	}

	return l, start
}

// parseEnclosedExpression := '(' EXPRESSION ')'
func parseEnclosedExpression(l lexer.Lexer) (lexer.Lexer, Expression) {
	new, openParens := l.Next()

	if openParens.Type != lexer.LPAREN {
		return l, nil
	}

	new, expression := parseExpression(new)
	if expression == nil {
		return l, nil
	}

	new, closeParens := new.Next()
	if closeParens.Type != lexer.RPAREN {
		return l, nil
	}

	return new, expression

}

func parseInfix(l lexer.Lexer, lhs Expression, expectOp lexer.TokenType, buildExp func(Expression, Expression) Expression) (lexer.Lexer, Expression) {
	new, operator := l.Next()

	if operator.Type != expectOp {
		return l, nil
	}

	new, rhs := parseExpressionStart(new)
	if rhs == nil {
		return l, nil
	}
	return new, buildExp(lhs, rhs)
}

// add := "+", start
func parseAddition(l lexer.Lexer, lhs Expression) (lexer.Lexer, Expression) {
	return parseInfix(l, lhs, lexer.PLUS, func(lhs, rhs Expression) Expression {
		return &AddExpr{lhs, rhs}
	})
}

// sub := "-", start
func parseSubtraction(l lexer.Lexer, lhs Expression) (lexer.Lexer, Expression) {
	return parseInfix(l, lhs, lexer.MINUS, func(lhs, rhs Expression) Expression {
		return &SubExpr{lhs, rhs}
	})
}

// lt := "<", start
func parseLessThan(l lexer.Lexer, lhs Expression) (lexer.Lexer, Expression) {
	return parseInfix(l, lhs, lexer.LT, func(lhs, rhs Expression) Expression {
		return &LessThanExpr{lhs, rhs}
	})
}

// gt := ">", start
func parseGreaterThan(l lexer.Lexer, lhs Expression) (lexer.Lexer, Expression) {
	return parseInfix(l, lhs, lexer.GT, func(lhs, rhs Expression) Expression {
		return &GreaterThanExpr{lhs, rhs}
	})
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

	errorMsg := fmt.Sprintln("[parse err] unexpected statement type")
	return l, nil, errors.New(errorMsg)
}

func ParseProgram(l lexer.Lexer) (*Program, error) {
	var program Program
	for {

		// If we reached the end of the input, return
		if _, tok := l.Next(); tok.Type == lexer.EOF {
			return &program, nil
		}

		var stmt Statement
		var err error
		l, stmt, err = ParseStatement(l)
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
	}
}
