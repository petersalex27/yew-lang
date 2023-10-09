package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
)

type DataNode struct {
	data expr.Expression[token.Token]
}

func (node DataNode) getExpression() ExpressionNode { 
	return ExpressionNode{node.data} 
}

func (n DataNode) Equals(a ast.Ast) bool {
	n2, ok := a.(DataNode)
	if !ok {
		return false
	}

	return n.data.Equals(glb_cxt.exprCxt, n2.data)
}

func (n DataNode) NodeType() ast.Type { return Data }

func (n DataNode) InOrderTraversal(f func(itoken.Token)) {
	for _, dat := range n.data.Collect() {
		f(dat)
	}
}

/*
data          ::= constructor
                  | data expr
                  | '(' data ')'
*/

func dataFromConstructorReduction(nodes ...ast.Ast) ast.Ast {
	nodePtr := new(DataNode)
	nodePtr.data = nil
	toData := func(tok itoken.Token) {
		dat := expr.Const[token.Token]{Name: tok.(token.Token)}
		if nodePtr.data == nil {
			nodePtr.data = dat
		} else {
			nodePtr.data = expr.Apply[token.Token](nodePtr.data, dat)
		}
	}

	getConstructor(nodes[0]).InOrderTraversal(toData)
	return *nodePtr
}

func dataAppendExprReduction(nodes ...ast.Ast) ast.Ast {
	data := nodes[0].(DataNode)
	ex := getExpression(nodes[1]).Expression
	data.data = expr.Apply[token.Token](data.data, ex)
	return data
}

var data__constructor_r = parser.
	Get(dataFromConstructorReduction).From(Constructor)

var data__data_expr_r = parser.
	Get(dataAppendExprReduction).From(Data, Expr)

var data__enclosed_r = parser.
	Get(grab_enclosed).From(LeftParen, Data, RightParen)
