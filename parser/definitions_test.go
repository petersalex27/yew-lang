package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestDefinitions(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			definitions__Indent_function_r,
			definitions__Indent_funcDef_r,
			definitions__Indent_function_definitions_r,
			definitions__Indent_funcDef_definitions_r,
		))

	intToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intToken)

	// Int -> Int
	intToIntType := types.Apply[token.Token](arrowConst, intType, intType)

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)

	fToken := makeIdToken_test("f", 1, 1)

	indent := ast.TokenNode(makeToken_test(token.Indent, "", 1, 1))

	fDef := FunctionDefNode{
		head: FunctionHeadNode{fToken, []expr.Expression[token.Token]{xVar}},
		typing: intToIntType,
	}

	f := FunctionNode{
		def:  fDef,
		body: xVar,
	}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{indent, f},
			parser.MakeSource("test/parser/definitions", "f x = x"),
			ast.AstRoot{
				DefinitionsNode{
					[]FunctionNode{f},
					[]FunctionDefNode{},
				},
			},
		},
		{
			[]ast.Ast{indent, fDef},
			parser.MakeSource("test/parser/definitions", "f: Int -> Int"),
			ast.AstRoot{
				DefinitionsNode{
					[]FunctionNode{},
					[]FunctionDefNode{fDef},
				},
			},
		},
		{
			[]ast.Ast{
				indent,
				f,
				DefinitionsNode{
					[]FunctionNode{},
					[]FunctionDefNode{fDef},
				},
			},
			parser.MakeSource("test/parser/definitions",
				"f x = x",
				"f: Int -> Int",
			),
			ast.AstRoot{
				DefinitionsNode{
					[]FunctionNode{f},
					[]FunctionDefNode{fDef},
				},
			},
		},
		{
			[]ast.Ast{
				indent,
				fDef,
				DefinitionsNode{
					[]FunctionNode{f},
					[]FunctionDefNode{},
				},
			},
			parser.MakeSource("test/parser/definitions", 
				"f: Int -> Int",
				"f x = x",
			),
			ast.AstRoot{
				DefinitionsNode{
					[]FunctionNode{f},
					[]FunctionDefNode{fDef},
				},
			},
		},
		{
			[]ast.Ast{
				indent,
				f,
				DefinitionsNode{
					[]FunctionNode{f},
					[]FunctionDefNode{},
				},
			},
			parser.MakeSource("test/parser/definitions",
				"f x = x",
				"f x = x",
			),
			ast.AstRoot{
				DefinitionsNode{
					[]FunctionNode{f, f},
					[]FunctionDefNode{},
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
