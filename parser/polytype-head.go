package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
	itoken "github.com/petersalex27/yew-packages/token"
)

/*
polyHead      ::= 'forall' ID
                  | polyHead ID
*/

// == polytype head reduction rules ===========================================

// polyHead <- 'forall' ID
var polyHead__Forall_Id_r = parser. 
	Get(initPolyHeadReduction).From(Forall, Id)

// polyHead <- polyHead ID
var polyHead__polyHead_Id_r = parser.
	Get(appendVarReduction).From(PolyHead, Id)

type PolyHeadNode struct {
	readyToUse bool
	vars       []types.Variable[token.Token]
}

// == polytype head reductions ================================================

func initPolyHeadReduction(nodes ...ast.Ast) ast.Ast {
	const _, varIndex int = 0, 1
	head := PolyHeadNode{
		readyToUse: false, 
		vars: make([]types.Variable[token.Token],1),
	}
	v := types.Var[token.Token](nodes[varIndex].(ast.Token).Token.(token.Token))
	head.vars[0] = v
	return head
}

func appendVarReduction(nodes ...ast.Ast) ast.Ast {
	const headIndex, varIndex int = 0, 1
	head := nodes[headIndex].(PolyHeadNode)
	v := types.Var[token.Token](nodes[varIndex].(ast.Token).Token.(token.Token))
	head.vars = append(head.vars, v)
	return head
}

// TODO =================================

func (h PolyHeadNode) Equals(a ast.Ast) bool {
	h2, ok := a.(PolyHeadNode)
	if !ok {
		return false
	}
	if h.readyToUse != h2.readyToUse || len(h.vars) != len(h2.vars) {
		return false
	}

	for i, v := range h.vars {
		if !v.Equals(h2.vars[i]) {
			return false
		}
	}
	return true
}

func (h PolyHeadNode) NodeType() ast.Type {
	if h.readyToUse {
		return PolyBinders
	}
	return PolyHead
}

func (h PolyHeadNode) InOrderTraversal(f func(itoken.Token)) {
	for _, v := range h.vars {
		tok := v.Collect()
		f(tok[0])
	}
}