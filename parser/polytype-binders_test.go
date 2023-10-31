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

func TestPolyBinders(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			polyBinders__polyHead_r,
			polyBinders__enclosed_r,
		))
	aToken := makeIdToken_test("a",1,1)
	aVar := types.Var[token.Token](aToken)
	rparen := ast.TokenNode(makeToken_test(token.RightParen,")",1,1))
	lparen := ast.TokenNode(makeToken_test(token.LeftParen,"(",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{PolyHeadNode{false, []types.Variable[token.Token]{aVar}}},
			parser.MakeSource("test/parser/type", "forall a"),
			ast.AstRoot{PolyHeadNode{true, []types.Variable[token.Token]{aVar}}},
		},
		{
			[]ast.Ast{lparen, PolyHeadNode{true, []types.Variable[token.Token]{aVar}}, rparen},
			parser.MakeSource("test/parser/type", "(forall a)"),
			ast.AstRoot{PolyHeadNode{true, []types.Variable[token.Token]{aVar}}},
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