package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
)

type LiteralNode expr.Const[token.Token]

func (node LiteralNode) asConstant() expr.Const[token.Token] {
	return expr.Const[token.Token](node)
}

func (node LiteralNode) Equals(a ast.Ast) bool {
	lit, ok := a.(LiteralNode)
	if !ok {
		return false
	}

	return expr.Const[token.Token](node).Equals(glb_cxt.exprCxt, expr.Const[token.Token](lit))
}

func (node LiteralNode) InOrderTraversal(f func(itoken.Token)) {
	f(node.Name)
}

func (node LiteralNode) NodeType() ast.Type { return Literal }

func literalReduction(nodes ...ast.Ast) ast.Ast {
	return LiteralNode(expr.Const[token.Token]{Name: GetToken(nodes[0])})
}

var literal__IntValue_r = parser.Get(literalReduction).From(IntValue)

var literal__CharValue_r = parser.Get(literalReduction).From(CharValue)

var literal__StringValue_r = parser.Get(literalReduction).From(StringValue)

var literal__FloatValue_r = parser.Get(literalReduction).From(FloatValue)