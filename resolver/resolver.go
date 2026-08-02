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
			s.resolveExpr(v.Initializer)
			s.define(v.Name)
		}

	case *dt.BlockStmt:
		s.startScope()

		for _, stmt := range v.Statements {
			s.resolveStmt(stmt)
		}
		s.endScope()

	case *dt.AssignStmt:
		token := v.GetAssignable()
		s.checkPos(v.NameExpr, token)
		s.resolveExpr(v.Expression)

	case *dt.FunStmt:
		s.declare(v.Name)
		s.define(v.Name)
		s.startScope()
		for _, param := range v.ParamList {
			s.declare(param)
			s.define(param)
		}
		s.resolveStmt(v.Code) // This results as a double enclosing, which leads to un intuitive hop numbers when debugging.
		s.endScope()

	case *dt.ExprStmt:
		s.resolveExpr(v.Expression)

	case *dt.PrintStmt:
		s.resolveExpr(v.Expression)

	case *dt.IfStmt:
		s.resolveExpr(v.Condition)
		s.resolveStmt(v.ThenStmt)
		if v.ElseStmt != nil {
			s.resolveStmt(v.ElseStmt)
		}

	case *dt.WhileStmt:
		s.resolveExpr(v.Condition)
		s.resolveStmt(v.Code)

	case *dt.ForStmt:
		if v.InitStmt != nil {
			s.resolveStmt(v.InitStmt)
		}

		s.resolveExpr(v.Condition)

		if v.UpdStmt != nil {
			s.resolveStmt(v.UpdStmt)
		}

		s.resolveStmt(v.Code)

	case *dt.ReturnStmt:
		s.resolveExpr(v.RetExpr)
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
		s.resolveExpr(v.Name)
		if v.ArgList != nil {
			for _, expr := range v.ArgList {
				s.resolveExpr(expr)
			}
		}

	case *dt.Binary:
		s.resolveExpr(v.Left)
		s.resolveExpr(v.Right)

	case *dt.Literal:

	case *dt.Unary:
		s.resolveExpr(v.Right)

	case *dt.Grouping:
		s.resolveExpr(v.Inner)
	}
}
