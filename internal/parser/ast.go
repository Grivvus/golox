package parser

import "golox/internal/scanner"

type Visitor[T string | any] interface {
	VisitUnaryExpr(*UnaryExpr) T
	VisitBinaryExpr(*BinaryExpr) T
	VisitGroupingExpr(*GroupingExpr) T
	VisitLiteralExpr(*LiteralExpr) T
	VisitVarExpr(*VarExpr) T
	VisitAssignExpr(*AssignExpr) T
	VisitLogicalExpr(*LogicalExpr) T
	VisitCallExpr(*CallExpr) T
	VisitGetExpr(*GetExpr) T
	VisitSetExpr(*SetExpr) T
	VisitThisExpr(*ThisExpr) T
	VisitSuperExpr(*SuperExpr) T
	VisitArrayDeclExpr(*ArrayDeclExpr) T
	VisitSubscriptExpr(*SubscriptExpr) T
}

type Expr interface {
	Print(v Visitor[string]) string
	Accept(v Visitor[any]) any
}

type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

func NewUnaryExpr(operator scanner.Token, right Expr) *UnaryExpr {
	return &UnaryExpr{
		Operator: operator,
		Right:    right,
	}
}

func (unary *UnaryExpr) Print(v Visitor[string]) string {
	return v.VisitUnaryExpr(unary)
}

func (unary *UnaryExpr) Accept(v Visitor[any]) any {
	return v.VisitUnaryExpr(unary)
}

type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func NewBinaryExpr(left Expr, operator scanner.Token, right Expr) *BinaryExpr {
	bin := new(BinaryExpr)
	bin.Left = left
	bin.Right = right
	bin.Operator = operator
	return bin
}

func (binary *BinaryExpr) Print(v Visitor[string]) string {
	return v.VisitBinaryExpr(binary)
}

func (binary *BinaryExpr) Accept(v Visitor[any]) any {
	return v.VisitBinaryExpr(binary)
}

type GroupingExpr struct {
	Expr Expr
}

func NewGroupingExpr(expr Expr) *GroupingExpr {
	return &GroupingExpr{
		Expr: expr,
	}
}

func (grouping *GroupingExpr) Print(v Visitor[string]) string {
	return v.VisitGroupingExpr(grouping)
}

func (grouping *GroupingExpr) Accept(v Visitor[any]) any {
	return v.VisitGroupingExpr(grouping)
}

type LiteralExpr struct {
	Value any
}

func NewLiteralExpr(value any) *LiteralExpr {
	return &LiteralExpr{
		Value: value,
	}
}

func (literal *LiteralExpr) Print(v Visitor[string]) string {
	return v.VisitLiteralExpr(literal)
}

func (literal *LiteralExpr) Accept(v Visitor[any]) any {
	return v.VisitLiteralExpr(literal)
}

type VarExpr struct {
	Name scanner.Token
}

func NewVarExpr(name scanner.Token) *VarExpr {
	return &VarExpr{
		Name: name,
	}
}

func (var_ *VarExpr) Print(v Visitor[string]) string {
	return v.VisitVarExpr(var_)
}

func (var_ *VarExpr) Accept(v Visitor[any]) any {
	return v.VisitVarExpr(var_)
}

type AssignExpr struct {
	Name  scanner.Token
	Value Expr
}

func NewAssignExpr(name scanner.Token, value Expr) *AssignExpr {
	a := new(AssignExpr)
	a.Name = name
	a.Value = value
	return a
}

func (assign *AssignExpr) Print(v Visitor[string]) string {
	return v.VisitAssignExpr(assign)
}

func (assign *AssignExpr) Accept(v Visitor[any]) any {
	return v.VisitAssignExpr(assign)
}

type LogicalExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func NewLogicalExpr(left Expr, operator scanner.Token, right Expr) *LogicalExpr {
	l := new(LogicalExpr)
	l.Left = left
	l.Right = right
	l.Operator = operator
	return l
}

func (logical *LogicalExpr) Accept(v Visitor[any]) any {
	return v.VisitLogicalExpr(logical)
}

func (logical *LogicalExpr) Print(v Visitor[string]) string {
	return v.VisitLogicalExpr(logical)
}

type CallExpr struct {
	CaleeToken scanner.Token
	Callee     Expr
	Args       []Expr
}

func NewCallExpr(caleeToken scanner.Token, callee Expr, args []Expr) *CallExpr {
	c := new(CallExpr)
	c.Args = args
	c.CaleeToken = caleeToken
	c.Callee = callee
	return c
}

func (call *CallExpr) Accept(v Visitor[any]) any {
	return v.VisitCallExpr(call)
}

func (call *CallExpr) Print(v Visitor[string]) string {
	return v.VisitCallExpr(call)
}

type GetExpr struct {
	Object Expr
	Name   scanner.Token
}

func NewGetExpr(object Expr, name scanner.Token) *GetExpr {
	return &GetExpr{
		Object: object,
		Name:   name,
	}
}

func (get *GetExpr) Accept(v Visitor[any]) any {
	return v.VisitGetExpr(get)
}

func (get *GetExpr) Print(v Visitor[string]) string {
	return v.VisitGetExpr(get)
}

type SetExpr struct {
	Object Expr
	Name   scanner.Token
	Value  Expr
}

func NewSetExpr(object Expr, name scanner.Token, value Expr) *SetExpr {
	return &SetExpr{
		Object: object,
		Name:   name,
		Value:  value,
	}
}

func (set *SetExpr) Accept(v Visitor[any]) any {
	return v.VisitSetExpr(set)
}

func (set *SetExpr) Print(v Visitor[string]) string {
	return v.VisitSetExpr(set)
}

type ThisExpr struct {
	Keyword scanner.Token
}

func NewThisExpr(keyword scanner.Token) *ThisExpr {
	return &ThisExpr{
		Keyword: keyword,
	}
}

func (this *ThisExpr) Accept(v Visitor[any]) any {
	return v.VisitThisExpr(this)
}

func (this *ThisExpr) Print(v Visitor[string]) string {
	return v.VisitThisExpr(this)
}

type SuperExpr struct {
	Keyword scanner.Token
	Method  scanner.Token
}

func NewSuperExpr(keyword, method scanner.Token) *SuperExpr {
	return &SuperExpr{
		Keyword: keyword,
		Method:  method,
	}
}

func (super *SuperExpr) Accept(v Visitor[any]) any {
	return v.VisitSuperExpr(super)
}

func (super *SuperExpr) Print(v Visitor[string]) string {
	return v.VisitSuperExpr(super)
}

type ArrayDeclExpr struct {
	Elements []Expr
}

func NewArrayDeclExpr(elements []Expr) *ArrayDeclExpr {
	return &ArrayDeclExpr{Elements: elements}
}

func (arr *ArrayDeclExpr) Accept(v Visitor[any]) any {
	return v.VisitArrayDeclExpr(arr)
}

func (arr *ArrayDeclExpr) Print(v Visitor[string]) string {
	return v.VisitArrayDeclExpr(arr)
}

type SubscriptExpr struct {
	ObjectToken scanner.Token
	Object      Expr
	IndexToken  scanner.Token
	Index       Expr
}

func NewSubscriptExpr(object, index Expr, objectT, indexT scanner.Token) *SubscriptExpr {
	return &SubscriptExpr{
		Object:      object,
		Index:       index,
		ObjectToken: objectT,
		IndexToken:  indexT,
	}
}

func (sub *SubscriptExpr) Accept(v Visitor[any]) any {
	return v.VisitSubscriptExpr(sub)
}

func (sub *SubscriptExpr) Print(v Visitor[string]) string {
	return v.VisitSubscriptExpr(sub)
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *Print)
	VisitExpressionStmt(stmt *Expression)
	VisitVarStmt(stmt *Var)
	VisitBlockStmt(stmt *Block)
	VisitIfStmt(stmt *If)
	VisitWhileStmt(stmt *While)
	VisitClassStmt(stmt *Class)
	VisitFunctionStmt(stmt *Function)
	VisitReturnStmt(stmt *Return)
}

type Stmt interface {
	Accept(StmtVisitor)
}

type Print struct {
	Expr Expr
}

func NewPrint(expr Expr) *Print {
	p := new(Print)
	p.Expr = expr
	return p
}

func (p *Print) Accept(vis StmtVisitor) {
	vis.VisitPrintStmt(p)
}

type Expression struct {
	Expr Expr
}

func NewExpression(expr Expr) *Expression {
	e := new(Expression)
	e.Expr = expr
	return e
}

func (e *Expression) Accept(vis StmtVisitor) {
	vis.VisitExpressionStmt(e)
}

type Var struct {
	VarName  scanner.Token
	VarValue Expr
}

func NewVar(varName scanner.Token, varValue Expr) *Var {
	return &Var{
		VarName:  varName,
		VarValue: varValue,
	}
}

func (v *Var) Accept(vis StmtVisitor) {
	vis.VisitVarStmt(v)
}

type Block struct {
	Stmts []Stmt
}

func NewBlock(stmts []Stmt) *Block {
	b := new(Block)
	b.Stmts = stmts
	return b
}

func (b *Block) Accept(vis StmtVisitor) {
	vis.VisitBlockStmt(b)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenBranch, elseBranch Stmt) *If {
	return &If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (i *If) Accept(vis StmtVisitor) {
	vis.VisitIfStmt(i)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func NewWhile(condition Expr, body Stmt) *While {
	return &While{
		Condition: condition,
		Body:      body,
	}
}

func (w *While) Accept(vis StmtVisitor) {
	vis.VisitWhileStmt(w)
}

type Function struct {
	Name      scanner.Token
	Arguments []scanner.Token
	Body      *Block
}

func NewFunction(name scanner.Token, arguments []scanner.Token, body *Block) *Function {
	return &Function{
		Name:      name,
		Arguments: arguments,
		Body:      body,
	}
}

func (fn *Function) Accept(vis StmtVisitor) {
	vis.VisitFunctionStmt(fn)
}

type Return struct {
	RetKeyWord scanner.Token
	Value      Expr
}

func NewReturn(retKeyWord scanner.Token, value Expr) *Return {
	return &Return{
		Value:      value,
		RetKeyWord: retKeyWord,
	}
}

func (ret *Return) Accept(vis StmtVisitor) {
	vis.VisitReturnStmt(ret)
}

type Class struct {
	Name       scanner.Token
	Superclass *VarExpr
	Methods    []*Function
}

func NewClass(name scanner.Token, superclass *VarExpr, methods []*Function) *Class {
	return &Class{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}
}

func (cls *Class) Accept(vis StmtVisitor) {
	vis.VisitClassStmt(cls)
}
