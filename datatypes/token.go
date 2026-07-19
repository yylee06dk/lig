package datatypes

import (
	"fmt"
)

type Token struct {
	Type Tokentype
	Name string
	Value any
}

func (t Token) String() string {
	switch t.Type{
		case Number:
			return fmt.Sprintf("{num[%v]}", t.Value)
		case String:
			val := t.Value.(string)
			return fmt.Sprintf("{string[%s]}", val)
		case Identifier:
			name := t.Name
			return fmt.Sprintf("{identifier[%s]}", name)
		case Error:
			return "{Error}"
		default:
			return fmt.Sprintf("{%v}", t.Type.String())
	}
}

