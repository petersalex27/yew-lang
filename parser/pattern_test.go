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

func TestPatternC(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			patternC__constructor_r,
		))

	nameToken := makeTypeIdToken_test("Name", 1, 1)

	nameConstructor := Node{Constructor, nameToken}
	namePatternC := Node{PatternC, nameToken}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{nameConstructor},
			parser.MakeSource("test/parser/patternC", "Name"),
			ast.AstRoot{namePatternC},
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

func TestPattern(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			pattern__patternC_r,
			pattern__literal_r,
			pattern__funcName_r,
			pattern__pattern_pattern_r,
			pattern__enclosed_r,
		))

	nameToken := makeTypeIdToken_test("Name", 1, 1)
	nameConst := Const(nameToken)
	namePattern := SomeExpression{Pattern, nameConst}

	namePatternC := Node{PatternC, nameToken}

	intLitToken := makeSymbolToken_test("1", 1, 1)
	intExpr := Const(intLitToken)
	literal := LiteralNode{intExpr}

	fToken := makeIdToken_test("f", 1, 1)
	f := Node{FuncName, fToken}
	fConst := Const(fToken)

	lparen := ast.TokenNode(token.LeftParen.Make())
	rparen := ast.TokenNode(token.RightParen.Make())

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{namePatternC},
			parser.MakeSource("test/parser/pattern", "Name"),
			ast.AstRoot{namePattern},
		},
		{
			[]ast.Ast{literal},
			parser.MakeSource("test/parser/pattern", "1"),
			ast.AstRoot{SomeExpression{Pattern, intExpr}},
		},
		{
			[]ast.Ast{f},
			parser.MakeSource("test/parser/pattern", "f"),
			ast.AstRoot{SomeExpression{Pattern, fConst}},
		},
		{
			[]ast.Ast{namePattern, namePattern},
			parser.MakeSource("test/parser/pattern", "(Name) Name"),
			ast.AstRoot{
				SomeExpression{
					Pattern,
					expr.Apply[token.Token](nameConst, nameConst),
				},
			},
		},
		{
			[]ast.Ast{lparen, namePattern, rparen},
			parser.MakeSource("test/parser/data", "( Name )"),
			ast.AstRoot{namePattern},
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
