// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package cons

import (
	typename "github.com/petersalex27/yew-lang/parser/type-name"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type TypeConstructorNode struct {
	Data bridge.Data[token.Token]
}

func MakeTypeConstructorNode(name token.Token, members []*typename.TypeNameNode) (constructor *TypeConstructorNode) {
	constructor = &TypeConstructorNode{}
	constructor.Data = bridge.MakeData(expr.MakeConst(name))
	return
}

// does nothing
func (constructor *TypeConstructorNode) Visit(cxt *inf.Context[token.Token]) {}

func (constructor *TypeConstructorNode) Equals(node ast.Ast) bool {
	constructor2, ok := node.(*TypeConstructorNode)
	if !ok {
		return false
	}
	return constructor.Data.StrictEquals(constructor2.Data)
}

func (*TypeConstructorNode) NodeType() ast.Type { return Constructor }

func (constructor *TypeConstructorNode) InOrderTraversal(action func(itoken.Token)) {
	for _, token := range constructor.Data.Collect() {
		action(token)
	}
}
