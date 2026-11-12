package parser

import (
	"errors"
	"fmt"
	"golox/internal/scanner"
	"strings"
)

type ParsingError struct {
	msg string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("error happend during parsing stage: %w", e.msg)
}

type Parser struct {
	tokens       []scanner.Token
	Exprs        []Expr
	Errs         []error
	ExitCode     int
	currentIndex int
}

func New(tokens []scanner.Token) *Parser {
	p := new(Parser)
	p.tokens = tokens
	p.Exprs = make([]Expr, 0)
	p.Errs = make([]error, 0)
	p.currentIndex = 0
	p.ExitCode = 0
	return p
}

func (p *Parser) ParseExprs() ([]Expr, error) {
	for !p.isAtEnd() {
		expr := p.nextExpr()
		if len(p.Errs) != 0 {
			return nil, errors.Join(p.Errs...)
		}
		if expr == nil {
			panic("invalid state: expr is nil while no errors")
		}
		p.Exprs = append(p.Exprs, expr)
	}
	return p.Exprs, nil
}

func (p *Parser) ParseStmts() ([]Stmt, error) {
	var stmts []Stmt
	for !p.isAtEnd() {
		newStmt := p.declaration()
		if len(p.Errs) != 0 {
			return nil, errors.Join(p.Errs...)
		}
		if newStmt == nil {
			panic("invalid state: stmt is nil while no errors")
		}
		stmts = append(stmts, newStmt)
	}

	return stmts, nil
}

func (p *Parser) declaration() Stmt {
	if p.match(scanner.CLASS) {
		return p.classDeclaration()
	} else if p.match(scanner.FUN) {
		return p.funStatement("function")
	} else if p.match(scanner.VAR) {
		return p.varStatement()
	}
	return p.statement()
}

func (p *Parser) statement() Stmt {
	if p.match(scanner.FOR) {
		return p.forStatement()
	} else if p.match(scanner.IF) {
		return p.ifStatement()
	} else if p.match(scanner.PRINT) {
		return p.printStatement()
	} else if p.match(scanner.RETURN) {
		return p.returnStatement()
	} else if p.match(scanner.WHILE) {
		return p.whileStatement()
	} else if p.match(scanner.LEFT_BRACE) {
		return p.blockStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) blockStatement() Stmt {
	var stmts []Stmt
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	if p.GetCurrent().Token != scanner.RIGHT_BRACE {
		p.error("Expect '}' after block statement")
		return nil
	}
	p.incrIndex()
	return NewBlock(stmts)
}

func (p *Parser) classDeclaration() Stmt {
	name := p.GetCurrent()
	var superclass *VarExpr = nil
	p.incrIndex()
	if name.Token != scanner.IDENTIFIER {
		p.error("Expect class name.")
		return nil
	}
	if p.match(scanner.LESS) {
		if !p.match(scanner.IDENTIFIER) {
			p.error("Expect superclass name.")
			return nil
		}
		superclass = NewVarExpr(p.getPrev())
	}
	if !p.match(scanner.LEFT_BRACE) {
		p.error("Expect '{' before class body.")
		return nil
	}
	methods := make([]*Function, 0)
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		method, ok := p.funStatement("method").(*Function)
		if !ok {
			panic("never")
		}
		methods = append(methods, method)
	}
	if !p.match(scanner.RIGHT_BRACE) {
		p.error("Expected '}' after class body.")
		return nil
	}

	return NewClass(name, superclass, methods)
}

func (p *Parser) funStatement(kind string) Stmt {
	var parameters []scanner.Token = make([]scanner.Token, 0)
	name := p.GetCurrent()
	p.currentIndex++
	if name.Token != scanner.IDENTIFIER {
		p.error(fmt.Sprintf("Expect %v name", kind))
		return nil
	}
	if !p.match(scanner.LEFT_PAREN) {
		p.error("Expect '(' before condition expression")
		return nil
	}

	if !p.check(scanner.RIGHT_PAREN) {
		param := p.GetCurrent()
		p.currentIndex++
		parameters = append(parameters, param)
		for p.match(scanner.COMMA) {
			if len(parameters) > 255 {
				p.error(fmt.Sprintf("too many arguments %v, expect no more than 255", len(parameters)))
				return nil
			}
			param := p.GetCurrent()
			p.currentIndex++
			if param.Token != scanner.IDENTIFIER {
				p.error("Expect parametr name")
				return nil
			}
			parameters = append(parameters, param)
		}
	}
	if !p.match(scanner.RIGHT_PAREN) {
		p.error("Expect '(' after parameters")
		return nil
	}
	if !p.match(scanner.LEFT_BRACE) {
		p.error("Expect '{' before function body")
		return nil
	}
	body := p.blockStatement()
	return NewFunction(name, parameters, body.(*Block))
}

func (p *Parser) returnStatement() Stmt {
	retKeyWord := p.getPrev()
	var value Expr = nil
	if !p.check(scanner.SEMICOLON) {
		value = p.nextExpr()
	}

	if !p.match(scanner.SEMICOLON) {
		p.error("Expect ';' after expression")
		return nil
	}
	return NewReturn(retKeyWord, value)
}

func (p *Parser) varStatement() Stmt {
	varName := p.incrIndex()
	var varValue Expr = nil
	if varName.Token != scanner.IDENTIFIER {
		p.error("Expect variable name.")
		return nil
	}
	if p.match(scanner.EQUAL) {
		varValue = p.nextExpr()
	}
	if p.GetCurrent().Token != scanner.SEMICOLON {
		p.error("Expect ';' after variable declaration")
		return nil
	}
	p.incrIndex()
	return NewVar(varName, varValue)
}

func (p *Parser) printStatement() Stmt {
	expr := p.nextExpr()
	if !p.match(scanner.SEMICOLON) {
		p.error("Expect ';' after expression")
		return nil
	}
	return NewPrint(expr)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.nextExpr()
	if !p.match(scanner.SEMICOLON) {
		p.error("Expect ';' after expression")
		return nil
	}
	return NewExpression(expr)
}

func (p *Parser) ifStatement() Stmt {
	if !p.match(scanner.LEFT_PAREN) {
		p.error("Expect '(' before condition expression")
		return nil
	}
	condition := p.nextExpr()
	if !p.match(scanner.RIGHT_PAREN) {
		p.error("Expect ')' after condition expression")
		return nil
	}
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(scanner.ELSE) {
		elseBranch = p.statement()
	} else {
		elseBranch = nil
	}
	return NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) whileStatement() Stmt {
	if !p.match(scanner.LEFT_PAREN) {
		p.error("Expect '(' before condition expression")
		return nil
	}
	condition := p.nextExpr()
	if !p.match(scanner.RIGHT_PAREN) {
		p.error("Expect ')' after condition expression")
		return nil
	}
	body := p.statement()

	return NewWhile(condition, body)
}

func (p *Parser) forStatement() Stmt {
	if !p.match(scanner.LEFT_PAREN) {
		p.error("Expect '(' before condition expression")
		return nil
	}

	var initializer Stmt
	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer = p.varStatement()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr = nil
	if !p.check(scanner.SEMICOLON) {
		condition = p.nextExpr()
	}
	if !p.match(scanner.SEMICOLON) {
		p.error("Expect ';' after loop condition")
		return nil
	}

	var increment Expr = nil
	if !p.check(scanner.RIGHT_PAREN) {
		increment = p.nextExpr()
	}
	if !p.match(scanner.RIGHT_PAREN) {
		p.error("Expect ')' after for clauses")
		return nil
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
		condition = &LiteralExpr{Value: true}
	}
	loop := While{Condition: condition, Body: body}

	if initializer != nil {
		return &Block{
			Stmts: []Stmt{initializer, &loop},
		}
	}

	return &loop
}

func (p *Parser) isAtEnd() bool {
	return p.GetCurrent().Token == scanner.EOF
}

func (p *Parser) check(t scanner.TokenType) bool {
	return p.GetCurrent().Token == t
}

func (p *Parser) match(tts ...scanner.TokenType) bool {
	for _, tt := range tts {
		if p.check(tt) {
			p.incrIndex()
			return true
		}
	}
	return false
}

func (p *Parser) GetCurrent() scanner.Token {
	return p.tokens[p.currentIndex]
}

func (p *Parser) getPrev() scanner.Token {
	return p.tokens[p.currentIndex-1]
}

func (p *Parser) incrIndex() scanner.Token {
	if !p.isAtEnd() {
		p.currentIndex++
	}
	return p.getPrev()
}

func (p *Parser) literalExpr() Expr {
	if p.match(scanner.TRUE) {
		return NewLiteralExpr(true)
	} else if p.match(scanner.FALSE) {
		return NewLiteralExpr(false)
	} else if p.match(scanner.NIL) {
		return NewLiteralExpr(nil)
	} else if p.match(scanner.STRING, scanner.NUMBER) {
		return NewLiteralExpr(p.getPrev().Literal)
	} else if p.match(scanner.SUPER) {
		token := p.getPrev()
		if !p.match(scanner.DOT) {
			p.error("Expect '.' after 'super'.")
			return nil
		}
		method := p.incrIndex()
		if method.Token != scanner.IDENTIFIER {
			p.error("Expect method name after 'super'")
			return nil
		}
		return NewSuperExpr(token, method)
	} else if p.match(scanner.THIS) {
		return NewThisExpr(p.getPrev())
	} else if p.match(scanner.IDENTIFIER) {
		return NewVarExpr(p.getPrev())
	} else if p.match(scanner.LEFT_SQUARE_BRACKET) {
		return p.arrayLiteral()
	}
	p.currentIndex++
	p.error("Expect expression")
	return nil
}

func (p *Parser) arrayLiteral() Expr {
	elements := make([]Expr, 0)
	if p.check(scanner.RIGHT_SQUARE_BRACKET) {
		p.currentIndex++
		return NewArrayDeclExpr(elements)
	} else {
		elements = append(elements, p.nextExpr())
	}
	for p.match(scanner.COMMA) {
		if p.GetCurrent().Token != scanner.RIGHT_SQUARE_BRACKET {
			elements = append(elements, p.nextExpr())
		}
	}
	if !p.check(scanner.RIGHT_SQUARE_BRACKET) {
		p.error("Expect ']' after array declaration")
		return nil
	}
	p.currentIndex++
	return NewArrayDeclExpr(elements)
}

func (p *Parser) group() Expr {
	if p.match(scanner.LEFT_PAREN) {
		expr := p.nextExpr()
		if expr == nil {
			return nil
		}
		if !p.check(scanner.RIGHT_PAREN) {
			p.error("Expect ) after expression")
			return nil
		} else {
			p.currentIndex++
		}
		return NewGroupingExpr(expr)
	}
	return p.literalExpr()
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := make([]Expr, 0)
	if !p.check(scanner.RIGHT_PAREN) {
		arguments = append(arguments, p.nextExpr())
		for p.match(scanner.COMMA) {
			arguments = append(arguments, p.nextExpr())
		}
		if len(arguments) > 255 {
			p.error(fmt.Sprintf("too many arguments %v, expect no more than 255", len(arguments)))
			return nil
		}
	}

	rp := p.GetCurrent()
	if rp.Token != scanner.RIGHT_PAREN {
		p.error("expected ')' after arguments")
		return nil
	}
	p.currentIndex++
	return NewCallExpr(p.getPrev(), callee, arguments)
}

func (p *Parser) finishSubscript(object Expr) Expr {
	index := p.nextExpr()
	if p.GetCurrent().Token != scanner.RIGHT_SQUARE_BRACKET {
		p.error("Expect ']' after array subscription")
		return nil
	}
	indexToken := p.getPrev()
	objectToken := p.tokens[p.currentIndex-2]
	p.currentIndex++
	return NewSubscriptExpr(object, index, objectToken, indexToken)
}

func (p *Parser) call() Expr {
	expr := p.group()
	for true {
		if p.match(scanner.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(scanner.LEFT_SQUARE_BRACKET) {
			expr = p.finishSubscript(expr)
		} else if p.match(scanner.DOT) {
			name := p.GetCurrent()
			p.currentIndex++
			if name.Token != scanner.IDENTIFIER {
				p.error("Expect property name after '.'")
				return nil
			}
			expr = NewGetExpr(expr, name)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(scanner.MINUS, scanner.BANG) {
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
	for p.match(scanner.STAR, scanner.SLASH, scanner.PERCENT) {
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
	for p.match(scanner.PLUS, scanner.MINUS) {
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
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
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
	for p.match(scanner.EQUAL_EQUAL, scanner.BANG_EQUAL) {
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
	for p.match(scanner.OR) {
		operator := p.getPrev()
		right := p.equality()
		expr = NewLogicalExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(scanner.AND) {
		operator := p.getPrev()
		right := p.and()
		expr = NewLogicalExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(scanner.EQUAL) {
		value := p.assignment()
		switch expr := expr.(type) {
		case *VarExpr:
			name := expr.Name
			return NewAssignExpr(name, value)
		case *GetExpr:
			get := expr
			return NewSetExpr(get.Object, get.Name, value)
		default:
			p.error("Invalid assignment target")
			return nil
		}
	}
	return expr
}

func (p *Parser) nextExpr() Expr {
	return p.assignment()
}

func (p *Parser) error(msg string) {
	p.Errs = append(p.Errs, &ParsingError{msg: strings.ToLower(msg)})
}
