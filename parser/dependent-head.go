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
dependentHead ::= 'mapall' varJudgement
                  | dependentHead varJudgement
*/

type DependHeadNode struct {
	readyToUse bool
	params []types.TypeJudgement[token.Token, expr.Variable[token.Token]]
}

// dependHead <- 'mapall' varJudgement
var dependHead__Mapall_varJudgement_r = parser. 
	Get(initDependHeadReduction).From(Mapall, VarJudgement)

// dependHead <- dependHead varJudgement
var dependHead__dependHead_varJudgement_r = parser.
	Get(appendVarJudgementReduction).From(DependHead, VarJudgement)

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