package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

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