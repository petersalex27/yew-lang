// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package rec

import (
	"github.com/petersalex27/yew-lang/parser/expression"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

type RecNode struct {
	Names       []token.Token
	Assignments []expression.ExpressionNode
	Expression  expression.ExpressionNode
	inf.Conclusion[token.Token, expr.RecIn[token.Token], types.Monotyped[token.Token]]
}

func (rec *RecNode) GetJudgement() types.TypedJudgement[token.Token, expr.Expression[token.Token], types.Monotyped[token.Token]] {
	e, ty := rec.Conclusion.Judgement().GetExpressionAndType()
	mono := ty.(types.Monotyped[token.Token])
	return types.TypedJudge[token.Token](e, mono)
}

func (rec *RecNode) Visit(cxt *inf.Context[token.Token]) {
	dischargePrime := cxt.Rec(rec.Names)
	assignments := fun.FMap(
		rec.Assignments,
		func(e expression.ExpressionNode) inf.TypeJudgement[token.Token] {
			return expression.Visit(e, cxt)
		},
	)
	dischargePrimePrime := dischargePrime(assignments)
	expression := expression.Visit(rec.Expression, cxt)
	rec.Conclusion = dischargePrimePrime(expression)
}

func (rec *RecNode) Equals(a ast.Ast) bool {
	rec2, ok := a.(*RecNode)
	if !ok {
		return false
	}

	if len(rec.Names) != len(rec2.Names) || len(rec.Assignments) != len(rec2.Assignments) {
		return false
	}

	for i, token := range rec.Names {
		if !utils.EquateTokens(token, rec2.Names[i]) {
			return false
		}
	}

	for i, assignment := range rec.Assignments {
		if !assignment.Equals(rec2.Assignments[i]) {
			return false
		}
	}

	return rec.Expression.Equals(rec2.Expression) && expression.JudgeEquals(rec, rec2)
}

func (*RecNode) NodeType() ast.Type { return RecExpr }

func (rec *RecNode) traverseNamesAndAssignments(f func(itoken.Token)) {
	if len(rec.Assignments) != len(rec.Names) {
		panic("len(Assignments) != len(Names)")
	}

	for i, token := range rec.Names {
		f(token)
		rec.Assignments[i].InOrderTraversal(f)
	}
}

func (rec *RecNode) traverseExpression(f func(itoken.Token)) {
	rec.Expression.InOrderTraversal(f)
}

func (rec *RecNode) InOrderTraversal(f func(itoken.Token)) {
	rec.traverseNamesAndAssignments(f)
	rec.traverseExpression(f)
}
