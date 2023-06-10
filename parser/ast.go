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

// Boolean ----------------------------
type BoolExpr struct {
	Value bool
}

func (*BoolExpr) isExpression() {}
func (*BoolExpr) isNode()       {}

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

// Less than --------------------------
type LessThanExpr struct {
	Lhs Expression
	Rhs Expression
}

func (*LessThanExpr) isExpression() {}
func (*LessThanExpr) isNode()       {}

// Greater than -----------------------
type GreaterThanExpr struct {
	Lhs Expression
	Rhs Expression
}

func (*GreaterThanExpr) isExpression() {}
func (*GreaterThanExpr) isNode()       {}

// Assignment -------------------------
type AssignStmt struct {
	Lhs  string
	Tipe string
	Rhs  Expression
}

func (*AssignStmt) isStatement() {}
func (*AssignStmt) isNode()      {}
