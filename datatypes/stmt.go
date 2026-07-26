package datatypes

import (
	"fmt"
)

type Stmt interface {
	dummyStmt()
}

type ExprStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type VarStmt struct {
	Name        Token
	Initializer Expr // Allows nil! (Since initializer is optional)
}

type AssignStmt struct {
	Name       Token
	Expression Expr
}

func (e *ExprStmt) String() string {
	return fmt.Sprintf("Stmt[ %v ]", e.Expression)
}

func (p *PrintStmt) String() string {
	return fmt.Sprintf("Print[ %v ]", p.Expression)
}

func (v *VarStmt) String() string {
	return fmt.Sprintf("Var[ %s, %v ]", v.Name.Name, v.Initializer)
}

func (a *AssignStmt) String() string {
	return fmt.Sprintf("Assign[ %s, %v ]", a.Name.Name, a.Expression)
}

func (p *PrintStmt) dummyStmt()  {}
func (e *ExprStmt) dummyStmt()   {}
func (v *VarStmt) dummyStmt()    {}
func (a *AssignStmt) dummyStmt() {}
