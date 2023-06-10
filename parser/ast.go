package parser

type (
	Node interface {
		Position() int
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
	Pos  int
}

func (*IdentExpr) isExpression() {}
func (*IdentExpr) isNode()       {}
func (e *IdentExpr) Position() int {
	return e.Pos
}

// Integer ----------------------------
type IntExpr struct {
	Value int
	Pos   int
}

func (*IntExpr) isExpression() {}
func (*IntExpr) isNode()       {}
func (e *IntExpr) Position() int {
	return e.Pos
}

// Boolean ----------------------------
type BoolExpr struct {
	Value bool
	Pos   int
}

func (*BoolExpr) isExpression() {}
func (*BoolExpr) isNode()       {}
func (e *BoolExpr) Position() int {
	return e.Pos
}

// Addition ---------------------------

type AddExpr struct {
	Lhs Expression
	Rhs Expression
	Pos int
}

func (*AddExpr) isExpression() {}
func (*AddExpr) isNode()       {}
func (e *AddExpr) Position() int {
	return e.Pos
}

// Subtraction ------------------------

type SubExpr struct {
	Lhs Expression
	Rhs Expression
	Pos int
}

func (*SubExpr) isExpression() {}
func (*SubExpr) isNode()       {}
func (e *SubExpr) Position() int {
	return e.Pos
}

// Less than --------------------------
type LessThanExpr struct {
	Lhs Expression
	Rhs Expression
	Pos int
}

func (*LessThanExpr) isExpression() {}
func (*LessThanExpr) isNode()       {}
func (e *LessThanExpr) Position() int {
	return e.Pos
}

// Greater than -----------------------
type GreaterThanExpr struct {
	Lhs Expression
	Rhs Expression
	Pos int
}

func (*GreaterThanExpr) isExpression() {}
func (*GreaterThanExpr) isNode()       {}
func (e *GreaterThanExpr) Position() int {
	return e.Pos
}

// Assignment -------------------------
type AssignStmt struct {
	Lhs  string
	Tipe string
	Rhs  Expression
	Pos  int
}

func (*AssignStmt) isStatement() {}
func (*AssignStmt) isNode()      {}
func (s *AssignStmt) Position() int {
	return s.Pos
}
