package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/token"
)

// function ::= functionDef '=' expr

// function instances
type FunctionNode struct {
	def FunctionDefNode
	prologue prologueSemiNode
	body expr.Expression[token.Token]
}

type deconstructionInstruction byte
const (
	moveLeft deconstructionInstruction = iota
	moveRight
	skipParam
	extractLeft
	extractRight
)

type prologueSemiNode struct {
	/*
	params = a (Data a (Thing (Thing _) a))
	[][]deconstructionInstruction{
		// a
		{skipParam},
		// (MyData a (Thing (Thing _) a))
		//					  /\
		//					_/  \_
		//				 /\    /\
		//     Data  a  /\ a
		//				 Thing /\
		//					Thing  *
		{moveLeft, extractRight, moveRight, extractRight},
	}
	*/
	
	deconstruct [][]deconstructionInstruction
}

func (f FunctionNode) deconstruct() FunctionNode {
	f.prologue.deconstruct = recursiveDeconstruction(f.def.head.params)
	return f
}

func functionReduction(nodes ...ast.Ast) ast.Ast {
	const funcDefIndex, _, bodyIndex int = 0, 1, 2
	funcDef := nodes[funcDefIndex].(FunctionDefNode)
	body := nodes[bodyIndex].(ExpressionNode).Expression
	return FunctionNode{
		def: funcDef,
		body: body,
	}
}

// function <- functionDef '=' expr
var function__functionDef_Assign_expr_r = parser. 
	Get(functionReduction).
	From(FunctionDefinition, Assign, Expr)

func (f FunctionNode) Equals(a ast.Ast) bool {
	f2, ok := a.(FunctionNode)
	if !ok {
		return false
	}

	return f.def.Equals(f2.def) && f.body.StrictEquals(f2.body)
}

// returns Function
func (f FunctionNode) NodeType() ast.Type { return Function }

func (f FunctionNode) InOrderTraversal(g func(itoken.Token)) {
	f.def.InOrderTraversal(g)
	for _, token := range f.body.Collect() {
		g(token)
	}
}