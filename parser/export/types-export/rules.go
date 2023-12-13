package typesexport

import (
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

// initializes type export
func produceInitialTypeExport(nodes ...ast.Ast) ast.Ast {
	const typeNameIndex, _ int = 0, 1
	typeName := nodes[typeNameIndex].(ast.Token).Token.(token.Token)
	typeExportNode := new(TypeExportNode)
	typeExportNode.TypeName = typeName
	typeExportNode.ConstructorNames = []token.Token{}
	return typeExportNode
}

// adds constructor to type export
func produceTypeExportConstructor(nodes ...ast.Ast) ast.Ast {
	const typeExportIndex, consNameIndex int = 0, 1
	typeExport := nodes[typeExportIndex].(*TypeExportNode)
	consName := nodes[consNameIndex].(ast.Token).Token.(token.Token)
	typeExport.ConstructorNames = append(typeExport.ConstructorNames, consName)
	return typeExport
}

// returns first arg, drops all remaining args
func produceTypeExportDrop(nodes ...ast.Ast) ast.Ast {
	const typeExportIndex, _ int = 0, 1
	return nodes[typeExportIndex]
}

var typeExportInitRule = parser.Get(produceInitialTypeExport).From(TypeId, LeftParen)
var typeExportConstructorRule = parser.Get(produceTypeExportConstructor).From(TypeExport, TypeId)
var typeExportCommaRule = parser.Get(produceTypeExportDrop).From(TypeExport, Comma)
var typeExportFinishRule = parser.Get(produceTypeExportDrop).From(TypeExport, RightParen)


// productions for type export list initialization
var typeExportInitProductions = parser.Order(typeExportInitRule)

// productions for ending parsing of type export list
var typeExportEndProductions = parser.Order(typeExportFinishRule)

// productions for adding a constructor to type export list
var typeExportProductions = parser.Order(typeExportConstructorRule)

// productions for handling comma sep. type export members
var typeExportCommaProductions = parser.Order(typeExportCommaRule)