package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-lang/token"
)

/*
dependBinders ::= dependentHead                   # when l.a. is '.'
                  | '(' dependBinders ')'         #   //
dependent     ::= dependBinders '.' monotype
*/

// == dependent type reduction rules ==========================================

var dependent__dependBinders_Dot_monotype_r = parser.
	Get(dependTypeReduction).From(DependBinders, Dot, Monotype)

// == dependent type reductions ===============================================

func dependTypeReduction(nodes ...ast.Ast) ast.Ast {
	const bindersIndex, _, monoIndex int = 0, 1, 2
	binders := nodes[bindersIndex].(DependHeadNode).params
	mono := nodes[monoIndex].(TypeNode).Type.(types.Monotyped[token.Token])
	return TypeNode{
		Dependtype,
		types.MakeDependentType[token.Token](
			binders,
			types.Apply[token.Token](mono.(types.ReferableType[token.Token])),
		),
	}
}