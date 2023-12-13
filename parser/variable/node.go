// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package variable

import (
	"github.com/petersalex27/yew-lang/parser/expression"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

// nodes that use the [Var] rule
type VarNode struct {
	Name token.Token
	inf.Conclusion[token.Token, expr.Const[token.Token], types.Monotyped[token.Token]]
}

func (v *VarNode) GetJudgement() types.TypedJudgement[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]] {
	e, ty := v.Conclusion.Judgement().GetExpressionAndType()
	mono := ty.(types.Monotyped[token.Token])
	return types.TypedJudge[token.Token](e, mono)
}

func (v *VarNode) Visit(cxt *inf.Context[token.Token]) {
	x := expr.Const[token.Token]{Name: v.Name}
	v.Conclusion = cxt.Var(x)
	v.Conclusion.Judgement()
}

func (v *VarNode) Equals(a ast.Ast) bool {
	v2, ok := a.(*VarNode)
	if !ok {
		return false
	}

	return utils.EquateTokens(v.Name, v2.Name) && expression.JudgeEquals(v, v2)
}

func (v *VarNode) NodeType() ast.Type { return Variable }

func (v *VarNode) InOrderTraversal(f func(itoken.Token)) {
	f(v.Name)
}
