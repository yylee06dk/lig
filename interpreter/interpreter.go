package interpreter

import (
	"fmt"
	"lig/datatypes"
)

var dummy any

func Interpret(expr datatypes.Expr) (any, error) {
	switch v := expr.(type) {
		case datatypes.Binary:
			res, err := binary(v)
			if err != nil {
				return dummy, fmt.Errorf("RuntimeError: %w", err)
			}
			return res, nil

		case datatypes.Literal:
			res, err := literal(v)
			if err != nil {
				return dummy, fmt.Errorf("RuntimeError: %w", err)
			}
			return res, nil

		default:
			// Temporary solve
			return dummy, nil
	}
}

type RuntimeError struct {
	CurExpr datatypes.Expr
	Msg string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("In expression %v, error occured: %s", e.CurExpr, e.Msg)
}

func binary(expr datatypes.Binary) (any, error) {
	operator := expr.Operator

	left, leftErr := Interpret(expr.Left)
	if leftErr != nil {
		return dummy, leftErr
	}

	right, rightErr := Interpret(expr.Right)
	if rightErr != nil {
		return dummy, rightErr
	}

	leftVal, rightVal, runtimeErr := runtimeCheck(left, right, expr)
	if runtimeErr != nil {
		
	}

	switch operator {
		case datatypes.Add:
			return leftVal + rightVal, nil
		case datatypes.Sub:
			return leftVal - rightVal, nil
		case datatypes.Mult:
			return leftVal * rightVal, nil
		case datatypes.Div:
			return leftVal / rightVal, nil
		default:
			// Temporary fix
			return dummy, nil
	}
}

func runtimeCheck(left any, right any, expr datatypes.Binary) (int, int, error){
	leftVal, okLeft := left.(int)
	rightVal, okRight := right.(int)
	if !okLeft || !okRight {
		return 0, 0, &RuntimeError{expr, fmt.Sprintf("Operands of %v must be type of int, received: %T, %T", expr.Operator, left, right)}
	}
	return leftVal, rightVal, nil
}


func literal(expr datatypes.Literal) (any, error) {
	switch v := expr.Value.(type) {
		case int:
			return v, nil
		default:
			return 0, &RuntimeError{expr, fmt.Sprintf("Expected type int, received type: %T", v)}
	}
}
