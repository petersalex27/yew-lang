package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
)

func astToCase(a ast.Ast) CaseNode { return a.(CaseNode) }

func astToPattern(a ast.Ast) SomeExpression { return a.(SomeExpression) }

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
case          ::= case data '->' expr
									| data '->' expr
*/

// pattern <- expr 'when' case
var pattern__expr_When_case_r = parser. 
	Get(patternReduction).From(Expr, When, Case)

// case <- case data '->' expr
var case__case_data_Arrow_expr_r = parser. 
	Get(caseJoinReduction).From(Case, Data, Arrow, Expr)

// case <- data '->' expr
var case__data_Arrow_expr_r = parser. 
	Get(caseReduction).From(Data, Arrow, Expr)

// for each case: collect vars declared in `data`, using them to bind
// corr. free vars in `expr`

var arbitraryExpressionVariable = expr.Var(
	token.Id.Make().
		AddValue("_").
		SetLineChar(1,1).(token.Token),
)

func patternReduction(nodes ...ast.Ast) ast.Ast {
	const exprIndex, _, caseIndex int = 0, 1, 2
	e := astToExpression(nodes[exprIndex]).Expression
	cs := astToCase(nodes[caseIndex])
	return SomeExpression{
		Pattern,
		expr.Select[token.Token](e, cs...),
	}
}

func caseReduction(nodes ...ast.Ast) ast.Ast {
	const dataIndex, _, exprIndex int = 0, 1, 2
	data := astToData(nodes[dataIndex]).Expression
	expression := astToExpression(nodes[exprIndex]).Expression
	var binders expr.BindersOnly[token.Token] = 
		data.ExtractFreeVariables(arbitraryExpressionVariable)
	expression.Bind(binders)
	return CaseNode{binders.InCase(data, expression),}
}

func caseJoinReduction(nodes ...ast.Ast) ast.Ast {
	const caseIndex, dataIndex, _, _ int = 0, 1, 2, 3
	rightCase := caseReduction(nodes[dataIndex:]...).(CaseNode)
	leftCase := nodes[caseIndex].(CaseNode)
	leftCase = append(leftCase, rightCase...)
	return leftCase
}