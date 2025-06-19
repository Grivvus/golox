package main

import (
	"fmt"
	"os"
)

type Resolver struct {
	interpreter *Interpreter
	scopes      []map[string]bool
}

func NewResolver(i *Interpreter) *Resolver {
	r := new(Resolver)
	r.interpreter = i
	r.scopes = make([]map[string]bool, 0)
	return r
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) currentScope() map[string]bool {
	if len(r.scopes) == 0 {
		return nil
	}
	return r.scopes[len(r.scopes)-1]
}

func (r Resolver) resolveStmts(stmts []Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r Resolver) resolveStmt(stmt Stmt) {
	stmt.accept(r)
}

func (r Resolver) resolveExpr(expr Expr) {
	expr.accept(r)
}

func (r *Resolver) resolveFunction(stmt *Function) {
	r.beginScope()
	for _, param := range stmt.arguments {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(stmt.body.stmts)
	r.endScope()
}

func (r *Resolver) define(name Token) {
	if scope := r.currentScope(); scope != nil {
		scope[name.Lexeme] = true
	}
}

func (r *Resolver) declare(name Token) {
	if scope := r.currentScope(); scope != nil {
		if _, found := scope[name.Lexeme]; found {
			fmt.Fprintf(os.Stderr, "variable %v already exist in this scope, [line %v]", name.Lexeme, name.line)
			os.Exit(70)
		}
		scope[name.Lexeme] = false
	}
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, found := r.scopes[i][name.Lexeme]; found {
			depth := len(r.scopes) - 1 - i
			r.interpreter.resolve(expr, depth)
			return
		}
	}
}

func (r Resolver) visitBlockStmt(stmt *Block) {
	// fmt.Println("Entering block scope")
	r.beginScope()
	r.resolveStmts(stmt.stmts)
	r.endScope()
	// fmt.Println("Exiting block scope")
}

func (r Resolver) visitVarStmt(stmt *Var) {
	// fmt.Printf("Declaring variable: %s\n", stmt.varName.Lexeme)
	r.declare(stmt.varName)
	if stmt.varValue != nil {
		r.resolveExpr(stmt.varValue)
	}
	r.define(stmt.varName)
}

func (r Resolver) visitFunctionStmt(stmt *Function) {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt)
}

func (r Resolver) visitExpressionStmt(stmt *Expression) {
	r.resolveExpr(stmt.expr)
}

func (r Resolver) visitIfStmt(stmt *If) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(stmt.elseBranch)
	}
}

func (r Resolver) visitPrintStmt(stmt *Print) {
	r.resolveExpr(stmt.expr)
}

func (r Resolver) visitReturnStmt(stmt *Return) {
	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
}

func (r Resolver) visitWhileStmt(stmt *While) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
}

func (r Resolver) visitVarExpr(expr *VarExpr) any {
	if len(r.scopes) == 0 {
		defined, exist := r.currentScope()[expr.name.Lexeme]
		if exist && !defined {
			r.error(fmt.Sprintf("Can't read local variable %v from it's own initializer. line %v", expr.name.Lexeme, expr.name.line))
		}
	}
	r.resolveLocal(expr, expr.name)

	return nil
}

func (r Resolver) visitAssignExpr(expr *AssignExpr) any {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r Resolver) visitBinaryExpr(expr *BinaryExpr) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r Resolver) visitCallExpr(expr *CallExpr) any {
	r.resolveExpr(expr.callee)

	for _, arg := range expr.args {
		r.resolveExpr(arg)
	}

	return nil
}

func (r Resolver) visitGroupingExpr(expr *GroupingExpr) any {
	r.resolveExpr(expr.expr)
	return nil
}

func (r Resolver) visitLiteralExpr(expr *LiteralExpr) any {
	return nil
}

func (r Resolver) visitLogicalExpr(expr *LogicalExpr) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r Resolver) visitUnaryExpr(expr *UnaryExpr) any {
	r.resolveExpr(expr.right)
	return nil
}

func (r Resolver) error(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(70)
}
