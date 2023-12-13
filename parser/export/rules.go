// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content:
// rules for export list parsing
//
// Notes: -
// =============================================================================
package export

import (
	typesexport "github.com/petersalex27/yew-lang/parser/export/types-export"
	"github.com/petersalex27/yew-lang/parser/module"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

// exports a non-abstract type
func produceExportTypeExport(nodes ...ast.Ast) ast.Ast {
	const moduleIndex, typeExportIndex int = 0, 1
	mod := nodes[moduleIndex].(*module.ModuleSourceNode)
	ty := nodes[typeExportIndex].(*typesexport.TypeExportNode)
	mod.TypeNames = append(mod.TypeNames, ty.TypeName)
	mod.ConstructorNames = append(mod.ConstructorNames, ty.ConstructorNames)
	return mod
}

// exports an abstract type
func produceExportType(nodes ...ast.Ast) ast.Ast {
	const moduleIndex, typeNameIndex int = 0, 1
	modNode := nodes[moduleIndex]
	typeName := nodes[typeNameIndex].(ast.Token).Token.(token.Token)
	typeExport := typesexport.TypeExportNode{
		TypeName:         typeName,
		ConstructorNames: []token.Token{}, // empty slice--exports no constructors
	}
	// produce export
	return produceExportTypeExport(modNode, &typeExport)
}

// exports type and all its constructors
func produceImplicitFullType(nodes ...ast.Ast) ast.Ast {
	const moduleIndex, typeNameIndex, _, dotDotIndex, _ = 0, 1, 2, 3, 4
	modNode := nodes[moduleIndex]
	typeName := nodes[typeNameIndex].(ast.Token).Token.(token.Token)
	dotDotToken := nodes[dotDotIndex].(ast.Token).Token.(token.Token)
	dotDot := token.TypeId.Make().
		AddValue("..").
		SetLength(dotDotToken.GetLength()).
		SetLineChar(dotDotToken.GetLineChar()).(token.Token)
	typeExport := typesexport.TypeExportNode{
		TypeName:         typeName,
		ConstructorNames: []token.Token{dotDot},
	}
	// produce export
	return produceExportTypeExport(modNode, &typeExport)
}

// add name to export list
func produceExportName(nodes ...ast.Ast) ast.Ast {
	const moduleIndex, nameIndex int = 0, 1
	mod := nodes[moduleIndex].(*module.ModuleSourceNode)
	name := nodes[nameIndex].(ast.Token).Token.(token.Token)
	mod.FunctionNames = append(mod.FunctionNames, name)
	return mod
}

// export name rules
//
//	export ::= export ID
//	           | export INFIXED
var exportIdRule = parser.Get(produceExportName).From(ExportList, Id)
var exportInfixedRule = parser.Get(produceExportName).From(ExportList, Infixed)

// export ::= export typeExport
var exportTypeExportRule = parser.Get(produceExportTypeExport).From(ExportList, TypeExport)
var exportTypeRule = parser.Get(produceExportType).From(ExportList, TypeId)
var exportExplicitAbstractTypeRule = parser.Get(produceExportType).From(ExportList, TypeId, LeftParen, RightParen)
var exportImplicitFullTypeRule = parser.Get(produceImplicitFullType).From(ExportList, TypeId, LeftParen, DotDot, RightParen)

// export ::= export ','
var exportDropCommaRule = parser.Get(produceExportDrop).From(ExportList, Comma)

// productions for export members
var exportProductions = parser.Order(
	exportIdRule,
	exportInfixedRule,
	exportTypeExportRule,
	exportTypeRule,
	exportExplicitAbstractTypeRule,
	exportImplicitFullTypeRule,
)

// productions for handling comma sep. export members
var exportCommaProductions = parser.Order(
	exportDropCommaRule,
)
