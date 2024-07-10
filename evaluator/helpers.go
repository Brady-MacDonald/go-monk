package evaluator

import (
	"fmt"
	"monkey/object"
)

func newError(formatStr string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(formatStr, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

// All objects are truthy expect for NULL/FALSE
func isTruthy(condition object.Object) bool {
	switch condition {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

// Get the instance of the Boolean object (TRUE/FALSE)
func getBoolObj(val bool) *object.Boolean {
	if val {
		return TRUE
	}

	return FALSE
}

// Evaluate the operand of a - prefixExpression.
// Operand must be an object.Integer
func evalMinusPrefix(operand object.Object) object.Object {
	intLit, ok := operand.(*object.Integer)
	if !ok {
		return newError("Invalid operand type: -%s", operand.Type())
	}

	return &object.Integer{Value: -intLit.Value}
}

// Operand should only be an object.Integer or object.Boolean
func evalBangPrefix(operand object.Object) object.Object {
	switch operand {
	case TRUE:
		return FALSE

	case FALSE:
		return TRUE

	case NULL:
		return TRUE

	default:
		return FALSE
	}
}

func evalStringInfix(operator string, left, right object.Object) object.Object {
	leftStr, ok := left.(*object.String)
	if !ok {
		return newError("Left is not a String, Got=%T", left)
	}

	rightStr, ok := right.(*object.String)
	if !ok {
		return newError("Right is not a String, Got=%T", left)
	}

	if operator != "+" {
		return newError("Strings only support the + operator, Got=%s", operator)
	}

	return &object.String{Value: leftStr.Value + rightStr.Value}
}

// Evaluate infix expression if both operands are Object.Integer
func evalIntegerInfix(operator string, left, right object.Object) object.Object {
	leftInt, ok := left.(*object.Integer)
	if !ok {
		return newError("Left infix object not Integer: Got=%T", left)
	}

	rightInt, ok := right.(*object.Integer)
	if !ok {
		return newError("Right infix object not Integer: Got=%T", right)
	}

	switch operator {
	/* Integer Producing Infix Expressions */

	case "+":
		return &object.Integer{Value: leftInt.Value + rightInt.Value}
	case "-":
		return &object.Integer{Value: leftInt.Value - rightInt.Value}
	case "*":
		return &object.Integer{Value: leftInt.Value * rightInt.Value}
	case "/":
		return &object.Integer{Value: leftInt.Value / rightInt.Value}

		/* Boolean Producing Infix Expressions */

	case "<":
		return getBoolObj(leftInt.Value < rightInt.Value)
	case ">":
		return getBoolObj(leftInt.Value > rightInt.Value)
	case "==":
		return getBoolObj(leftInt.Value == rightInt.Value)
	case "!=":
		return getBoolObj(leftInt.Value != rightInt.Value)

	default:
		return newError("Unknown infix operator: %s %s %s", leftInt.Type(), operator, rightInt.Type())
	}
}
