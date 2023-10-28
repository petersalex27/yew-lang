package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-lang/token"
)

// monoList <- monotype Comma monoTail
var monoList__mono_Comma_monoTail_r = parser.
	Get(monoList__mono_Comma_monoTail).
	From(Monotype, Comma, MonoTail)
	
func monoList__mono_Comma_monoTail(nodes ...ast.Ast) ast.Ast {
	head := GetMonotype(nodes[0])
	comma := types.MakeInfixConst(GetToken(nodes[1]))
	tail := GetMonotype(nodes[2])
	ls := types.Apply[token.Token](comma, head, tail)
	return TypeNode{MonoList, ls}
}

// monoTail <- monotype Comma
var monoTail__mono_Comma_r = parser.
	Get(monoTail__mono_Comma).
	From(Monotype, Comma)

func monoTail__mono_Comma(nodes ...ast.Ast) ast.Ast {
	return TypeNode{MonoTail, GetMonotype(nodes[0])}
}

// monoTail <- monotype
var monoTail__mono_r = parser.
	Get(monoTail__mono).
	From(Monotype)

func monoTail__mono(nodes ...ast.Ast) ast.Ast {
	return TypeNode{MonoTail, GetMonotype(nodes[0])}
}


// monoList <- monotype Comma monoList
var monoList__mono_Comma_monoList_r = parser.
	Get(monoList__mono_Comma_monoList).
	From(Monotype, Comma, MonoList)

func monoList__mono_Comma_monoList(nodes ...ast.Ast) ast.Ast {
	head, tail := GetMonotype(nodes[0]), getApplicationType(nodes[2])
	comma := types.MakeInfixConst(GetToken(nodes[1]))
	return TypeNode{MonoList, infixType(head, comma, tail)}
}