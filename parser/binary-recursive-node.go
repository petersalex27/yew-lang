package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
)

type BinaryRecursiveNode interface {
	ast.Ast
	UpdateType(target, then ast.Type) BinaryRecursiveNode
	GetString(head string) string
	SplitNode() (left, right BinaryRecursiveNode)
	HasValue() (val Node, ok bool)
}

func getBinaryRecursiveNode(node ast.Ast) BinaryRecursiveNode {
	return node.(BinaryRecursiveNode)
}

// changes type of BinaryRecursiveNode
//
//	nodes[0].(BinaryRecursiveNode).NodeType() == t =>
//	rewrapReduction(t2)(nodes[0]) =>
//	nodes[0].(BinaryRecursiveNode).NodeType() == t2
func rewrapReduction(newType ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		brn := nodes[0].(BinaryRecursiveNode)
		// will always update b/c target = brn.NodeType() and brn is updated iff
		// target == brn.NodeType()
		return brn.UpdateType(brn.NodeType(), newType)
	}
}
