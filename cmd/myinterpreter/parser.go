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
	currentIndex int
}

func NewParser(tokens []Token) *Parser {
	p := new(Parser)
	p.tokens = tokens
	p.exprs = make([]Expr, 0, 1)
	p.errs = make([]error, 0, 1)
	p.currentIndex = 0
	return p
}

func (p *Parser) parse() []Expr {
	for !p.isAtEnd() {
		expr, err := p.nextExpr()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(0)
		}
		p.exprs = append(p.exprs, expr)
	}
	return p.exprs
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

func (p *Parser) literalExpr() (*LiteralExpr, error) {
	if p.match(TRUE) {
		return NewLiteralExpr(true), nil
	} else if p.match(FALSE) {
		return NewLiteralExpr(false), nil
	} else if p.match(NIL) {
		return NewLiteralExpr(nil), nil
	} else if p.match(STRING, NUMBER) {
		return NewLiteralExpr(p.getPrev().Literal), nil
	}
	return nil, errors.New("Unknown expr")
}

func (p *Parser) group() (Expr, error){
    if p.match(LEFT_PAREN) {
        expr, err := p.nextExpr()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        if !p.check(RIGHT_PAREN){
            fmt.Fprintln(os.Stderr, errors.New("Expect ) after expression"))
            os.Exit(1)
        }
        p.currentIndex++
        return NewGroupingExpr(expr), nil
    }
    return p.literalExpr()
}

func (p *Parser) unary() (Expr, error){
    if p.match(MINUS, BANG) {
        token := p.getPrev()
        expr, err := p.unary()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        return NewUnaryExpr(token, expr), nil
    }
    return p.group()
}

func (p *Parser) factor() (Expr, error){
    expr, err := p.unary()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    for p.match(STAR, SLASH) {
        operator := p.getPrev()
        right, _ := p.unary()
        expr = NewBinaryExpr(expr, operator, right)
    }
    return expr, nil
}

func (p *Parser) term() (Expr, error) {
    expr, err := p.factor()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    for p.match(PLUS, MINUS)  {
        operator := p.getPrev()
        right, _ := p.factor()
        expr = NewBinaryExpr(expr, operator, right)
    }
    return expr, nil
}

func (p *Parser) comparison() (Expr, error){
    expr, err := p.term()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
        operator := p.getPrev()
        right, _ := p.term()
        expr = NewBinaryExpr(expr, operator, right)
    }
    return expr, nil
}

func (p *Parser) equality() (Expr, error){
    expr, err := p.comparison()
    if err !=  nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    for p.match(EQUAL_EQUAL, BANG_EQUAL) {
        operator := p.getPrev()
        right, _ := p.comparison()
        expr = NewBinaryExpr(expr, operator, right)
    }
    return expr, nil
}

func (p *Parser) nextExpr() (Expr, error) {
    return p.equality()
}
