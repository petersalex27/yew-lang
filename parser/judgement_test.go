package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestJudgement(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			judgement__expr_Colon_type_r,
			judgement__varJudgement_r,
			judgement__enclosed_r,
		))

	intNameToken := makeTypeIdToken_test("Int", 1, 1)
	intType := TypeNode{Type, types.MakeConst[token.Token](intNameToken)}
	intExpr := ExpressionNode{Const(makeToken_test(token.IntValue, "1", 1, 1))}

	aNameToken := makeIdToken_test("a", 1, 1)
	aVar := expr.Var(aNameToken)

	varJudge := VariableJudgement(types.Judgement[token.Token, expr.Variable[token.Token]](
		aVar, intType.Type,
	))
	judge := JudgementNode(types.Judgement[token.Token, expr.Expression[token.Token]](
		intExpr.Expression, intType.Type,
	))
	judgeFromVarJudge := JudgementNode(types.Judgement[token.Token, expr.Expression[token.Token]](
		aVar, intType.Type,
	))
	colon := ast.TokenNode(makeToken_test(token.Typing, ":", 1, 1))
	rparen := ast.TokenNode(makeToken_test(token.RightParen, ")", 1, 1))
	lparen := ast.TokenNode(makeToken_test(token.LeftParen, "(", 1, 1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{intExpr, colon, intType},
			parser.MakeSource("test/parser/judgement", "1: Int"),
			ast.AstRoot{judge},
		},
		{
			[]ast.Ast{varJudge},
			parser.MakeSource("test/parser/judgement", "a: Int"),
			ast.AstRoot{judgeFromVarJudge},
		},
		{
			[]ast.Ast{lparen, judge, rparen},
			parser.MakeSource("test/parser/expression", "(1: Int)"),
			ast.AstRoot{judge},
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
