package interpreter

import (
	"fmt"
	"golox/internal/parser"
	"golox/internal/scanner"
)

type _functionType struct{}

var FunctionType = _functionType{}

func (ft _functionType) None() int {
	return 0
}

func (ft _functionType) Function() int {
	return 1
}

func (ft _functionType) Method() int {
	return 2
}

func (ft _functionType) Initializer() int {
	return 4
}

type _classType struct{}

var ClassType = _classType{}

func (ct _classType) None() int {
	return 0
}

func (ct _classType) Class() int {
	return 1
}

func (ct _classType) Subclass() int {
	return 2
}

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction int
	currentClass    int
}

func NewResolver(i *Interpreter) *Resolver {
	r := new(Resolver)
	r.interpreter = i
	r.scopes = make([]map[string]bool, 0)
	r.currentFunction = FunctionType.None()
	r.currentClass = ClassType.None()
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

func (r Resolver) ResolveStmts(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r Resolver) resolveStmt(stmt parser.Stmt) {
	stmt.Accept(r)
}

func (r Resolver) resolveExpr(expr parser.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(stmt *parser.Function, type_ int) {
	enclosingFunctionType := r.currentFunction
	defer func() { r.currentFunction = enclosingFunctionType }()
	r.currentFunction = type_
	r.beginScope()
	for _, param := range stmt.Arguments {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStmts(stmt.Body.Stmts)
	r.endScope()
}

func (r *Resolver) define(name scanner.Token) {
	if scope := r.currentScope(); scope != nil {
		scope[name.Lexeme] = true
	}
}

func (r *Resolver) declare(name scanner.Token) {
	if scope := r.currentScope(); scope != nil {
		if _, found := scope[name.Lexeme]; found {
			r.error(fmt.Sprintf("[line %v] Error at '%v': variable already exist in this scope", name.Line, name.Lexeme))
		}
		scope[name.Lexeme] = false
	}
}

func (r *Resolver) resolveLocal(expr parser.Expr, name scanner.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, found := r.scopes[i][name.Lexeme]; found {
			depth := len(r.scopes) - 1 - i
			r.interpreter.resolve(expr, depth)
			return
		}
	}
}

func (r Resolver) VisitBlockStmt(stmt *parser.Block) {
	r.beginScope()
	r.ResolveStmts(stmt.Stmts)
	r.endScope()
}

func (r Resolver) VisitVarStmt(stmt *parser.Var) {
	r.declare(stmt.VarName)
	if stmt.VarValue != nil {
		r.resolveExpr(stmt.VarValue)
	}
	r.define(stmt.VarName)
}

func (r Resolver) VisitClassStmt(stmt *parser.Class) {
	enclosingClass := r.currentClass
	r.currentClass = ClassType.Class()
	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil {
		r.currentClass = ClassType.Subclass()
		if stmt.Superclass.Name.Lexeme == stmt.Name.Lexeme {
			r.error("Can't inherit from itself")
		}
		r.resolveExpr(stmt.Superclass)
		r.beginScope()
		defer func() { r.endScope() }()
		r.currentScope()["super"] = true
	}

	r.beginScope()
	r.currentScope()["this"] = true

	for _, method := range stmt.Methods {
		if method.Name.Lexeme == "init" {
			r.resolveFunction(method, FunctionType.Initializer())
		} else {
			r.resolveFunction(method, FunctionType.Method())
		}
	}

	r.endScope()

	r.currentClass = enclosingClass
}

func (r Resolver) VisitFunctionStmt(stmt *parser.Function) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FunctionType.Function())
}

func (r Resolver) VisitExpressionStmt(stmt *parser.Expression) {
	r.resolveExpr(stmt.Expr)
}

func (r Resolver) VisitIfStmt(stmt *parser.If) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
}

func (r Resolver) VisitPrintStmt(stmt *parser.Print) {
	r.resolveExpr(stmt.Expr)
}

func (r Resolver) VisitReturnStmt(stmt *parser.Return) {
	if r.currentFunction == FunctionType.None() {
		r.error(fmt.Sprintf("[line %v] Error at '%v': Can't return from top-level code", stmt.RetKeyWord.Line, stmt.RetKeyWord.Lexeme))
	}
	if stmt.Value != nil {
		if r.currentFunction == FunctionType.Initializer() {
			r.error(fmt.Sprintf("[line %v] Error at '%v': Can't return a value from initializer", stmt.RetKeyWord.Line, stmt.RetKeyWord.Lexeme))
		}
		r.resolveExpr(stmt.Value)
	}
}

func (r Resolver) VisitWhileStmt(stmt *parser.While) {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
}

func (r Resolver) VisitVarExpr(expr *parser.VarExpr) any {
	if scope := r.currentScope(); scope != nil {
		if val, exists := scope[expr.Name.Lexeme]; exists && !val {
			r.error(fmt.Sprintf("[Line %d] Error at '%v': Can't read local variable in its own initializer", expr.Name.Line, expr.Name.Lexeme))
		}
	}
	r.resolveLocal(expr, expr.Name)

	return nil
}

func (r Resolver) VisitAssignExpr(expr *parser.AssignExpr) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r Resolver) VisitBinaryExpr(expr *parser.BinaryExpr) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r Resolver) VisitCallExpr(expr *parser.CallExpr) any {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Args {
		r.resolveExpr(arg)
	}

	return nil
}

func (r Resolver) VisitSubscriptExpr(expr *parser.SubscriptExpr) any {
	r.resolveExpr(expr.Object)
	r.resolveExpr(expr.Index)
	return nil
}

func (r Resolver) VisitGetExpr(expr *parser.GetExpr) any {
	r.resolveExpr(expr.Object)
	return nil
}

func (r Resolver) VisitSetExpr(expr *parser.SetExpr) any {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}

func (r Resolver) VisitSuperExpr(expr *parser.SuperExpr) any {
	if r.currentClass == ClassType.None() {
		r.error(fmt.Sprintf("[line %v] Can't use 'super' outside of class", expr.Method.Line))
	}
	if r.currentClass != ClassType.Subclass() {
		r.error(fmt.Sprintf("[line %v] Can't use 'super' in a class with no superclass", expr.Method.Line))
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r Resolver) VisitThisExpr(expr *parser.ThisExpr) any {
	if r.currentClass == ClassType.None() {
		r.error(fmt.Sprintf("[line %v] Can't use 'this' outside of a class", expr.Keyword.Line))
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r Resolver) VisitGroupingExpr(expr *parser.GroupingExpr) any {
	r.resolveExpr(expr.Expr)
	return nil
}

func (r Resolver) VisitLiteralExpr(expr *parser.LiteralExpr) any {
	return nil
}

func (r Resolver) VisitArrayDeclExpr(expr *parser.ArrayDeclExpr) any {
	for _, element := range expr.Elements {
		r.resolveExpr(element)
	}
	return nil
}

func (r Resolver) VisitLogicalExpr(expr *parser.LogicalExpr) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r Resolver) VisitUnaryExpr(expr *parser.UnaryExpr) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r Resolver) error(msg string) {
	panic(msg)
}
