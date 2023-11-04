package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
moduleDec    ::= INDENT 'module' ID
*/

var moduleDec__Indent_Module_Id_r = parser. 
	Get(moduleDecReduction). 
	From(Indent, Module, Id)

func moduleDecReduction(nodes ...ast.Ast) ast.Ast {
	const _, _, idIndex int = 0, 1, 2
	id := GetToken(nodes[idIndex])
	return ModuleNode{
		Type: ModuleDeclaration,
		name: id,
		exportList: []exportToken{},
		DefinitionsNode: DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}
}