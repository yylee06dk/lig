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

type BlockStmt struct {
	Statements []Stmt
}

type IfStmt struct {
	Condition Expr
	ThenStmt  Stmt
	ElseStmt  Stmt // Can be nil
}

type WhileStmt struct {
	Condition Expr
	Code      Stmt
}

type ForStmt struct {
	InitStmt  Stmt
	Condition Expr
	UpdStmt   Stmt
	Code      Stmt
}

type FunStmt struct {
	Name      Token
	ParamList []Token
	Code      *BlockStmt
}

type ReturnStmt struct {
	RetExpr Expr
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

func (b *BlockStmt) String() string {
	return fmt.Sprintf("Block[ %v ]", b.Statements)
}

func (i *IfStmt) String() string {
	return fmt.Sprintf("If[ %v, %v, %v ]", i.Condition, i.ThenStmt, i.ElseStmt)
}

func (w *WhileStmt) String() string {
	return fmt.Sprintf("While[ %v, %v ]", w.Condition, w.Code)
}

func (f *ForStmt) String() string {
	return fmt.Sprintf("For[ %v, %v, %v, %v ]", f.InitStmt, f.Condition, f.UpdStmt, f.Code)
}

func (f *FunStmt) String() string {
	return fmt.Sprintf("Fun[ %v | %v | %v", f.Name.Name, f.ParamList, f.Code)
}

func (r *ReturnStmt) String() string {
	return fmt.Sprintf("Return [ %v ]", r.RetExpr)
}

func (p *PrintStmt) dummyStmt()  {}
func (e *ExprStmt) dummyStmt()   {}
func (v *VarStmt) dummyStmt()    {}
func (a *AssignStmt) dummyStmt() {}
func (b *BlockStmt) dummyStmt()  {}
func (i *IfStmt) dummyStmt()     {}
func (w *WhileStmt) dummyStmt()  {}
func (f *ForStmt) dummyStmt()    {}
func (f *FunStmt) dummyStmt()    {}
func (r *ReturnStmt) dummyStmt() {}
