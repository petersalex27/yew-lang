// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package let

import (
	"github.com/petersalex27/yew-lang/parser/expression"
	"github.com/petersalex27/yew-lang/parser/function/bound"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

type LetNode struct {
	*bound.BoundFunctionNode
	Expression expression.ExpressionNode
	inf.Conclusion[token.Token, expr.NameContext[token.Token], types.Monotyped[token.Token]]
}

func (let *LetNode) GetJudgement() types.TypedJudgement[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]] {
	e, ty := let.Conclusion.Judgement().GetExpressionAndType()
	mono := ty.(types.Monotyped[token.Token])
	return types.TypedJudge[token.Token](e, mono)
}

func (let *LetNode) Visit(cxt *inf.Context[token.Token]) {
	let.BoundFunctionNode.Visit(cxt)
	assignment := let.BoundFunctionNode.TypedJudgement
	// bind assignment to name, name = assignment
	discharge := cxt.Let(let.Name, assignment)
	// with assignment bound to name, get judgement for expression
	expression := expression.Visit(let.Expression, cxt)
	// unbind name from assignment, concluding judgement from expression
	let.Conclusion = discharge(expression)
}

func (let *LetNode) Equals(a ast.Ast) bool {
	let2, ok := a.(*LetNode)
	if !ok {
		return false
	}

	return utils.EquateTokens(let.Name, let2.Name) &&
		let.Assignment.Equals(let2.Assignment) &&
		let.Expression.Equals(let2.Expression) &&
		expression.JudgeEquals(let, let2)
}

func (let *LetNode) NodeType() ast.Type { return LetExpr }

func (let *LetNode) InOrderTraversal(f func(itoken.Token)) {
	f(let.Name)
	let.Assignment.InOrderTraversal(f)
	let.Expression.InOrderTraversal(f)
}
