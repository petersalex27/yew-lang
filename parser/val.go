package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/token"
)

type ValNode ExpressionNode

func (node ValNode) getExpression() ExpressionNode {
	return ExpressionNode(node)
}

func (ValNode) NodeType() ast.Type {
	return Val
}

func (node ValNode) InOrderTraversal(f func(token.Token)) {
	node.getExpression().InOrderTraversal(f)
} 

func (node ValNode) Equals(a ast.Ast) bool {
	if val, ok := a.(ValNode); ok {
		return node.Expression.Equals(glb_cxt.exprCxt, val.Expression)
	}
	return false
}

/*
val           ::= literal
                  | array
*/

func valFromExpressionTyped(nodes ...ast.Ast) ast.Ast {
	e := nodes[0].(expressionNodeTypes).getExpression()
	return ValNode(e)
}

// val <- literal
var val__literal_r = parser.Get(valFromExpressionTyped).From(Literal)

// val <- array
var val__array_r = parser.Get(valFromExpressionTyped).From(Array)