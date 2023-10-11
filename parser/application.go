package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

// ApplicationNode = SomeExpression{Application, _}

func applicationReduction(nodes ...ast.Ast) ast.Ast {
	const leftIndex, rightIndex int = 0, 1
	left := getExpression(nodes[leftIndex])
	right := getExpression(nodes[rightIndex])
	return SomeExpression{
		Application,
		expr.Apply(left.Expression, right.Expression),
	}
}

// application <- expr expr
var application__expr_expr_r = parser. 
	Get(applicationReduction).From(Expr, Expr)