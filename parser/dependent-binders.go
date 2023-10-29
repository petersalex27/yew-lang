package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
)

/*
dependBinders ::= dependentHead                   # when l.a. is '.'
                  | '(' dependBinders ')'         #   //
*/

// == dependent binders reduction rules =======================================

var dependBinders__dependHead_r = parser.
	Get(dependentBindersFromHeadReduction).From(DependHead)

var dependBinders__enclosed_r = parser.
	Get(parenEnclosedReduction).From(LeftParen, DependBinders, RightParen)

// == dependent binders reductions ============================================

func dependentBindersFromHeadReduction(nodes ...ast.Ast) ast.Ast {
	const dependHeadIndex int = 0
	params := nodes[dependHeadIndex].(DependHeadNode).params
	return DependHeadNode{readyToUse: true, params: params}
}

// == dependent binders utils =================================================

func variableJudgements(judgements ...types.TypeJudgement[token.Token, expr.Variable[token.Token]]) []types.TypeJudgement[token.Token, expr.Variable[token.Token]] {
	return judgements
}
