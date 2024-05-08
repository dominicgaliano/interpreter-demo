package ast

import "github.com/dominicgaliano/interpreter-demo/token"

// Node represents any node in the abstract syntax tree (AST).
// It defines a debug method to retrieve the literal value of the token
// associated with the node.
type Node interface {
	TokenLiteral() string
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
// Ex. 2 + 2
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

// ReturnStatement represents a return statement in the AST.
// The value of the return statement is an expression.
type ReturnStatement struct {
	Token token.Token // the token.RETURN token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
