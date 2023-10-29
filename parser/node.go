package parser

import (
	"fmt"

	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

type Node struct {
	ty ast.Type
	token.Token
}

// == node reduction rule generators ==========================================

// given an ast type `ty`, generate the following reduction:
//
//	Node{ty, someToken} ::= ast.Token{someToken}
func giveTypeToTokenReductionGen(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return Node{ty, GetToken(nodes[0])}
	}
}

// given an ast type `ty`, generate the following reduction:
//
//	Node{ty, someToken} ::= Node{ty2, someToken}
func retypeNodeReductionGen(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return Node{ty, nodes[0].(Node).Token}
	}
}

func (n Node) UpdateType(target, then ast.Type) BinaryRecursiveNode {
	if n.ty == target {
		n.ty = then
	}
	return n
}

func (n Node) String() string {
	return fmt.Sprintf("Node{ty: %v, Token: %v}", n.ty, n.Token)
}

func (n Node) GetString(head string) string {
	return fmt.Sprintf("Node{\n%s\tty:%v,\n%s\tToken:%v,\n%s}\n",
		head, n.ty,
		head, n.Token,
		head,
	)
}

func asType(node Node) types.ReferableType[token.Token] {
	if node.Token.GetType() == uint(token.TypeId) {
		return types.MakeConst[token.Token](node.Token)
	} else if node.Token.GetType() == uint(token.Id) {
		return types.Var[token.Token](node.Token)
	}
	panic("illegal operation: node.Token's type is not a valid type")
}

func (n Node) Equals(a ast.Ast) bool {
	n2, ok := a.(Node)
	if !ok {
		return false
	}

	return n.ty == n2.ty && EqualsToken(n.Token, n2.Token)
}

func (n Node) NodeType() ast.Type { return n.ty }

func (n Node) InOrderTraversal(f func(itoken.Token)) { f(n.Token) }

func (node Node) SplitNode() (left, right BinaryRecursiveNode) { return nil, nil }

func (node Node) HasValue() (val Node, ok bool) { return node, true }
