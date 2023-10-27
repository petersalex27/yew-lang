package parser

import (
	"fmt"
	"sync"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

// prefix for compiler generated parameter names
const parameterStringPrefix string = "$p"

// prefix for compiler generated data member variables
const memberStringPrefix string = "$m"

// prefix for compiler generated type names
const typeStringPrefix string = "$t"

// indexes must be sorted from low to high
func sliceRule(rule func(...ast.Ast) ast.Ast, indexes ...int) func(nodes ...ast.Ast) ast.Ast {
	n := len(indexes)
	return func(nodes ...ast.Ast) ast.Ast {
		buff := make([]ast.Ast, n)
		for i, index := range indexes {
			buff[i] = nodes[index]
		}
		return rule(buff...)
	}
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

type glb_cxt_t struct {
	typeMutex sync.Mutex
	exprMutex sync.Mutex
	typeCxt   *types.Context[token.Token]
	exprCxt   *expr.Context[token.Token]
}

var glb_cxt *glb_cxt_t

func reInit() {
	glb_cxt = new(glb_cxt_t)
	glb_cxt.exprCxt = expr.NewContext[token.Token]().SetNameMaker(
		func(s string) token.Token {
			return token.Id.Make().AddValue(s)
		},
	)
	glb_cxt.typeCxt = types.NewContext[token.Token]().SetNameMaker(
		func(s string) token.Token {
			return token.Id.Make().AddValue(s)
		},
	)
}

func init() {
	reInit()
}

// creates a new type variable
func newTypeVar() types.Variable[token.Token] {
	glb_cxt.typeMutex.Lock()
	defer glb_cxt.typeMutex.Unlock()

	return glb_cxt.typeCxt.NewVar()
}

func lockType() {
	glb_cxt.typeMutex.Lock()
}

func unlockType() {
	glb_cxt.typeMutex.Unlock()
}

// creates a new kind variable
func newKindVar() expr.Variable[token.Token] {
	glb_cxt.exprMutex.Lock()
	defer glb_cxt.exprMutex.Unlock()

	return glb_cxt.exprCxt.NewVar()
}

func generateGetConst(constantName string) func() types.Constant[token.Token] {
	return func() types.Constant[token.Token] {
		tok := token.TypeId.Make().AddValue(constantName)
		return types.MakeConst(tok)
	}
}

var getUint = generateGetConst("Uint")
var getInt = generateGetConst("Int")
var getAny = generateGetConst("@any")
var getString = generateGetConst("String")
var getChar = generateGetConst("Char")
var getBool = generateGetConst("Bool")

func makeFreeJudgementOf(ty TypeNode) JudgementNode {
	varNode := ExpressionNode{newKindVar()}
	return makeJudgement(varNode, ty)
}

func GetToken(a ast.Ast) token.Token {
	tmp, _ := a.(ast.Token)
	tok, _ := tmp.Token.(token.Token)
	return tok
}

// a <- LeftParen a RightParen
func grab_enclosed(nodes ...ast.Ast) ast.Ast {
	return nodes[1]
}

func EqualsToken(a, b token.Token) bool {
	lineA, charA := a.GetLineChar()
	lineB, charB := b.GetLineChar()
	tyA, tyB := a.GetType(), b.GetType()
	valA, valB := a.GetValue(), b.GetValue()
	return lineA == lineB &&
		charA == charB &&
		tyA == tyB &&
		valA == valB
}

type Node struct {
	ty ast.Type
	token.Token
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

func getNode(node ast.Ast) Node {
	return node.(Node)
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

// ty -> (TokenNode{Token} -> Node{ty, Token})
func simpleNodeRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return Node{ty, GetToken(nodes[0])}
	}
}

// ty -> (Node{ty2, Token} -> Node{ty, Token})
func rewrapNodeRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return Node{ty, nodes[0].(Node).Token}
	}
}

func monoSelect(rule func(nodes ...ast.Ast) ast.Ast, at int) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return rule(nodes[at])
	}
}

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

func getBinaryNode(node ast.Ast) BinaryNode {
	return node.(BinaryNode)
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

type Wrapper struct {
	ast.Type
	ast.Ast
}

func (w Wrapper) Equals(a ast.Ast) bool {
	w2, ok := a.(Wrapper)
	if !ok {
		return false
	}

	return w.Type == w2.Type && w.Ast.Equals(w2.Ast)
}

func (w Wrapper) NodeType() ast.Type { return w.Type }

func (w Wrapper) InOrderTraversal(f func(itoken.Token)) { w.Ast.InOrderTraversal(f) }

func wrapRule(ty ast.Type) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast { return Wrapper{ty, nodes[0]} }
}
