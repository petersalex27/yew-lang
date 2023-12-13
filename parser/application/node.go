// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package application

import (
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/parser/expression"
	"github.com/petersalex27/yew-lang/parser/node"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

type ApplicationNode struct {
	Apply []expression.ExpressionNode
	inf.Conclusion[token.Token, expr.Application[token.Token], types.Monotyped[token.Token]]
}

func (a *ApplicationNode) GetJudgement() types.TypedJudgement[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]] {
	e, ty := a.Conclusion.Judgement().GetExpressionAndType()
	mono := ty.(types.Monotyped[token.Token])
	return types.TypedJudge[token.Token](e, mono)
}

func (a *ApplicationNode) visitThenApplyExpressions(cxt *inf.Context[token.Token]) {
	// visit first node
	left := expression.Visit(a.Apply[0], cxt)

	for _, app := range a.Apply[1:] {
		// visit applicant
		right := expression.Visit(app, cxt)
		// apply "applyee" to applicant
		a.Conclusion = cxt.App(left, right)
		if a.Conclusion.NotOk() {
			// TODO: ? report or something
			return
		}

		// make result new applicant (or ignore on final iteration)
		e, t := a.Conclusion.Judgement().GetExpressionAndType()
		left = types.TypedJudge[token.Token](e, t.(types.Monotyped[token.Token]))
	}
}

func (a *ApplicationNode) Visit(cxt *inf.Context[token.Token]) {
	if len(a.Apply) < 2 {
		panic("illegal node structure")
	}

	a.visitThenApplyExpressions(cxt)
}

func (a *ApplicationNode) Equals(ast ast.Ast) bool {
	a2, ok := ast.(*ApplicationNode)
	if !ok {
		return false
	}

	return node.NodesEquals(a.Apply, a2.Apply)
}

func (*ApplicationNode) NodeType() ast.Type { return Application }

func (a *ApplicationNode) InOrderTraversal(action func(itoken.Token)) {
	for _, app := range a.Apply {
		app.InOrderTraversal(action)
	}
}
