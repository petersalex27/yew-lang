package parser

import (
	"testing"

	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestFunctionDeclaration(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			functionDecl__Indent_funcName_r,
		))

	fToken := makeIdToken_test("f", 1, 1)
	fName := Node{FuncName, fToken}

	indent := ast.TokenNode(makeToken_test(token.Indent, "", 1, 1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{indent, fName},
			parser.MakeSource("test/parser/function-decl", "f"),
			ast.AstRoot{
				indent, 
				Node{FunctionDecl, fToken},
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
