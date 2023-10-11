package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

/*
polyHead      ::= 'forall' Id
                  | polyHead Id
polyBinders   ::= polyHead                        # when l.a. is '.'
                  | '(' polyBinders ')'           #   //
*/

// polyHead <- 'forall' ID
var polyHead__Forall_Id_r = parser. 
	Get(initPolyHeadReduction).From(Forall, Id)

// polyHead <- polyHead ID
var polyHead__polyHead_Id_r = parser.
	Get(appendVarReduction).From(PolyHead, Id)

// polyBinders <- polyHead
var polyBinders__polyHead_r = parser.
	Get(func(nodes ...ast.Ast) ast.Ast { 
		return PolyHeadNode{readyToUse: true, vars: nodes[0].(PolyHeadNode).vars}
	}).
	From(PolyHead)

// polyBinders <- '(' polyBinders ')'
var polyBinders__enclosed_r = parser.
	Get(grab_enclosed).From(LeftParen, PolyBinders, RightParen)

// polytype <- polyBinders '.' dependTyped
var polytype__polyBinders_Dot_dependTyped_r = parser.
	Get(polytypeReduction).From(PolyBinders, Dot, Dependtyped)

type PolyHeadNode struct {
	readyToUse bool
	vars       []types.Variable[token.Token]
}

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

func polytypeReduction(nodes ...ast.Ast) ast.Ast {
	const bindersIndex, _, dependTypedIndex int = 0, 1, 2
	binders := nodes[bindersIndex].(PolyHeadNode).vars
	dependTyped := nodes[dependTypedIndex].(TypeNode).Type.(types.DependentTyped[token.Token])
	return TypeNode{
		Polytype,
		types.Forall[token.Token](binders...).Bind(dependTyped),
	}
}

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

// polytype <- forallHead Dot depend
var polytype__forallHead_Dot_depend_r = parser.
	Get(polytype__forallHead_Dot_depend).
	From(PolyHead, Dot, Dependtype)

func polytype__forallHead_Dot_depend(nodes ...ast.Ast) ast.Ast {
	head := nodes[0].(PolyHeadNode).vars
	dep := getDependentTyped(nodes[2])
	res := types.Forall(head...).Bind(dep)
	return TypeNode{Polytype, res}
}

// polytype <- forallHead Dot monotype
var polytype__forallHead_Dot_mono_r = parser.
	Get(polytype__forallHead_Dot_mono).
	From(PolyHead, Dot, Monotype)

func polytype__forallHead_Dot_mono(nodes ...ast.Ast) ast.Ast {
	head := nodes[0].(PolyHeadNode)
	dep := getDependentTyped(nodes[2])
	res := types.Forall(head.vars...).Bind(dep)
	return TypeNode{Polytype, res}
}

// forallHead <- Forall var
var forallHead__Forall_var_r = parser.
	Get(forallHead__Forall_var).
	From(Forall, FreeVar)

func forallHead__Forall_var(nodes ...ast.Ast) ast.Ast {
	return PolyHeadNode{false, []types.Variable[token.Token]{getVariable(nodes[1])}}
}

// forallHead <- forallHead var
var forallHead__forallHead_var_r = parser.
	Get(forallHead__forallHead_var).
	From(PolyHead, FreeVar)

func forallHead__forallHead_var(nodes ...ast.Ast) ast.Ast {
	left := nodes[0].(PolyHeadNode)
	right := getVariable(nodes[1])
	return PolyHeadNode{false, append(left.vars, right)}
}
