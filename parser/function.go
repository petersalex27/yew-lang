package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
)

// function ::= functionDef '=' expr

// function instances
type FunctionNode struct {
	def FunctionDefNode
	body expr.Expression[token.Token]
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