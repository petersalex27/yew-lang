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

func TestLiteral(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			literal__IntValue_r,
			literal__CharValue_r,
			literal__FloatValue_r,
			literal__StringValue_r,
			literal__literalArray_r,
		))

	intToken := makeToken_test(token.IntValue, "1", 1, 1)
	charToken := makeToken_test(token.CharValue, "a", 1, 1)
	floatToken := makeToken_test(token.IntValue, "1.1", 1, 1)
	stringToken := makeToken_test(token.IntValue, "hello, world!", 1, 1)

	intLit := LiteralNode{expr.Const[token.Token]{Name: intToken}}
	charLit := LiteralNode{expr.Const[token.Token]{Name: charToken}}
	floatLit := LiteralNode{expr.Const[token.Token]{Name: floatToken}}
	stringLit := LiteralNode{expr.Const[token.Token]{Name: stringToken}}
	litArray := ArrayNode{true, expr.List[token.Token]{intLit.Expression, intLit.Expression}}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{ast.TokenNode(intToken)},
			parser.MakeSource("test/parser/literal", "1"),
			ast.AstRoot{intLit},
		},
		{
			[]ast.Ast{ast.TokenNode(charToken)},
			parser.MakeSource("test/parser/literal", "'a'"),
			ast.AstRoot{charLit},
		},
		{
			[]ast.Ast{ast.TokenNode(floatToken)},
			parser.MakeSource("test/parser/literal", "1.1"),
			ast.AstRoot{floatLit},
		},
		{
			[]ast.Ast{ast.TokenNode(stringToken)},
			parser.MakeSource("test/parser/literal", `"hello, world!"`),
			ast.AstRoot{stringLit},
		},
		{
			[]ast.Ast{litArray},
			parser.MakeSource("test/parser/literal", "[1,1]"),
			ast.AstRoot{LiteralNode{litArray.List}},
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