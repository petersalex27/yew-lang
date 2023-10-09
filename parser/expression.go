package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
)

type ExpressionNode struct{ expr.Expression[token.Token] }

type SomeExpression struct {
	ty ast.Type
	expr.Expression[token.Token]
}

func asExpression(node ast.Ast) expr.Expression[token.Token] {
	return getExpression(node).Expression
}

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
