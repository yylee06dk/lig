package interpreter

import (
	"fmt"
	dt "lig/datatypes"
)

type Environment struct {
	Parent  *Environment
	Current map[string]any
}

func newEnvironment(parent *Environment) *Environment {
	return &Environment{Parent: parent, Current: make(map[string]any)}
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

func (e *Environment) define(name string, value any) {
	e.Current[name] = value
}

func (e *Environment) assign(name dt.Token, value any) *RuntimeError {
	if _, ok := e.Current[name.Name]; ok {
		e.define(name.Name, value)
		return nil
	}
	if e.Parent != nil {
		e.Parent.assign(name, value)
	}
	return &RuntimeError{name, fmt.Sprintf("Undeclared variable referred: %s", name.Name)}
}
