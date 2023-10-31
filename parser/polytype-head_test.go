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

func TestPolyHead(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			polyHead__Forall_Id_r,
			polyHead__polyHead_Id_r,
		))
	aToken := makeIdToken_test("a",1,1)
	bToken := makeIdToken_test("b",1,1)
	a := ast.TokenNode(aToken)
	b := ast.TokenNode(bToken)
	aVar := types.Var[token.Token](aToken)
	bVar := types.Var[token.Token](bToken)
	forall := ast.TokenNode(makeToken_test(token.Forall,"forall",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{forall, a},
			parser.MakeSource("test/parser/polytype", "forall a"),
			ast.AstRoot{PolyHeadNode{false, []types.Variable[token.Token]{aVar}}},
		},
		{
			[]ast.Ast{PolyHeadNode{false, []types.Variable[token.Token]{aVar}}, b},
			parser.MakeSource("test/parser/type", "forall a b"),
			ast.AstRoot{PolyHeadNode{false, []types.Variable[token.Token]{aVar, bVar}}},
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