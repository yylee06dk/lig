package resolver

import (
	"fmt"
	dt "lig/datatypes"
	"lig/interpreter"
)

// --------------- Stack Implementation ----------------

type scopeStack struct {
	items []map[string]bool
}

func newscopeStack() *scopeStack {
	return &scopeStack{items: nil}
}

func (s *scopeStack) push(val map[string]bool) {
	s.items = append(s.items, val)
}

func (s *scopeStack) pop() map[string]bool {
	if s.isEmpty() {
		return nil
	}
	l := len(s.items) - 1
	val := s.items[l]
	s.items = s.items[:l]
	return val
}

func (s *scopeStack) peek() map[string]bool {
	if s.isEmpty() {
		return nil
	}
	return s.items[len(s.items)-1]
}

func (s *scopeStack) isEmpty() bool {
	return len(s.items) == 0
}

// --------------- Stack Implementation ----------------

func Resolve(ast []dt.Stmt) error {
	s := newscopeStack()
	s.startScope() // Global. We need to define some native functions

	s.declare(dt.Token{Name: "clock"})
	s.define(dt.Token{Name: "clock"})

	err := s.resolve(ast)
	if err != nil {
		return err
	}
	s.endScope()
	return nil
}

func (s *scopeStack) resolve(ast []dt.Stmt) error {
	for _, stmt := range ast {
		err := s.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *scopeStack) declare(token dt.Token) {
	if !s.isEmpty() {
		s.peek()[token.Name] = false
	}
}

func (s *scopeStack) define(token dt.Token) {
	if !s.isEmpty() {
		s.peek()[token.Name] = true
	}
}

func (s *scopeStack) startScope() {
	s.push(make(map[string]bool))
}

func (s *scopeStack) endScope() {
	_ = s.pop()
}

func (s *scopeStack) checkPos(expr dt.Expr, token dt.Token) error {
	l := len(s.items)
	var pos = l - 1

	var checkFound bool = false
	for ; pos >= 0; pos -= 1 {
		if val, exists := s.items[pos][token.Name]; exists {
			if val == false {
				return fmt.Errorf("Uninitialized variable referred: %s", token.Name)
			} else {
				checkFound = true
				interpreter.SetLocals(expr, l-1-pos)
			}
		}
	}
	if !checkFound {
		return fmt.Errorf("Undeclared variable referred: %s", token.Name)
	}
	return nil
}

func (s *scopeStack) resolveStmt(stmt dt.Stmt) error {
	switch v := stmt.(type) {
	case *dt.VarStmt:
		s.declare(v.Name)

		if v.Initializer != nil {
			err := s.resolveExpr(v.Initializer)
			if err != nil {
				return err
			}
			s.define(v.Name)
		}

	case *dt.BlockStmt:
		s.startScope()

		for _, stmt := range v.Statements {
			err := s.resolveStmt(stmt)
			if err != nil {
				return err
			}
		}
		s.endScope()

	case *dt.AssignStmt:
		token := v.GetAssignable()
		s.checkPos(v.NameExpr, token)
		err := s.resolveExpr(v.Expression)
		if err != nil {
			return err
		}

	case *dt.FunStmt:
		s.declare(v.Name)
		s.define(v.Name)
		s.startScope()
		for _, param := range v.ParamList {
			s.declare(param)
			s.define(param)
		}
		err := s.resolveStmt(v.Code) // This results as a double enclosing, which leads to un intuitive hop numbers when debugging.
		if err != nil {
			return err
		}
		s.endScope()

	case *dt.ExprStmt:
		err := s.resolveExpr(v.Expression)
		if err != nil {
			return err
		}

	case *dt.PrintStmt:
		err := s.resolveExpr(v.Expression)
		if err != nil {
			return err
		}

	case *dt.IfStmt:
		err := s.resolveExpr(v.Condition)
		if err != nil {
			return err
		}
		err = s.resolveStmt(v.ThenStmt)
		if err != nil {
			return err
		}
		if v.ElseStmt != nil {
			err := s.resolveStmt(v.ElseStmt)
			if err != nil {
				return err
			}
		}

	case *dt.WhileStmt:
		err := s.resolveExpr(v.Condition)
		if err != nil {
			return err
		}
		err = s.resolveStmt(v.Code)
		if err != nil {
			return err
		}

	case *dt.ForStmt:
		if v.InitStmt != nil {
			err := s.resolveStmt(v.InitStmt)
			if err != nil {
				return err
			}
		}

		err := s.resolveExpr(v.Condition)
		if err != nil {
			return err
		}

		if v.UpdStmt != nil {
			err := s.resolveStmt(v.UpdStmt)
			if err != nil {
				return err
			}
		}

		err = s.resolveStmt(v.Code)
		if err != nil {
			return err
		}

	case *dt.ReturnStmt:
		err := s.resolveExpr(v.RetExpr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *scopeStack) resolveExpr(expr dt.Expr) error {
	switch v := expr.(type) {
	case *dt.Variable:
		if !s.isEmpty() {
			if val, ok := s.peek()[v.Name.Name]; ok && val == false {
				return fmt.Errorf("Uninitialized variable referred in its initializer")
			}
		}
		s.checkPos(v, v.Name)

	case *dt.Call:
		err := s.resolveExpr(v.Name)
		if err != nil {
			return err
		}
		if v.ArgList != nil {
			for _, expr := range v.ArgList {
				err := s.resolveExpr(expr)
				if err != nil {
					return err
				}
			}
		}

	case *dt.Binary:
		err := s.resolveExpr(v.Left)
		if err != nil {
			return err
		}
		err = s.resolveExpr(v.Right)
		if err != nil {
			return err
		}

	case *dt.Literal:

	case *dt.Unary:
		err := s.resolveExpr(v.Right)
		if err != nil {
			return err
		}

	case *dt.Grouping:
		err := s.resolveExpr(v.Inner)
		if err != nil {
			return err
		}
	}
	return nil
}
