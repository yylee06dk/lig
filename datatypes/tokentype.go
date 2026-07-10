package datatypes

import (
	"fmt"
)

type Tokentype int

const (
	Number Tokentype = iota
	Add
	Sub
	Mult
	Div
	Group
)

func (t Tokentype) String() string {
	switch t {
	case Number:
		return "Number"
	case Add:
		return "+"
	case Sub:
		return "-"
	case Mult:
		return "*"
	case Div:
		return "/"
	case Group:
		return "group"
	default:
		return fmt.Sprintf("UnknownTokentype(%d)", t)
	}
}
