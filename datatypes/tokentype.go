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
	EOF
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
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("UnknownTokentype(%d)", t)
	}
}
