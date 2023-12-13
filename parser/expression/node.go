package expression

import (
	"github.com/petersalex27/yew-lang/parser/node"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/types"
)

type ExpressionNode interface {
	node.Node
	GetJudgement() types.TypedJudgement[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]]
}

func Visit(e ExpressionNode, cxt *inf.Context[token.Token]) types.TypedJudgement[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]] {
	e.Visit(cxt)
	return e.GetJudgement()
}

func JudgeEquals(a, b ExpressionNode) bool {
	return types.JudgementEquals[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]](
		a.GetJudgement(), 
		b.GetJudgement(),
	)
}