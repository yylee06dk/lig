package parser

import (
	"fmt"
	dt "lig/datatypes"
)

// Top-level declaration.
 
// Dummy Expression (for error returning)

var dummy dt.Expr = dt.Literal{0}


type Parser struct {
	Tokens []dt.Token
	cur int
}

func New(tokens []dt.Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() (dt.Expr, error) {
	if len(p.Tokens) == 1 { return dt.End{}, nil }
	res, err := p.expression()
	if err != nil {
		return res, fmt.Errorf("ParseError: %w", err)
	}

	if p.peek().Type != dt.EOF {
		endErr := &ParseError{p.peek(), fmt.Sprintf("Expected EOF, received %v", p.peek())}
		return res, fmt.Errorf("ParseError: %w", endErr)
	}
	return res, nil
}

type ParseError struct {
	Source dt.Token
	Msg string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("In token %v, Error occured: %s", e.Source, e.Msg)
}

func (p *Parser) expression() (dt.Expr, error) {
	res, err := p.equality()
	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *Parser) equality() (dt.Expr, error) {
	left, leftErr := p.comparison()
	if leftErr != nil {
		return dummy, leftErr
	}

	for !p.isAtEnd() && (p.match(dt.EqualEqual) || p.match(dt.BangEqual)) {
		operator := p.previous()
		right, rightErr := p.comparison()
		if rightErr != nil {
			return left, rightErr
		}
		left = dt.Binary{left, operator.Type, right}
	}

	return left, nil
}


func (p *Parser) comparison() (dt.Expr, error) {
	left, leftErr := p.term()
	if leftErr != nil {
		return dummy, leftErr
	}

	if p.match(dt.Greater) || p.match(dt.GreaterEqual) || p.match(dt.Less) || p.match(dt.LessEqual) {
		operator := p.previous()
		right, rightErr := p.term() // No chaining!
		if rightErr != nil {
			return left, rightErr
		}
		left = dt.Binary{left, operator.Type, right}
	}

	return left, nil
}

func (p *Parser) term() (dt.Expr, error) {
	left, leftErr := p.concat()
	if leftErr != nil {
		return dummy, leftErr
	}

	for !p.isAtEnd() && (p.match(dt.Sub) || p.match(dt.Add)) {
		operator := p.previous()
		right, rightErr := p.concat()
		if rightErr != nil {
			return left, rightErr
		}
		left = dt.Binary{left, operator.Type, right}
	}

	return left, nil
}

func (p *Parser) concat() (dt.Expr, error) {
	left, leftErr := p.factor()
	if leftErr != nil {
		return dummy, leftErr
	}

	for !p.isAtEnd() && p.match(dt.AddAdd) {
		operator := p.previous()
		right, rightErr := p.factor()
		if rightErr != nil {
			return left, rightErr
		}
		left = dt.Binary{left, operator.Type, right}
	}

	return left, nil
}

func (p *Parser) factor() (dt.Expr, error) {
	left, leftErr := p.unary()
	if leftErr != nil {
		return dummy, leftErr
	}

	for !p.isAtEnd() && (p.match(dt.Mult) || p.match(dt.Div)) {
		operator := p.previous()
		right, rightErr := p.unary()
		if rightErr != nil {
			return left, rightErr
		}
		left = dt.Binary{left, operator.Type, right}
	}

	return left, nil
}

func (p *Parser) unary() (dt.Expr, error) {
	if !p.isAtEnd() && (p.match(dt.Bang) || p.match(dt.Sub)) {
		operator := p.previous()
		right, rightErr := p.unary()
		if rightErr != nil {
			return dummy, rightErr
		}
		right = dt.Unary{operator.Type, right}
		return right, nil
	}

	prim, primErr := p.primary()
	if primErr != nil {
		return dummy, primErr
	}

	return prim, nil
}

func (p *Parser) primary() (dt.Expr, error) {
	if p.isAtEnd() {
		return dummy, &ParseError{p.previous(), "Expected literal, received nothing"}
	}
	if p.match(dt.Number) {
		return dt.Literal{p.previous().Value}, nil
	} else if p.match(dt.String) {
		return dt.Literal{p.previous().Value}, nil
	} else {
		return dummy, &ParseError{p.previous(), fmt.Sprintf("Expected literal, received %v", p.peek())}
	}
}

func (p *Parser) advance() dt.Token {
	p.cur += 1
	return p.Tokens[p.cur - 1]
}

func (p *Parser) isAtEnd() bool {
	// Due to EOF token
	return len(p.Tokens) - 1 <= p.cur
}

func (p *Parser) peek() dt.Token {
	return p.Tokens[p.cur]
}

func (p *Parser) previous() dt.Token {
	return p.Tokens[p.cur-1]
}

func (p *Parser) match(expectType dt.Tokentype) bool {
	if p.peek().Type == expectType {
		p.cur += 1 
		return true
	}
	return false
}
