package parser

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

// judgement ::= expr ':' type

// Wraps types.TypeJudgement
type JudgementNode types.TypeJudgement[token.Token, expr.Expression[token.Token]]

// judgement <- expr Colon type
var judgement__expr_Colon_type_r = parser.
	Get(judgementReduction).
	From(Expr, Colon, Type)

// judgement <- expr Colon type
func judgementReduction(nodes ...ast.Ast) ast.Ast {
	// ignore Colon, i.e., element nodes[1]
	const exprIndex, _, typeIndex = 0, 1, 2
	return JudgementNode(types.Judgement[token.Token, expr.Expression[token.Token]](
		getExpression(nodes[exprIndex]).Expression,
		GetType(nodes[typeIndex]).Type,
	))
}

// judgement <- '(' judgement ')'
var judgement__enclosed_r = parser.Get(grab_enclosed).From(LeftParen, TypeJudgement, RightParen)

func judgementToExpression(nodes ...ast.Ast) ast.Ast {
	return ExpressionNode{
		bridge.JudgementAsExpression[token.Token, expr.Expression[token.Token]](
			nodes[0].(JudgementNode),
		),
	}
}

type VariableJudgement types.TypeJudgement[token.Token, expr.Variable[token.Token]]

type ApplicationJudgement types.TypeJudgement[token.Token, expr.Application[token.Token]]

type AnonFuncJudgement types.TypeJudgement[token.Token, expr.Function[token.Token]]

type SomeJudgementNode struct {
	ty ast.Type
	JudgementNode
}

func getJudgement(node ast.Ast) types.TypeJudgement[token.Token, expr.Expression[token.Token]] {
	return types.TypeJudgement[token.Token, expr.Expression[token.Token]](node.(JudgementNode))
}

func getFreeJudgement(node ast.Ast) types.FreeJudgement[token.Token, expr.Expression[token.Token]] {
	return types.FreeJudgement[token.Token, expr.Expression[token.Token]](getJudgement(node))
}

func getExprJudgement_ty(node ast.Ast) types.ExpressionJudgement[token.Token, expr.Expression[token.Token]] {
	return types.ExpressionJudgement[token.Token, expr.Expression[token.Token]](getJudgement(node))
}

func getExprJudgement_free(node ast.Ast) types.ExpressionJudgement[token.Token, expr.Expression[token.Token]] {
	return types.ExpressionJudgement[token.Token, expr.Expression[token.Token]](getFreeJudgement(node))
}

func getAppJudgement(nodes ast.Ast) types.TypeJudgement[token.Token, expr.Application[token.Token]] {
	return types.TypeJudgement[token.Token, expr.Application[token.Token]](nodes.(ApplicationJudgement))
}

func getAnonJudgement(nodes ast.Ast) types.TypeJudgement[token.Token, expr.Function[token.Token]] {
	return types.TypeJudgement[token.Token, expr.Function[token.Token]](nodes.(AnonFuncJudgement))
}

func getSomeJudgement(nodes ast.Ast) SomeJudgementNode {
	return nodes.(SomeJudgementNode)
}

func getVariableJudgement(node ast.Ast) types.TypeJudgement[token.Token, expr.Variable[token.Token]] {
	return types.TypeJudgement[token.Token, expr.Variable[token.Token]](node.(VariableJudgement))
}

func mkFreeJudge(ex expr.Expression[token.Token], ty types.Type[token.Token]) types.ExpressionJudgement[token.Token, expr.Expression[token.Token]] {
	return (types.FreeJudgement[token.Token, expr.Expression[token.Token]]{}).MakeJudgement(ex, ty)
}

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

func (j VariableJudgement) Equals(a ast.Ast) bool {
	j2, ok := a.(VariableJudgement)
	if !ok {
		return false
	}
	ja := getVariableJudgement(j)
	exp, ty := ja.GetExpression(), ja.GetType()
	jb := getVariableJudgement(j2)
	exp2, ty2 := jb.GetExpression(), jb.GetType()
	return exp.Equals(glb_cxt.exprCxt, expr.Expression[token.Token](exp2)) && ty.Equals(ty2)
}

func (j VariableJudgement) NodeType() ast.Type { return VarJudgement }

func (j VariableJudgement) InOrderTraversal(f func(itoken.Token)) {
	tmp := getVariableJudgement(j).Collect()
	for _, a := range tmp {
		f(a)
	}
}

func (j ApplicationJudgement) Equals(a ast.Ast) bool {
	j2, ok := a.(ApplicationJudgement)
	if !ok {
		return false
	}
	ja := getAppJudgement(j)
	exp, ty := ja.GetExpression(), ja.GetType()
	jb := getAppJudgement(j2)
	exp2, ty2 := jb.GetExpression(), jb.GetType()
	return exp.Equals(glb_cxt.exprCxt, expr.Expression[token.Token](exp2)) && ty.Equals(ty2)
}

func (j ApplicationJudgement) NodeType() ast.Type { return AppJudgement }

func (j ApplicationJudgement) InOrderTraversal(f func(itoken.Token)) {
	tmp := getAppJudgement(j).Collect()
	for _, a := range tmp {
		f(a)
	}
}

func (j AnonFuncJudgement) Equals(a ast.Ast) bool {
	j2, ok := a.(AnonFuncJudgement)
	if !ok {
		return false
	}
	ja := getAnonJudgement(j)
	exp, ty := ja.GetExpression(), ja.GetType()
	jb := getAnonJudgement(j2)
	exp2, ty2 := jb.GetExpression(), jb.GetType()
	return exp.Equals(glb_cxt.exprCxt, expr.Expression[token.Token](exp2)) && ty.Equals(ty2)
}

func (j AnonFuncJudgement) NodeType() ast.Type { return AnonJudgement }

func (j AnonFuncJudgement) InOrderTraversal(f func(itoken.Token)) {
	tmp := getAnonJudgement(j).Collect()
	for _, a := range tmp {
		f(a)
	}
}

func (j SomeJudgementNode) Equals(a ast.Ast) bool {
	j2, ok := a.(SomeJudgementNode)
	if !ok {
		return false
	}
	return j.ty == j2.ty && j.JudgementNode.Equals(j2.JudgementNode)
}

func (j SomeJudgementNode) NodeType() ast.Type { return j.ty }

func (j SomeJudgementNode) InOrderTraversal(f func(itoken.Token)) {
	tmp := getJudgement(j.JudgementNode).Collect()
	for _, a := range tmp {
		f(a)
	}
}

func makeJudgement(ex ExpressionNode, type_ TypeNode) JudgementNode {
	return JudgementNode(types.Judgement(ex.Expression, type_.Type))
}

func makeSomeJudgement[E expr.Expression[token.Token]](ty ast.Type, e E, type_ TypeNode) SomeJudgementNode {
	return SomeJudgementNode{
		ty,
		JudgementNode(types.Judgement[token.Token, expr.Expression[token.Token]](e, type_.Type)),
	}
}

// judgement <- varJudgement
var judgement__varJudgement_r = parser.
	Get(judgement__varJudgement).
	From(VarJudgement)

func judgement__varJudgement(nodes ...ast.Ast) ast.Ast {
	vj := getVariableJudgement(nodes[0])
	ty, exp := vj.GetType(), vj.GetExpression()
	j := types.Judgement(expr.Expression[token.Token](exp), ty)
	return JudgementNode(j)
}

// judgement <- appJudgement
var judgement__appJudgement_r = parser.
	Get(judgement__appJudgement).
	From(AppJudgement)

func judgement__appJudgement(nodes ...ast.Ast) ast.Ast {
	return nodes[0].(SomeJudgementNode).JudgementNode
}

// varJudgement <- Id Colon monotype
var varJudgement__Id_Colon_monotype_r = parser.
	Get(varJudgement__Id_Colon_monotype).
	From(Id, Colon, Monotype)

func varJudgement__Id_Colon_monotype(nodes ...ast.Ast) ast.Ast {
	// var
	v := expr.Var(GetToken(nodes[0]))
	// ignore Colon, i.e., element nodes[1]
	type_ := GetType(nodes[2])
	return makeSomeJudgement(VarJudgement, v, type_)
}

// varJudgement <- LeftParen varJudgement RightParen
var varJudgement__LeftParen_varJudgement_RightParen_r = parser.
	Get(varJudgement__LeftParen_varJudgement_RightParen).
	From(LeftParen, VarJudgement, RightParen)

func varJudgement__LeftParen_varJudgement_RightParen(nodes ...ast.Ast) ast.Ast {
	return nodes[1].(SomeJudgementNode)
}

// appJudgement <- app Colon monotype
var appJudgement__app_Colon_mono_r = parser.
	Get(appJudgement__app_Colon_mono).
	From(Application, Colon, Monotype)

func appJudgement__app_Colon_mono(nodes ...ast.Ast) ast.Ast {
	app := getApplication(nodes[0])
	// ignore Colon, i.e., element nodes[1]
	type_ := GetType(nodes[2])
	return makeSomeJudgement(AppJudgement, app, type_)
}

// appJudgement <- LeftParen appJudgement RightParen
var appJudgement__LeftParen_appJudgement_RightParen_r = parser.
	Get(appJudgement__LeftParen_appJudgement_RightParen).
	From(LeftParen, AppJudgement, RightParen)

func appJudgement__LeftParen_appJudgement_RightParen(nodes ...ast.Ast) ast.Ast {
	return nodes[1].(SomeJudgementNode)
}

// TODO:
// tupleJudgement <- tuple Colon monotype
func tupleJudgement__tuple_Colon_monotype(nodes ...ast.Ast) ast.Ast {
	app := getApplication(nodes[0])
	// ignore Colon, i.e., element nodes[1]
	type_ := GetType(nodes[2])
	return makeSomeJudgement(AppJudgement, app, type_)
}

// TODO:
// tupleJudgement <- LeftParen tupleJudgement RightParen
func tupleJudgement__LeftParen_tupleJudgement_RightParen(nodes ...ast.Ast) ast.Ast {
	return nodes[1].(SomeJudgementNode)
}

// TODO:
// anonJudgement <- anonymousFunction Colon monotype

// TODO:
// anonJudgement <- LeftParen anonJudgement RightParen
