package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

func TestEmpty(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(empty__LeftParen_RightParen_r))

	lparen := ast.TokenNode(token.LeftParen.Make())
	rparen := ast.TokenNode(token.RightParen.Make())
	expected := ast.AstRoot{Node{Empty, lparen.Token.(token.Token)}}

	p := parser.NewParser().
			LA(1).
			UsingReductionTable(table).Load(
			[]itoken.Token{},
			parser.MakeSource("test/parser/empty", "()"),
			nil, nil,
		).InitialStackPush(lparen, rparen)

		actual := p.Parse()
		if p.HasErrors() {
			es := p.GetErrors()
			errors.PrintErrors(es...)
			t.Fatal(testutil.TestFail2("errors", nil, es, 0))
		}

		if !actual.Equals(expected) {
			t.Fatal(testutil.TestFail2("equality", expected, actual, 0))
		}
}