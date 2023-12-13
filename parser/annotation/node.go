// =============================================================================
// Author-Date: Alex Peters - November 27, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package annotation

import (
	"github.com/petersalex27/yew-lang/parser/utils"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type AnnotationNode struct {
	// TODO: finish
	token.Token
}

func (*AnnotationNode) Visit(cxt *inf.Context[token.Token]) {}

func (annot *AnnotationNode) Equals(node ast.Ast) bool {
	annot2, ok := node.(*AnnotationNode)
	if !ok {
		return false
	}
	return utils.EquateTokens(annot.Token, annot2.Token)
}

func (annot *AnnotationNode) NodeType() ast.Type { return Annotation }

func (annot *AnnotationNode) InOrderTraversal(action func(itoken.Token)) {
	action(annot.Token)
}