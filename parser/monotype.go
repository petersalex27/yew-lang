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
	Get(grab_enclosed).
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

// TODO: monotype <- monotype SemiColon expr

func infixType(m1 types.Monotyped[token.Token], c types.InfixConst[token.Token], m2 types.Monotyped[token.Token]) types.Application[token.Token] {
	return types.Application[token.Token](types.Apply[token.Token](c, m1, m2))
}

func GetType(a ast.Ast) TypeNode {
	return a.(TypeNode)
}

func getVariable(a ast.Ast) types.Variable[token.Token] {
	return a.(TypeNode).Type.(types.Variable[token.Token])
}

func GetMonotype(a ast.Ast) types.Monotyped[token.Token] {
	ty := a.(TypeNode)
	return ty.Type.(types.Monotyped[token.Token])
}

func getApplicationType(a ast.Ast) types.Application[token.Token] {
	return a.(TypeNode).Type.(types.Application[token.Token])
}

// monotype <- monotype Arrow monotype
var monotype__mono_Arrow_mono_r = parser.
	Get(monotype__mono_Arrow_mono).
	From(Monotype, Arrow, Monotype)

func monotype__mono_Arrow_mono(nodes ...ast.Ast) ast.Ast {
	m1, m2 := GetMonotype(nodes[0]), GetMonotype(nodes[2])
	arrow := types.MakeInfixConst(GetToken(nodes[1]))
	return TypeNode{Monotype, infixType(m1, arrow, m2)}
}

// typeApp <- monotype monotype
var typeApp__mono_mono_r = parser.
	Get(typeApp__mono_mono).
	From(Monotype, Monotype)

func typeApp__mono_mono(nodes ...ast.Ast) ast.Ast {
	// all monotypes are referable
	m1 := GetMonotype(nodes[0]).(types.ReferableType[token.Token])
	m2 := GetMonotype(nodes[1])
	return TypeNode{TypeApp, types.Apply[token.Token](m1, m2)}
}

// typeApp <- TypeId monotype
var typeApp__TypeId_mono_r = parser.
	Get(typeApp__TypeId_mono).
	From(TypeId, Monotype)

func typeApp__TypeId_mono(nodes ...ast.Ast) ast.Ast {
	c := types.MakeConst(GetToken(nodes[0]))
	m := GetMonotype(nodes[1])
	return TypeNode{TypeApp, types.Apply[token.Token](c, m)}
}

// typeApp <- typeApp monotype
var typeApp__typeApp_mono_r = parser.
	Get(typeApp__typeApp_mono).
	From(TypeApp, Monotype)

func typeApp__typeApp_mono(nodes ...ast.Ast) ast.Ast {
	app := types.Merge(getApplicationType(nodes[0]), GetMonotype(nodes[1]))
	return TypeNode{TypeApp, app}
}

// monotype <- typeApp
var monotype__typeApp_r = parser.
	Get(monotype__typeApp).
	From(TypeApp)

func monotype__typeApp(nodes ...ast.Ast) ast.Ast {
	return TypeNode{Monotype, GetType(nodes[0]).Type}
}

// monotype <- LeftParen monoTail RightParen
var monotype__LeftParen_monoTail_RightParen_r = parser.
	Get(monotype__LeftParen_monoTail_RightParen).
	From(LeftParen, MonoTail, RightParen)

func monotype__LeftParen_monoTail_RightParen(nodes ...ast.Ast) ast.Ast {
	m := GetMonotype(nodes[1])
	return TypeNode{Monotype, m}
}

// monotype <- LeftParen monoList RightParen
var monotype__LeftParen_monoList_RightParen_r = parser.
	Get(monotype__LeftParen_monoList_RightParen).
	From(LeftParen, MonoList, RightParen)

func monotype__LeftParen_monoList_RightParen(nodes ...ast.Ast) ast.Ast {
	return TypeNode{Monotype, GetMonotype(nodes[1])}
}

// monotype <- var
var monotype__var_r = parser.
	Get(monotype__var).
	From(FreeVar)

func monotype__var(nodes ...ast.Ast) ast.Ast {
	v := GetType(nodes[0])
	return TypeNode{Monotype, v.Type}
}

// monotype <- dependInstance
var monotype__dependInstance_r = parser.
	Get(monotype__dependInstance).
	From(DependInstance)

func monotype__dependInstance(nodes ...ast.Ast) ast.Ast {
	depInst := getDependInstance(nodes[0])
	return TypeNode{Monotype, depInst}
}

// var <- Id
var var__Id_r = parser.Get(var__Id).From(Id)

func var__Id(nodes ...ast.Ast) ast.Ast {
	v := GetToken(nodes[0])
	return TypeNode{FreeVar, types.Var(v)}
}
