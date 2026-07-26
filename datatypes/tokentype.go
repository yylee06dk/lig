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

	Group

	// Stage 1.2 - user convenience
	Error

	LeftParen
	RightParen

	Skip // A hack to skipping characters

	// Stage 2
	Semicolon
	Print

	// Stage 2.2
	Var

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
	case EqualEqual:
		return "=="
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
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case Semicolon:
		return ";"
	case Print:
		return "Print"
	case Var:
		return "Var"
	default:
		return fmt.Sprintf("UnknownTokentype(%d)", t)
	}
}
