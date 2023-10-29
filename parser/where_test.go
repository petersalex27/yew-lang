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

func TestWhere(t *testing.T) {
	//reInit()

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			whereExpr__expr_Where_function_r,
		))

	whereToken := token.Where.Make()
	where := ast.TokenNode(whereToken)

	intToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intToken)

	// Int -> Int
	intToIntType := types.Apply[token.Token](arrowConst, intType, intType)

	// Int -> Int -> Int
	intToIntToIntType := types.Apply[token.Token](arrowConst, intType, intToIntType)

	fToken := makeIdToken_test("f", 1, 1)
	gToken := makeIdToken_test("g", 1, 1)
	hToken := makeIdToken_test("h", 1, 1)

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)
	xExpr := ExpressionNode{xVar}

	yToken := makeIdToken_test("y", 1, 1)
	yVar := expr.Var(yToken)
	yExpr := ExpressionNode{yVar}

	fDecl := FunctionHeadNode{fToken, []expr.Expression[token.Token]{xVar}}
	fDef := FunctionDefNode{fDecl, intToIntType}
	f := FunctionNode{
		def:  fDef,
		body: xExpr.Expression,
	}

	p0, p1 := exprVar("$p0"), exprVar("$p1")
	empty := expr.Bind[token.Token]()

	// let f = (\$p0 -> ($p0 when x -> x))
	whereF := expr.Where[token.Token](
		yVar,
		Const(fToken),
		expr.Bind(p0).In(
			expr.Select[token.Token](p0, empty.InCase(xVar, xVar)),
		),
	)

	// let g = (\$p0 $p1 -> (($p0 $p1) when (x y) -> x))
	whereG := expr.Where[token.Token](
		yVar,
		Const(gToken),
		expr.Bind(p0, p1).In(
			expr.Select[token.Token](
				expr.Apply[token.Token](p0, p1), 
				empty.InCase(expr.Apply[token.Token](xVar, yVar), xVar),
			),
		),
	)

	whereH := expr.Where[token.Token](yVar, Const(hToken), xVar)

	gDecl := FunctionHeadNode{gToken, []expr.Expression[token.Token]{xVar, yVar}}
	gDef := FunctionDefNode{gDecl, intToIntToIntType}
	g := FunctionNode{
		def: gDef,
		body: xExpr.Expression,
	}

	hDecl := FunctionHeadNode{hToken, []expr.Expression[token.Token]{}}
	hDef := FunctionDefNode{hDecl, intType}
	h := FunctionNode{
		def: hDef,
		body: xExpr.Expression,
	}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{yExpr, where, h},
			parser.MakeSource("test/parser/let", "y where h = x"),
			ast.AstRoot{SomeExpression{WhereExpr, whereH}},
		},
		{
			[]ast.Ast{yExpr, where, f},
			parser.MakeSource("test/parser/let", "y where f x = x"),
			ast.AstRoot{SomeExpression{WhereExpr, whereF}},
		},
		{
			[]ast.Ast{yExpr, where, g},
			parser.MakeSource("test/parser/let", "y where g x y = x"),
			ast.AstRoot{SomeExpression{WhereExpr, whereG}},
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