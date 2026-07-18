package datatypes

import (
	"fmt"
)

type Token struct {
	Type Tokentype
	Value any
}

func (t Token) String() string {
	switch t.Type{
		case Number:
			return fmt.Sprintf("{num[%v]}", t.Value)
		case String:
			val := t.Value.(string)
			return fmt.Sprintf("string[%s]", val)
		default: // +, -, *, /
			return fmt.Sprintf("{%v}", t.Type.String())
	}
}

