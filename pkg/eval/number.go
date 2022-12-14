package eval

import (
	"math"

	"github.com/qiushiyan/qlang/pkg/object"
)

func evalNumberLiteral(value float64) object.Object {
	return &object.Number{Value: value}
}

func evalNumberInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value

	switch operator {
	case "+":
		return &object.Number{Value: leftVal + rightVal}
	case "-":
		return &object.Number{Value: leftVal - rightVal}
	case "*":
		return &object.Number{Value: leftVal * rightVal}
	case "/":
		return &object.Number{Value: leftVal / rightVal}
	case "%":
		return &object.Number{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return evalBoolean(leftVal < rightVal)
	case "<=":
		return evalBoolean(leftVal <= rightVal)
	case ">":
		return evalBoolean(leftVal > rightVal)
	case ">=":
		return evalBoolean(leftVal >= rightVal)
	case "==":
		return evalBoolean(leftVal == rightVal)
	case "!=":
		return evalBoolean(leftVal != rightVal)
	default:
		return newInfixError(left, operator, right)
	}
}
