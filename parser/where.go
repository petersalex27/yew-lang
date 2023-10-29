package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

// whereExpr ::= expr 'where' function

var whereExpr__expr_Where_function_r = parser.
	Get(whereReduction).From(Expr, Where, Function)

func whereReduction(nodes ...ast.Ast) ast.Ast {
	const exprIndex, whereIndex, funcIndex int = 0, 1, 2
	let, in := nodes[whereIndex], nodes[whereIndex]
	res := letReduction(let, nodes[funcIndex], in, nodes[exprIndex]).
		(SomeExpression).
		Expression.
		(expr.NameContext[token.Token])
	whereExpr := expr.Where[token.Token](res.GetContextualized(), res.GetName(), res.GetAssignment())
	return SomeExpression{WhereExpr, whereExpr,}
}