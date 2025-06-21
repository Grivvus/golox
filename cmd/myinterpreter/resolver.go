package main

import (
	"fmt"
	"os"
)

type functionType int

const (
	none functionType = iota
	function
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction functionType
}

func NewResolver(i *Interpreter) *Resolver {
	r := new(Resolver)
	r.interpreter = i
	r.scopes = make([]map[string]bool, 0)
	r.currentFunction = none
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

func (r *Resolver) resolveFunction(stmt *Function, type_ functionType) {
	enclosingFunctionType := r.currentFunction
	defer func() { r.currentFunction = enclosingFunctionType }()
	r.currentFunction = type_
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
			r.error(fmt.Sprintf("[line %v] Error at '%v': variable already exist in this scope", name.line, name.Lexeme))
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
	r.beginScope()
	r.resolveStmts(stmt.stmts)
	r.endScope()
}

func (r Resolver) visitVarStmt(stmt *Var) {
	r.declare(stmt.varName)
	if stmt.varValue != nil {
		r.resolveExpr(stmt.varValue)
	}
	r.define(stmt.varName)
}

func (r Resolver) visitClassStmt(stmt *Class) {
	r.declare(stmt.name)
	r.define(stmt.name)
}

func (r Resolver) visitFunctionStmt(stmt *Function) {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt, function)
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
	if r.currentFunction == none {
		r.error(fmt.Sprintf("[line %v] Error at '%v': Can't return from top-level code", stmt.retKeyWord.line, stmt.retKeyWord.Lexeme))
	}
	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
}

func (r Resolver) visitWhileStmt(stmt *While) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
}

func (r Resolver) visitVarExpr(expr *VarExpr) any {
	if scope := r.currentScope(); scope != nil {
		if val, exists := scope[expr.name.Lexeme]; exists && !val {
			r.error(fmt.Sprintf("[Line %d] Error at '%v': Can't read local variable in its own initializer", expr.name.line, expr.name.Lexeme))
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
	os.Exit(65)
}
