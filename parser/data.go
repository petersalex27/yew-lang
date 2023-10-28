package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-lang/token"
)

/*
data          ::= patternC
                  | data expr
                  | '(' data ')'
*/

// Ast -> (Data==SomeExpression)
func astToData(a ast.Ast) SomeExpression { return a.(SomeExpression) }

func dataFromPatternCReduction(nodes ...ast.Ast) ast.Ast {
	const patternCIndex int = 0
	// patternC == constructor.(BinaryRecursiveNode).
	//		UpdateType(constructor.(BinaryRecursiveNode).NodeType(), PatternC)
	constructorExpr := constructorToExpression(nodes[patternCIndex])
	return SomeExpression{Data, constructorExpr}
}

func dataAppendExprReduction(nodes ...ast.Ast) ast.Ast {
	data := nodes[0].(SomeExpression)
	ex := getExpression(nodes[1]).Expression
	data.Expression = expr.Apply[token.Token](data.Expression, ex)
	return data
}

var data__patternC_r = parser.
	Get(dataFromPatternCReduction).From(Pattern)

var data__data_expr_r = parser.
	Get(dataAppendExprReduction).From(Data, Expr)

var data__enclosed_r = parser.
	Get(grab_enclosed).From(LeftParen, Data, RightParen)
