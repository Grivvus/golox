package main

import (

)

type stmtVisitor interface {
    visitPrintStmt(stmt Print)
    visitExpressionStmt(stmt Expression)
}

type Stmt interface {
    accept(stmtVisitor)
}

// type Block[T any] struct {
//     statements []Stmt[T]
// }
//
// func NewBlock[T any](statements []Stmt[T]) *Block[T]{
//     stmt := new(Block[T])
//     stmt.statements = statements
//     return stmt
// }
//
// func (b *Block[T]) accept(vis stmtVisitor[T]) T {
//     return vis.visitBlockStmt(*b)
// }

type Print struct {
    expr Expr
}

func NewPrint (expr Expr) *Print {
    p := new(Print)
    p.expr = expr
    return p
}

func (p Print) accept(vis stmtVisitor){
    vis.visitPrintStmt(p)
}

type Expression struct {
    expr Expr
}

func NewExpression (expr Expr) *Expression {
    e := new(Expression)
    e.expr = expr
    return e
}

func (e Expression) accept(vis stmtVisitor){
    vis.visitExpressionStmt(e)
}
