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

func TestDependtype(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			dependent__dependBinders_Dot_monotype_r,
		))

	intNameToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intNameToken)

	mono := TypeNode{
		Monotype, 
		// Int
		types.MakeConst[token.Token](intNameToken),
	}

	judge := types.TypeJudgement[token.Token, expr.Variable[token.Token]](
		types.Judgement(exprVar("a"), mono.Type),
	)

	dep := TypeNode{
		Dependtype,
		// mapall (a: Int) . Int
		types.MakeDependentType[token.Token](
			variableJudgements(judge),
			types.Apply[token.Token](intType),
		),
	}

	dot := ast.TokenNode(makeToken_test(token.Dot,".",1,1))

	dependHead := DependHeadNode{
		readyToUse: true, 
		params: variableJudgements(judge),
	}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{dependHead, dot, mono},
			parser.MakeSource("test/parser/dependent-type", "mapall (a: Int) . Int"),
			ast.AstRoot{dep},
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