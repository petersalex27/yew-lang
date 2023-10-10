package parser

import (
	"testing"

	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
	"yew.lang/main/errors"
	"yew.lang/main/token"
)

func TestTypeDecl(t *testing.T) {
	reInit() // reset global context

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			typeDecl__TypeId_r,
			typeDecl__typeDecl_Id_r,
		))

	aToken := makeToken_test(token.Id, "a", 1, 1)
	typeToken := makeTypeIdToken_test("Type", 1, 1)

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{ast.TokenNode(typeToken)},
			parser.MakeSource("test/parser/literal", "Type"),
			ast.AstRoot{Node{TypeDecl, typeToken}},
		},
		{
			[]ast.Ast{Node{TypeDecl, typeToken}, ast.TokenNode(aToken)},
			parser.MakeSource("test/parser/literal", "Type a"),
			ast.AstRoot{BinaryNode{TypeDecl, Node{TypeDecl, typeToken}, Node{Name, aToken}}},
		},
	}

	for i, test := range tests {
		// reset global context 
		reInit()

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