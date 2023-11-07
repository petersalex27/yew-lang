package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type FunctionNodeInstances []FunctionNode

// =============================================================================
// functions production rules
// =============================================================================

var functions__exprBlock_function_r = parser.
	Get(functionsInitialProduction).
	When(IndentExprBlock).From(Function)

var functions__Indent_function_r = parser.
	Get(functionsInitialProduction).
	When(Indent).From(Function)

var functions__exprBlock_function_Indent_functions_r = 
	parser.Precondition(
		parser.Get(functionsPrependProduction).
			When(IndentExprBlock).
			From(Function, Indent, FunctionInstances),
		indentMatchPrecondition_0_2(IndentExprBlock, Indent),
	)

var functions__Indent_function_Indent_functions_r = 
	parser.Precondition(
		parser.Get(functionsPrependProduction).
			When(Indent).
			From(Function, Indent, FunctionInstances),
		indentMatchPrecondition_0_2(Indent, Indent),
	)

// =============================================================================
// functions production functions
// =============================================================================

// creates a function instances node
func functionsInitialProduction(nodes ...ast.Ast) ast.Ast {
	const functionIndex int = 0
	function := nodes[functionIndex].(FunctionNode)
	return FunctionNodeInstances{function}
}

// prepends function to slice of functions
func functionsPrependProduction(nodes ...ast.Ast) ast.Ast {
	const functionIndex, _, functionsIndex = 0, 1, 2
	function := nodes[functionIndex].(FunctionNode)
	functions := nodes[functionsIndex].(FunctionNodeInstances)
	functions = append(FunctionNodeInstances{function}, functions...)
	return functions
}

// =============================================================================
// precondition check function
// =============================================================================

func exprBlockGetParam(eb ast.Ast) string {
	return token.Token(eb.(ExprBlockStart)).GetValue()
}

func indentGetParam(tok ast.Ast) string {
	return GetToken(tok).GetValue()
}

func indentGetParamGen(indentType ast.Type) func(ast.Ast) string {
	if indentType == IndentExprBlock {
		return exprBlockGetParam
	}
	return indentGetParam
}

func indentMatchPreconditionGen(first, second int) func(ast.Type, ast.Type) func(...ast.Ast) bool { 
	return func(left, right ast.Type) func(nodes ...ast.Ast) bool {
		return func(nodes ...ast.Ast) bool {
			indent1 := nodes[first]
			indent2 := nodes[second]
			return indentGetParamGen(left)(indent1) == indentGetParamGen(right)(indent2) 
		}
	}
}

var indentMatchPrecondition_0_2 = indentMatchPreconditionGen(0, 2)

// =============================================================================
// FunctionNodeInstances's ast.Ast implementation
// =============================================================================



func (f FunctionNodeInstances) Equals(a ast.Ast) bool {
	f2, ok := a.(FunctionNodeInstances)
	if !ok {
		return false
	}

	ok = len(f) == len(f2)
	if !ok {
		return false
	}

	for i, function := range f {
		if !function.Equals(f2[i]) {
			return false
		}
	}

	return true
}

// returns Function
func (FunctionNodeInstances) NodeType() ast.Type { return FunctionInstances }

func (fs FunctionNodeInstances) InOrderTraversal(g func(itoken.Token)) {
	for _, f := range fs {
		f.InOrderTraversal(g)
	}
}