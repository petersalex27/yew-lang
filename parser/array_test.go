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

func TestArray(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			arrayValHead__LeftBracket_expr_r,
			arrayValHead__arrayValHead_Comma_expr_r,
			arrayValHead__litArrHead_Comma_expr_r,

			array__arrayValHead_RightBracket_r,
			array__arrayValHead_Comma_RightBracket_r,
		))

	litToken := makeToken_test(token.IntValue, "1", 1, 1)
	idToken := makeIdToken_test("a", 1, 1)
	commaToken := makeToken_test(token.Comma, ",", 1, 1)
	leftToken := makeToken_test(token.LeftBracket, "[", 1, 1)
	rightToken := makeToken_test(token.RightBracket, "]", 1, 1)

	lit := LiteralNode{Const(litToken)}
	//litExpr := ExpressionNode{expr.Expression[token.Token](expr.Const[token.Token]{litToken})}
	idExpr := ExpressionNode{Const(idToken)}
	leftBracket := ast.TokenNode(leftToken)
	rightBracket := ast.TokenNode(rightToken)
	comma := ast.TokenNode(commaToken)

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{leftBracket, idExpr,},
			parser.MakeSource("test/parser/array", "[a"),
			ast.AstRoot{
				NodeSequence{ArrayValHead, []ast.Ast{idExpr}},
			},
		},
		{
			[]ast.Ast{
				NodeSequence{ArrayValHead, []ast.Ast{idExpr}}, 
				comma,
				idExpr,
			},
			parser.MakeSource("test/parser/array", "[a,a"),
			ast.AstRoot{
				NodeSequence{ArrayValHead, []ast.Ast{idExpr, idExpr}},
			},
		},
		{
			[]ast.Ast{
				NodeSequence{LitArrHead, []ast.Ast{lit,lit}}, 
				comma,
				idExpr,
			},
			parser.MakeSource("test/parser/array", "[1,1,a"),
			ast.AstRoot{
				NodeSequence{ArrayValHead, []ast.Ast{lit, lit, idExpr}},
			},
		},
		{
			[]ast.Ast{
				NodeSequence{ArrayValHead, []ast.Ast{idExpr, idExpr}}, 
				rightBracket,
			},
			parser.MakeSource("test/parser/array", "[a,a]"),
			ast.AstRoot{
				ArrayNode{false, expr.List[token.Token]{idExpr.Expression, idExpr.Expression}},
			},
		},
		{
			[]ast.Ast{
				NodeSequence{ArrayValHead, []ast.Ast{idExpr, idExpr}}, 
				comma,
				rightBracket,
			},
			parser.MakeSource("test/parser/array", "[a,a,]"),
			ast.AstRoot{
				ArrayNode{false, expr.List[token.Token]{idExpr.Expression, idExpr.Expression}},
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

func TestLiteralArray(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			litArrHead__LeftBracket_literal_r,
			litArrHead__litArrHead_Comma_literal_r,

			literalArray__litArrHead_RightBracket_r,
			literalArray__litArrHead_Comma_RightBracket_r,
		))

	litToken := makeToken_test(token.IntValue, "1", 1, 1)
	commaToken := makeToken_test(token.Comma, ",", 1, 1)
	leftToken := makeToken_test(token.LeftBracket, "[", 1, 1)
	rightToken := makeToken_test(token.RightBracket, "]", 1, 1)

	lit := LiteralNode{Const(litToken)}
	leftBracket := ast.TokenNode(leftToken)
	rightBracket := ast.TokenNode(rightToken)
	comma := ast.TokenNode(commaToken)

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{leftBracket, lit,},
			parser.MakeSource("test/parser/literal-array", "[1"),
			ast.AstRoot{
				NodeSequence{LitArrHead, []ast.Ast{lit}},
			},
		},
		{
			[]ast.Ast{
				NodeSequence{LitArrHead, []ast.Ast{lit}}, 
				comma,
				lit,
			},
			parser.MakeSource("test/parser/literal-array", "[1,1"),
			ast.AstRoot{
				NodeSequence{LitArrHead, []ast.Ast{lit, lit}},
			},
		},
		{
			[]ast.Ast{
				NodeSequence{LitArrHead, []ast.Ast{lit, lit}}, 
				rightBracket,
			},
			parser.MakeSource("test/parser/literal-array", "[1,1]"),
			ast.AstRoot{
				ArrayNode{true, expr.List[token.Token]{lit.Expression, lit.Expression}},
			},
		},
		{
			[]ast.Ast{
				NodeSequence{LitArrHead, []ast.Ast{lit, lit}}, 
				comma,
				rightBracket,
			},
			parser.MakeSource("test/parser/literal-array", "[1,1,]"),
			ast.AstRoot{
				ArrayNode{true, expr.List[token.Token]{lit.Expression, lit.Expression}},
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