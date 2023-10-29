package parser

import (
	"testing"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestPatternMatch(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			patternMatch__expr_When_case_r,
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
			parser.MakeSource("test/parser/pattern-match", "x when Data a -> a"),
			ast.AstRoot{
				SomeExpression{
					PatternMatch, 
					expr.Select[token.Token](x.Expression, caseData...),
				},
			},
		},
		{
			[]ast.Ast{x, when, caseData_caseName},
			parser.MakeSource("test/parser/pattern-match", 
				"x when",
				"  Data a -> a",
				"  Name b -> b",
			),
			ast.AstRoot{
				SomeExpression{
					PatternMatch, 
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
			case__case_pattern_Arrow_expr_r,
			case__pattern_Arrow_expr_r,
		))

	// foo
	fooToken := makeTypeIdToken_test("Foo", 1, 1)
	fooConst := Const(fooToken)

	// bar 
	barToken := makeTypeIdToken_test("Bar", 1, 1)
	barConst := Const(barToken)

	// vars
	a := exprVar("a")
	b := exprVar("b")

	foo := SomeExpression{
		Pattern,
		expr.Apply[token.Token](fooConst, a),
	}

	bar := SomeExpression{
		Pattern,
		expr.Apply[token.Token](barConst, b),
	}

	caseFoo := CaseNode{(expr.BindersOnly[token.Token]{a}).InCase(foo.Expression, a)}

	caseFoo_caseBar := CaseNode{
		caseFoo[0],
		(expr.BindersOnly[token.Token]{b}).InCase(bar.Expression, b),
	}

	arrow := ast.TokenNode(makeToken_test(token.Arrow,"->",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{caseFoo, bar, arrow, ExpressionNode{b}},
			parser.MakeSource("test/parser/case",
				"Foo a -> a",
				"Bar b -> b",
			),
			ast.AstRoot{caseFoo_caseBar},
		},
		{
			[]ast.Ast{foo, arrow, ExpressionNode{a}},
			parser.MakeSource("test/parser/case", "Foo a -> a"),
			ast.AstRoot{caseFoo},
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
