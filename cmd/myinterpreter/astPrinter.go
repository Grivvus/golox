package main

import (
	"fmt"
	"strings"
)

type astPrinter struct{}

func NewPrinter() *astPrinter {
	return new(astPrinter)
}

func (printer *astPrinter) print(e Expr) {
	fmt.Println(e.print(printer))
}

func (printer *astPrinter) parenthesize(name string, exprs ...Expr) string {
	builder := strings.Builder{}
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.print(printer))
	}
	builder.WriteString(")")
	return builder.String()
}

func (printer astPrinter) visitUnaryExpr(expr *UnaryExpr) string {
	return printer.parenthesize(expr.operator.Lexeme, expr.right)
}

func (printer astPrinter) visitBinaryExpr(expr *BinaryExpr) string {
	return printer.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (printer astPrinter) visitGroupingExpr(expr *GroupingExpr) string {
	return printer.parenthesize("group", expr.expr)
}

func (printer astPrinter) visitLiteralExpr(expr *LiteralExpr) string {
	if expr.value == nil {
		return "nil"
	}
	switch v := expr.value.(type) {
	case float64:
		if expr.value == float64(int(v)) {
			return fmt.Sprintf("%.1f", expr.value)
		}
	}
	return fmt.Sprint(expr.value)
}

func (printer astPrinter) visitVarExpr(expr *VarExpr) string {
	return fmt.Sprintf("var %v\n", expr.name.Lexeme)
}

func (printer astPrinter) visitAssignExpr(expr *AssignExpr) string {
	return fmt.Sprintf("var %v=%v\n", expr.name, expr.value)
}

func (printer astPrinter) visitLogicalExpr(expr *LogicalExpr) string {
	return printer.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (printer astPrinter) visitCallExpr(expr *CallExpr) string {
	return "callable"
}

func (printer astPrinter) visitGetExpr(expr *GetExpr) string {
	return fmt.Sprintf("get on %v\n", expr.object.print(printer))
}

func (printer astPrinter) visitSetExpr(expr *SetExpr) string {
	return fmt.Sprintf("set on %v, value=%v\n", expr.object.print(printer), expr.value.print(printer))
}
