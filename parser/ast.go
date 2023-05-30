package parser

type (
	Node interface {
		isNode()
	}

	// Expressions produce values
	Expression interface {
		Node
		isExpression()
	}

	// Statements do not produce values
	Statement interface {
		Node
		isStatement()
	}

	Program struct {
		Statements []Statement
	}
)

// Identifier -------------------------
type IdentExpr struct {
	Name string
}

func (*IdentExpr) isExpression() {}
func (*IdentExpr) isNode()       {}

// Integer ----------------------------
type IntExpr struct {
	Value int
}

func (*IntExpr) isExpression() {}
func (*IntExpr) isNode()       {}

// Addition ---------------------------

type AddExpr struct {
	Lhs Expression
	Rhs Expression
}

func (*AddExpr) isExpression() {}
func (*AddExpr) isNode()       {}

// Subtraction ------------------------

type SubExpr struct {
	Lhs Expression
	Rhs Expression
}

func (*SubExpr) isExpression() {}
func (*SubExpr) isNode()       {}

// Assignment -------------------------
type AssignStmt struct {
	Lhs  string
	Tipe string
	Rhs  Expression
}

func (*AssignStmt) isStatement() {}
func (*AssignStmt) isNode()      {}
