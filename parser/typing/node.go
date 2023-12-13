// =============================================================================
// Author-Date: Alex Peters - December 09, 2023
// =============================================================================
package typing

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/types"
)

type TypeNode struct {
	TypeKind ast.Type
	types.Type[token.Token]
}

func (ty *TypeNode) Equals(node ast.Ast) bool {
	ty2, ok := node.(*TypeNode)
	if !ok {
		return false
	}

	if ty.TypeKind != ty2.TypeKind {
		return false
	}

	return ty.Equals(ty2)
}

func (ty *TypeNode) NodeType() ast.Type { return ty.TypeKind }

func (ty *TypeNode) InOrderTraversal(action func(itoken.Token)) {
	for _, token := range ty.Collect() {
		action(token)
	}
}
