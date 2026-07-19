package interpreter

import (
	"strings"
	"fmt"
	dt "lig/datatypes"
)

var dummy any

func Interpret(expr dt.Expr) (any, error) {
	if (expr == dt.End{}) {
		return nil, nil
	}
	switch v := expr.(type) {
		case dt.Binary:
			res, err := binary(v)
			if err != nil {
				return dummy, fmt.Errorf("RuntimeError: %w", err)
			}
			return res, nil

		case dt.Literal:
			// Need to deal with keywords
			res, err := literal(v)
			if err != nil {
				return dummy, fmt.Errorf("RuntimeError: %w", err)
			}
			return res, nil

		case dt.Unary:
			res, err := unary(v)
			if err != nil {
				return dummy, fmt.Errorf("RuntimeError: %w", err)
			}
			return res, nil

		default:
			return dummy, &RuntimeError{expr, fmt.Sprintf("RuntimeError: Unexpected expression")}
	}
}

type RuntimeError struct {
	CurExpr dt.Expr
	Msg string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("In expression %v, error occured: %s", e.CurExpr, e.Msg)
}

func binary(expr dt.Binary) (any, error) {
	operator := expr.Operator

	left, leftErr := Interpret(expr.Left)
	if leftErr != nil {
		return dummy, leftErr
	}

	right, rightErr := Interpret(expr.Right)
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
							 strings.Compare(leftVal, rightVal) == 0 , nil
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
							 strings.Compare(leftVal, rightVal) == 0 , nil
			} else {
				rightVal := rightVal.(int)

				return temp <= rightVal, nil
			}

		default:
			return dummy, &RuntimeError{expr, fmt.Sprintf("Unknown Operator: [%v]", operator)}
	}
}

func literal(expr dt.Literal) (any, error) {
	switch v := expr.Value.(type) {
		case int, string:
			return v, nil
		default:
			return 0, &RuntimeError{expr, fmt.Sprintf("Expected type int or string, received type: %T", v)}
	}
	//return 0, &RuntimeError{expr, fmt.Sprintf("UNREACHABLE!!! in literal: %v", expr)}
}

func unary(expr dt.Unary) (any, error) {
	right, rightErr := Interpret(expr.Right)
	if rightErr != nil {
		return dummy, rightErr
	}

	rightVal, runtimeErr := runtimeCheckUnary(right, expr)
	if runtimeErr != nil {
		return dummy, runtimeErr
	}

	switch expr.Operator {
		case dt.Sub:
			rightVal := rightVal.(int)
			return -rightVal, nil
		case dt.Bang:
			return !isTruthy(rightVal), nil
		default:
			return dummy, &RuntimeError{expr, fmt.Sprintf("Unrechable, or you maybe added another unary operator: %v", expr)}
	}
	//return 0, &RuntimeError{expr, fmt.Sprintf("UNREACHABLE!!! in unary: %v", expr)}
}

func runtimeCheckBinary(left any, right any, expr dt.Binary) (any, any, error){
	var leftVal any
	var rightVal any

	switch expr.Operator {
		case dt.Add, dt.Sub, dt.Mult, dt.Div:
			ltemp, okLeft := left.(int)
			rtemp, okRight := right.(int)
			if !okLeft || !okRight {
				return 0, 0, &RuntimeError{expr, fmt.Sprintf("Operands of %v must be type of int, received: %T, %T", expr.Operator, left, right)}
			}

			leftVal = ltemp
			rightVal = rtemp

		case dt.AddAdd, dt.BangEqual, dt.EqualEqual, dt.Greater, dt.GreaterEqual, dt.Less, dt.LessEqual:
			ltempInt, okLeft := left.(int)
			rtempInt, okRight := right.(int)
			if !okLeft || !okRight {
				return 0, 0, &RuntimeError{expr, fmt.Sprintf("Operands of %v must be type of int or string, received: %T, %T", expr.Operator, left, right)}
			} else {
				leftVal = ltempInt
				rightVal = rtempInt
				break
			}

			ltempStr, okLeft := left.(string)
			rtempStr, okRight := right.(string)
			if !okLeft || !okRight {
				return 0, 0, &RuntimeError{expr, fmt.Sprintf("Operands of %v must be type of int or string, received: %T, %T", expr.Operator, left, right)}
			}
			leftVal = ltempStr
			rightVal = rtempStr

		default:
			return 0, 0, &RuntimeError{expr, fmt.Sprintf("Unknown or not yet implemented binary operator: %v", expr.Operator)}
	}
	return leftVal, rightVal, nil
}

func runtimeCheckUnary(right any, expr dt.Unary) (any, error){
	var rightVal any
	switch expr.Operator {
		case dt.Sub:
			rtemp, okRight := right.(int)
			if !okRight {
				return 0, &RuntimeError{expr, fmt.Sprintf("Operands of %v must be type of int, received: %T", expr.Operator, right)}
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
			if v == 0 { return false }
			return true
		case string:
			if strings.Compare(v, "") == 0 { return false }
			return true
		default:
			fmt.Println("Need to add more to this!(isTruthy)")
			return false
	}
}
