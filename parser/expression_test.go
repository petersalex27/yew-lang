package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestExpression(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			expr__val_r,
			expr__data_r,
			expr__funcName_r,
			expr__letExpr_r,
			expr__exprWhere_r,
			expr__application_r,
			expr__pattern_r,
			expr__judgement_r,
			expr__enclosed_r,
			expr__LeftParen_expr_Indent_r,
		))

	typeNameToken := makeTypeIdToken_test("Name", 1, 1)
	thingToken := makeIdToken_test("thing", 1, 1)
	aToken := makeIdToken_test("a", 1, 1)
	xVar := expr.Var(makeIdToken_test("x", 1, 1))
	aConst := expr.Const[token.Token]{Name: aToken}
	thingVal := expr.Var(thingToken)
	intVal := SomeExpression{Val, Const(makeToken_test(token.IntValue,"1",1,1))}
	nameConst := Const(typeNameToken)
	data := SomeExpression{Data, nameConst}
	app := SomeExpression{Application, expr.Apply[token.Token](Const(typeNameToken), Const(typeNameToken))}
	judge := JudgementNode(types.Judgement[token.Token, expr.Expression[token.Token]](
		thingVal, 
		types.MakeConst(typeNameToken),
	))
	rparen := ast.TokenNode(makeToken_test(token.RightParen,")",1,1))
	lparen := ast.TokenNode(makeToken_test(token.LeftParen,"(",1,1))
	indent := ast.TokenNode(makeToken_test(token.Indent, "", 1, 1))
	pattern := expr.Select[token.Token](
		thingVal, 
		(expr.BindersOnly[token.Token]{}).InCase(
			nameConst, 
			intVal.Expression,
		),
	)


	// let a = x in a
	let := SomeExpression{
		LetExpr,
		expr.Let[token.Token](
			/* let */ aConst, /* = */ xVar, /* in */ aConst,
		),
	}

	// a where a = x
	where := SomeExpression{
		WhereExpr,
		expr.Where[token.Token](
			aConst, /* where */ aConst, /* = */ xVar,
		),
	}

	tests := []struct {
		desc string
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			"expr ::= val",
			[]ast.Ast{intVal},
			parser.MakeSource("test/parser/expression", "1"),
			ast.AstRoot{ExpressionNode{intVal.Expression}},
		},
		{
			"expr ::= data",
			[]ast.Ast{data},
			parser.MakeSource("test/parser/expression", "Name"),
			ast.AstRoot{ExpressionNode{data.Expression}},
		},
		{
			"expr ::= funcName",
			[]ast.Ast{data},
			parser.MakeSource("test/parser/expression", "thing"),
			ast.AstRoot{ExpressionNode{data.Expression}},
		},
		{
			"expr ::= letExpr",
			[]ast.Ast{let},
			parser.MakeSource("test/parser/expression", "let a = x in a"),
			ast.AstRoot{ExpressionNode{let.Expression}},
		},
		{
			"expr ::= exprWhere",
			[]ast.Ast{where},
			parser.MakeSource("test/parser/expression", "a where a = x"),
			ast.AstRoot{ExpressionNode{where.Expression}},
		},
		{
			"expr ::= app",
			[]ast.Ast{app},
			parser.MakeSource("test/parser/expression", "Name Name"),
			ast.AstRoot{ExpressionNode{app.Expression}},
		},
		{
			"expr ::= pattern",
			[]ast.Ast{SomeExpression{Pattern, pattern}},
			parser.MakeSource("test/parser/expression", "match thing in Name -> 1"),
			ast.AstRoot{ExpressionNode{pattern}},
		},
		{
			"expr ::= judgement",
			[]ast.Ast{judge},
			parser.MakeSource("test/parser/expression", "thing: Name"),
			ast.AstRoot{ExpressionNode{bridge.JudgementAsExpression[token.Token, expr.Expression[token.Token]](judge)}},
		},
		{
			"expr ::= '(' expr ')'",
			[]ast.Ast{lparen, ExpressionNode{intVal.Expression}, rparen},
			parser.MakeSource("test/parser/expression", "( 1 )"),
			ast.AstRoot{ExpressionNode{intVal.Expression}},
		},
		{
			"expr ::= ( '(' ) expr INDENT(_)",
			[]ast.Ast{lparen, ExpressionNode{intVal.Expression}, indent},
			parser.MakeSource(
				"test/parser/expression", 
				"( 1",
				"",
			),
			ast.AstRoot{lparen, ExpressionNode{intVal.Expression}},
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
			t.Fatal(testutil.Testing("errors", test.desc).FailMessage(nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(testutil.Testing("equality", test.desc).FailMessage(test.expect, actual, i))
		}
	}
}
