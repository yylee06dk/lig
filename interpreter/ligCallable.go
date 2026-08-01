package interpreter

import (
	dt "lig/datatypes"
)

type callable interface {
	arity() int
	call(i *Interpreter, argList []any) (any, *RuntimeError)
}

// Functions ------------------------------------------------------

type closure struct { // Like an anonymous function
	paramList []dt.Token
	code      *dt.BlockStmt
	env       *Environment
}

func (c *closure) arity() int {
	return len(c.paramList)
}

func (c *closure) call(interp *Interpreter, argList []any) (any, *RuntimeError) {
	nI := NewInterpreter(c.env)
	newEnv := nI.Env
	for idx, value := range argList {
		name := c.paramList[idx]
		newEnv.define(name.Name, value)
	}
	err := nI.interpretStmt(c.code)
	if err != nil && err.ErrToken.Type == dt.Returning {
		return err.ErrToken.Value, nil
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}

// Functions ------------------------------------------------------
