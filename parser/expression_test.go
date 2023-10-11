package parser

import (
	"testing"

	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
	"yew.lang/main/errors"
	"yew.lang/main/token"
)

func TestExpression(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			expr__val_r,
			expr__data_r,
			//expr__funcName_r,
			//expr__letIn_r,
			expr__application_r,
			//expr__pattern_r,
			//expr__exprWhere_r,
			expr__judgement_r,
			expr__enclosed_r,
		))

	typeNameToken := makeTypeIdToken_test("Name", 1, 1)
	thingToken := makeIdToken_test("thing", 1, 1)
	thingVal := expr.Var(thingToken)
	intVal := SomeExpression{Val, Const(makeToken_test(token.IntValue,"1",1,1))}
	data := SomeExpression{Data, Const(typeNameToken)}
	app := SomeExpression{Application, expr.Apply[token.Token](Const(typeNameToken), Const(typeNameToken))}
	judge := JudgementNode(types.Judgement[token.Token, expr.Expression[token.Token]](
		thingVal, 
		types.MakeConst(typeNameToken),
	))
	rparen := ast.TokenNode(makeToken_test(token.RightParen,")",1,1))
	lparen := ast.TokenNode(makeToken_test(token.LeftParen,"(",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{intVal},
			parser.MakeSource("test/parser/expression", "1"),
			ast.AstRoot{ExpressionNode{intVal.Expression}},
		},
		{
			[]ast.Ast{data},
			parser.MakeSource("test/parser/expression", "Name"),
			ast.AstRoot{ExpressionNode{data.Expression}},
		},
		/*
		{ // TODO: funcName
			[]ast.Ast{data},
			parser.MakeSource("test/parser/expression", "thing"),
			ast.AstRoot{ExpressionNode{data.Expression}},
		},
		{
			// TODO: letIn
			[]ast.Ast{data},
			parser.MakeSource("test/parser/expression", "let thing = 1 in thing"),
			ast.AstRoot{ExpressionNode{data.Expression}},
		},
		*/
		{
			[]ast.Ast{app},
			parser.MakeSource("test/parser/expression", "Name Name"),
			ast.AstRoot{ExpressionNode{app.Expression}},
		},
		/*
		{ // TODO: pattern
			[]ast.Ast{data},
			parser.MakeSource("test/parser/expression", "thing when Name -> 1"),
			ast.AstRoot{ExpressionNode{data.Expression}},
		},
		{
			// TODO: exprWhere
			[]ast.Ast{data},
			parser.MakeSource("test/parser/expression", "thing where thing = 1"),
			ast.AstRoot{ExpressionNode{data.Expression}},
		},
		*/
		{
			[]ast.Ast{judge},
			parser.MakeSource("test/parser/expression", "thing: Name"),
			ast.AstRoot{ExpressionNode{bridge.JudgementAsExpression[token.Token, expr.Expression[token.Token]](judge)}},
		},
		{
			[]ast.Ast{lparen, ExpressionNode{intVal.Expression}, rparen},
			parser.MakeSource("test/parser/expression", "( 1 )"),
			ast.AstRoot{ExpressionNode{intVal.Expression}},
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
