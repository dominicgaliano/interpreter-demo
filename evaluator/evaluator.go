package evaluator

import (
	"fmt"

	"github.com/dominicgaliano/interpreter-demo/ast"
	"github.com/dominicgaliano/interpreter-demo/object"
	"github.com/dominicgaliano/interpreter-demo/token"
)

// Singleton instances of NULL, TRUE, and FALSE to ensure a single shared
// instance of these values. Provides a performance increase by reducing
// the need to instantiate these instances more than once.
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
    case *ast.LetStatement:
        val := Eval(node.Value, env)
        if isError(val) {
            return val
        }
        env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		operator := node.Operator
		return evalInfixExpression(operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)
    case *ast.Identifier:
        return evalIdentifier(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		// stop evaluating statements when return statement has been evaluated
		// since this is called at the program level, upwrap return statement
		// same for errors
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result

		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		// stop evaluating statements when return statement has been evaluated
		// or an error has occurred.
		// Do not unwrap value as it will be used in lower level nests
		// ex:
		// if (true) {
		//     if (true) {
		//         // this wil create a RETURN_VALUE_OBJ that must be passed
		//         // to the parent if statement. This will allow the parent
		//         // if statement to know that it should also stop evaluating
		//         return 1;
		//     }
		//     return 0;
		// }

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
	operator string,
	left object.Object,
	right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(operator, left, right)
	// The following equality checks are only intended for comparing object.BOOLEAN_OBJ
	// All other comparisons with return False
	case operator == token.EQ:
		return nativeBoolToBooleanObject(left == right)
	case operator == token.NOT_EQ:
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalInfixIntegerExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case token.PLUS:
		return &object.Integer{Value: leftValue + rightValue}
	case token.MINUS:
		return &object.Integer{Value: leftValue - rightValue}
	case token.ASTERISK:
		return &object.Integer{Value: leftValue * rightValue}
	case token.SLASH:
		return &object.Integer{Value: leftValue / rightValue}
	case token.GT:
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case token.LT:
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case token.EQ:
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	// determine if node.Condition evaluates to a truthy value
	// if it does, evaluate and return node.Consequence
	// if it does not, check if node.Alternative is not null
	// if node.Alternative is not null, eval and return node.Alternative

	condition := Eval(node.Condition, env)
    if isError(condition) {
        return condition
    }

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	}

	if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	// Returns true if an object is "truthy"
	// All objects are truthy expect the following:
	// FALSE, NULL, INTERGER_OBJ with value = 0

	switch obj.Type() {
	case object.BOOLEAN_OBJ:
		return obj == TRUE
	case object.NULL_OBJ:
		return false
	case object.INTEGER_OBJ:
		int_obj := obj.(*object.Integer)
		return int_obj.Value != 0
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    val, ok := env.Get(node.Value)
    if !ok {
        return newError("identifier not found: " + node.Value)
    }
    return val
}

