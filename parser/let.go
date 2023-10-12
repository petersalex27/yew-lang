package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

// letIn ::= 'let' function 'in' expr

var letIn__Let_function_In_expr_r = parser. 
	Get(letReduction).From(Let, FunctionDecl, In, Expr)

func letReduction(nodes ...ast.Ast) ast.Ast {
	const _, functionIndex, _, exprIndex int = 0, 1, 2, 3
	nodes[functionIndex].
}


// exprWhere ::= expr 'where' function