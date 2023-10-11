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

func TestData(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			data__constructor_r,
			data__data_expr_r,
			data__enclosed_r,
		))

	nameToken := makeTypeIdToken_test("Name",1,1)
	intLitToken := makeSymbolToken_test("1",1,1)

	tests := []struct{
		nodes []ast.Ast
		src source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{Node{Constructor, nameToken}},
			parser.MakeSource("test/parser/data", "Name"),
			ast.AstRoot{SomeExpression{Data, Const(nameToken)}},
		},
		{
			[]ast.Ast{
				SomeExpression{Data, Const(nameToken)}, 
				ExpressionNode{Const(intLitToken)},
			},
			parser.MakeSource("test/parser/data", "Name 1"),
			ast.AstRoot{
				SomeExpression{Data, expr.Apply[token.Token](Const(nameToken), Const(intLitToken))},
			},
		},
		{
			[]ast.Ast{
				SomeExpression{Data, expr.Apply[token.Token](Const(nameToken), Const(intLitToken))}, 
				ExpressionNode{expr.Apply[token.Token](Const(nameToken), Const(intLitToken))},
			},
			parser.MakeSource("test/parser/data", "Name 1 (Name 1)"),
			ast.AstRoot{
				SomeExpression{Data, expr.Apply[token.Token](
					expr.Apply[token.Token](Const(nameToken), Const(intLitToken)), 
					expr.Apply[token.Token](Const(nameToken), Const(intLitToken)),
				)},
			},
		},
		{
			[]ast.Ast{
				ast.TokenNode(token.LeftParen.Make()), 
				SomeExpression{Data, Const(nameToken)},
				ast.TokenNode(token.RightParen.Make()),
			},
			parser.MakeSource("test/parser/data", "( Name )"),
			ast.AstRoot{SomeExpression{Data, Const(nameToken)}},
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