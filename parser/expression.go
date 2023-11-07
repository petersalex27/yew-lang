// =============================================================================
// Author-Date: Alex Peters - 2023
//
// Content:
// this file contains expression production rules and functions, and it contains
// expression and some-expression nodes and ast.Ast implementations. It also has
// a few expression-related utility functions--mostly wrappers for type
// assertions.
//
// Grammar:
//
//	expr ::= val
//					 | data
//					 | funcName
//					 | letExpr
//					 | application
//					 | pattern
//					 | whereExpr
//					 | judgement
//					 | '(' expr ')'
//					 | ( '(' ) expr INDENT(_)
//
// Notes: -
// =============================================================================
package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

// represents an expression: one that can be used by the 
// <whatever>/yew-packages/expr module
type ExpressionNode struct{ expr.Expression[token.Token] }

// represents an node that is not yet an "expression" node but has been 
// converted (and wrapped inside this ast node) into a form that can be
// used by the <whatever>/yew-packages/expr module
type SomeExpression struct {
	ty ast.Type
	expr.Expression[token.Token]
}

// things that hold expressions useable by the <whatever>/yew-packages/expr
// module
type expressionNodeTypes interface {
	getExpression() ExpressionNode
}

// =============================================================================
// production rules
// =============================================================================

// expr ::= val
var expr__val_r = parser.
	Get(unwrapSomeExpression).
	From(Val)

// expr ::= data
var expr__data_r = parser.
	Get(unwrapSomeExpression).
	From(Data)

// expr ::= funcName
var expr__funcName_r = parser.
	Get(constantExpressionProduction).
	From(FuncName)

// expr ::= letExpr
var expr__letExpr_r = parser.
	Get(unwrapSomeExpression).
	From(LetExpr)

// expr ::= exprWhere
var expr__exprWhere_r = parser.
	Get(unwrapSomeExpression).
	From(WhereExpr)

// expr ::= application
var expr__application_r = parser.
	Get(unwrapSomeExpression).
	From(Application)

// expr ::= pattern
var expr__pattern_r = parser.
	Get(unwrapSomeExpression).
	From(Pattern)

// expr ::= judgement
var expr__judgement_r = parser.
	Get(judgementToExpressionProduction).
	From(TypeJudgement)

// expr ::= '(' expr ')'
var expr__enclosed_r = parser.
	Get(parenEnclosedProduction).
	From(LeftParen, Expr, RightParen)

var expr__LeftParen_expr_Indent_r = parser.
	Get(grabInitialProduction).
	When(LeftParen).
	From(Expr, Indent)

// =============================================================================
// production functions
// =============================================================================

// transforms a judgment into an expression
func judgementToExpressionProduction(nodes ...ast.Ast) ast.Ast {
	return ExpressionNode{
		bridge.JudgementAsExpression[token.Token, expr.Expression[token.Token]](
			nodes[0].(JudgementNode),
		),
	}
}

// transforms a some-expression node to an expression node
func unwrapSomeExpression(nodes ...ast.Ast) ast.Ast {
	return ExpressionNode{nodes[0].(SomeExpression).Expression}
}

// creates a name (i.e., expr.Const[token.Token]) and wraps it in an 
// expression node
func constantExpressionProduction(nodes ...ast.Ast) ast.Ast {
	const nameIndex int = 0
	name := GetToken(nodes[nameIndex])
	nameConst := expr.Const[token.Token]{Name: name}
	return ExpressionNode{nameConst}
}

// =============================================================================
// expression node ast.Ast implementation
// =============================================================================

func (e1 ExpressionNode) Equals(a ast.Ast) bool {
	e2, ok := a.(ExpressionNode)
	if !ok {
		return false
	}

	return e1.Expression.Equals(globalContext__.exprCxt, e2.Expression)
}

func (e ExpressionNode) NodeType() ast.Type { return Expr }

func (e ExpressionNode) InOrderTraversal(f func(itoken.Token)) {
	elems := e.Expression.Collect()
	for _, elem := range elems {
		f(elem)
	}
}

// =============================================================================
// some expression node ast.Ast implementation
// =============================================================================

func (e SomeExpression) NodeType() ast.Type { return e.ty }

func (e SomeExpression) InOrderTraversal(f func(itoken.Token)) {
	elems := e.Expression.Collect()
	for _, elem := range elems {
		f(elem)
	}
}

func (e SomeExpression) Equals(a ast.Ast) bool {
	e2, ok := a.(SomeExpression)
	if !ok {
		return false
	}
	return e.ty == e2.ty && e.Expression.Equals(globalContext__.exprCxt, e2.Expression)
}

// =============================================================================
// utils
// =============================================================================

// casts an ast.Ast to an ExpressionNode
func astToExpression(a ast.Ast) ExpressionNode { return a.(ExpressionNode) }

// returns expression node
func (e ExpressionNode) getExpression() ExpressionNode { return e }

// returns expression node
func getExpression(node ast.Ast) ExpressionNode {
	return node.(ExpressionNode)
}

