package main

type Interpreter struct {}

func NewInterpreter() *Interpreter{
    return new(Interpreter)
}

func (i Interpreter) visitLiteralExpr(expr LiteralExpr) any{
    return expr.value
}

func (i Interpreter) visitGroupingExpr(expr GroupingExpr) any {
    return i.evaluate(expr.expr)
}

func (i Interpreter) visitUnaryExpr(expr UnaryExpr) any {
    return 1
}

func (i Interpreter) visitBinaryExpr(expr BinaryExpr) any {
    return 1
}

func (i Interpreter) evaluate(expr Expr) any {
    return expr.accept(i)
}
