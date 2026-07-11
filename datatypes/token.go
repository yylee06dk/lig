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
		default: // +, -, *, /
			return fmt.Sprintf("{%v}", t.Type.String())
	}
}

