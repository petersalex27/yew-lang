// =============================================================================
// Author-Date: Alex Peters - November 25, 2023
//
// Content:
// exported type rules
//
// Notes: -
// =============================================================================
package typesexport

import (
	"github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type TypeExportNode struct {
	TypeName         token.Token
	ConstructorNames []token.Token
}

func (typeExportNode *TypeExportNode) Visit(cxt *inf.Context[token.Token]) {}

func (typeExportNode *TypeExportNode) Equals(node ast.Ast) bool {
	typeExportNode2, ok := node.(*TypeExportNode)
	if !ok {
		return false
	}

	return utils.EquateTokens(typeExportNode.TypeName, typeExportNode2.TypeName) && 
		utils.TokensEquals(typeExportNode.ConstructorNames, typeExportNode2.ConstructorNames)
}

func (typeExportNode *TypeExportNode) NodeType() ast.Type { return types.TypeExport }

func (typeExportNode *TypeExportNode) InOrderTraversal(action func(itoken.Token)) {
	action(typeExportNode.TypeName)
	for _, token := range typeExportNode.ConstructorNames {
		action(token)
	}
}