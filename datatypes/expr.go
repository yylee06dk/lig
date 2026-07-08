package datatypes

type Expr interface {
	exprOnly() int
}

type Binary struct {
	Left Expr
	Operator Tokentype
	Right Expr
}

func (_ Binary) exprOnly() int {
	return 1
}

type Literal struct {
	Value any
}

func (_ Literal) exprOnly() int {
	return 1
}
