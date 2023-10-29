package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

/*
judgement ::= expr ':' type
							| varJudgement
						  | '(' judgement ')'
*/

// Wraps types.TypeJudgement
type JudgementNode types.TypeJudgement[token.Token, expr.Expression[token.Token]]

// == judgement reduction rules ===============================================

var judgement__expr_Colon_type_r = parser.
	Get(judgementReduction).
	From(Expr, Colon, Type)

var judgement__enclosed_r = parser.
	Get(parenEnclosedReduction).
	From(LeftParen, TypeJudgement, RightParen)

var judgement__varJudgement_r = parser.
	Get(judgementFromVarJudgementReduction).
	From(VarJudgement)

// == judgement reductions ====================================================

func judgementReduction(nodes ...ast.Ast) ast.Ast {
	const exprIndex, _, typeIndex = 0, 1, 2
	expression := getExpression(nodes[exprIndex]).Expression
	someType := GetType(nodes[typeIndex]).Type
	return JudgementNode(types.Judgement[token.Token, expr.Expression[token.Token]](
		expression, someType, // expression: someType
	))
}

func judgementFromVarJudgementReduction(nodes ...ast.Ast) ast.Ast {
	const varJudgementIndex int = 0
	varJudgment := getVariableJudgement(nodes[varJudgementIndex])
	ty, variable := varJudgment.GetType(), varJudgment.GetExpression()
	expression := expr.Expression[token.Token](variable)
	judgment := types.Judgement(expression, ty)
	return JudgementNode(judgment)
}

// == judgment node implementation of ast.Ast =================================

func (j JudgementNode) Equals(a ast.Ast) bool {
	j2, ok := a.(JudgementNode)
	if !ok {
		return false
	}
	judge := types.TypeJudgement[token.Token, expr.Expression[token.Token]](j)
	judge2 := types.TypeJudgement[token.Token, expr.Expression[token.Token]](j2)
	e, t := judge.GetExpression(), judge.GetType()
	e2, t2 := judge2.GetExpression(), judge2.GetType()
	return e.StrictEquals(e2) && t.Equals(t2)
}

func (j JudgementNode) NodeType() ast.Type { return TypeJudgement }

func (j JudgementNode) InOrderTraversal(f func(itoken.Token)) {
	tmp := getJudgement(j).Collect()
	for _, a := range tmp {
		f(a)
	}
}

// == judgement utils =========================================================

func getJudgement(node ast.Ast) types.TypeJudgement[token.Token, expr.Expression[token.Token]] {
	return types.TypeJudgement[token.Token, expr.Expression[token.Token]](node.(JudgementNode))
}
