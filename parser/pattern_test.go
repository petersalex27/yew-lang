package parser

import (
	"testing"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
	"yew.lang/main/errors"
	"yew.lang/main/token"
)

func TestPattern(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			pattern__expr_When_case_r,
		))

	x := ExpressionNode{exprVar("x")}
	dataToken := makeTypeIdToken_test("Data", 1, 1)
	dataConst := Const(dataToken)
	a := exprVar("a")
	nameToken := makeTypeIdToken_test("Name", 1, 1)
	nameConst := Const(nameToken)
	b := exprVar("b")
	data := SomeExpression{
		Data,
		expr.Apply[token.Token](dataConst, a),
	}
	name := SomeExpression{
		Data,
		expr.Apply[token.Token](nameConst, b),
	}
	caseData := CaseNode{(expr.BindersOnly[token.Token]{a}).InCase(data.Expression, a)}
	caseData_caseName := CaseNode{
		caseData[0],
		(expr.BindersOnly[token.Token]{b}).InCase(name.Expression,b),
	}
	when := ast.TokenNode(makeToken_test(token.When,"when",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{x, when, caseData},
			parser.MakeSource("test/parser/pattern", "x when Data a -> a"),
			ast.AstRoot{
				SomeExpression{
					Pattern, 
					expr.Select[token.Token](x.Expression, caseData...),
				},
			},
		},
		{
			[]ast.Ast{x, when, caseData_caseName},
			parser.MakeSource("test/parser/pattern", 
				"x when",
				"  Data a -> a",
				"  Name b -> b",
			),
			ast.AstRoot{
				SomeExpression{
					Pattern, 
					expr.Select[token.Token](
						x.Expression, 
						caseData_caseName...,
					),
				},
			},
		},
	}

	for i, test := range tests {
		p := parser.NewParser().
			LA(1).
			UsingReductionTable(table).Load(
			[]itoken.Token{},
			test.src,
			nil, nil,
		).InitialStackPush(test.nodes...)

		actual := p.Parse()
		if p.HasErrors() {
			es := p.GetErrors()
			errors.PrintErrors(es...)
			t.Fatal(testutil.TestFail2("errors", nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(testutil.TestFail(test.expect, actual, i))
		}
	}
}

func TestCase(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			case__case_data_Arrow_expr_r,
			case__data_Arrow_expr_r,
		))

	dataToken := makeTypeIdToken_test("Data", 1, 1)
	dataConst := Const(dataToken)
	a := exprVar("a")
	nameToken := makeTypeIdToken_test("Name", 1, 1)
	nameConst := Const(nameToken)
	b := exprVar("b")
	data := SomeExpression{
		Data,
		expr.Apply[token.Token](dataConst, a),
	}
	name := SomeExpression{
		Data,
		expr.Apply[token.Token](nameConst, b),
	}
	caseData := CaseNode{(expr.BindersOnly[token.Token]{a}).InCase(data.Expression, a)}
	caseData_caseName := CaseNode{
		caseData[0],
		(expr.BindersOnly[token.Token]{b}).InCase(name.Expression,b),
	}
	arrow := ast.TokenNode(makeToken_test(token.Arrow,"->",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{caseData, name, arrow, ExpressionNode{b}},
			parser.MakeSource("test/parser/case",
				"Data a -> a",
				"Name b -> b",
			),
			ast.AstRoot{caseData_caseName},
		},
		{
			[]ast.Ast{data, arrow, ExpressionNode{a}},
			parser.MakeSource("test/parser/case", "Data a -> a"),
			ast.AstRoot{caseData},
		},
	}

	for i, test := range tests {
		p := parser.NewParser().
			LA(1).
			UsingReductionTable(table).Load(
			[]itoken.Token{},
			test.src,
			nil, nil,
		).InitialStackPush(test.nodes...)

		actual := p.Parse()
		if p.HasErrors() {
			es := p.GetErrors()
			errors.PrintErrors(es...)
			t.Fatal(testutil.TestFail2("errors", nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(testutil.TestFail(test.expect, actual, i))
		}
	}
}
