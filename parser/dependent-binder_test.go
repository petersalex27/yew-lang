package parser

import (
	"testing"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestDependBinders(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			dependBinders__dependHead_r,
			dependBinders__enclosed_r,
		))

	aType := types.MakeConst[token.Token](makeTypeIdToken_test("A",1,1))
	monoA := TypeNode{Monotype, aType}
	rparen := ast.TokenNode(makeToken_test(token.RightParen,")",1,1))
	lparen := ast.TokenNode(makeToken_test(token.LeftParen,"(",1,1))
	judgeA := types.TypeJudgement[token.Token, expr.Variable[token.Token]](
		types.Judgement(exprVar("a"), monoA.Type),
	)
	head := DependHeadNode{false, []types.TypeJudgement[token.Token, expr.Variable[token.Token]]{judgeA}}
	binder := DependHeadNode{true, []types.TypeJudgement[token.Token, expr.Variable[token.Token]]{judgeA}}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{head},
			parser.MakeSource("test/parser/dependent", "mapall (a: A)"),
			ast.AstRoot{binder},
		},
		{
			[]ast.Ast{lparen, binder, rparen},
			parser.MakeSource("test/parser/dependent", "(mapall (a: A))"),
			ast.AstRoot{binder},
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