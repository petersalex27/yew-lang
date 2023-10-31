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

func TestLet(t *testing.T) {
	//reInit()

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			letExpr__Let_function_In_expr_r,
		))

	letToken := token.Let.Make()
	let := ast.TokenNode(letToken)

	inToken := token.In.Make()
	in := ast.TokenNode(inToken)

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
	letF := expr.Let[token.Token](
		/* let */ Const(fToken), /* = */ 
		expr.Bind(p0).In(
			expr.Select[token.Token](p0, empty.InCase(xVar, xVar)),
		),
		/* in */ yVar,
	)

	// let g = (\$p0 $p1 -> (($p0 $p1) when (x y) -> x))
	letG := expr.Let[token.Token](
		/* let */ Const(gToken), /* = */ 
		expr.Bind(p0, p1).In(
			expr.Select[token.Token](
				expr.Apply[token.Token](p0, p1), 
				empty.InCase(expr.Apply[token.Token](xVar, yVar), xVar),
			),
		), /* in */ yVar,
	)

	letH := expr.Let[token.Token](/* let */ Const(hToken), /* = */ xVar, /* in */ yVar,)

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
			[]ast.Ast{let, h, in, yExpr},
			parser.MakeSource("test/parser/let", "let h = x in y"),
			ast.AstRoot{SomeExpression{LetExpr, letH}},
		},
		{
			[]ast.Ast{let, f, in, yExpr},
			parser.MakeSource("test/parser/let", "let f x = x in y"),
			ast.AstRoot{SomeExpression{LetExpr, letF}},
		},
		{
			[]ast.Ast{let, g, in, yExpr},
			parser.MakeSource("test/parser/let", "let g x y = x in y"),
			ast.AstRoot{SomeExpression{LetExpr, letG}},
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