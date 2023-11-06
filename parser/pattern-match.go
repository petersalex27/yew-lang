package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

func astToCase(a ast.Ast) CaseNode { return a.(CaseNode) }

type CaseNode []expr.Case[token.Token]

func (cs CaseNode) Equals(a ast.Ast) bool {
	cs2, ok := a.(CaseNode)
	if !ok {
		return false
	}

	if len(cs) != len(cs2) {
		return false
	}

	for i, c := range cs {
		if !c.StrictEquals(cs2[i]) {
			return false
		}
	}

	return true
}

func (c CaseNode) NodeType() ast.Type { return Case }

func (cs CaseNode) InOrderTraversal(f func(itoken.Token)) {
	for _, c := range cs {
		for _, tok := range c.Collect() {
			f(tok)
		}
	}
}

/*
pattern       ::= expr 'when' case
case          ::= case pattern '->' expr
									| pattern '->' expr
*/

// patternMatch <- expr 'when' case
var patternMatch__expr_Match_case_r = parser.
	Get(patternMatchReduction).From(Expr, Match, Case)

// case <- case pattern '->' expr
var case__case_pattern_Arrow_expr_r = parser.
	Get(caseJoinReduction).From(Case, Pattern, Arrow, Expr)

// case <- pattern '->' expr
var case__pattern_Arrow_expr_r = parser.
	Get(caseReduction).From(Pattern, Arrow, Expr)

// for each case: collect vars declared in `data`, using them to bind
// corr. free vars in `expr`

var arbitraryExpressionVariable = expr.Var(
	token.Id.Make().
		AddValue("_").
		SetLineChar(1, 1).(token.Token),
)

func patternMatchReduction(nodes ...ast.Ast) ast.Ast {
	const exprIndex, _, caseIndex int = 0, 1, 2
	e := astToExpression(nodes[exprIndex]).Expression
	cs := astToCase(nodes[caseIndex])
	return SomeExpression{
		PatternMatch,
		expr.Select[token.Token](e, cs...),
	}
}

func caseReduction(nodes ...ast.Ast) ast.Ast {
	const patternIndex, _, exprIndex int = 0, 1, 2
	pattern := nodes[patternIndex].(SomeExpression).Expression
	expression := astToExpression(nodes[exprIndex]).Expression
	var binders expr.BindersOnly[token.Token] = pattern.ExtractFreeVariables(arbitraryExpressionVariable)
	expression.Bind(binders)
	return CaseNode{binders.InCase(pattern, expression)}
}

func caseJoinReduction(nodes ...ast.Ast) ast.Ast {
	const caseIndex, rightCaseStartIndex, _, _ int = 0, 1, 2, 3

	// right case from index 1, 2, and 3 (which is why index 2 and 3 need not
	// have associated const identifiers)
	rightCase := caseReduction(nodes[rightCaseStartIndex:]...).(CaseNode)

	// left case is at index 0
	leftCase := nodes[caseIndex].(CaseNode)
	// case ::= case pattern '->' expr
	leftCase = append(leftCase, rightCase...)
	return leftCase
}
