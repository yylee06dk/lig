package interpreter

import (
	"fmt"
	dt "lig/datatypes"
	"strings"
	"time"
)

var dummy any

type Interpreter struct {
	Env *Environment
}

type clockFunc struct{}

func (c *clockFunc) arity() int {
	return 0
}

func (c *clockFunc) call(i *Interpreter, argList []any) (any, *RuntimeError) {
	return float64(time.Now().UnixMilli()) / 1000.0, nil
}

func NewInterpreter(parent *Environment) *Interpreter {
	env := newEnvironment(parent)
	if parent == nil { // defining global
		c := clockFunc{}
		env.define("clock", &c)
	}
	return &Interpreter{env}
}

func (i *Interpreter) Interpret(stmtSlice []dt.Stmt) *RuntimeError {
	// Empty stmtSlice
	if stmtSlice == nil {
		return nil
	}

	for _, stmt := range stmtSlice {
		err := i.interpretStmt(stmt)
		if err != nil {
			return err
		} // We don't check if the error is a returning error since we won't run function body here
	}
	return nil
}

func (i *Interpreter) interpretStmt(stmt dt.Stmt) *RuntimeError {
	switch v := stmt.(type) {
	case *dt.VarStmt:
		name := v.Name
		initVal := v.Initializer
		if initVal == nil {
			i.Env.define(name.Name, initVal)
		}
		res, err := i.interpretExpr(initVal)
		if err != nil {
			return err
		}
		i.Env.define(name.Name, res)

	case *dt.FunStmt:
		name := v.Name
		function := &closure{v.ParamList, v.Code, i.Env}
		i.Env.define(name.Name, function)

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
		// Checking if assignment target is defined
		curVal := i.Env.get(v.Name)
		if _, ok := curVal.(*RuntimeError); ok {
			return &RuntimeError{v.Name, fmt.Sprintf("Undeclared variable[%s] tried to be assigned", v.Name.Name)}
		}
		rightVal, err := i.interpretExpr(v.Expression)
		if err != nil {
			return err
		}
		i.Env.assign(v.Name, rightVal)

	case *dt.BlockStmt:
		newInterp := NewInterpreter(i.Env)
		err := newInterp.Interpret(v.Statements)
		if err != nil {
			return err
		}

	case *dt.IfStmt:
		condition, err := i.interpretExpr(v.Condition)
		if err != nil {
			return err
		}
		if isTruthy(condition) {
			return i.interpretStmt(v.ThenStmt)
		}
		if v.ElseStmt != nil && !isTruthy(condition) {
			return i.interpretStmt(v.ElseStmt)
		}

	case *dt.WhileStmt:
		condition, err := i.interpretExpr(v.Condition)
		if err != nil {
			return err
		}
		for ; isTruthy(condition); condition, err = i.interpretExpr(v.Condition) {
			if err != nil {
				return err
			}
			runtimeErr := i.interpretStmt(v.Code)
			if runtimeErr != nil {
				return runtimeErr
			}
		}

	case *dt.ForStmt:
		var err *RuntimeError
		if v.InitStmt != nil {
			err = i.interpretStmt(v.InitStmt)
			if err != nil {
				return err
			}
		}
		condition, condErr := i.interpretExpr(v.Condition)
		if condErr != nil {
			return condErr
		}

		updExist := v.UpdStmt != nil

		for ; isTruthy(condition); condition, condErr = i.interpretExpr(v.Condition) {
			if condErr != nil {
				return condErr
			}
			runtimeErr := i.interpretStmt(v.Code)
			if runtimeErr != nil {
				return runtimeErr
			}
			if updExist {
				runtimeErr = i.interpretStmt(v.UpdStmt)
				if runtimeErr != nil {
					return runtimeErr
				}
			}
		}

	case *dt.ReturnStmt:
		retVal, err := i.interpretExpr(v.RetExpr)
		if err != nil {
			return err
		}
		return &RuntimeError{dt.Token{Type: dt.Returning, Value: retVal}, ""}

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

	case *dt.Call:
		name, err := i.interpretExpr(v.Name) // Returns closure!! (when given good code)
		if err != nil {
			return nil, err
		}

		var res []any
		for _, expr := range v.ArgList {
			value, err := i.interpretExpr(expr)
			if err != nil {
				return nil, err
			}
			res = append(res, value)
		}
		if closure, ok := name.(callable); !ok {
			return nil, &RuntimeError{v.Paren, "Expected callable name"}
		} else {
			retVal, err := closure.call(i, res)
			if err != nil {
				return nil, err
			}
			return retVal, nil
		}

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
