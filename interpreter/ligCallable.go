package interpreter

import (
	dt "lig/datatypes"
)

type callable interface {
	arity() int
	call(argList []any) any
}

type returnCheck interface {
	dummyReturn()
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

func (c *closure) call(interp *Interpreter, argList []dt.Expr) (any, *RuntimeError) {
	newEnv := NewEnvironment(c.env)
	for idx, expr := range argList {
		name := c.paramList[idx]
		value, err := interp.interpretExpr(expr)
		if err != nil {
			return nil, err
		}
		newEnv.define(name, value)
	}
	nI := NewInterpreter(newEnv)
	err := nI.interpretStmt(c.code)
	if err != nil && err.ErrToken.Type == dt.Returning {
		return err.ErrToken.Value, nil
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}

// Functions ------------------------------------------------------
