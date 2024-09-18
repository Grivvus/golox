package main

import (

)

type stmtVisitor interface {
    visitPrintStmt(stmt Print)
    visitExpressionStmt(stmt Expression)
    visitVarStmt(stmt Var)
}

type Stmt interface {
    accept(stmtVisitor)
}

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

type Var struct {
    varName Token
    varValue Expr
}

func NewVar(varName Token, varValue Expr) *Var {
    v := new(Var)
    v.varName = varName
    v.varValue = varValue
    return v
}

func (v Var) accept(vis stmtVisitor){
    vis.visitVarStmt(v)
}
