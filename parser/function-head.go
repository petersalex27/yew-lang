package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/token"
)

/*
functionDecl  ::= INDENT funcName
functionHead  ::= functionDecl pattern
                  | functionDecl                # only when L.A. is '=' or ':'
functionDef   ::= functionHead ':' type
									| functionHead      					# only when L.A. is '='
*/

type FunctionDeclNode_ token.Token

// function declarations--declares a function refered to as `name` with an
// application pattern of `params[0] params[1] .. params[len(params)-1]` exists
type FunctionHeadNode struct {
	name   token.Token
	params []expr.Expression[token.Token]
}

// takes a parameter node and converts it into a Variable expression
//
// ASSUMPTION: `paramNode` has type `Node`
func paramToVarExpression(paramNode ast.Ast) expr.Variable[token.Token] {
	paramToken := paramAsToken(paramNode)
	return expr.Var(paramToken)
}

func (fd FunctionHeadNode) appendExpressionToFuncHead(ex expr.Expression[token.Token]) FunctionHeadNode {
	newParams := make([]expr.Expression[token.Token], len(fd.params)+1)
	copy(newParams, fd.params)
	newParams[len(fd.params)] = ex
	return FunctionHeadNode{
		name:   fd.name,
		params: newParams,
	}
}

// reduction: funcDecl <- funcName param
func funcHeadInitialParamReduction(nodes ...ast.Ast) ast.Ast {
	const funcNameIndex, paramIndex int = 0, 1
	name := funcNameAsToken(nodes[funcNameIndex])
	param := paramToVarExpression(nodes[paramIndex])
	return FunctionHeadNode{
		name:   name,
		params: []expr.Expression[token.Token]{param},
	}
}

// reduction: funcHead <- funcDecl pattern
func funcHeadPatternReduction(nodes ...ast.Ast) ast.Ast {
	const funcDeclIndex, patternIndex int = 0, 1
	name := nodes[funcDeclIndex].(Node).Token
	pattern := nodes[patternIndex].(SomeExpression).Expression.(expr.List[token.Token])
	return FunctionHeadNode{
		name:   name,
		params: pattern,
	}
}

// reduction: funcHead <- funcDecl
func funcHeadFromFuncDeclReduction(nodes ...ast.Ast) ast.Ast {
	const funcDeclIndex int = 0
	name := nodes[funcDeclIndex].(Node).Token
	var pattern expr.List[token.Token] = nil
	return FunctionHeadNode{
		name:   name,
		params: pattern,
	}
}

// funcHead <- funcDecl pattern
var funcHead__funcDecl_pattern_r = parser.
	Get(funcHeadPatternReduction).
	From(FunctionDecl, Pattern)

// funcHead <- funcDecl
var funcHead__funcDecl_r = parser.
	Get(funcHeadFromFuncDeclReduction).
	From(FunctionDecl)

func (fd FunctionHeadNode) Equals(a ast.Ast) bool {
	fd2, ok := a.(FunctionHeadNode)
	if !ok {
		return false
	}

	if !token.TokenEquals(fd.name, fd2.name) {
		return false
	}

	if len(fd.params) != len(fd2.params) {
		return false
	}

	for i, param := range fd.params {
		if !param.StrictEquals(fd2.params[i]) {
			return false
		}
	}
	return true
}

func (fd FunctionHeadNode) NodeType() ast.Type { return FunctionHead }

func (fd FunctionHeadNode) InOrderTraversal(f func(itoken.Token)) {
	f(fd.name)
	for _, param := range fd.params {
		for _, token := range param.Collect() {
			f(token)
		}
	}
}
