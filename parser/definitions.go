package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

/*
definitions   ::= INDENT function
                  | INDENT functionDef
                  | INDENT function definitions
                  | INDENT functionDef definitions
*/

type DefinitionsNode struct {
	functions []FunctionNode
	functionDefs []FunctionDefNode
}

// == reduction rules =========================================================

var definitions__Indent_function_r = parser.
	Get(defFunctionReduction).From(Indent, Function)

var definitions__Indent_funcDef_r = parser.
	Get(defFuncDefReduction).From(Indent, FunctionDefinition)

var definitions__Indent_function_definitions_r = parser.
	Get(defFunctionAppendReduction).From(Indent, Function, Definitions)

var definitions__Indent_funcDef_definitions_r = parser.
	Get(defFuncDefAppendReduction).From(Indent, FunctionDefinition, Definitions)

// == reductions ==============================================================

func defFunctionReduction(nodes ...ast.Ast) ast.Ast {
	const _, funcIndex int = 0, 1
	function := nodes[funcIndex].(FunctionNode)
	return DefinitionsNode{
		functions: []FunctionNode{function},
		functionDefs: []FunctionDefNode{},
	}
}

func defFuncDefReduction(nodes ...ast.Ast) ast.Ast {
	const _, funcDefIndex int = 0, 1
	funcDef := nodes[funcDefIndex].(FunctionDefNode)
	return DefinitionsNode{
		functions: []FunctionNode{},
		functionDefs: []FunctionDefNode{funcDef},
	}
}

func defFunctionAppendReduction(nodes ...ast.Ast) ast.Ast {
	const _, funcIndex, defsIndex int = 0, 1, 2
	defs := nodes[defsIndex].(DefinitionsNode)
	function := nodes[funcIndex].(FunctionNode)
	return DefinitionsNode{
		functions: append(defs.functions, function),
		functionDefs: defs.functionDefs,
	}
}

func defFuncDefAppendReduction(nodes ...ast.Ast) ast.Ast {
	const _, funcDefIndex, defsIndex int = 0, 1, 2
	defs := nodes[defsIndex].(DefinitionsNode)
	funcDef := nodes[funcDefIndex].(FunctionDefNode)
	return DefinitionsNode{
		functions: defs.functions,
		functionDefs: append(defs.functionDefs, funcDef),
	}
}

// == DefinitionsNode implementation of ast.Ast ===============================

func (defs DefinitionsNode) Equals(a ast.Ast) bool {
	defs2, ok := a.(DefinitionsNode)
	if !ok {
		return false
	}

	ok = len(defs.functionDefs) == len(defs2.functionDefs) &&
		len(defs.functions) == len(defs2.functions) 
	if !ok {
		return false
	}

	for i, function := range defs.functions {
		if !function.Equals(defs2.functions[i]) {
			return false
		}
	}

	for i, funcDef := range defs.functionDefs {
		if !funcDef.Equals(defs2.functionDefs[i]) {
			return false
		}
	}

	return true
}

// returns Definitions
func (defs DefinitionsNode) NodeType() ast.Type { return Definitions }

func (defs DefinitionsNode) InOrderTraversal(f func(itoken.Token)) {
	for _, function := range defs.functions {
		function.InOrderTraversal(f)
	}

	for _, funcDef := range defs.functionDefs {
		funcDef.InOrderTraversal(f)
	}
}