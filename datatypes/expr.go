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
	switch v := l.Value.(type){
		case int:
			return fmt.Sprintf("num[%v]", l.Value)
		case string:
			return fmt.Sprintf("string[%s]", v)
		default:
			return fmt.Sprintf("UnknownLiteral: %v",l. Value)
	}
}

func (_ Binary) dummy() {}
func (_ Literal) dummy() {}
