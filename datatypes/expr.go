package datatypes

import (
	"fmt"
)

type Expr interface {
	dummy()
}

type Binary struct {
	Left Expr
	Operator Tokentype
	Right Expr
}

func (b Binary) String() string {
	return fmt.Sprintf("[ {%v} {%v} {%v} ]", b.Left, b.Operator, b.Right)
}

type Literal struct {
	Value any
}

func (l Literal) String() string {
	return fmt.Sprintf("num[%v]", l.Value)
}

func (_ Binary) dummy() {}
func (_ Literal) dummy() {}
