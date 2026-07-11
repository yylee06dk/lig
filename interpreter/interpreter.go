package interpreter

import (
	"fmt"
	"lig/datatypes"
)

var dummy int

func Interpret(expr datatypes.Expr) (int, error) {
	switch v := expr.(type) {
		case datatypes.Binary:
			res, err := binary(v)
			if err != nil {
				return dummy, fmt.Errorf("InterpError: %w", err)
			}
			return res, nil

		case datatypes.Literal:
			res, err := literal(v)
			if err != nil {
				return dummy, fmt.Errorf("InterpError: %w", err)
			}
			return res, nil

		default:
			// Temporary solve
			return dummy, nil
	}
}

type InterpError struct {
	CurExpr datatypes.Expr
	Msg string
}

func (e *InterpError) Error() string {
	return fmt.Sprintf("In expression %v, error occured: %s", e.CurExpr, e.Msg)
}

func binary(expr datatypes.Binary) (int, error) {
	operator := expr.Operator

	left, leftErr := Interpret(expr.Left)
	if leftErr != nil {
		return dummy, leftErr
	}

	right, rightErr := Interpret(expr.Right)
	if rightErr != nil {
		return dummy, rightErr
	}

	// Need to implement runtime error check
	//left, right = runtimeCheck(left, right)

	switch operator {
		case datatypes.Add:
			return left + right, nil
		case datatypes.Sub:
			return left - right, nil
		case datatypes.Mult:
			return left * right, nil
		case datatypes.Div:
			return left / right, nil
		default:
			// Temporary fix
			return 0, nil
	}
}
/*
func runtimeCheck(left any, right any) (int, int){
	return left, right
}
*/

func literal(expr datatypes.Literal) (int, error) {
	switch v := expr.Value.(type) {
		case int:
			return v, nil
		default:
			return 0, &InterpError{expr, "Expected type int."}
	}
}
