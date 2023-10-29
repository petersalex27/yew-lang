package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
)

/*
monotype      ::= monotype monotype
									| monotype '->' monotype
                  | '(' monotype ')'
                  | TYPE_ID
                  | '(' monotype tupleType
*/

func nodeAsMonotype(node ast.Ast) types.Monotyped[token.Token] {
	return node.(TypeNode).Type.(types.Monotyped[token.Token])
}

// == monotype reduction rules ================================================

var monotype__monotype_monotype_r = parser.
	Get(monotypeApplyReduction).
	From(Monotype, Monotype)

var monotype__monotype_Arrow_monotype_r = parser.
	Get(arrowTypeReduction).
	From(Monotype, Arrow, Monotype)

var monotype__enclosed_r = parser.
	Get(parenEnclosedReduction).
	From(LeftParen, Monotype, RightParen)

var monotype__TypeId_r = parser.
	Get(monotypeConstReduction).
	From(TypeId)

var monotype__Id_r = parser.
	Get(monotypeVarReduction).
	From(Id)

var monotype__LeftParen_monotype_tupleType_r = parser.
	Get(monotypeTupleReduction).
	From(LeftParen, Monotype, TupleType)

// == monotype reduction functions ============================================

func monotypeApplyReduction(nodes ...ast.Ast) ast.Ast {
	const leftIndex, rightIndex int = 0, 1
	left := nodeAsMonotype(nodes[leftIndex]).(types.ReferableType[token.Token])
	right := nodeAsMonotype(nodes[rightIndex])
	return TypeNode{
		Monotype,
		types.Apply[token.Token](left, right),
	}
}

func arrowTypeReduction(nodes ...ast.Ast) ast.Ast {
	const leftIndex, arrowIndex, rightIndex int = 0, 1, 2
	left, right := nodeAsMonotype(nodes[leftIndex]), nodeAsMonotype(nodes[rightIndex])
	arrow := types.MakeInfixConst(GetToken(nodes[arrowIndex]))
	return TypeNode{
		Monotype,
		types.Apply[token.Token](arrow, left, right),
	}
}

func monotypeConstReduction(nodes ...ast.Ast) ast.Ast {
	const typeIdIndex int = 0
	typeIdToken := GetToken(nodes[typeIdIndex])
	return TypeNode{Monotype, types.MakeConst(typeIdToken)}
}

func monotypeVarReduction(nodes ...ast.Ast) ast.Ast {
	const typeIdIndex int = 0
	typeIdToken := GetToken(nodes[typeIdIndex])
	return TypeNode{Monotype, types.Var(typeIdToken)}
}

func monotypeTupleReduction(nodes ...ast.Ast) ast.Ast {
	const _, monoIndex, tupleIndex int = 0, 1, 2
	head := nodeAsMonotype(nodes[monoIndex])
	tupleNodes := nodes[tupleIndex].(NodeSequence).nodes
	comma := tupleNodes[0].(TypeNode).Type.(types.InfixConst[token.Token])
	tail := tupleNodes[1].(TypeNode).Type.(types.Monotyped[token.Token])
	return TypeNode{
		Monotype,
		types.Apply[token.Token](comma, head, tail),
	}
}

func GetType(a ast.Ast) TypeNode {
	return a.(TypeNode)
}

func GetMonotype(a ast.Ast) types.Monotyped[token.Token] {
	ty := a.(TypeNode)
	return ty.Type.(types.Monotyped[token.Token])
}

// func getApplicationType(a ast.Ast) types.Application[token.Token] {
// 	return a.(TypeNode).Type.(types.Application[token.Token])
// }
