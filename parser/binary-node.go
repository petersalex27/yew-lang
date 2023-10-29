package parser

import (
	"fmt"

	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type BinaryNode struct {
	ty          ast.Type
	left, right BinaryRecursiveNode
}

func (n BinaryNode) UpdateType(target, then ast.Type) BinaryRecursiveNode {
	if n.ty == target {
		n.ty = then
	}
	n.left = n.left.UpdateType(target, then)
	n.right = n.right.UpdateType(target, then)
	return n
}

func (n BinaryNode) GetString(head string) string {
	return fmt.Sprintf("BinaryNode{\n%s\tty:%v,\n%s\tleft:%s,\n%s\tright:%s\n%s}\n",
		head, n.ty,
		head, n.left.GetString(head+"\t"),
		head, n.right.GetString(head+"\t"),
		head,
	)
}

func (node BinaryNode) SplitNode() (left, right BinaryRecursiveNode) {
	return node.left, node.right
}

func (node BinaryNode) HasValue() (val Node, ok bool) { return val, false }

func (n BinaryNode) Equals(a ast.Ast) bool {
	n2, ok := a.(BinaryNode)
	if !ok {
		return false
	}

	if n.ty != n2.ty {
		return false
	}

	left, right := n.SplitNode()
	left2, right2 := n2.SplitNode()

	if left == nil || left2 == nil {
		if left != left2 {
			return false
		}
	}

	if right == nil || right2 == nil {
		if right != right2 {
			return false
		}
	}

	if left == nil && right == nil {
		return true
	} else if left == nil {
		return right.Equals(right2)
	} else if right == nil {
		return left.Equals(left2)
	}

	return left.Equals(left2) && right.Equals(right2)
}

func (n BinaryNode) NodeType() ast.Type { return n.ty }

func (n BinaryNode) InOrderTraversal(f func(itoken.Token)) {
	left, right := n.SplitNode()

	if left != nil {
		left.InOrderTraversal(f)
	}

	if right != nil {
		right.InOrderTraversal(f)
	}
}

// ty -> ((BinRec1, BinRec2) -> BinaryNode{ty, BinRec1, BinRec2})
func simpleBinaryNodeRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		left := nodes[0].(BinaryRecursiveNode)
		right := nodes[1].(BinaryRecursiveNode)
		return BinaryNode{ty, left, right}
	}
}

func binarySelect(rule func(nodes ...ast.Ast) ast.Ast, first, second int) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return rule(nodes[first], nodes[second])
	}
}