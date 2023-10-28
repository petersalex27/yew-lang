package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/token"
)

func astToExpression(a ast.Ast) ExpressionNode { return a.(ExpressionNode) }

/*
expr          ::= val
                  | data
                  | funcName
                  | letIn
                  | application
                  | pattern
                  | exprWhere
                  | judgement
                  | '(' expr ')'
*/

// expr <- val
var expr__val_r = parser.Get(unwrapSomeExpression).From(Val)

// expr <- data
var expr__data_r = parser.Get(unwrapSomeExpression).From(Data)

// expr <- funcName
// var expr__funcName_r = parser.Get().From(funcName)

// expr <- letIn
// var expr__letIn_r = parser.Get().From(LetIn)

// expr <- application
var expr__application_r = parser.Get(unwrapSomeExpression).From(Application)

// expr <- pattern
var expr__pattern_r = parser.Get(unwrapSomeExpression).From(Pattern)

// expr <- exprWhere
// var expr__exprWhere_r = parser.Get().From(Val)

// expr <- judgement
var expr__judgement_r = parser. 
	Get(judgementToExpression).From(TypeJudgement)

// expr <- '(' expr ')'
var expr__enclosed_r = parser. 
	Get(grab_enclosed).From(LeftParen, Expr, RightParen)

type expressionNodeTypes interface {
	getExpression() ExpressionNode
}

type ExpressionNode struct{ expr.Expression[token.Token] }

type SomeExpression struct {
	ty ast.Type
	expr.Expression[token.Token]
}

func unwrapSomeExpression(nodes ...ast.Ast) ast.Ast {
	return ExpressionNode{nodes[0].(SomeExpression).Expression}
}

func (e ExpressionNode) getExpression() ExpressionNode { return e }

func getExpression(node ast.Ast) ExpressionNode {
	return node.(ExpressionNode)
}

func getApplication(node ast.Ast) expr.Application[token.Token] {
	return getExpression(node).Expression.(expr.Application[token.Token])
}

func (e1 ExpressionNode) Equals(a ast.Ast) bool {
	e2, ok := a.(ExpressionNode)
	if !ok {
		return false
	}

	return e1.Expression.Equals(glb_cxt.exprCxt, e2.Expression)
}

func (e ExpressionNode) NodeType() ast.Type { return Expr }

func (e ExpressionNode) InOrderTraversal(f func(itoken.Token)) {
	elems := e.Expression.Collect()
	for _, elem := range elems {
		f(elem)
	}
}

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
	return e.ty == e2.ty && e.Expression.Equals(glb_cxt.exprCxt, e2.Expression)
}
