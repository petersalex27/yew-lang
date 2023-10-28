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

func TestType(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			type__monotype_r,
			type__polytype_r,
			type__dependent_r,
			type__enclosed_r,
		))

	intNameToken := makeTypeIdToken_test("Int", 1, 1)
	aToken := makeIdToken_test("a",1,1)
	mono := TypeNode{
		Monotype, 
		// Int
		types.MakeConst[token.Token](intNameToken),
	}
	poly := TypeNode{
		Polytype, 
		// forall a . a
		types.Forall[token.Token](types.Var[token.Token](aToken)).Bind(types.Var[token.Token](aToken)),
	}
	dep := TypeNode{
		Dependtype,
		// mapall (a: Int) . Int
		types.MakeDependentType[token.Token](
			[]types.TypeJudgement[token.Token, expr.Variable[token.Token]]{
				types.Judgement(exprVar("a"), mono.Type),
			},
			types.Apply[token.Token](types.MakeConst(intNameToken)),
		),
	}
	rparen := ast.TokenNode(makeToken_test(token.RightParen,")",1,1))
	lparen := ast.TokenNode(makeToken_test(token.LeftParen,"(",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{mono},
			parser.MakeSource("test/parser/type", "Int"),
			ast.AstRoot{TypeNode{Type, mono.Type}},
		},
		{
			[]ast.Ast{poly},
			parser.MakeSource("test/parser/type", "forall a . a"),
			ast.AstRoot{TypeNode{Type, poly.Type}},
		},
		{
			[]ast.Ast{dep},
			parser.MakeSource("test/parser/type", "mapall (a: Int) . Int"),
			ast.AstRoot{TypeNode{Type, dep.Type}},
		},
		{
			[]ast.Ast{lparen, TypeNode{Type, mono.Type}, rparen},
			parser.MakeSource("test/parser/type", "(Int)"),
			ast.AstRoot{TypeNode{Type, mono.Type}},
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
