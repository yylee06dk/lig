package parser

import (
	"lig/datatypes"
)

func Parse(tokens []datatypes.Token) datatypes.Expr {
	return expression(tokens, 0)
}

func expression(tokens []datatypes.Token, cur int) datatypes.Expr {
	return term(tokens, cur)
}

func term(tokens []datatypes.Token, cur int) datatypes.Expr {
	left := factor(tokens, cur)
	cur = advance(cur)

	for !isAtEnd(tokens, cur) && (match(tokens, cur, datatypes.Sub) || match(tokens, cur, datatypes.Add)) {
		operator := peek(tokens, cur)
		cur = advance(cur)
		right := factor(tokens, cur)
		cur = advance(cur)
		left = datatypes.Binary{left, operator.Type, right}
	}

	return left
}

func factor(tokens []datatypes.Token, cur int) datatypes.Expr {
	var operator datatypes.Token
	var right datatypes.Expr
	left := literal(tokens, cur)
	cur = advance(cur)

	for !isAtEnd(tokens, cur) && (match(tokens, cur, datatypes.Mult) || match(tokens, cur, datatypes.Div)) {
		operator = peek(tokens, cur)
		cur = advance(cur)
		right = literal(tokens, cur)
		cur = advance(cur)
		left = datatypes.Binary{left, operator.Type, right}	}

	return left
}

func literal(tokens []datatypes.Token, cur int) datatypes.Expr {
	return datatypes.Literal{peek(tokens, cur).Value}
}

func advance(cur int) int {
	return cur+1
}

func isAtEnd(tokens []datatypes.Token, cur int) bool {
	return len(tokens) <= cur
}

func peek(tokens []datatypes.Token, cur int) datatypes.Token {
	return tokens[cur]
}

func match(tokens []datatypes.Token, cur int, expectType datatypes.Tokentype) bool {
	if peek(tokens, cur).Type == expectType { return true }
	return false
}
