package ast

import (
	"bytes"

	"github.com/dominicgaliano/interpreter-demo/token"
)

// Node represents any node in the abstract syntax tree (AST).
// It defines a debug method, TokenLiteral and String, to retrieve the literal
// value of the token associated with the node and to print the node, respectively.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents any statement node in the AST.
// A statement is expresses some action to be carried out.
// Ex. foo := 10
// It embeds the Node interface and defines a method specific to statement
// nodes that can help the Go compiler distinguish it from the Expression
// interface.
type Statement interface {
	Node
	statementNode()
}

// Expression represents any expression node in the AST.
// An expression is an entity that may be evaluated to determine its value.
// Expression with prefix operators
// Ex. -5, !true
// Expression with infix operators
// Ex. 5 + 5, 5 * 5
// Expression with comparison operators
// Ex. 5 < 10, 5 == 5
// Expression with function calls
// Ex. add(5, 5), add(5 + 5, 5 * 5)
// It embeds the Node interface and defines a method specific to  statement
// nodes that can help the Go compiler distinguish it from the Expression
// interface.
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement represents a variable assignment statement in the AST.
// The LET token is stored to be used in the TokenLiteral() method required
// for all AST nodes.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// Identifier represent a variable identifier in the AST.
// Ex. let x = 5; => Token = token.IDENT, Value = x
// Identifiers can be used as expressions in some cases, so it implements the
// Expression interface.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// ReturnStatement represents a return statement in the AST.
// The value of the return statement is an expression.
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement represents an expression statement in the AST.
// The expression is stored as a field and can be any valid expression.
// ExpressionStatement is a wrapper around an expression that allows it to
// be used as a statement.
// Ex. 5 + 5;
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral represents an integer literal in the AST.
// The literal value is stored in the Value field.
// IntegerLiteral implements the Expression interface.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return il.Token.Literal }

