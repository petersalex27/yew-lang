// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package module

import (
	//tokennode "github.com/petersalex27/yew-lang/parser/token-node"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

// =============================================================================
// helper functions
// =============================================================================

// creates initial module node w/ given type
func produceInitialHelper(ty ast.Type, exportAll bool, nodes ...ast.Ast) *ModuleSourceNode {
	const _, nameIndex int = 0, 1
	nameNode := nodes[nameIndex]
	name := nameNode.(ast.Token).Token.(token.Token)
	mod := new(ModuleSourceNode)
	mod.Name = name
	mod.Type = ty
	mod.exportAll = exportAll
	return mod
}

// =============================================================================
// production functions
// =============================================================================

// makes export node
func produceInitialExport(nodes ...ast.Ast) ast.Ast {
	return produceInitialHelper(ExportList, false, nodes...)
}

// makes module source node
func produceModuleHead(nodes ...ast.Ast) ast.Ast {
	return produceInitialHelper(ModuleDef, false, nodes...)
}

func produceModuleAndExportList(nodes ...ast.Ast) ast.Ast {
	const exportListIndex, _ int = 0, 1
	exportList := nodes[exportListIndex].(*ModuleSourceNode)
	exportList.Type = ModuleDef
	return exportList
}

func produceExportEverythingModule(nodes ...ast.Ast) ast.Ast {
	return produceInitialHelper(ModuleDef, true, nodes...)
}

// func produceInitializedExport(nodes ...ast.Ast) ast.Ast {
// 	mod := produceModuleHead(nodes...).(*ModuleSourceNode)
// 	mod.Type = ExportList
// 	return mod
// }

// // export all defs in file
// func produceExportAll(nodes ...ast.Ast) ast.Ast {
// 	const moduleIndex, _, _ int = 0, 1, 2
// 	mod := nodes[moduleIndex].(*ModuleSourceNode)
// 	mod.Type = ModuleDef
// 	mod.exportAll = true
// 	return mod
// }

// =============================================================================
// rules
// =============================================================================

// module ::= 'module' ID '(' ')'
var moduleExportNothingRule = parser.Get(produceModuleHead).From(Module, Id, LeftParen, RightParen)

// module ::= 'module' ID '(' .. ')'
var moduleExportEverythingRule = parser.Get(produceExportEverythingModule).From(Module, Id, LeftParen, DotDot, RightParen)

// export ::= 'module' ID '('
var beginExportRule = parser.Get(produceInitialExport).From(Module, Id, LeftParen)

// module ::= 'module' ID
var emptyModuleExportRule = parser.Get(produceModuleHead).From(Module, Id)

// module ::= export ')'
var attachExportRule = parser.Get(produceModuleAndExportList).From(ExportList, RightParen)

// =============================================================================
// production orders
// =============================================================================

var parseModuleProductions = parser.Order(beginExportRule, emptyModuleExportRule, moduleExportEverythingRule, moduleExportNothingRule)

var attachExportListProductions = parser.Order(attachExportRule)
