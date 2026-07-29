package parser

import (
	"fmt"
	dt "lig/datatypes"
)

type Parser struct {
	Tokens []dt.Token
	cur    int
}

func New(tokens []dt.Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() ([]dt.Stmt, *ParseError) {
	var res []dt.Stmt
	if p.isAtEnd() {
		return nil, nil
	}

	for !p.isAtEnd() {
		temp, err := p.declaration()
		if err != nil {
			return res, err
		}
		res = append(res, temp)
	}

	if !p.isAtEnd() {
		endErr := &ParseError{p.peek(), fmt.Sprintf("Expected EOF, received %v", p.peek())}
		return res, endErr
	}
	return res, nil
}

type ParseError struct {
	Token dt.Token
	Msg   string
}

func (e *ParseError) Error() string {
	return e.Msg
}

func (p *Parser) declaration() (dt.Stmt, *ParseError) {
	if p.match(dt.Var) {
		return p.varDeclaration()
	} else {
		return p.statement()
	}
}

func (p *Parser) varDeclaration() (dt.Stmt, *ParseError) {
	var initVal dt.Expr = nil
	err := p.consume(dt.Identifier, "Expect variable name")
	if err != nil {
		return nil, err
	}
	name := p.previous()

	if p.match(dt.Equal) {
		initVal, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	err = p.consume(dt.Semicolon, "Expect ';' after variable declaration")

	return &dt.VarStmt{name, initVal}, err
}

func (p *Parser) statement() (dt.Stmt, *ParseError) {
	if p.match(dt.Print) {
		return p.printStatement()
	} else if p.match(dt.LeftBrace) {
		return p.blockStatement()
	} else if p.match(dt.If) {
		return p.ifStatement()
	} else if p.match(dt.While) {
		return p.whileStatement()
	} else if p.match(dt.For) {
		return p.forStatement()
	}
	return p.assignOrExprStatement()
}

func (p *Parser) printStatement() (dt.Stmt, *ParseError) {
	res, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(dt.Semicolon, "Expect ';' after print Statement")
	return &dt.PrintStmt{res}, err
}

func (p *Parser) blockStatement() (dt.Stmt, *ParseError) {
	var res []dt.Stmt
	for !p.isAtEnd() && p.peek().Type != dt.RightBrace {
		temp, err := p.declaration()
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	err := p.consume(dt.RightBrace, "Expect closing braces")
	if err != nil {
		return nil, err
	}
	return &dt.BlockStmt{res}, nil
}

func (p *Parser) ifStatement() (dt.Stmt, *ParseError) {
	var elseStmt dt.Stmt = nil
	var err *ParseError
	err = p.consume(dt.LeftParen, "Expect '(' after if")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(dt.RightParen, "Expect ')' after condition of if statement")
	if err != nil {
		return nil, err
	}

	thenStmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	if p.match(dt.Else) { //Eagerly look for else-branch
		elseStmt, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &dt.IfStmt{condition, thenStmt, elseStmt}, nil
}

func (p *Parser) whileStatement() (dt.Stmt, *ParseError) {
	var err *ParseError
	err = p.consume(dt.LeftParen, "Expect '(' after while")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(dt.RightParen, "Expect ')' after condition of while statement")
	if err != nil {
		return nil, err
	}

	code, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &dt.WhileStmt{condition, code}, nil
}

func (p *Parser) forStatement() (dt.Stmt, *ParseError) {
	var initStmt dt.Stmt = nil
	var condition dt.Expr = &dt.Literal{true}
	var updStmt dt.Stmt = nil
	var err *ParseError
	err = p.consume(dt.LeftParen, "Expect '(' after for")
	if err != nil {
		return nil, err
	}
	if p.peek().Type != dt.Semicolon { // Initializer declaration exists
		if p.match(dt.Var) {
			initStmt, err = p.varDeclaration()
		} else {
			initStmt, err = p.assignOrExprStatement()
		}
		if err != nil {
			return nil, err
		}
	}

	p.consume(dt.Semicolon, "Expect ';' before condition of for")
	if p.peek().Type != dt.Semicolon { // Condition exists
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(dt.Semicolon, "Expect ';' after condition of for")
	if p.peek().Type != dt.RightBrace { // Exists a update statement
		updStmt, err = p.assignOrExprStatement()
		if err != nil {
			return nil, err
		}
	}

	p.consume(dt.RightBrace, "Expect ')' after update statement of for")
	code, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &dt.ForStmt{initStmt, condition, updStmt, code}, nil
}

func (p *Parser) assignOrExprStatement() (dt.Stmt, *ParseError) {
	res, err := p.expression()
	if err != nil {
		return nil, err
	}

	if p.match(dt.Equal) { // Assignment Statement located
		if varExpr, ok := res.(*dt.Variable); ok {
			rightExpr, err := p.expression()
			err = p.consume(dt.Semicolon, "Expect ';' after statement")
			return &dt.AssignStmt{varExpr.Name, rightExpr}, err
		}
		return nil, &ParseError{p.previous(), "Invalid assignment target"}
	}

	err = p.consume(dt.Semicolon, "Expect ';' after statement")
	return &dt.ExprStmt{res}, err
}

func (p *Parser) expression() (dt.Expr, *ParseError) {
	res, err := p.equality()
	if err != nil {
		return res, err
	}
	return res, nil
}

func (p *Parser) equality() (dt.Expr, *ParseError) {
	left, leftErr := p.comparison()
	if leftErr != nil {
		return nil, leftErr
	}

	for !p.isAtEnd() && (p.match(dt.EqualEqual) || p.match(dt.BangEqual)) {
		operator := p.previous()
		right, rightErr := p.comparison()
		if rightErr != nil {
			return left, rightErr
		}
		left = &dt.Binary{left, operator, right}
	}

	return left, nil
}

func (p *Parser) comparison() (dt.Expr, *ParseError) {
	left, leftErr := p.term()
	if leftErr != nil {
		return nil, leftErr
	}

	if p.match(dt.Greater) || p.match(dt.GreaterEqual) || p.match(dt.Less) || p.match(dt.LessEqual) {
		operator := p.previous()
		right, rightErr := p.term() // No chaining!
		if rightErr != nil {
			return left, rightErr
		}
		left = &dt.Binary{left, operator, right}
	}

	return left, nil
}

func (p *Parser) term() (dt.Expr, *ParseError) {
	left, leftErr := p.concat()
	if leftErr != nil {
		return nil, leftErr
	}

	for !p.isAtEnd() && (p.match(dt.Sub) || p.match(dt.Add)) {
		operator := p.previous()
		right, rightErr := p.concat()
		if rightErr != nil {
			return left, rightErr
		}
		left = &dt.Binary{left, operator, right}
	}

	return left, nil
}

func (p *Parser) concat() (dt.Expr, *ParseError) {
	left, leftErr := p.factor()
	if leftErr != nil {
		return nil, leftErr
	}

	for !p.isAtEnd() && p.match(dt.AddAdd) {
		operator := p.previous()
		right, rightErr := p.factor()
		if rightErr != nil {
			return left, rightErr
		}
		left = &dt.Binary{left, operator, right}
	}

	return left, nil
}

func (p *Parser) factor() (dt.Expr, *ParseError) {
	left, leftErr := p.unary()
	if leftErr != nil {
		return nil, leftErr
	}

	for !p.isAtEnd() && (p.match(dt.Mult) || p.match(dt.Div)) {
		operator := p.previous()
		right, rightErr := p.unary()
		if rightErr != nil {
			return left, rightErr
		}
		left = &dt.Binary{left, operator, right}
	}

	return left, nil
}

func (p *Parser) unary() (dt.Expr, *ParseError) {
	if !p.isAtEnd() && (p.match(dt.Bang) || p.match(dt.Sub)) {
		operator := p.previous()
		right, rightErr := p.unary()
		if rightErr != nil {
			return nil, rightErr
		}
		right = &dt.Unary{operator, right}
		return right, nil
	}

	prim, primErr := p.primary()
	if primErr != nil {
		return nil, primErr
	}

	return prim, nil
}

func (p *Parser) primary() (dt.Expr, *ParseError) {
	if p.isAtEnd() {
		return nil, &ParseError{p.previous(), "Expected literal, received nothing"}
	}

	if p.match(dt.Number) {
		return &dt.Literal{p.previous().Value}, nil
	} else if p.match(dt.String) {
		return &dt.Literal{p.previous().Value}, nil
	} else if p.match(dt.Identifier) {
		return &dt.Variable{p.previous()}, nil
	} else if p.match(dt.True) {
		return &dt.Literal{true}, nil
	} else if p.match(dt.False) {
		return &dt.Literal{false}, nil
	} else if p.match(dt.LeftParen) {
		exprPtr, err := p.expression()
		if err != nil {
			return nil, err
		}
		err = p.consume(dt.RightParen, "Expected closing parentheses")
		if err != nil {
			return nil, err
		}
		return exprPtr, nil
	} else {
		return nil, &ParseError{p.peek(), fmt.Sprintf("Expected literal, received %v", p.peek())}
		// we need to peek since no match has occured!
	}
}

func (p *Parser) advance() dt.Token {
	p.cur += 1
	return p.Tokens[p.cur-1]
}

func (p *Parser) isAtEnd() bool {
	// Due to EOF token
	return len(p.Tokens)-1 <= p.cur
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

func (p *Parser) consume(expect dt.Tokentype, msg string) *ParseError {
	if p.peek().Type == expect {
		p.cur += 1
		return nil
	}
	if p.cur == 0 {
		return &ParseError{p.peek(), "private function consume has been used at the start causing error"}
	}
	return &ParseError{p.previous(), msg}
}
