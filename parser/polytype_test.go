package parser

import (
	"testing"
	
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestPolytype(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			polytype__polyBinders_Dot_dependTyped_r,
		))
	aToken := makeIdToken_test("a",1,1)
	aVar := types.Var[token.Token](aToken)
	poly := TypeNode{
		Polytype, 
		// forall a . a
		types.Forall[token.Token](aVar).Bind(aVar),
	}
	dot := ast.TokenNode(makeToken_test(token.Dot,".",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{PolyHeadNode{true, []types.Variable[token.Token]{aVar}}, dot, TypeNode{Dependtyped, aVar}},
			parser.MakeSource("test/parser/polytype", "forall a . a"),
			ast.AstRoot{poly},
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