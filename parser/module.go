package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

/*
module        ::= moduleHead
                  | export ')'
*/

type ModuleNode struct {
	// ModuleDeclaration, ExportList, ExportHead, ModuleDefinition, or Source
	ast.Type
	name token.Token
	exportList []token.Token
	DefinitionsNode
}

// == module reduction rules ==================================================

var module__moduleDec_r = parser. 
	Get(rewriteModuleTypeReduction(ModuleDefinition)). 
	From(ModuleDeclaration)

var module__export_RightParen_r = parser. 
	Get(rewriteModuleTypeReduction(ModuleDefinition)). 
	From(ExportList, RightParen)

// == module reduction functions ==============================================

func rewriteModuleTypeReduction(ty ast.Type) func (nodes ...ast.Ast) ast.Ast {
	return func (nodes ...ast.Ast) ast.Ast {
		const moduleBasedIndex int = 0
		moduleBased := nodes[moduleBasedIndex].(ModuleNode)
		moduleBased.Type = ty
		return moduleBased
	}
}

// == ModuleNode implementation of ast.Ast ====================================

func (mod ModuleNode) Equals(a ast.Ast) bool {
	mod2, ok := a.(ModuleNode)
	if !ok {
		return false
	}

	if mod.Type != mod2.Type || !EqualsToken(mod.name, mod2.name) {
		return false
	}

	if len(mod.exportList) != len(mod2.exportList) {
		return false
	}

	for i, item := range mod.exportList {
		if !EqualsToken(item, mod2.exportList[i]) {
			return false
		}
	}

	return mod.DefinitionsNode.Equals(mod2.DefinitionsNode)
}

// returns mod.Type
func (mod ModuleNode) NodeType() ast.Type { return mod.Type }

func (mod ModuleNode) InOrderTraversal(f func(itoken.Token)) {
	f(mod.name)
	for _, token := range mod.exportList {
		f(token)
	}
	mod.DefinitionsNode.InOrderTraversal(f)
}