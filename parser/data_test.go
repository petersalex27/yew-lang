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
			data__patternC_r,
			data__data_expr_r,
			data__enclosed_r,
		))

	nameToken := makeTypeIdToken_test("Name",1,1)
	nameConst := Const(nameToken)
	namePatternC := Node{PatternC, nameToken}

	intLitToken := makeSymbolToken_test("1",1,1)
	intExpr := Const(intLitToken)

	name1Expr := expr.Apply[token.Token](nameConst, intExpr)

	lparen := ast.TokenNode(token.LeftParen.Make())
	rparen := ast.TokenNode(token.RightParen.Make())

	tests := []struct{
		nodes []ast.Ast
		src source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{namePatternC},
			parser.MakeSource("test/parser/data", "Name"),
			ast.AstRoot{SomeExpression{Data, nameConst}},
		},
		{
			[]ast.Ast{
				SomeExpression{Data, nameConst}, 
				ExpressionNode{intExpr},
			},
			parser.MakeSource("test/parser/data", "Name 1"),
			ast.AstRoot{
				SomeExpression{Data, name1Expr},
			},
		},
		{
			[]ast.Ast{
				SomeExpression{Data, name1Expr}, 
				ExpressionNode{name1Expr},
			},
			parser.MakeSource("test/parser/data", "Name 1 (Name 1)"),
			ast.AstRoot{
				SomeExpression{
					Data, 
					expr.Apply[token.Token](name1Expr, name1Expr),
				},
			},
		},
		{
			[]ast.Ast{
				lparen,
				SomeExpression{Data, nameConst},
				rparen,
			},
			parser.MakeSource("test/parser/data", "( Name )"),
			ast.AstRoot{SomeExpression{Data, nameConst}},
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