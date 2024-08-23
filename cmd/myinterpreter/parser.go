package main

import (
	"errors"
)

type Parser struct {
	tokens       []Token
	currentIndex int
}

func NewParser(tokens []Token) *Parser {
	p := new(Parser)
	p.tokens = tokens
	p.currentIndex = 0
	return p
}

func (p *Parser) nextExpr() (Expr, error) {
	if p.isAtEnd() {
		return nil, errors.New("EOF")
	}
	token := p.tokens[p.currentIndex]
	p.currentIndex++
	if token.Token == TRUE || token.Token == FALSE || token.Token == NIL || token.Token == NUMBER || token.Token == STRING {
		return literalExpr(token)
	} else if token.Token == LEFT_PAREN {
		e, err := p.nextExpr()
		if err != nil {
			return nil, err
		}
		if p.tokens[p.currentIndex].Token != RIGHT_PAREN {
            return nil, errors.New("Error: Unmatched parentheses.")
		} else {
			p.currentIndex++
		}
		return groupExpr(e)
	}
	return nil, errors.New("unknown expr")
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.currentIndex].Token == EOF
}

func literalExpr(token Token) (*LiteralExpr, error) {
	if token.Token == TRUE {
		return NewLiteralExpr(true), nil
	} else if token.Token == FALSE {
		return NewLiteralExpr(false), nil
	} else if token.Token == NIL {
		return NewLiteralExpr(nil), nil
	} else if token.Token == NUMBER {
		return NewLiteralExpr(token.Literal), nil
	} else if token.Token == STRING {
		return NewLiteralExpr(token.Literal), nil
	}
	return nil, errors.New("can't parse literal expr")
}

func groupExpr(expr Expr) (*GroupingExpr, error) {
	return NewGroupingExpr(expr), nil
}
