package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

/*
dependentHead ::= 'mapall' varJudgement
                  | dependentHead varJudgement
dependBinders ::= dependentHead                   # when l.a. is '.'
                  | '(' dependBinders ')'         #   //
dependent     ::= dependBinders '.' monotype
*/

// dependHead <- 'mapall' varJudgement
var dependHead__Mapall_varJudgement_r = parser. 
	Get(initDependHeadReduction).From(Mapall, VarJudgement)

// dependHead <- dependHead varJudgement
var dependHead__dependHead_varJudgement_r = parser.
	Get(appendVarJudgementReduction).From(DependHead, VarJudgement)

// dependBinders <- dependHead
var dependBinders__dependHead_r = parser. 
	Get(func(nodes ...ast.Ast) ast.Ast {
		params := nodes[0].(DependHeadNode).params
		return DependHeadNode{
			readyToUse: true,
			params: params,
		}
	}).From(DependHead)

// dependBinders <- '(' dependBinders ')'
var dependBinders__enclosed_r = parser. 
	Get(grab_enclosed).From(LeftParen, DependBinders, RightParen)

// dependent <- dependBinders '.' monotype
var dependent__dependBinders_Dot_monotype_r = parser.
	Get(dependTypeReduction).From(DependBinders, Dot, Monotype)

func initDependHeadReduction(nodes ...ast.Ast) ast.Ast {
	const _, judgeIndex int = 0, 1
	judge := getVariableJudgement(nodes[judgeIndex])
	return DependHeadNode{
		readyToUse: false,
		params: []types.TypeJudgement[token.Token, expr.Variable[token.Token]]{
			judge,
		},
	}
}

func appendVarJudgementReduction(nodes ...ast.Ast) ast.Ast {
	const headIndex, judgeIndex int = 0, 1
	params := nodes[headIndex].(DependHeadNode).params
	judge := getVariableJudgement(nodes[judgeIndex])
	return DependHeadNode{
		readyToUse: false,
		params: append(params, judge),
	}
}

func dependTypeReduction(nodes ...ast.Ast) ast.Ast {
	const bindersIndex, _, monoIndex int = 0, 1, 2
	binders := nodes[bindersIndex].(DependHeadNode).params
	mono := nodes[monoIndex].(TypeNode).Type.(types.Monotyped[token.Token])
	return TypeNode{
		Dependtype,
		types.MakeDependentType[token.Token](
			binders,
			types.Apply[token.Token](mono.(types.ReferableType[token.Token])),
		),
	}
}

type DependHeadNode struct {
	readyToUse bool
	params []types.TypeJudgement[token.Token, expr.Variable[token.Token]]
}

func (m DependHeadNode) Equals(a ast.Ast) bool {
	m2, ok := a.(DependHeadNode)
	if !ok {
		return false
	}
	if m.readyToUse != m2.readyToUse || len(m.params) != len(m2.params) {
		return false
	}

	for i, judge := range m.params {
		if !VariableJudgement(judge).Equals(VariableJudgement(m2.params[i])) {
			return false
		}
	}
	return true
}

func (m DependHeadNode) NodeType() ast.Type { 
	if m.readyToUse {
		return DependBinders
	}
	return DependHead 
}

func (m DependHeadNode) InOrderTraversal(f func(itoken.Token)) {
	for _, judge := range m.params {
		VariableJudgement(judge).InOrderTraversal(f)
	}
}

// depend <- mapallHead Dot typeApp
var depend__mapallHead_Dot_typeApp_r = parser.
	Get(depend__mapallHead_Dot_typeApp).
	From(DependHead, Dot, TypeApp)

func depend__mapallHead_Dot_typeApp(nodes ...ast.Ast) ast.Ast {
	app := getApplicationType(nodes[2])
	head := []types.TypeJudgement[token.Token, expr.Variable[token.Token]](nodes[0].(DependHeadNode).params)
	return TypeNode{Dependtype, types.MakeDependentType[token.Token](head, app)}
}

// depend <- mapallHead Dot TypeId
var depend__mapallHead_Dot_TypeId_r = parser.
	Get(depend__mapallHead_Dot_TypeId).
	From(DependHead, Dot, TypeId)

func depend__mapallHead_Dot_TypeId(nodes ...ast.Ast) ast.Ast {
	ty := types.MakeConst(GetToken(nodes[2]))
	app := types.Apply[token.Token](ty)
	head := []types.TypeJudgement[token.Token, expr.Variable[token.Token]](nodes[0].(DependHeadNode).params)
	return TypeNode{Dependtype, types.MakeDependentType[token.Token](head, app)}
}

// mapallHead <- Mapall varJudgement
var mapallHead__Mapall_varJudgement_r = parser.
	Get(mapallHead__Mapall_varJudgement).
	From(Mapall, VarJudgement)

func mapallHead__Mapall_varJudgement(nodes ...ast.Ast) ast.Ast {
	return DependHeadNode{
		false, 
		[]types.TypeJudgement[token.Token, expr.Variable[token.Token]]{
			getVariableJudgement(nodes[1]),
		},
	}
}

// mapallHead <- mapallHead varJudgement
var mapallHead__mapallHead_varJudgement_r = parser.
	Get(mapallHead__mapallHead_varJudgement).
	From(DependHead, VarJudgement)

func mapallHead__mapallHead_varJudgement(nodes ...ast.Ast) ast.Ast {
	left := nodes[0].(DependHeadNode).params
	right := getVariableJudgement(nodes[1])
	return DependHeadNode{false, append(left, right)}
}
