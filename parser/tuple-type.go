package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
)

var tupleType__Comma_monotype_RightParen_r = parser. 
	Get(initialTupleTypeReduction).
	From(Comma, Monotype, RightParen)

var tupleType__Comma_monotype_tupleType_r = parser.
	Get(rightRecursiveTupleTypeReduction). 
	From(Comma, Monotype, TupleType)

func initialTupleTypeReduction(nodes ...ast.Ast) ast.Ast {
	const commaIndex, monoIndex, _ int = 0, 1, 2
	comma := types.MakeInfixConst(GetToken(nodes[commaIndex]))
	monotype := nodeAsMonotype(nodes[monoIndex])
	return NodeSequence{
		TupleType, 
		[]ast.Ast{
			TypeNode{Monotype, comma}, 
			TypeNode{Monotype, monotype},
		},
	}
}

func rightRecursiveTupleTypeReduction(nodes ...ast.Ast) ast.Ast {
	const commaIndex, monoIndex, savedIndex int = 0, 1, 2
	commaNew := types.MakeInfixConst(GetToken(nodes[commaIndex]))
	left := nodeAsMonotype(nodes[monoIndex])
	saved := nodes[savedIndex].(NodeSequence)
	comma := nodeAsMonotype(saved.nodes[0]).(types.InfixConst[token.Token])
	right := nodeAsMonotype(saved.nodes[1])
	return NodeSequence{
		TupleType, 
		[]ast.Ast{
			TypeNode{Monotype, commaNew}, 
			TypeNode{Monotype, types.Apply[token.Token](comma, left, right)},
		},
	}
}