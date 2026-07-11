package parser

import (
	"fmt"
	"lig/datatypes"
)

// Top-level declaration.
 
// Dummy Expression (for error returning)

var dummy datatypes.Expr = datatypes.Literal{0}


type Parser struct {
	Tokens []datatypes.Token
	cur int
}

func New(tokens []datatypes.Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() (datatypes.Expr, error) {
	res, err := p.expression()
	if err != nil {
		return res, fmt.Errorf("ParseError: %w", err)
	}

	if p.peek().Type != datatypes.EOF {
		endErr := &ParseError{p.peek(), fmt.Sprintf("Expected EOF, received %v", p.peek())}
		return res, fmt.Errorf("ParseError: %w", endErr)
	}
	return res, nil
}

type ParseError struct {
	Source datatypes.Token
	Msg string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("In token %v, Error occured: %s", e.Source, e.Msg)
}

func (p *Parser) expression() (datatypes.Expr, error) {
	res, err := p.term()
	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *Parser) term() (datatypes.Expr, error) {
	left, leftErr := p.factor()
	if leftErr != nil {
		return dummy, leftErr
	}

	for !p.isAtEnd() && (p.match(datatypes.Sub) || p.match(datatypes.Add)) {
		operator := p.previous()
		right, rightErr := p.factor()
		if rightErr != nil {
			return left, rightErr
		}
		left = datatypes.Binary{left, operator.Type, right}
	}

	return left, nil
}

func (p *Parser) factor() (datatypes.Expr, error) {
	left, leftErr := p.literal()
	if leftErr != nil {
		return dummy, leftErr
	}

	for !p.isAtEnd() && (p.match(datatypes.Mult) || p.match(datatypes.Div)) {
		operator := p.previous()
		right, rightErr := p.literal()
		if rightErr != nil {
			return left, rightErr
		}
		left = datatypes.Binary{left, operator.Type, right}
	}

	return left, nil
}

func (p *Parser) literal() (datatypes.Expr, error) {
	if !p.isAtEnd() && p.match(datatypes.Number) {
		return datatypes.Literal{p.previous().Value}, nil
	} else {
		return dummy, &ParseError{p.previous(), fmt.Sprintf("Expected literal, received %v", p.peek())}
	}
}

func (p *Parser) advance() datatypes.Token {
	p.cur += 1
	return p.Tokens[p.cur - 1]
}

func (p *Parser) isAtEnd() bool {
	// Due to EOF token
	return len(p.Tokens) - 1 <= p.cur
}

func (p *Parser) peek() datatypes.Token {
	return p.Tokens[p.cur]
}

func (p *Parser) previous() datatypes.Token {
	return p.Tokens[p.cur-1]
}

func (p *Parser) match(expectType datatypes.Tokentype) bool {
	if p.peek().Type == expectType {
		p.cur += 1 
		return true
	}
	return false
}
