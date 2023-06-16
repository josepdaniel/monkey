package parser

import "strings"

type (

	// Expressions produce values
	Expression interface {
		isExpression()
	}

	// Statements do not produce values
	Statement interface {
		isStatement()
		Position() int
	}

	Program struct {
		Statements []Statement
	}

	// Type expressions produce types
	TypeExpression interface {
		isTypeExpression()
		Render() string
	}
)

// Identifier -------------------------
type IdentExpr struct {
	Name string
}

func (*IdentExpr) isExpression() {}

// Integer ----------------------------
type IntExpr struct {
	Value int
}

func (*IntExpr) isExpression() {}

// Boolean ----------------------------
type BoolExpr struct {
	Value bool
}

func (*BoolExpr) isExpression() {}

// Addition ---------------------------

type AddExpr struct {
	Lhs Expression
	Rhs Expression
}

func (*AddExpr) isExpression() {}

// Subtraction ------------------------

type SubExpr struct {
	Lhs Expression
	Rhs Expression
}

func (*SubExpr) isExpression() {}

// Less than --------------------------
type LessThanExpr struct {
	Lhs Expression
	Rhs Expression
}

func (*LessThanExpr) isExpression() {}

// Greater than -----------------------
type GreaterThanExpr struct {
	Lhs Expression
	Rhs Expression
}

func (*GreaterThanExpr) isExpression() {}

type FunctionParameter struct {
	Name IdentExpr
	Tipe TypeExpression
}

type BlockBodyExpr struct {
	Statements []Statement
	Final      Expression
}

func (*BlockBodyExpr) isExpression() {}

type LambdaExpr struct {
	Parameters []FunctionParameter
	Returns    TypeExpression
	Body       BlockBodyExpr
}

func (*LambdaExpr) isExpression() {}

// Assignment -------------------------
type AssignStmt struct {
	Lhs  string
	Tipe TypeExpression
	Rhs  Expression
	Pos  int
}

func (*AssignStmt) isStatement() {}
func (s *AssignStmt) Position() int {
	return s.Pos
}

// Literal types ----------------------
type LiteralType struct {
	Name string
}

func (*LiteralType) isTypeExpression() {}
func (l *LiteralType) Render() string {
	return l.Name
}

// Function types ---------------------
type ArrowType struct {
	Parameters []TypeExpression
	Returns    TypeExpression
}

func (*ArrowType) isTypeExpression() {}
func (a *ArrowType) Render() string {
	var s = strings.Builder{}
	s.WriteString("(")
	for i, te := range a.Parameters {
		s.WriteString(te.Render())
		if i < len(a.Parameters)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString(")")
	s.WriteString(" -> ")
	s.WriteString(a.Returns.Render())
	return s.String()
}
