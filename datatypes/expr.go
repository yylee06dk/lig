package datatypes

import (
	"fmt"
)

type Expr interface {
	dummyExpr()
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator Token
	Right    Expr
}

type Variable struct {
	Name Token
}

type Grouping struct {
	Inner Expr
}

type Call struct {
	Name    Expr
	Paren   Token
	ArgList []Expr
}

func (b *Binary) String() string {
	return fmt.Sprintf("[ {%v} {%v} {%v} ]", b.Left, b.Operator, b.Right)
}

func (l *Literal) String() string {
	switch v := l.Value.(type) {
	case int:
		return fmt.Sprintf("num[%v]", l.Value)
	case string:
		return fmt.Sprintf("string[%s]", v)
	case bool:
		return fmt.Sprintf("boolean[%v]", v)
	default:
		return fmt.Sprintf("UnknownLiteral: %v", l.Value)
	}
}

func (v *Variable) String() string {
	return fmt.Sprintf("Var[%s]", v.Name)
}

func (u *Unary) String() string {
	return fmt.Sprintf("[ {%v} {%v} ]", u.Operator, u.Right)
}

func (g *Grouping) String() string {
	return fmt.Sprintf("Group[ {%v} ]", g.Inner)
}

func (c *Call) String() string {
	return fmt.Sprintf("Call [ %s | %v ]", c.Name, c.ArgList)
}

func (_ *Binary) dummyExpr()   {}
func (_ *Literal) dummyExpr()  {}
func (_ *Unary) dummyExpr()    {}
func (_ *Variable) dummyExpr() {}
func (_ *Grouping) dummyExpr() {}
func (_ *Call) dummyExpr()     {}
