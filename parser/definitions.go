// =============================================================================
// Author-Date: Alex Peters - 2023
//
// Content:
// rules and ast node (DefinitionsNode) related to collection blocks of
// functions and function definitions
//
// Grammar Rules:
//
//	definitions ::= ( INDENT(n) ) functionDefs INDENT(n) definitions		# grp(1)
//	                | ( INDENT(n) ) functionDefs INDENT(n) definitions  #  //
//	                | ( indent(n) ) functions INDENT(n) definitions			# grp(2)
//	                | ( indent(n) ) functionDefs INDENT(n) definitions	#  //
//	                | ( INDENT(n) ) functions														# grp(3)
//	                | ( INDENT(n) ) functionDefs												#  //
//	                | ( indent(n) ) functions														# grp(4)
//	                | ( indent(n) ) functionDefs												#  //
//
// Notes:
// definitions's grammar alternatives can be broken into four groups:
// (1) right-recursive, (2) terminal rules of right-recursive rules,
// (3) initial-cases of right-recursive rules, and (4) singleton rules. They are
// labeled accordingly above with `grp(n)` for group (n)
// =============================================================================
package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

// represents a collection of definitions / instances
type DefinitionsNode struct {
	functions    FunctionNodeInstances
	functionDefs FunctionNodeDefs
}

// =============================================================================
// definitions production rules
// =============================================================================

// one of the two right recursive defs rules
//
//	definitions ::= ( INDENT(n) ) functions INDENT(n) definitions
var definitions__Indent_functions_Indent_definitions_r = parser.Precondition(
	parser.
		Get(defFunctionsPrependProduction).
		When(Indent).
		From(FunctionInstances, Indent, Definitions),
	indentMatchPrecondition_0_2(Indent, Indent),
)

// one of the two right recursive defs rules
//
//	definitions ::= ( INDENT(n) ) functionDefs INDENT(n) definitions
var definitions__Indent_functionDefs_Indent_definitions_r = parser.Precondition(
	parser.
		Get(defFuncDefsPrependProduction).
		When(Indent).
		From(FunctionDefinitions, Indent, Definitions),
	indentMatchPrecondition_0_2(Indent, Indent),
)

// one of two rules that terminate right recursive rules
//
//	definitions ::= ( indent(n) ) functions INDENT(n) definitions
var definitions__exprBlock_functions_Indent_definitions_r = parser.Precondition(
	parser.
		Get(defFunctionsPrependProduction).
		When(IndentExprBlock).
		From(FunctionInstances, Indent, Definitions),
	indentMatchPrecondition_0_2(IndentExprBlock, Indent),
)

// one of two rules that terminate right recursive rules
//
//	definitions ::= ( indent(n) ) functionDefs INDENT(n) definitions
var definitions__exprBlock_functionDefs_Indent_definitions_r = parser.Precondition(
	parser.
		Get(defFuncDefsPrependProduction).
		When(IndentExprBlock).
		From(FunctionDefinitions, Indent, Definitions),
	indentMatchPrecondition_0_2(IndentExprBlock, Indent),
)

// one of the two base cases for right recursive rules
//
//	definitions ::= ( INDENT(n) ) functions
var definitions__Indent_functions_r = parser.
	Get(defFunctionsSingleProduction).
	When(Indent).
	From(FunctionInstances)

// one of the two base cases for right recursive rules
//
//	definitions ::= ( INDENT(n) ) functionDefs
var definitions__Indent_functionDefs_r = parser.
	Get(defFuncDefsSingleProduction).
	When(Indent).
	From(FunctionDefinitions)

// only used when there is one sequence of functions in the expr block
//
//	definitions ::= ( indent(n) ) functions
var definitions__exprBlock_functions_r = parser.
	Get(defFunctionsSingleProduction).
	When(IndentExprBlock).
	From(FunctionInstances)

// only used when there is on sequence of function defs in the expr block
//
//	definitions ::= ( indent(n) ) functionDefs
var definitions__exprBlock_functionDefs_r = parser.
	Get(defFuncDefsSingleProduction).
	When(IndentExprBlock).
	From(FunctionDefinitions)

// =============================================================================
// production functions
// =============================================================================

func defFunctionsSingleProduction(nodes ...ast.Ast) ast.Ast {
	const funcsIndex int = 0
	functions := nodes[funcsIndex].(FunctionNodeInstances)
	return DefinitionsNode{
		functions:    functions,
		functionDefs: []FunctionDefNode{},
	}
}

func defFuncDefsSingleProduction(nodes ...ast.Ast) ast.Ast {
	const funcDefsIndex int = 0
	funcDefs := nodes[funcDefsIndex].(FunctionNodeDefs)
	return DefinitionsNode{
		functions:    []FunctionNode{},
		functionDefs: funcDefs,
	}
}

func defFunctionsPrependProduction(nodes ...ast.Ast) ast.Ast {
	const funcsIndex, _, defsIndex int = 0, 1, 2
	functions := nodes[funcsIndex].(FunctionNodeInstances)
	defs := nodes[defsIndex].(DefinitionsNode)
	return DefinitionsNode{
		functions:    append(functions, defs.functions...),
		functionDefs: defs.functionDefs,
	}
}

func defFuncDefsPrependProduction(nodes ...ast.Ast) ast.Ast {
	const funcDefsIndex, _, defsIndex int = 0, 1, 2
	funcDefs := nodes[funcDefsIndex].(FunctionNodeDefs)
	defs := nodes[defsIndex].(DefinitionsNode)
	return DefinitionsNode{
		functions:    defs.functions,
		functionDefs: append(funcDefs, defs.functionDefs...),
	}
}

// == DefinitionsNode implementation of ast.Ast ===============================

func (defs DefinitionsNode) Equals(a ast.Ast) bool {
	defs2, ok := a.(DefinitionsNode)
	if !ok {
		return false
	}

	return defs.functionDefs.Equals(defs2.functionDefs) &&
		defs.functions.Equals(defs2.functions)
}

// returns Definitions
func (defs DefinitionsNode) NodeType() ast.Type { return Definitions }

func (defs DefinitionsNode) InOrderTraversal(f func(itoken.Token)) {
	defs.functions.InOrderTraversal(f)
	defs.functionDefs.InOrderTraversal(f)
}
