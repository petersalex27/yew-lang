package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

type MapHeadNode []types.TypeJudgement[token.Token, expr.Variable[token.Token]]

func (m MapHeadNode) Equals(a ast.Ast) bool {
	m2, ok := a.(MapHeadNode)
	if !ok {
		return false
	}
	if len(m) != len(m2) {
		return false
	}

	for i, judge := range m {
		if !VariableJudgement(judge).Equals(VariableJudgement(m2[i])) {
			return false
		}
	}
	return true
}

func (m MapHeadNode) NodeType() ast.Type { return MapHead }

func (m MapHeadNode) InOrderTraversal(f func(itoken.Token)) {
	for _, judge := range m {
		VariableJudgement(judge).InOrderTraversal(f)
	}
}

// depend <- mapallHead Dot typeApp
var depend__mapallHead_Dot_typeApp_r = parser. 
	Get(depend__mapallHead_Dot_typeApp). 
	From(MapHead, Dot, TypeApp)

func depend__mapallHead_Dot_typeApp(nodes ...ast.Ast) ast.Ast {
	app := getApplicationType(nodes[2])
	head := []types.TypeJudgement[token.Token, expr.Variable[token.Token]](nodes[0].(MapHeadNode))
	return TypeNode{Dependtype, types.MakeDependentType[token.Token](head, app)}
}

// depend <- mapallHead Dot TypeId
var depend__mapallHead_Dot_TypeId_r = parser. 
	Get(depend__mapallHead_Dot_TypeId). 
	From(MapHead, Dot, TypeId)

func depend__mapallHead_Dot_TypeId(nodes ...ast.Ast) ast.Ast {
	ty := types.MakeConst(GetToken(nodes[2]))
	app := types.Apply[token.Token](ty)
	head := []types.TypeJudgement[token.Token, expr.Variable[token.Token]](nodes[0].(MapHeadNode))
	return TypeNode{Dependtype, types.MakeDependentType[token.Token](head, app)}
}

// mapallHead <- Mapall varJudgement
var mapallHead__Mapall_varJudgement_r = parser. 
	Get(mapallHead__Mapall_varJudgement).
	From(Mapall, VarJudgement)

func mapallHead__Mapall_varJudgement(nodes ...ast.Ast) ast.Ast {
	return MapHeadNode{getVariableJudgement(nodes[1])}
}

// mapallHead <- mapallHead varJudgement
var mapallHead__mapallHead_varJudgement_r = parser. 
	Get(mapallHead__mapallHead_varJudgement). 
	From(MapHead, VarJudgement)

func mapallHead__mapallHead_varJudgement(nodes ...ast.Ast) ast.Ast {
	left := nodes[0].(MapHeadNode)
	right := getVariableJudgement(nodes[1])
	return MapHeadNode(append(left, right))
}