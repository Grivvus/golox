package main 

import(
    "strings"
    "fmt"
)

type visitor interface{
    visitUnaryExpr(UnaryExpr) string
    visitBinaryExpr(BinaryExpr) string
    visitGroupExpr(GroupingExpr) string
    visitLiteralExpr(LiteralExpr) string
}

type astPrinter struct{

}

func NewPrinter() *astPrinter{
    return new(astPrinter)
}

func (printer *astPrinter) print(e Expr){
    fmt.Println(e.accept(printer))
}

func (printer *astPrinter) parenthesize(name string, exprs... Expr) string{
    builder := strings.Builder{}
    builder.WriteString("(")
    builder.WriteString(name)
    for _, expr := range exprs{
        builder.WriteString(" ")
        builder.WriteString(expr.accept(printer))
    }
    builder.WriteString(")")
    return builder.String()
}

func (printer astPrinter) visitUnaryExpr(expr UnaryExpr) string{
    return printer.parenthesize(expr.operator.Lexeme, expr.right)
}

func (printer astPrinter) visitBinaryExpr(expr BinaryExpr) string{
    return printer.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (printer astPrinter) visitGroupExpr(expr GroupingExpr) string{
    return printer.parenthesize("group", expr.expr)
}

func (printer astPrinter) visitLiteralExpr(expr LiteralExpr) string{
    if expr.value == nil{
        return "nil"
    }
    switch v := expr.value.(type) {
    case float64:
        if expr.value == float64(int(v)){
            return fmt.Sprintf("%.1f", expr.value)
        }
    }
    return fmt.Sprint(expr.value) 
}
