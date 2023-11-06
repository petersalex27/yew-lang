package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
source        ::= module 'where' indent(n) definitions
                  | module
*/

// =============================================================================
// source production rules
// =============================================================================

var source__module_Where_exprBlock_definitions_r = parser. 
	Get(defineModuleSourceProduction). 
	From(ModuleDefinition, Where, IndentExprBlock, Definitions)

var source_module_r = parser. 
	Get(rewriteModuleTypeReduction(Source)). 
	From(ModuleDefinition)

// =============================================================================
// source production functions
// =============================================================================

func defineModuleSourceProduction(nodes ...ast.Ast) ast.Ast {
	const moduleIndex, _, _, defsIndex int = 0, 1, 2, 3
	module := nodes[moduleIndex].(ModuleNode)
	defs := nodes[defsIndex].(DefinitionsNode)
	module.DefinitionsNode = defs
	module.Type = Source 
	return module
}