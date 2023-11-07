package parser

import (
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
functionDefs  ::= ( INDENT(n) ) functionDef INDENT(n) functionDefs
                  | ( INDENT(n) ) functionDef
                  | ( indent(n) ) functionDef INDENT(n) functionDefs
                  | ( indent(n) ) functionDef
*/

type FunctionNodeDefs []FunctionDefNode

// =============================================================================
// function defs productions
// =============================================================================

var functionDefs__Indent_functionDef_Indent_functionDefs_r = 
	parser.Precondition(
		parser.
			Get(functionDefsPrependProduction).
			When(Indent).
			From(FunctionDefinition, Indent, FunctionDefinitions),
		indentMatchPrecondition_0_2(Indent, Indent),
	)

var functionDefs__Indent_functionDef_r = parser.
	Get(functionDefsInitialProduction).
	When(Indent).
	From(FunctionDefinition)

var functionDefs__exprBlock_functionDef_Indent_functionDefs_r = 
	parser.Precondition(
		parser.
			Get(functionDefsPrependProduction).
			When(IndentExprBlock).
			From(FunctionDefinition, Indent, FunctionDefinitions),
		indentMatchPrecondition_0_2(IndentExprBlock, Indent),
	)

var functionDefs__exprBlock_functionDef_r = parser.
	Get(functionDefsInitialProduction).
	When(IndentExprBlock).
	From(FunctionDefinition)


// =============================================================================
// functions production functions
// =============================================================================

// creates a function instances node
func functionDefsInitialProduction(nodes ...ast.Ast) ast.Ast {
	const defIndex int = 0
	def := nodes[defIndex].(FunctionDefNode)
	return FunctionNodeDefs{def}
}

// prepends function to slice of functions
func functionDefsPrependProduction(nodes ...ast.Ast) ast.Ast {
	const defIndex, _, defsIndex = 0, 1, 2
	def := nodes[defIndex].(FunctionDefNode)
	defs := nodes[defsIndex].(FunctionNodeDefs)
	defs = append(FunctionNodeDefs{def}, defs...)
	return defs
}

// =============================================================================
// FunctionNodeInstances's ast.Ast implementation
// =============================================================================

func (f FunctionNodeDefs) Equals(a ast.Ast) bool {
	f2, ok := a.(FunctionNodeDefs)
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
func (FunctionNodeDefs) NodeType() ast.Type { return FunctionDefinitions }

func (fs FunctionNodeDefs) InOrderTraversal(g func(itoken.Token)) {
	for _, f := range fs {
		f.InOrderTraversal(g)
	}
}