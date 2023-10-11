package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
val           ::= literal
                  | array
*/

func valFromExpressionTyped(nodes ...ast.Ast) ast.Ast {
	e := nodes[0].(expressionNodeTypes).getExpression().Expression
	return SomeExpression{Val, e}
}

// val <- literal
var val__literal_r = parser.Get(valFromExpressionTyped).From(Literal)

// val <- array
var val__array_r = parser.Get(valFromExpressionTyped).From(Array)