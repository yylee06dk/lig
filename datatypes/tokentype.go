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

	// Stage 2.3
	LeftBrace
	RightBrace

	// Stage 3
	If
	Else
	While
	For

	// Stage 4
	Fun
	Comma
	Return
	Returning

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
		return "false"
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
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"
	case If:
		return "if"
	case Else:
		return "else"
	case While:
		return "while"
	case For:
		return "for"
	case Fun:
		return "fun"
	case Comma:
		return ","
	case Return:
		return "return"
	case Error:
		return "Lexical Error"
	default:
		return fmt.Sprintf("UnknownTokentype(%d)", t)
	}
}
