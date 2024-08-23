package main

type Expr interface {
	accept(v visitor) string
}

type UnaryExpr struct {
	operator Token
	right    Expr
}

func NewUnaryExpr(operator Token, right Expr) *UnaryExpr {
	unary := new(UnaryExpr)
	unary.right = right
	unary.operator = operator
	return unary
}

func (unary *UnaryExpr) accept(v visitor) string {
	return v.visitUnaryExpr(*unary)
}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func NewBinaryExpr(left Expr, operator Token, right Expr) *BinaryExpr {
	bin := new(BinaryExpr)
	bin.left = left
	bin.right = right
	bin.operator = operator
	return bin
}

func (binary *BinaryExpr) accept(v visitor) string {
	return v.visitBinaryExpr(*binary)
}

type GroupingExpr struct {
	expr Expr
}

func NewGroupingExpr(expr Expr) *GroupingExpr {
	g := new(GroupingExpr)
	g.expr = expr
	return g
}

func (grouping *GroupingExpr) accept(v visitor) string {
	return v.visitGroupingExpr(*grouping)
}

type LiteralExpr struct {
	value any
}

func NewLiteralExpr(value any) *LiteralExpr {
	lit := new(LiteralExpr)
	lit.value = value
	return lit
}

func (literal *LiteralExpr) accept(v visitor) string {
	return v.visitLiteralExpr(*literal)
}
