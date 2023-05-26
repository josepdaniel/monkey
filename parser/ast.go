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
		Nodes []Node
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

// Assignment -------------------------
type AssignStmt struct {
	Lhs  string
	Tipe string
	Rhs  Expression
}

func (*AssignStmt) isStatement() {}
func (*AssignStmt) isNode()      {}
