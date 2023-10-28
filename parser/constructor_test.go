package parser

import (
	//"fmt"
	"testing"

	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestConstructor(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			constructor__TypeId_r,
			constructor__typeDecl_r,
			constructor__constructor_name_r,
			constructor__constructor_constructor_r,
			constructor__enclosed_r,
		))

	typeNameToken := makeTypeIdToken_test("Name", 1, 1)
	typeName2Token := makeTypeIdToken_test("Name2", 1, 1)
	thingToken := makeIdToken_test("thing", 1, 1)

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{ast.TokenNode(typeNameToken)},
			parser.MakeSource("test/parser/constructor", "Name"),
			ast.AstRoot{Node{Constructor, typeNameToken}},
		},
		{
			[]ast.Ast{BinaryNode{TypeDecl, BinaryNode{TypeDecl, Node{Name, typeNameToken}, Node{Name, thingToken}}, Node{Name, thingToken}}},
			parser.MakeSource("test/parser/constructor", "Name thing thing"),
			ast.AstRoot{BinaryNode{Constructor, BinaryNode{Constructor, Node{Name, typeNameToken}, Node{Name, thingToken}}, Node{Name, thingToken}}},
		},
		{
			[]ast.Ast{
				Node{Constructor, typeNameToken},
				Node{Name, thingToken},
			},
			parser.MakeSource("test/parser/constructor", "Name thing"),
			ast.AstRoot{
				BinaryNode{
					Constructor,
					Node{Constructor, typeNameToken},
					Node{Name, thingToken},
				},
			},
		},
		{
			[]ast.Ast{
				Node{Constructor, typeNameToken},
				Node{Constructor, typeName2Token},
			},
			parser.MakeSource("test/parser/constructor", "Name Name2"),
			ast.AstRoot{
				BinaryNode{
					Constructor,
					Node{Constructor, typeNameToken},
					Node{Constructor, typeName2Token},
				},
			},
		},
		{
			[]ast.Ast{
				ast.TokenNode(token.LeftParen.Make()),
				Node{Constructor, typeNameToken},
				ast.TokenNode(token.RightParen.Make()),
			},
			parser.MakeSource("test/parser/constructor", "( Name )"),
			ast.AstRoot{Node{Constructor, typeNameToken}},
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
