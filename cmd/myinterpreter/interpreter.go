package main

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return new(Interpreter)
}

func (i Interpreter) visitLiteralExpr(expr LiteralExpr) any {
	return expr.value
}

func (i Interpreter) visitGroupingExpr(expr GroupingExpr) any {
	return i.evaluate(expr.expr)
}

func (i Interpreter) visitUnaryExpr(expr UnaryExpr) any {
	right := i.evaluate(expr.right)

	switch expr.operator.Token {
	case BANG:
		return !booleanCast(right)
	case MINUS:
		return -(right.(float64))
	}

    return nil
}

func (i Interpreter) visitBinaryExpr(expr BinaryExpr) any {
    left := i.evaluate(expr.left)
    right := i.evaluate(expr.right)

    switch expr.operator.Token {
    case STAR:
        return left.(float64) * right.(float64)
    case SLASH:
        return left.(float64) / right.(float64)
    case PLUS:
        return left.(float64) + right.(float64)
    case MINUS:
        return left.(float64) - right.(float64)
    }

    return nil
}

func (i Interpreter) evaluate(expr Expr) any {
	return expr.accept(i)
}

func booleanCast(expr any) bool {
	if expr == nil {
		return false
	}
	switch expr.(type) {
	case bool:
		return expr.(bool)
	default:
		return true
	}
}
