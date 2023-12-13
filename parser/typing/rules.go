// =============================================================================
// Author-Date: Alex Peters - December 09, 2023
// =============================================================================
package typing

import (
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
)

// =============================================================================
// production functions
// =============================================================================

func produceRightAssoc(nodes ...ast.Ast) ast.Ast {
	const leftIndex, infixTokenIndex, rightIndex int = 0, 1, 2
	left := nodes[leftIndex].(*TypeNode).Type.(types.Monotyped[token.Token])
	right := nodes[rightIndex].(*TypeNode).Type.(types.Monotyped[token.Token])
	infix := nodes[infixTokenIndex].(ast.Token).Token.(token.Token)
	infixConst := types.MakeInfixConst(infix)
	res := types.Apply[token.Token](infixConst, left, right)
	return &TypeNode{Monotype, res}
}

func produceConst(nodes ...ast.Ast) ast.Ast {
	const tokenIndex int = 0
	constTok := nodes[tokenIndex].(ast.Token).Token.(token.Token)
	constant := types.MakeConst(constTok)
	return &TypeNode{Monotype, constant}
}

func produceVar(nodes ...ast.Ast) ast.Ast {
	const tokenIndex int = 0
	varTok := nodes[tokenIndex].(ast.Token).Token.(token.Token)
	variable := types.Var(varTok)
	return &TypeNode{Monotype, variable}
}

func produceApplication(nodes ...ast.Ast) ast.Ast {
	const leftIndex, rightIndex int = 0, 1
	left := nodes[leftIndex].(*TypeNode).Type.(types.Monotyped[token.Token])
	right := nodes[rightIndex].(*TypeNode).Type.(types.Monotyped[token.Token])
	res := types.Apply(left, right)
	return &TypeNode{Monotype, res}
}

// just returns second arg
func produceSecond(nodes ...ast.Ast) ast.Ast {
	return nodes[1]
}

// =============================================================================
// rules
// =============================================================================

var functionRule = parser.Get(produceRightAssoc).From(Monotype, Arrow, Monotype)

var tupleRule = parser.Get(produceRightAssoc).From(Monotype, Comma, Monotype)

var closeTupleRule = parser.Get(produceSecond).From(LeftParen, Monotype)

var constRule = parser.Get(produceConst).From(TypeId)

var varRule = parser.Get(produceVar).From(Id)

var applicationRule = parser.Get(produceApplication).From(Monotype, Monotype)

// =============================================================================
// production orders
// =============================================================================

var singletonProductions = parser.Order(constRule, varRule)

var applicationProductions = parser.Order(applicationRule)

var beforeCommaProductions = parser.Order(functionRule)

var loopedProductions = parser.Order(functionRule, tupleRule)

var closeTupleProductions = parser.Order(closeTupleRule)