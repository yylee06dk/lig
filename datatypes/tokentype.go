package datatypes

import (
	"fmt"
)

type Tokentype int

const (

	// Stage 1
	Number Tokentype = iota
	Add
	Sub
	Mult
	Div

	// Stage 1.1

	String
	Identifier
	True
	False

	AddAdd // String concat
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	And
	Or

	Error
	
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
	case String:
		return "string"
	case Identifier:
		return "identifier"
	case True:
		return "true"
	case False:
		return "False"
	case AddAdd:
		return "++"
	case Bang:
		return "!"
	case BangEqual:
		return "!="
	case Equal:
		return "="
	case Greater:
		return ">"
	case GreaterEqual:
		return ">="
	case Less:
		return "<"
	case LessEqual:
		return "<="
	case And:
		return "&&"
	case Or:
		return "||"
	case Error:
		return "Lexical Error"
	default:
		return fmt.Sprintf("UnknownTokentype(%d)", t)
	}
}
