package datatypes

import (
	"fmt"
)

type Expr interface {
	dummy()
}

type End struct {}

type Binary struct {
	Left Expr
	Operator Tokentype
	Right Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator Tokentype
	Right Expr
}

func (e *End) String() string {
	return fmt.Sprintf("END")
}

func (b *Binary) String() string {
	return fmt.Sprintf("[ {%v} {%v} {%v} ]", b.Left, b.Operator, b.Right)
}

func (l *Literal) String() string {
	switch v := l.Value.(type){
		case int:
			return fmt.Sprintf("num[%v]", l.Value)
		case string:
			return fmt.Sprintf("string[%s]", v)
		default:
			return fmt.Sprintf("UnknownLiteral: %v",l. Value)
	}
}

func (u *Unary) String() string {
	return fmt.Sprintf("[ {%v} {%v} ]", u.Operator, u.Right)
}

func (_ *End) dummy() {}
func (_ *Binary) dummy() {}
func (_ *Literal) dummy() {}
func (_ *Unary) dummy() {}
