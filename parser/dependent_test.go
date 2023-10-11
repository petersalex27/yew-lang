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
	"yew.lang/main/errors"
	"yew.lang/main/token"
)

func TestDependtype(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			dependent__dependBinders_Dot_monotype_r,
		))

	intNameToken := makeTypeIdToken_test("Int", 1, 1)
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
			[]types.TypeJudgement[token.Token, expr.Variable[token.Token]]{judge},
			types.Apply[token.Token](types.MakeConst(intNameToken)),
		),
	}
	dot := ast.TokenNode(makeToken_test(token.Dot,".",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{DependHeadNode{true, []types.TypeJudgement[token.Token, expr.Variable[token.Token]]{judge}}, dot, mono},
			parser.MakeSource("test/parser/polytype", "mapall (a: Int) . Int"),
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

func TestDependHead(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			dependHead__Mapall_varJudgement_r,
			dependHead__dependHead_varJudgement_r,
		))
	
	aType := types.MakeConst[token.Token](makeTypeIdToken_test("A",1,1))
	bType := types.MakeConst[token.Token](makeTypeIdToken_test("B",1,1))
	mapall := ast.TokenNode(makeToken_test(token.Mapall,"mapall",1,1))
	monoA := TypeNode{Monotype, aType}
	monoB := TypeNode{Monotype, bType}
	judgeA := types.TypeJudgement[token.Token, expr.Variable[token.Token]](
		types.Judgement(exprVar("a"), monoA.Type),
	)
	judgeB := types.TypeJudgement[token.Token, expr.Variable[token.Token]](
		types.Judgement(exprVar("b"), monoB.Type),
	)
	head := DependHeadNode{false, []types.TypeJudgement[token.Token, expr.Variable[token.Token]]{judgeA}}
	headAB := DependHeadNode{false, []types.TypeJudgement[token.Token, expr.Variable[token.Token]]{judgeA, judgeB}}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{mapall, VariableJudgement(judgeA)},
			parser.MakeSource("test/parser/dependent", "mapall (a: A)"),
			ast.AstRoot{head},
		},
		{
			[]ast.Ast{head, VariableJudgement(judgeB)},
			parser.MakeSource("test/parser/dependent", "mapall (a: A) (b: B)"),
			ast.AstRoot{headAB},
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

func TestDependBinders(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
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
