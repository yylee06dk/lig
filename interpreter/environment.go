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

func (e *Environment) ancestor(hops int) *Environment {
	var targetEnv *Environment = e
	for i := 0; i < hops; i += 1 {
		if e.Parent == nil {
			panic("error ... with calculating hop counts")
		}
		targetEnv = targetEnv.Parent
	}
	return targetEnv
}

func (e *Environment) global() *Environment {
	for e.Parent != nil {
		e = e.Parent
	}
	return e
}

func (e *Environment) getAt(name dt.Token, hops int) (any, *RuntimeError) {
	targetEnv := e.ancestor(hops)
	return targetEnv.Current[name.Name], nil
}

func (e *Environment) assignAt(name dt.Token, hops int, value any) *RuntimeError {
	targetEnv := e.ancestor(hops)
	targetEnv.define(name.Name, value)
	return nil
}

func (e *Environment) define(name string, value any) {
	e.Current[name] = value
}

func (e *Environment) PrintAll(id string) {
	fmt.Println("start", id)
	for e.Parent != nil {
		fmt.Println(e.Current)
		e = e.Parent
	}
	fmt.Println(e.Current)
	fmt.Println("end", id)
}

/*

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

*/
