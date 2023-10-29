package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
source        ::= module 'where' definitions
                  | module
*/

// == source reduction rules ==================================================

var source__module_Where_definitions_r = parser. 
	Get(defineModuleSourceReduction). 
	From(ModuleDefinition, Where, Definitions)

var source_module_r = parser. 
	Get(rewriteModuleTypeReduction(Source)). 
	From(ModuleDefinition)

// == source reduction functions ==============================================

func defineModuleSourceReduction(nodes ...ast.Ast) ast.Ast {
	const moduleIndex, _, defsIndex int = 0, 1, 2
	module := nodes[moduleIndex].(ModuleNode)
	defs := nodes[defsIndex].(DefinitionsNode)
	module.DefinitionsNode = defs
	module.Type = Source 
	return module
}