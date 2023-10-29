package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

type VariableJudgement types.TypeJudgement[token.Token, expr.Variable[token.Token]]

// == variable judgement reduction rules ======================================

var varJudgement__Id_Colon_monotype_r = parser.
	Get(varJudgementReduction).From(Id, Colon, Monotype)

var varJudgement__enclosed_r = parser.
	Get(parenEnclosedReduction).
	From(LeftParen, VarJudgement, RightParen)

// == variable judgement reductions ===========================================

func varJudgementReduction(nodes ...ast.Ast) ast.Ast {
	const idIndex, _, monotypeIndex int = 0, 1, 2
	v := expr.Var(GetToken(nodes[idIndex]))
	someType := GetType(nodes[monotypeIndex]).Type
	judgment := types.Judgement[token.Token, expr.Variable[token.Token]](v, someType)
	return VariableJudgement(judgment)
}

// == variable judgement implementation of ast.Ast ============================

func (j VariableJudgement) Equals(a ast.Ast) bool {
	j2, ok := a.(VariableJudgement)
	if !ok {
		return false
	}
	ja := getVariableJudgement(j)
	exp, ty := ja.GetExpression(), ja.GetType()
	jb := getVariableJudgement(j2)
	exp2, ty2 := jb.GetExpression(), jb.GetType()
	return exp.Equals(globalContext__.exprCxt, expr.Expression[token.Token](exp2)) && ty.Equals(ty2)
}

func (j VariableJudgement) NodeType() ast.Type { return VarJudgement }

func (j VariableJudgement) InOrderTraversal(f func(itoken.Token)) {
	tmp := getVariableJudgement(j).Collect()
	for _, a := range tmp {
		f(a)
	}
}

// == variable judgement utils ================================================

func getVariableJudgement(node ast.Ast) types.TypeJudgement[token.Token, expr.Variable[token.Token]] {
	return types.TypeJudgement[token.Token, expr.Variable[token.Token]](node.(VariableJudgement))
}