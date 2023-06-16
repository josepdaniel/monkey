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
// start := enclosedExpression | ident | int | bool | lambdaExpr | blockExpr | Nothing
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

	new, lambda := parseLambdaExpr(l)
	if lambda != nil {
		return new, lambda
	}

	new, block := parseBlockBody(l)
	if block != nil {
		return new, block
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

// assignment := "let", IDENT, ":", typeExpr, "=", expression
func parseAssignment(l lexer.Lexer) (lexer.Lexer, Statement) {

	if new, toks := allOf(l, lexer.LET, lexer.IDENT, lexer.ASSIGN_T); toks != nil {
		if new, tipe := parseTypeExpr(new); tipe != nil {
			if new, tok := new.Next(); tok.Type == lexer.ASSIGN {
				if new, rhs := parseExpression(new); rhs != nil {
					return new, &AssignStmt{Lhs: toks[1].Lexeme, Tipe: tipe, Rhs: rhs, Pos: l.Position}
				}
			}
		}
	}
	return l, nil
}

// typeExpr := literalType | arrowType
func parseTypeExpr(l lexer.Lexer) (lexer.Lexer, TypeExpression) {
	new, ident := parseLiteralType(l)
	if ident != nil {
		return new, ident
	}

	new, arrow := parseArrowType(l)
	if arrow != nil {
		return new, arrow
	}

	return l, nil
}

// literalType := [a-zA-Z]
func parseLiteralType(l lexer.Lexer) (lexer.Lexer, TypeExpression) {

	if new, ident := l.Next(); ident.Type == lexer.IDENT {
		return new, &LiteralType{ident.Lexeme}
	}

	return l, nil
}

// arrowType := "(", {typeExpr}, ")", "->", typeExpr
func parseArrowType(l lexer.Lexer) (lexer.Lexer, TypeExpression) {
	new, tok := l.Next()

	if tok.Type != lexer.LPAREN {
		return l, nil
	}

	var types []TypeExpression
	for {
		// If we reached an RPAREN, we're done
		if _, tok := new.Next(); tok.Type == lexer.RPAREN {
			new, _ = new.Next()
			break
		}

		// Expect a comma between each parameter
		if len(types) > 0 {
			new, tok = new.Next()
			if tok.Type != lexer.COMMA {
				return l, nil
			}
		}

		var tipe TypeExpression
		new, tipe = parseTypeExpr(new)

		if tipe == nil {
			return l, nil
		}
		types = append(types, tipe)
	}

	if new, tok = new.Next(); tok.Type != lexer.ARROW {
		return l, nil
	}

	if new, tipe := parseTypeExpr(new); tipe != nil {
		return new, &ArrowType{types, tipe}
	}

	return l, nil

}

// a block is many statements followed by a single expression
// block := '{' {statement} expression '}'
func parseBlockBody(l lexer.Lexer) (lexer.Lexer, *BlockBodyExpr) {
	var statements []Statement = make([]Statement, 0)
	new, tok := l.Next()

	if tok.Type != lexer.LBRACE {
		return l, nil
	}

	for {
		newer, stmt, err := ParseStatement(new)
		if err != nil {
			break
		}
		new = newer
		statements = append(statements, stmt)
	}

	new, expr := parseExpression(new)
	if expr == nil {
		return l, nil
	}
	new, rbrace := new.Next()
	if rbrace.Type != lexer.RBRACE {
		return l, nil
	}
	return new, &BlockBodyExpr{
		Statements: statements,
		Final:      expr,
	}
}

func parseFuncParams(l lexer.Lexer) (lexer.Lexer, []FunctionParameter) {
	new, tok := l.Next()
	if tok.Type != lexer.LPAREN {
		return l, nil
	}

	var action func(lexer.Lexer, []FunctionParameter) (lexer.Lexer, []FunctionParameter)

	action = func(l lexer.Lexer, params []FunctionParameter) (lexer.Lexer, []FunctionParameter) {
		new, tok := l.Next()
		switch tok.Type {
		case lexer.RPAREN:
			return new, params
		case lexer.COMMA:
			return action(new, params)
		case lexer.IDENT:
			new, assign_t := new.Next()
			if assign_t.Type != lexer.ASSIGN_T {
				return l, nil
			}
			new, tipe := parseTypeExpr(new)
			if tipe == nil {
				return l, nil
			}
			params = append(params, FunctionParameter{IdentExpr{tok.Lexeme}, tipe})
			return action(new, params)
		default:
			return l, nil
		}
	}
	var params []FunctionParameter
	new, params = action(new, params)
	return new, params
}

// lambda := "def" "(" {func_param} ")" "->" ident block
func parseLambdaExpr(l lexer.Lexer) (lexer.Lexer, Expression) {
	new, toks := allOf(l, lexer.FUNCTION)
	if toks == nil {
		return l, nil
	}
	new, params := parseFuncParams(new)
	if params == nil {
		return l, nil
	}
	new, tok := new.Next()
	if tok.Type != lexer.ARROW {
		return l, nil
	}
	new, tipe := parseTypeExpr(new)
	if tipe == nil {
		return l, nil
	}
	new, block := parseBlockBody(new)
	if block == nil {
		return l, nil
	}
	return new, &LambdaExpr{
		Parameters: params,
		Returns:    tipe,
		Body:       *block,
	}
}

func ParseStatement(l lexer.Lexer) (lexer.Lexer, Statement, error) {
	l, assign := parseAssignment(l)
	if assign != nil {
		return l, assign, nil
	}

	errorMsg := fmt.Sprintln("unexpected statement type")
	return l, nil, errors.New(errorMsg)
}

type ParseError struct {
	Position int
	Error    error
}

func (e ParseError) ToError(source string) error {
	return errors.New(fmt.Sprint("[parse err]: ", e.Error, "\nCulprit:\n>>> ", source))
}

func ParseProgram(l lexer.Lexer) (*Program, *ParseError) {
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
			return nil, &ParseError{
				Position: l.Position,
				Error:    err,
			}
		}
		program.Statements = append(program.Statements, stmt)
	}
}
