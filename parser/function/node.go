// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package function

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

type FunctionNode struct {
	Param      token.Token
	Expression expression.ExpressionNode
	inf.Conclusion[token.Token, expr.Function[token.Token], types.Monotyped[token.Token]]
}

func (f *FunctionNode) GetJudgement() types.TypedJudgement[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]] {
	e, ty := f.Conclusion.Judgement().GetExpressionAndType()
	mono := ty.(types.Monotyped[token.Token])
	return types.TypedJudge[token.Token](e, mono)
}

func (f *FunctionNode) Visit(cxt *inf.Context[token.Token]) {
	discharge := cxt.Abs(f.Param)
	f.Conclusion = discharge(expression.Visit(f.Expression, cxt))
}

func (f *FunctionNode) Equals(a ast.Ast) bool {
	f2, ok := a.(*FunctionNode)
	if !ok {
		return false
	}

	return utils.EquateTokens(f.Param, f2.Param) &&
		f.Expression.Equals(f2.Expression) &&
		expression.JudgeEquals(f, f2)
}

func (f *FunctionNode) NodeType() ast.Type { return Function }

func (f *FunctionNode) InOrderTraversal(g func(itoken.Token)) {
	g(f.Param)
	f.Expression.InOrderTraversal(g)
}
