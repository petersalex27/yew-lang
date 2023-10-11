package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
)

/*
data          ::= constructor
                  | data expr
                  | '(' data ')'
*/

func dataFromConstructorReduction(nodes ...ast.Ast) ast.Ast {
	nodePtr := new(SomeExpression)
	nodePtr.ty = Data
	nodePtr.Expression = nil
	toData := func(tok itoken.Token) {
		dat := expr.Const[token.Token]{Name: tok.(token.Token)}
		if nodePtr.Expression == nil {
			nodePtr.Expression = dat
		} else {
			nodePtr.Expression = expr.Apply[token.Token](nodePtr.Expression, dat)
		}
	}

	getConstructor(nodes[0]).InOrderTraversal(toData)
	return *nodePtr
}

func dataAppendExprReduction(nodes ...ast.Ast) ast.Ast {
	data := nodes[0].(SomeExpression)
	ex := getExpression(nodes[1]).Expression
	data.Expression = expr.Apply[token.Token](data.Expression, ex)
	return data
}

var data__constructor_r = parser.
	Get(dataFromConstructorReduction).From(Constructor)

var data__data_expr_r = parser.
	Get(dataAppendExprReduction).From(Data, Expr)

var data__enclosed_r = parser.
	Get(grab_enclosed).From(LeftParen, Data, RightParen)
