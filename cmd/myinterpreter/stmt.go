package main

type stmtVisitor interface {
	visitPrintStmt(stmt Print)
	visitExpressionStmt(stmt Expression)
	visitVarStmt(stmt Var)
	visitBlockStmt(stmt Block)
    visitIfStmt(stmt If)
    visitWhileStmt(stmt While)
    visitFunctionStmt(stmt Function)
}

type Stmt interface {
	accept(stmtVisitor)
}

type Print struct {
	expr Expr
}

func NewPrint(expr Expr) *Print {
	p := new(Print)
	p.expr = expr
	return p
}

func (p Print) accept(vis stmtVisitor) {
	vis.visitPrintStmt(p)
}

type Expression struct {
	expr Expr
}

func NewExpression(expr Expr) *Expression {
	e := new(Expression)
	e.expr = expr
	return e
}

func (e Expression) accept(vis stmtVisitor) {
	vis.visitExpressionStmt(e)
}

type Var struct {
	varName  Token
	varValue Expr
}

func NewVar(varName Token, varValue Expr) *Var {
	v := new(Var)
	v.varName = varName
	v.varValue = varValue
	return v
}

func (v Var) accept(vis stmtVisitor) {
	vis.visitVarStmt(v)
}

type Block struct {
	stmts []Stmt
}

func NewBlock(stmts []Stmt) *Block {
	b := new(Block)
    b.stmts = stmts
	return b
}

func (b Block) accept(vis stmtVisitor) {
	vis.visitBlockStmt(b)
}

type If struct {
    condition Expr
    thenBranch Stmt
    elseBranch Stmt
}

func NewIf(condition Expr, thenBranch, elseBranch Stmt) *If {
    i := new(If)
    i.condition = condition
    i.thenBranch = thenBranch
    i.elseBranch = elseBranch
    return i
}

func (i If) accept(vis stmtVisitor){
    vis.visitIfStmt(i)
}

type While struct {
    condition Expr
    body Stmt
}

func NewWhile(condition Expr, body Stmt) *While {
    w := new(While)
    w.body = body
    w.condition = condition
    return w
}

func (w While) accept(vis stmtVisitor) {
    vis.visitWhileStmt(w)
}

type Function struct {
    name Token
    arguments []Token
    body Block
}

func NewFunction(name Token, arguments []Token, body Block) *Function{
    fn := new(Function)
    fn.name = name
    fn.arguments = arguments
    fn.body = body
    return fn
}

func (fn Function) accept(vis stmtVisitor){
    vis.visitFunctionStmt(fn)
}
