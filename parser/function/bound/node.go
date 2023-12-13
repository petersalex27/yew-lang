// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package bound

import (
	"github.com/petersalex27/yew-lang/parser/expression"
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

type BoundFunctionNode struct {
	Name       token.Token
	Assignment expression.ExpressionNode
	types.TypedJudgement[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]]
}

func (f *BoundFunctionNode) Visit(cxt *inf.Context[token.Token]) {
	_ = expression.Visit(f.Assignment, cxt)
}

func (f *BoundFunctionNode) Equals(node ast.Ast) bool {
	f2, ok := node.(*BoundFunctionNode)
	if !ok {
		return false
	}

	return utils.EquateTokens(f.Name, f2.Name) && f.Assignment.Equals(f2.Assignment)
}

func (f *BoundFunctionNode) NodeType() ast.Type {
	panic("TODO: finish")
}

func (f *BoundFunctionNode) InOrderTraversal(action func(itoken.Token)) {
	panic("TODO: finish")
}
