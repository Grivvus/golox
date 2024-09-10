package main

import (
    "reflect"
)

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
        left_type := reflect.TypeOf(left)
        right_type := reflect.TypeOf(right)
        if left_type.Kind() == reflect.String && right_type.Kind() == reflect.String {
            return left.(string) + right.(string)
        } else if left_type.Kind() == reflect.Float64 && right_type.Kind() == reflect.Float64 {
            return left.(float64) + right.(float64)
        }
	case MINUS:
		return left.(float64) - right.(float64)
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case EQUAL_EQUAL:
        left_type := reflect.TypeOf(left)
        right_type := reflect.TypeOf(right)
        if left_type.Kind() != right_type.Kind() {
            return false
        } else if left_type.Kind() == reflect.String {
            return left.(string) == right.(string)
        } else if left_type.Kind() == reflect.Float64 { 
            return left.(float64) == right.(float64)
        }

    case BANG_EQUAL:
        left_type := reflect.TypeOf(left)
        right_type := reflect.TypeOf(right)
        if left_type.Kind() != right_type.Kind() {
            return true
        } else if left_type.Kind() == reflect.String {
            return left.(string) != right.(string)
        } else if left_type.Kind() == reflect.Float64 { 
            return left.(float64) != right.(float64)
        }
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
