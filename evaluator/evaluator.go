package evaluator

import (
	"github.com/dominicgaliano/interpreter-demo/ast"
	"github.com/dominicgaliano/interpreter-demo/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
        return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
        return Eval(node.Expression)

	// Expression
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
    var result object.Object

    for _, stmt := range statements {
        result = Eval(stmt)
    }

    return result
}
