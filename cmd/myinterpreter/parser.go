package main

import (
	"fmt"
	"os"
)

type Parser struct {
	tokens       []Token
	exprs        []Expr
	errs         []error
	exitCode     int
	currentIndex int
}

func NewParser(tokens []Token) *Parser {
	p := new(Parser)
	p.tokens = tokens
	p.exprs = make([]Expr, 0, 1)
	p.errs = make([]error, 0, 1)
	p.currentIndex = 0
	p.exitCode = 0
	return p
}

func (p *Parser) parseExprs() []Expr {
	for !p.isAtEnd() {
		expr := p.nextExpr()
		if expr != nil {
			p.exprs = append(p.exprs, expr)
		}
	}
	return p.exprs
}

func (p *Parser) parseStmts() []Stmt {
	var stmts []Stmt
	for !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}

	return stmts
}

func (p *Parser) declaration() Stmt {
	if p.match(CLASS) {
		return p.classDeclaration()
	} else if p.match(FUN) {
		return p.funStatement("function")
	} else if p.match(VAR) {
		return p.varStatement()
	}
	return p.statement()
}

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	} else if p.match(IF) {
		return p.ifStatement()
	} else if p.match(PRINT) {
		return p.printStatement()
	} else if p.match(RETURN) {
		return p.returnStatement()
	} else if p.match(WHILE) {
		return p.whileStatement()
	} else if p.match(LEFT_BRACE) {
		return p.blockStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) blockStatement() Stmt {
	var stmts []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	if p.getCurrent().Token != RIGHT_BRACE {
		p.error("Expect '}' after block statement")
	}
	p.incrIndex()
	return NewBlock(stmts)
}

func (p *Parser) classDeclaration() Stmt {
	name := p.getCurrent()
	p.incrIndex()
	if name.Token != IDENTIFIER {
		p.error("Expect class name.")
	}
	if !p.match(LEFT_BRACE) {
		p.error("Expect '{' before class body.")
	}
	methods := make([]*Function, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		method, ok := p.funStatement("method").(*Function)
		if !ok {
			panic("never")
		}
		methods = append(methods, method)
	}
	if !p.match(RIGHT_BRACE) {
		p.error("Expected '}' after class body.")
	}

	return NewClass(name, methods)
}

func (p *Parser) funStatement(kind string) Stmt {
	var parameters []Token = make([]Token, 0)
	name := p.getCurrent()
	p.currentIndex++
	if name.Token != IDENTIFIER {
		p.error(fmt.Sprintf("Expect %v name", kind))
	}
	if !p.match(LEFT_PAREN) {
		p.error("Expect '(' before condition expression")
	}

	if !p.check(RIGHT_PAREN) {
		param := p.getCurrent()
		p.currentIndex++
		parameters = append(parameters, param)
		for p.match(COMMA) {
			if len(parameters) > 255 {
				p.error(fmt.Sprintf("too many arguments %v, expect no more than 255", len(parameters)))
			}
			param := p.getCurrent()
			p.currentIndex++
			if param.Token != IDENTIFIER {
				p.error("Expect parametr name")
			}
			parameters = append(parameters, param)
		}
	}
	if !p.match(RIGHT_PAREN) {
		p.error("Expect '(' after parameters")
	}
	if !p.match(LEFT_BRACE) {
		p.error("Expect '{' before function body")
	}
	body := p.blockStatement()
	return NewFunction(name, parameters, body.(*Block))
}

func (p *Parser) returnStatement() Stmt {
	retKeyWord := p.getPrev()
	var value Expr = nil
	if !p.check(SEMICOLON) {
		value = p.nextExpr()
	}

	if !p.match(SEMICOLON) {
		p.error("Expect ';' after expression")
	}
	return NewReturn(retKeyWord, value)
}

func (p *Parser) varStatement() Stmt {
	varName := p.incrIndex()
	var varValue Expr = nil
	if varName.Token != IDENTIFIER {
		p.error("Expect variable name.")
	}
	if p.match(EQUAL) {
		varValue = p.nextExpr()
	}
	if p.getCurrent().Token != SEMICOLON {
		p.error("Expect ';' after variable declaration")
	}
	p.incrIndex()
	return NewVar(varName, varValue)
}

func (p *Parser) printStatement() Stmt {
	expr := p.nextExpr()
	if !p.match(SEMICOLON) {
		p.error("Expect ';' after expression")
	}
	return NewPrint(expr)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.nextExpr()
	if !p.match(SEMICOLON) {
		p.error("Expect ';' after expression")
	}
	return NewExpression(expr)
}

func (p *Parser) ifStatement() Stmt {
	if !p.match(LEFT_PAREN) {
		p.error("Expect '(' before condition expression")
	}
	condition := p.nextExpr()
	if !p.match(RIGHT_PAREN) {
		p.error("Expect ')' after condition expression")
	}
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	} else {
		elseBranch = nil
	}
	return NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) whileStatement() Stmt {
	if !p.match(LEFT_PAREN) {
		p.error("Expect '(' before condition expression")
		os.Exit(65)
	}
	condition := p.nextExpr()
	if !p.match(RIGHT_PAREN) {
		p.error("Expect ')' after condition expression")
	}
	body := p.statement()

	return NewWhile(condition, body)
}

func (p *Parser) forStatement() Stmt {
	if !p.match(LEFT_PAREN) {
		p.error("Expect '(' before condition expression")
	}

	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varStatement()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr = nil
	if !p.check(SEMICOLON) {
		condition = p.nextExpr()
	}
	if !p.match(SEMICOLON) {
		p.error("Expect ';' after loop condition")
	}

	var increment Expr = nil
	if !p.check(RIGHT_PAREN) {
		increment = p.nextExpr()
	}
	if !p.match(RIGHT_PAREN) {
		p.error("Expect ')' after for clauses")
	}

	body := p.statement()

	if increment != nil {
		body = &Block{
			[]Stmt{
				body,
				&Expression{increment},
			},
		}
	}

	if condition == nil {
		condition = &LiteralExpr{value: true}
	}
	loop := While{condition: condition, body: body}

	if initializer != nil {
		return &Block{
			stmts: []Stmt{initializer, &loop},
		}
	}

	return &loop
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.currentIndex].Token == EOF
}

func (p *Parser) check(t TokenType) bool {
	return p.getCurrent().Token == t
}

func (p *Parser) match(tts ...TokenType) bool {
	for _, tt := range tts {
		if p.check(tt) {
			p.incrIndex()
			return true
		}
	}
	return false
}

func (p *Parser) getCurrent() Token {
	return p.tokens[p.currentIndex]
}

func (p *Parser) getPrev() Token {
	return p.tokens[p.currentIndex-1]
}

func (p *Parser) incrIndex() Token {
	if !p.isAtEnd() {
		p.currentIndex++
	}
	return p.getPrev()
}

func (p *Parser) literalExpr() Expr {
	if p.match(TRUE) {
		return NewLiteralExpr(true)
	} else if p.match(FALSE) {
		return NewLiteralExpr(false)
	} else if p.match(NIL) {
		return NewLiteralExpr(nil)
	} else if p.match(STRING, NUMBER) {
		return NewLiteralExpr(p.getPrev().Literal)
	} else if p.match(IDENTIFIER) {
		return NewVarExpr(p.getPrev())
	}
	p.currentIndex++
	p.error("Expect expression")
	return nil
}

func (p *Parser) group() Expr {
	if p.match(LEFT_PAREN) {
		expr := p.nextExpr()
		if expr == nil {
			return nil
		}
		if !p.check(RIGHT_PAREN) {
			p.error("Expect ) after expression")
		} else {
			p.currentIndex++
		}
		return NewGroupingExpr(expr)
	}
	return p.literalExpr()
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := make([]Expr, 0, 0)
	if !p.check(RIGHT_PAREN) {
		arguments = append(arguments, p.nextExpr())
		for p.match(COMMA) {
			arguments = append(arguments, p.nextExpr())
		}
		if len(arguments) > 255 {
			p.error(fmt.Sprintf("too many arguments %v, expect no more than 255", len(arguments)))
		}
	}

	rp := p.getCurrent()
	if rp.Token != RIGHT_PAREN {
		p.error("expected ')' after arguments")
	}
	p.currentIndex++
	return NewCallExpr(callee, arguments)
}

func (p *Parser) call() Expr {
	expr := p.group()
	for true {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(MINUS, BANG) {
		token := p.getPrev()
		expr := p.unary()
		if expr == nil {
			return nil
		}
		return NewUnaryExpr(token, expr)
	}
	return p.call()
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	if expr == nil {
		return nil
	}
	for p.match(STAR, SLASH) {
		operator := p.getPrev()
		right := p.unary()
		if right == nil {
			return nil
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	if expr == nil {
		return nil
	}
	for p.match(PLUS, MINUS) {
		operator := p.getPrev()
		right := p.factor()
		if right == nil {
			return nil
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	if expr == nil {
		return nil
	}
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.getPrev()
		right := p.term()
		if right == nil {
			return nil
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	if expr == nil {
		return nil
	}
	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		operator := p.getPrev()
		right := p.comparison()
		if right == nil {
			return nil
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(OR) {
		operator := p.getPrev()
		right := p.equality()
		expr = NewLogicalExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(AND) {
		operator := p.getPrev()
		right := p.and()
		expr = NewLogicalExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		value := p.assignment()
		switch expr.(type) {
		case *VarExpr:
			name := expr.(*VarExpr).name
			return NewAssignExpr(name, value)
		default:
			p.error("Invalid assignment target")
		}
	}
	return expr
}

func (p *Parser) nextExpr() Expr {
	return p.assignment()
}

func (p *Parser) error(msg string) {
	fmt.Fprintf(os.Stderr, "[line %v] error at %v: %v\n", p.getCurrent().line, p.getCurrent().Lexeme, msg)
	os.Exit(65)
}
