// =============================================================================
// Author-Date: Alex Peters - December 13, 2023
// =============================================================================
package pattern

import (
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type PatternNode struct {
	head []token.Token
	tail *PatternNode
}

func (patternNode *PatternNode) Visit(cxt *inf.Context[token.Token]) {}

func (patternNode *PatternNode) Equals(node ast.Ast) bool {
	patternNode2, ok := node.(*PatternNode)
	if !ok {
		return false
	}
	if patternNode == nil {
		return patternNode2 == nil
	}

	return patternNode.tail.Equals(patternNode2.tail) && 
		utils.TokensEquals(patternNode.head, patternNode2.head)
}

func (patternNode *PatternNode) NodeType() ast.Type {
	return Pattern
}

func (patternNode *PatternNode) InOrderTraversal(action func(itoken.Token)) {
	if patternNode == nil {
		return
	}
	
	for _, token := range patternNode.head {
		action(token)
	}
	patternNode.tail.InOrderTraversal(action)
}