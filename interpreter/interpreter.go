package interpreter

import (
	"fmt"
	dt "lig/datatypes"
	"strings"
)

var dummy any

type Interpreter struct {
	CodeBlock []dt.Stmt
	Env       *Environment
}

func NewInterpreter(codeBlk []dt.Stmt, env *Environment) *Interpreter {
	return &Interpreter{codeBlk, env}
}

func (i *Interpreter) Interpret() *RuntimeError {
	// Empty stmtSlice
	if i.CodeBlock == nil {
		return nil
	}

	return i.interpretDecl()
}

func (i *Interpreter) interpretDecl() *RuntimeError {
	for _, stmt := range i.CodeBlock {
		switch v := stmt.(type) {
		case *dt.VarStmt:
			name := v.Name
			initVal := v.Initializer
			if initVal == nil {
				i.Env.define(name, initVal)
			}
			res, err := i.interpretExpr(initVal)
			if err != nil {
				return err
			}
			i.Env.define(name, res)

		default:
			err := i.interpretStmt(stmt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Interpreter) interpretStmt(stmt dt.Stmt) *RuntimeError {
	switch v := stmt.(type) {
	case *dt.PrintStmt:
		res, err := i.interpretExpr(v.Expression)
		if err != nil {
			return err
		}
		fmt.Println(res)

	case *dt.ExprStmt:
		_, err := i.interpretExpr(v.Expression)
		if err != nil {
			return err
		}

	case *dt.AssignStmt:
		rightVal, err := i.interpretExpr(v.Expression)
		if err != nil {
			return err
		}
		i.Env.define(v.Name, rightVal)

	default:
		fmt.Println("Implement more statement types")
		// Could deal with this better.. But since I'm gonna add more statements, keep it like this for now.
	}
	return nil
}

func (i *Interpreter) interpretExpr(expr dt.Expr) (any, *RuntimeError) {
	switch v := expr.(type) {
	case *dt.Binary:
		res, err := i.binary(v)
		if err != nil {
			return dummy, err
		}
		return res, nil

	case *dt.Literal:
		// Need to deal with keywords
		res, err := i.literal(v)
		if err != nil {
			return dummy, err
		}
		return res, nil

	case *dt.Unary:
		res, err := i.unary(v)
		if err != nil {
			return dummy, err
		}
		return res, nil

	case *dt.Grouping:
		return i.interpretExpr(v.Inner)

	case *dt.Variable:
		retVal := i.Env.get(v.Name)
		if err, ok := retVal.(*RuntimeError); ok {
			return dummy, err
		}
		return retVal, nil

	default:
		return dummy, &RuntimeError{dt.Token{}, fmt.Sprintf("RuntimeError: Unexpected expression")}
	}
}

type RuntimeError struct {
	ErrToken dt.Token
	Msg      string
}

func (e *RuntimeError) Error() string {
	return e.Msg
}

func (i *Interpreter) binary(expr *dt.Binary) (any, *RuntimeError) {
	operator := expr.Operator.Type

	left, leftErr := i.interpretExpr(expr.Left)
	if leftErr != nil {
		return dummy, leftErr
	}

	right, rightErr := i.interpretExpr(expr.Right)
	if rightErr != nil {
		return dummy, rightErr
	}

	leftVal, rightVal, runtimeErr := runtimeCheckBinary(left, right, expr)
	if runtimeErr != nil {
		return dummy, runtimeErr
	}

	switch operator {
	case dt.Add:
		leftVal := leftVal.(int)
		rightVal := rightVal.(int)
		return leftVal + rightVal, nil
	case dt.AddAdd:
		leftVal := leftVal.(string)
		rightVal := rightVal.(string)
		return leftVal + rightVal, nil
	case dt.Sub:
		leftVal := leftVal.(int)
		rightVal := rightVal.(int)
		return leftVal - rightVal, nil
	case dt.Mult:
		leftVal := leftVal.(int)
		rightVal := rightVal.(int)
		return leftVal * rightVal, nil
	case dt.Div:
		leftVal := leftVal.(int)
		rightVal := rightVal.(int)
		return leftVal / rightVal, nil
	case dt.BangEqual:
		temp, leftOk := leftVal.(int)
		if !leftOk {
			leftVal := leftVal.(string)
			rightVal := rightVal.(string)
			return strings.Compare(leftVal, rightVal) != 0, nil
		} else {
			rightVal := rightVal.(int)

			return temp != rightVal, nil
		}
	case dt.EqualEqual:
		temp, leftOk := leftVal.(int)
		if !leftOk {
			leftVal := leftVal.(string)
			rightVal := rightVal.(string)
			return strings.Compare(leftVal, rightVal) == 0, nil
		} else {
			rightVal := rightVal.(int)

			return temp == rightVal, nil
		}
	case dt.Greater:
		temp, leftOk := leftVal.(int)
		if !leftOk {
			leftVal := leftVal.(string)
			rightVal := rightVal.(string)
			return strings.Compare(leftVal, rightVal) == 1, nil
		} else {
			rightVal := rightVal.(int)

			return temp > rightVal, nil
		}
	case dt.GreaterEqual:
		temp, leftOk := leftVal.(int)
		if !leftOk {
			leftVal := leftVal.(string)
			rightVal := rightVal.(string)
			return strings.Compare(leftVal, rightVal) == 1 ||
				strings.Compare(leftVal, rightVal) == 0, nil
		} else {
			rightVal := rightVal.(int)

			return temp >= rightVal, nil
		}
	case dt.Less:
		temp, leftOk := leftVal.(int)
		if !leftOk {
			leftVal := leftVal.(string)
			rightVal := rightVal.(string)
			return strings.Compare(leftVal, rightVal) == -1, nil
		} else {
			rightVal := rightVal.(int)

			return temp < rightVal, nil
		}
	case dt.LessEqual:
		temp, leftOk := leftVal.(int)
		if !leftOk {
			leftVal := leftVal.(string)
			rightVal := rightVal.(string)
			return strings.Compare(leftVal, rightVal) == -1 ||
				strings.Compare(leftVal, rightVal) == 0, nil
		} else {
			rightVal := rightVal.(int)

			return temp <= rightVal, nil
		}

	default:
		return dummy, &RuntimeError{expr.Operator, fmt.Sprintf("Unknown binary operator: [%v]", operator)}
	}
}

func (i *Interpreter) literal(expr *dt.Literal) (any, *RuntimeError) {
	switch v := expr.Value.(type) {
	case int, string, bool:
		return v, nil
	default:
		return 0, &RuntimeError{dt.Token{}, fmt.Sprintf("Expected type int or string or bool, received type: %T", v)}
	}
	//return 0, &RuntimeError{expr, fmt.Sprintf("UNREACHABLE!!! in literal: %v", expr)}
}

func (i *Interpreter) unary(expr *dt.Unary) (any, *RuntimeError) {
	right, rightErr := i.interpretExpr(expr.Right)
	if rightErr != nil {
		return dummy, rightErr
	}

	rightVal, runtimeErr := runtimeCheckUnary(right, expr)
	if runtimeErr != nil {
		return dummy, runtimeErr
	}

	switch expr.Operator.Type {
	case dt.Sub:
		rightVal := rightVal.(int)
		return -rightVal, nil
	case dt.Bang:
		return !isTruthy(rightVal), nil
	default:
		return dummy, &RuntimeError{expr.Operator, fmt.Sprintf("Unreachable, or you maybe added another unary operator: %v", expr)}
	}
}

func runtimeCheckBinary(left any, right any, expr *dt.Binary) (any, any, *RuntimeError) {
	var leftVal any
	var rightVal any

	switch expr.Operator.Type {
	case dt.Add, dt.Sub, dt.Mult, dt.Div:
		ltemp, okLeft := left.(int)
		rtemp, okRight := right.(int)
		if !okLeft || !okRight {
			return 0, 0, &RuntimeError{expr.Operator, fmt.Sprintf("Operands of %v must be type of int, received: %T, %T", expr.Operator.Type, left, right)}
		}

		leftVal = ltemp
		rightVal = rtemp

	case dt.BangEqual, dt.EqualEqual, dt.Greater, dt.GreaterEqual, dt.Less, dt.LessEqual:
		ltempInt, okLeft := left.(int)
		rtempInt, okRight := right.(int)
		if !okLeft || !okRight {
			ltempStr, okLeft := left.(string)
			rtempStr, okRight := right.(string)
			if !okLeft || !okRight {
				return 0, 0, &RuntimeError{expr.Operator, fmt.Sprintf("Operands of %v must both be type of int or string, received: %T, %T", expr.Operator.Type, left, right)}
			}
			leftVal = ltempStr
			rightVal = rtempStr
		} else {
			leftVal = ltempInt
			rightVal = rtempInt
			break
		}

	case dt.AddAdd:
		ltempStr, okLeft := left.(string)
		rtempStr, okRight := right.(string)
		if !okLeft || !okRight {
			return 0, 0, &RuntimeError{expr.Operator, fmt.Sprintf("Operands of %v must be type of string, received: %T, %T", expr.Operator.Type, left, right)}
		}
		leftVal = ltempStr
		rightVal = rtempStr

	default:
		return 0, 0, &RuntimeError{expr.Operator, fmt.Sprintf("Unknown or not yet implemented binary operator: %v", expr.Operator.Type)}
	}
	return leftVal, rightVal, nil
}

func runtimeCheckUnary(right any, expr *dt.Unary) (any, *RuntimeError) {
	var rightVal any
	switch expr.Operator.Type {
	case dt.Sub:
		rtemp, okRight := right.(int)
		if !okRight {
			return 0, &RuntimeError{expr.Operator, fmt.Sprintf("Operands of %v must be type of int, received: %T", expr.Operator.Type, right)}
		}
		rightVal = rtemp

	case dt.Bang:
		rightVal = !isTruthy(right)
	}

	return rightVal, nil
}

func isTruthy(value any) bool {
	switch v := value.(type) {
	case int:
		if v == 0 {
			return false
		}
		return true
	case string:
		if v == "" {
			return false
		}
		return true
	case bool:
		return v
	default:
		// Hope its unreachable
		fmt.Println("Need to add more to this!(isTruthy)")
		return false
	}
}
