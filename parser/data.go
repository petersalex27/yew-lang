package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
data          ::= patternC
                  | data expr
                  | '(' data ')'
*/

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
	Get(dataFromPatternCReduction).From(PatternC)

var data__data_expr_r = parser.
	Get(dataAppendExprReduction).From(Data, Expr)

var data__enclosed_r = parser.
	Get(parenEnclosedReduction).From(LeftParen, Data, RightParen)
