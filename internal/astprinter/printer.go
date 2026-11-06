package astprinter

import (
	"fmt"
	"golox/internal/parser"
	"strings"
)

type astPrinter struct{}

func New() *astPrinter {
	return new(astPrinter)
}

func (p *astPrinter) Print(e parser.Expr) {
	fmt.Println(e.Print(p))
}

func (p *astPrinter) parenthesize(name string, exprs ...parser.Expr) string {
	builder := strings.Builder{}
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Print(p))
	}
	builder.WriteString(")")
	return builder.String()
}

func (p astPrinter) VisitUnaryExpr(expr *parser.UnaryExpr) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p astPrinter) VisitBinaryExpr(expr *parser.BinaryExpr) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p astPrinter) VisitGroupingExpr(expr *parser.GroupingExpr) string {
	return p.parenthesize("group", expr.Expr)
}

func (p astPrinter) VisitLiteralExpr(expr *parser.LiteralExpr) string {
	if expr.Value == nil {
		return "nil"
	}
	switch v := expr.Value.(type) {
	case float64:
		if expr.Value == float64(int(v)) {
			return fmt.Sprintf("%.1f", expr.Value)
		}
	}
	return fmt.Sprint(expr.Value)
}

func (p astPrinter) VisitVarExpr(expr *parser.VarExpr) string {
	return fmt.Sprintf("var %v\n", expr.Name.Lexeme)
}

func (p astPrinter) VisitAssignExpr(expr *parser.AssignExpr) string {
	return fmt.Sprintf("var %v=%v\n", expr.Name, expr.Value)
}

func (p astPrinter) VisitLogicalExpr(expr *parser.LogicalExpr) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p astPrinter) VisitCallExpr(expr *parser.CallExpr) string {
	return "callable"
}

func (p astPrinter) VisitGetExpr(expr *parser.GetExpr) string {
	return fmt.Sprintf("get on %v\n", expr.Object.Print(p))
}

func (p astPrinter) VisitSetExpr(expr *parser.SetExpr) string {
	return fmt.Sprintf("set on %v, value=%v\n", expr.Object.Print(p), expr.Value.Print(p))
}

func (p astPrinter) VisitThisExpr(expr *parser.ThisExpr) string {
	return "this\n"
}

func (p astPrinter) VisitSuperExpr(expr *parser.SuperExpr) string {
	return "super\n"
}

func (p astPrinter) VisitArrayDeclExpr(expr *parser.ArrayDeclExpr) string {
	b := strings.Builder{}
	b.WriteString("[ ")
	for _, element := range expr.Elements {
		b.WriteString(element.Print(p))
		b.WriteString(" ")
	}
	b.WriteString("]")
	return b.String()
}

func (p astPrinter) VisitSubscriptExpr(expr *parser.SubscriptExpr) string {
	return fmt.Sprintf("subscript index %v", expr.Index.Print(p))
}
