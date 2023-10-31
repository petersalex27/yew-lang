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

func TestFunctionHead(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			funcHead__funcDecl_pattern_r,
			funcHead__funcDecl_r,
		))

	fToken := makeIdToken_test("f", 1, 1)
	fDecl := Node{FunctionDecl, fToken}

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)

	pattern := SomeExpression{Pattern, expr.List[token.Token]{xVar}}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{fDecl, pattern},
			parser.MakeSource("test/parser/function-head", "f x"),
			ast.AstRoot{
				FunctionHeadNode{
					fToken, 
					pattern.Expression.(expr.List[token.Token])},
				},
		},
		{
			[]ast.Ast{fDecl},
			parser.MakeSource("test/parser/function-head", "f"),
			ast.AstRoot{
				FunctionHeadNode{
					name:   fToken,
					params: nil,
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
