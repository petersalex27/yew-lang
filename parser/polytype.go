package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-lang/token"
)

/*
polytype ::= polyBinders '.' dependTyped
*/

// == polytype reduction rules ================================================

var polytype__polyBinders_Dot_dependTyped_r = parser.
	Get(polytypeReduction).From(PolyBinders, Dot, Dependtyped)

// == polytype reductions =====================================================

func polytypeReduction(nodes ...ast.Ast) ast.Ast {
	const bindersIndex, _, dependTypedIndex int = 0, 1, 2
	binders := nodes[bindersIndex].(PolyHeadNode).vars
	dependTyped := nodes[dependTypedIndex].(TypeNode).Type.(types.DependentTyped[token.Token])
	return TypeNode{
		Polytype,
		types.Forall[token.Token](binders...).Bind(dependTyped),
	}
}