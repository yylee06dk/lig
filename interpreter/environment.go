package interpreter

import (
	"fmt"
	dt "lig/datatypes"
)

type Environment struct {
	Parent  *Environment
	Current map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{Parent: nil, Current: make(map[string]any)}
}

func (e *Environment) get(name dt.Token) any { //Bad design
	if res, ok := e.Current[name.Name]; !ok {
		if e.Parent == nil {
			return &RuntimeError{name, fmt.Sprintf("Undeclared variable referred: %s", name.Name)}
		}
		return e.Parent.get(name)
	} else {
		return res
	}
}

func (e *Environment) define(name dt.Token, value any) {
	e.Current[name.Name] = value
}
