package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
)

type LiteralNode struct{expr.Expression[token.Token]}

func (node LiteralNode) getExpression() ExpressionNode { 
	return ExpressionNode{node.Expression} 
}

func (node LiteralNode) asConstant() expr.Const[token.Token] {
	return node.Expression.(expr.Const[token.Token])
}

func (node LiteralNode) asList() expr.List[token.Token] {
	return node.Expression.(expr.List[token.Token])
}

func (node LiteralNode) Equals(a ast.Ast) bool {
	lit, ok := a.(LiteralNode)
	if !ok {
		return false
	}

	return node.Expression.Equals(glb_cxt.exprCxt, lit.Expression)
}

func (node LiteralNode) InOrderTraversal(f func(itoken.Token)) {
	for _, tok := range node.Expression.Collect() {
		f(tok)
	}
}

func (node LiteralNode) NodeType() ast.Type { return Literal }

func literalConstReduction(nodes ...ast.Ast) ast.Ast {
	return LiteralNode{expr.Const[token.Token]{Name: GetToken(nodes[0])}}
}

func literalFromLiteralArrayReduction(nodes ...ast.Ast) ast.Ast {
	e := nodes[0].(ArrayNode).getExpression().Expression
	return LiteralNode{e}
}

/*
literal       ::= INT_VALUE 
                  | CHAR_VALUE 
                  | STRING_VALUE 
                  | FLOAT_VALUE
                  | literalArray
*/

var literal__IntValue_r = parser.Get(literalConstReduction).From(IntValue)

var literal__CharValue_r = parser.Get(literalConstReduction).From(CharValue)

var literal__StringValue_r = parser.Get(literalConstReduction).From(StringValue)

var literal__FloatValue_r = parser.Get(literalConstReduction).From(FloatValue)

var literal__literalArray_r = parser.
	Get(literalFromLiteralArrayReduction).
	From(FloatValue)