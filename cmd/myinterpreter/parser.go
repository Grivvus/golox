package main

import (
	"errors"
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
    if p.match(VAR){
        return p.varStatement()
    }
    return p.statement()
}

func (p *Parser) statement() Stmt {
    if p.match(FOR){
        return p.forStatement()
    } else if p.match(IF){
		return p.ifStatement()
	} else if p.match(PRINT) {
		return p.printStatement()
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
		fmt.Fprintln(os.Stderr, "Expect } after block statement")
		os.Exit(65)
	}
	p.incrIndex()
	return *NewBlock(stmts)
}

func (p *Parser) varStatement() Stmt {
	varName := p.incrIndex()
	var varValue Expr = nil
	if varName.Token != IDENTIFIER {
		fmt.Fprintln(os.Stderr, "Expect variable name.")
		os.Exit(65)
	}
	if p.match(EQUAL) {
		varValue = p.nextExpr()
	}
	if p.getCurrent().Token != SEMICOLON {
		fmt.Fprintln(os.Stderr, "Expect ; after variable declaration")
        os.Exit(65)
	}
	p.incrIndex()
	return *NewVar(varName, varValue)
}

func (p *Parser) printStatement() Stmt {
	expr := p.nextExpr()
	if !p.match(SEMICOLON) {
		fmt.Fprintln(os.Stderr, "Expect ; after expression")
		os.Exit(65)
	}
	return *NewPrint(expr)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.nextExpr()
	if !p.match(SEMICOLON) {
		fmt.Fprintln(os.Stderr, "Expect ; after expression")
		os.Exit(65)
	}
	return *NewExpression(expr)
}

func (p *Parser) ifStatement() Stmt {
	if !p.match(LEFT_PAREN) {
		fmt.Fprintln(os.Stderr, "Expect ( before condition expression")
		os.Exit(65)
	}
	condition := p.nextExpr()
	if !p.match(RIGHT_PAREN) {
		fmt.Fprintln(os.Stderr, "Expect ) after condition expression")
		os.Exit(65)
	}
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	} else {
		elseBranch = nil
	}
	return *NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) whileStatement() Stmt {
	if !p.match(LEFT_PAREN) {
		fmt.Fprintln(os.Stderr, "Expect ( before condition expression")
		os.Exit(65)
	}
	condition := p.nextExpr()
	if !p.match(RIGHT_PAREN) {
		fmt.Fprintln(os.Stderr, "Expect ) after condition expression")
		os.Exit(65)
	}
    body := p.statement()

    return *NewWhile(condition, body)
}

func (p *Parser) forStatement() Stmt {
	if !p.match(LEFT_PAREN) {
		fmt.Fprintln(os.Stderr, "Expect ( before condition expression")
		os.Exit(65)
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
    if !p.check(SEMICOLON){
        condition = p.nextExpr()
    } 
    if !p.match(SEMICOLON) {
        fmt.Fprintln(os.Stderr, fmt.Sprintf("[line %v] Error at '%v' Expect ';' after expression", p.getCurrent().line, p.getCurrent().Lexeme))
        os.Exit(65)
    }

    var increment Expr = nil
    if !p.check(RIGHT_PAREN) {
        increment = p.nextExpr()
    }
    if !p.match(RIGHT_PAREN) {
		fmt.Fprintln(os.Stderr, "Expect ) after condition expression")
		os.Exit(65)
	}

    body := p.statement()

    if increment != nil {
        stmts := make([]Stmt, 2)
        stmts[0] = body
        stmts[1] = *NewExpression(increment)
        body = *NewBlock(stmts)
    }
    if condition ==  nil {
        condition = NewLiteralExpr(true)
    }
    body = *NewWhile(condition, body)
    if initializer != nil {
        stmts := make([]Stmt, 2)
        stmts[0] = initializer
        stmts[1] = body
        body = *NewBlock(stmts)
    }
    return body
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
	p.exitCode = 65
	fmt.Fprintln(os.Stderr, fmt.Sprintf("[line %v] Error at '%v': Expect expression.", p.getCurrent().line, p.getCurrent().Lexeme))
	return nil
}

func (p *Parser) group() Expr {
	if p.match(LEFT_PAREN) {
		expr := p.nextExpr()
		if expr == nil {
			return nil
		}
		if !p.check(RIGHT_PAREN) {
			p.exitCode = 65
			fmt.Fprintln(os.Stderr, errors.New("Expect ) after expression"))
			os.Exit(65)
		} else {
			p.currentIndex++
		}
		return NewGroupingExpr(expr)
	}
	return p.literalExpr()
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
	return p.group()
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
			fmt.Fprintln(os.Stderr, "Invalid assignment target")
			os.Exit(70)
		}
	}
	return expr
}

func (p *Parser) nextExpr() Expr {
	return p.assignment()
}
