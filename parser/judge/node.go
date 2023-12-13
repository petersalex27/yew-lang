// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package judge

import (
	"github.com/petersalex27/yew-lang/parser/expression"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

type JudgementNode struct {
	expression.ExpressionNode
	types.Monotyped[token.Token]
}

func (j *JudgementNode) Visit(cxt *inf.Context[token.Token]) {
	res := expression.Visit(j.ExpressionNode, cxt)

	status := cxt.Unify(res.GetType(), j.Monotyped)
	if status.NotOk() {
		// TODO

	}
}

func (j *JudgementNode) Equals(a ast.Ast) bool {
	j2, ok := a.(*JudgementNode)
	if !ok {
		return false
	}

	return j.Monotyped.Equals(j2.Monotyped) && expression.JudgeEquals(j, j2)
}

func (j *JudgementNode) NodeType() ast.Type { return Judgement }

func (j *JudgementNode) InOrderTraversal(f func(itoken.Token)) {
	j.ExpressionNode.InOrderTraversal(f)

	for _, token := range j.Monotyped.Collect() {
		f(token)
	}
}
