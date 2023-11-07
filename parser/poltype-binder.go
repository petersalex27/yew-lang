package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
polyBinders   ::= polyHead                        # when l.a. is '.'
                  | '(' polyBinders ')'           #   //
*/

// == polytype binders reduction rules ========================================

var polyBinders__polyHead_r = parser.
	Get(polyBindersFromPolyHeadReduction).From(PolyHead)

var polyBinders__enclosed_r = parser.
	Get(parenEnclosedProduction).From(LeftParen, PolyBinders, RightParen)

// == polytype binders reductions =============================================

func polyBindersFromPolyHeadReduction(nodes ...ast.Ast) ast.Ast {
	const polytypeHeadIndex int = 0
	binders := nodes[polytypeHeadIndex].(PolyHeadNode).vars
	return PolyHeadNode{readyToUse: true, vars: binders}
}
