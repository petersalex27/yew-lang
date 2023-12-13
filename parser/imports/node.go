// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package imports

import (
	"github.com/petersalex27/yew-lang/parser/utils"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type ImportElementNode struct {
	inf.QualificationType
	Name token.Token
	As   token.Token
	From string
	// TODO: allow only specific exported items from Name to be imported
	// Spec Stuff
}

type ImportNode []*ImportElementNode

func (importNode *ImportNode) Visit(cxt *inf.Context[token.Token]) {
	for _, imported := range *importNode {
		imported.Visit(cxt)
	}
}

func (elem *ImportElementNode) Visit(cxt *inf.Context[token.Token]) {
	// TODO
	cxt.Import(elem.QualificationType, elem.Name, elem.As)
}

func (elem *ImportElementNode) String() string {
	out := elem.Name.GetValue()
	if elem.QualificationType == inf.FullyQualified {
		out = "qualified " + out
	}

	if len(elem.From) > 0 {
		out = out + " from " + elem.From
	}
	return out
}

func (importNode *ImportNode) Equals(node ast.Ast) bool {
	importNode2, ok := node.(*ImportNode)
	if !ok {
		return false
	}

	if len(*importNode) != len(*importNode2) {
		return false
	}

	for i, imported := range *importNode {
		if imported.QualificationType != (*importNode2)[i].QualificationType {
			return false
		} else if !utils.EquateTokens(imported.Name, (*importNode2)[i].Name) {
			return false
		} else if !utils.EquateTokens(imported.As, (*importNode2)[i].As) {
			return false
		} else if imported.From != (*importNode2)[i].From {
			return false
		}
		// TODO:
	}

	return true
}

func (elem *ImportElementNode) Equals(node ast.Ast) bool {
	elem2, ok := node.(*ImportElementNode)
	if !ok {
		return false
	}

	if elem.QualificationType != elem2.QualificationType {
		return false
	} else if !utils.EquateTokens(elem.Name, elem2.Name) {
		return false
	} else if !utils.EquateTokens(elem.As, elem2.As) {
		return false
	} else if elem.From != elem2.From {
		return false
	}
	// TODO:

	return true
}

func (*ImportNode) NodeType() ast.Type { return ImportContext }

func (*ImportElementNode) NodeType() ast.Type { return ImportElement }

func (importNode *ImportNode) InOrderTraversal(action func(itoken.Token)) {
	for _, imported := range *importNode {
		imported.InOrderTraversal(action)
	}
}

func (elem *ImportElementNode) InOrderTraversal(action func(itoken.Token)) {
	action(elem.Name)
	action(elem.As)
	// TODO
}