package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

// polytype <- Forall freevars Dot monotype
// polytype <- Forall var Dot monotype
// polytype <- Forall freevars Dot monotype
// polytype <- Forall var Dot monotype

type ForallHeadNode []types.Variable[token.Token]

func (h ForallHeadNode) Equals(a ast.Ast) bool {
	h2, ok := a.(ForallHeadNode)
	if !ok {
		return false
	}
	if len(h) != len(h2) {
		return false
	}

	for i, v := range h {
		if !v.Equals(h2[i]) {
			return false
		}
	}
	return true
}

func (h ForallHeadNode) NodeType() ast.Type { return ForallHead }

func (h ForallHeadNode) InOrderTraversal(f func(itoken.Token)) {
	for _, v := range h {
		tok := v.Collect()
		f(tok[0])
	}
}

// polytype <- forallHead Dot depend
var polytype__forallHead_Dot_depend_r = parser. 
	Get(polytype__forallHead_Dot_depend). 
	From(ForallHead, Dot, Dependtype)

func polytype__forallHead_Dot_depend(nodes ...ast.Ast) ast.Ast {
	head := nodes[0].(ForallHeadNode)
	dep := getDependentTyped(nodes[2])
	res := types.Forall(head...).Bind(dep)
	return TypeNode{Polytype, res}
}

// polytype <- forallHead Dot monotype
var polytype__forallHead_Dot_mono_r = parser. 
	Get(polytype__forallHead_Dot_mono). 
	From(ForallHead, Dot, Monotype)

func polytype__forallHead_Dot_mono(nodes ...ast.Ast) ast.Ast {
	head := nodes[0].(ForallHeadNode)
	dep := getDependentTyped(nodes[2])
	res := types.Forall(head...).Bind(dep)
	return TypeNode{Polytype, res}
}

// forallHead <- Forall var
var forallHead__Forall_var_r = parser. 
	Get(forallHead__Forall_var).
	From(Forall, FreeVar)

func forallHead__Forall_var(nodes ...ast.Ast) ast.Ast {
	return ForallHeadNode{getVariable(nodes[1])}
}

// forallHead <- forallHead var
var forallHead__forallHead_var_r = parser. 
	Get(forallHead__forallHead_var). 
	From(ForallHead, FreeVar)

func forallHead__forallHead_var(nodes ...ast.Ast) ast.Ast {
	left := nodes[0].(ForallHeadNode)
	right := getVariable(nodes[1])
	return ForallHeadNode(append(left, right))
}