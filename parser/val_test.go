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

func TestVal(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			val__literal_r,
			val__array_r,
		))

	intToken := makeToken_test(token.IntValue, "1", 1, 1)

	lit := LiteralNode{expr.Const[token.Token]{Name: intToken}}
	array := ArrayNode{false, expr.List[token.Token]{lit.Expression, lit.Expression}}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{lit},
			parser.MakeSource("test/parser/val", "1"),
			ast.AstRoot{SomeExpression{Val, lit.getExpression().Expression}},
		},
		{
			[]ast.Ast{array},
			parser.MakeSource("test/parser/val", "[1,1]"),
			ast.AstRoot{SomeExpression{Val, array.getExpression().Expression}},
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