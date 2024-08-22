package main

import (
	"errors"
)

type Parser struct{
    tokens []Token
    currentIndex int
}

func NewParser(tokens []Token) *Parser{
    p := new(Parser)
    p.tokens = tokens
    p.currentIndex = 0
    return p
}

func (p *Parser) nextExpr() (Expr, error){
    if p.isAtEnd(){
        return nil, errors.New("EOF")
    }
    token := p.tokens[p.currentIndex]
    p.currentIndex++
    if token.Token == TRUE || token.Token == FALSE || token.Token == NIL{
        return literalExpr(token)
    }
    return nil, errors.New("unknown expr")
}

func (p *Parser) isAtEnd() bool{
    return p.tokens[p.currentIndex].Token == EOF
}

func literalExpr(token Token) (*LiteralExpr, error){
    if token.Token == TRUE{
        return NewLiteralExpr(true), nil
    } else if token.Token == FALSE{
        return NewLiteralExpr(false), nil
    } else if token.Token == NIL{
        return NewLiteralExpr(nil), nil
    }
    return nil, errors.New("can't parse literal expr")
}
