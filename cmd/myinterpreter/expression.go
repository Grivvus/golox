package main

type visitor[T string | any] interface {
	visitUnaryExpr(UnaryExpr) T
	visitBinaryExpr(BinaryExpr) T
	visitGroupingExpr(GroupingExpr) T
	visitLiteralExpr(LiteralExpr) T
	visitVarExpr(VarExpr) T
    visitAssignExpr(AssignExpr) T
    visitLogicalExpr(LogicalExpr) T
}

type Expr interface {
	print(v visitor[string]) string
	accept(v visitor[any]) any
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

func (unary *UnaryExpr) print(v visitor[string]) string {
	return v.visitUnaryExpr(*unary)
}

func (unary *UnaryExpr) accept(v visitor[any]) any {
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

func (binary *BinaryExpr) print(v visitor[string]) string {
	return v.visitBinaryExpr(*binary)
}

func (binary *BinaryExpr) accept(v visitor[any]) any {
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

func (grouping *GroupingExpr) print(v visitor[string]) string {
	return v.visitGroupingExpr(*grouping)
}

func (grouping *GroupingExpr) accept(v visitor[any]) any {
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

func (literal *LiteralExpr) print(v visitor[string]) string {
	return v.visitLiteralExpr(*literal)
}

func (literal *LiteralExpr) accept(v visitor[any]) any {
	return v.visitLiteralExpr(*literal)
}

type VarExpr struct {
	name Token
}

func NewVarExpr(name Token) *VarExpr {
	e := new(VarExpr)
	e.name = name
	return e
}

func (var_ *VarExpr) print(v visitor[string]) string {
	return v.visitVarExpr(*var_)
}

func (var_ *VarExpr) accept(v visitor[any]) any {
	return v.visitVarExpr(*var_)
}

type AssignExpr struct {
    name Token
    value Expr
}

func NewAssignExpr(name Token, value Expr) *AssignExpr{
    a := new(AssignExpr)
    a.name = name
    a.value = value
    return a
}

func (assign *AssignExpr) print(v visitor[string]) string {
    return v.visitAssignExpr(*assign)
}

func (assign *AssignExpr) accept(v visitor[any]) any {
    return v.visitAssignExpr(*assign)
}

type LogicalExpr struct {
    left Expr
    operator Token
    right Expr
}

func NewLogicalExpr(left Expr, operator Token, right Expr) *LogicalExpr {
    l := new(LogicalExpr)
    l.left = left
    l.right = right
    l.operator = operator
    return l
}

func (logical *LogicalExpr) accept(v visitor[any]) any {
    return v.visitLogicalExpr(*logical)
}

func (logical *LogicalExpr) print(v visitor[string]) string {
    return v.visitLogicalExpr(*logical)
}
