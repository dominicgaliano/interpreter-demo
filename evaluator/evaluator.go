package evaluator

import (
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

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
    case *ast.BlockStatement:
        return evalBlockStatement(node)
    case *ast.ReturnStatement:
        val := Eval(node.ReturnValue)
        return &object.ReturnValue{ Value: val }

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		operator := node.Operator
		return evalInfixExpression(operator, left, right)
    case *ast.IfExpression:
        return evalIfExpression(node)
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

        // stop evaluating statements when return statement has been evaluated
        // since this is called at the program level, upwrap return statement
        if returnValue, ok := result.(*object.ReturnValue); ok {
            return returnValue.Value
        }
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

        // stop evaluating statements when return statement has been evaluated
        // do not unwrap value as it will be used in lower level nests
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

        if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
            return result
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
		return NULL
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
		return NULL
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
    default:
        return NULL
    }

}

func evalInfixIntegerExpression(operator string, left, right object.Object) object.Object {
    leftValue := left.(*object.Integer).Value
    rightValue := right.(*object.Integer).Value

    switch operator {
    case token.PLUS:
        return &object.Integer{ Value: leftValue + rightValue }
    case token.MINUS:
        return &object.Integer{ Value: leftValue - rightValue }
    case token.ASTERISK:
        return &object.Integer{ Value: leftValue * rightValue }
    case token.SLASH:
        return &object.Integer{ Value: leftValue / rightValue }
    case token.GT:
        return nativeBoolToBooleanObject(leftValue > rightValue )
    case token.LT:
        return nativeBoolToBooleanObject(leftValue < rightValue )
    case token.EQ:
        return nativeBoolToBooleanObject(leftValue == rightValue )
    case token.NOT_EQ:
        return nativeBoolToBooleanObject(leftValue != rightValue )
    default:
        return NULL
    }
}

func evalIfExpression(node *ast.IfExpression) object.Object {
    // determine if node.Condition evaluates to a truthy value
    // if it does, evaluate and return node.Consequence
    // if it does not, check if node.Alternative is not null
    // if node.Alternative is not null, eval and return node.Alternative

    condition := Eval(node.Condition)
    if isTruthy(condition) {
        return Eval(node.Consequence)
    }

    if node.Alternative != nil {
        return Eval(node.Alternative)
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



