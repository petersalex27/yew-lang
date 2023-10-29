package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type NodeSequence struct {
	ty    ast.Type
	nodes []ast.Ast
}

func (n NodeSequence) GetString(head string) string {
	return fmt.Sprintf("NodeSequence{\n%s\tty:%v,\n%s\tnodes:%v,\n%s}\n",
		head, n.ty,
		head, n.nodes,
		head,
	)
}

func getNodeSequence(node ast.Ast) NodeSequence {
	return node.(NodeSequence)
}

func mergeNodeSequenceRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		left := getNodeSequence(nodes[0])
		left.ty = ty
		left.nodes = append(left.nodes, getNodeSequence(nodes[1]).nodes...)
		return left
	}
}

func reverseConsRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		left := getNodeSequence(nodes[0])
		left.ty = ty
		left.nodes = append(left.nodes, nodes[1])
		return left
	}
}

func consRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		left := NodeSequence{ty, []ast.Ast{nodes[0]}}
		left.nodes = append(left.nodes, getNodeSequence(nodes[1]).nodes...)
		return left
	}
}

func rewrapNodeSequenceRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		res := getNodeSequence(nodes[0])
		res.ty = ty
		return res
	}
}

func createNodeSequenceRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return NodeSequence{ty, []ast.Ast{nodes[0]}}
	}
}

func (n NodeSequence) Equals(a ast.Ast) bool {
	n2, ok := a.(NodeSequence)
	if !ok {
		return false
	}

	if n.ty != n2.ty || len(n.nodes) != len(n2.nodes) {
		return false
	}

	for i, node := range n.nodes {
		if !node.Equals(n2.nodes[i]) {
			return false
		}
	}

	return true
}

func (n NodeSequence) NodeType() ast.Type { return n.ty }

func (n NodeSequence) InOrderTraversal(f func(itoken.Token)) {
	for _, node := range n.nodes {
		node.InOrderTraversal(f)
	}
}