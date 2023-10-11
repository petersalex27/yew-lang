package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

/*
type          ::= polytype
                  | monotype
                  | dependent
                  | '(' type ')'
*/

func typeRewrapReduction(nodes ...ast.Ast) ast.Ast {
	return TypeNode{
		Type,
		nodes[0].(TypeNode).Type,
	}
}

// type <- polytype
var type__polytype_r = parser.
	Get(typeRewrapReduction).From(Polytype)

// type <- monotype
var type__monotype_r = parser.
	Get(typeRewrapReduction).From(Monotype)

// type <- dependent
var type__dependent_r = parser.
	Get(typeRewrapReduction).From(Dependtype)

// type <- '(' type ')'
var type__enclosed_r = parser. 
	Get(grab_enclosed).From(LeftParen, Type, RightParen)

type TypeNode struct{
	ty ast.Type 
	types.Type[token.Token]
}

func (t TypeNode) Equals(a ast.Ast) bool {
	t2, ok := a.(TypeNode)
	if !ok {
		return false
	}
	return t.ty == t2.ty && t.Type.Equals(t2.Type)
}

func (t TypeNode) NodeType() ast.Type { return t.ty }

func (t TypeNode) InOrderTraversal(f func(itoken.Token)) {
	tmp := t.Type.Collect()
	for _, a := range tmp {
		f(a)
	}
}