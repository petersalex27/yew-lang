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

func TestFunctionDef(t *testing.T) {
	reInit()

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			funcDef__funcHead_Colon_type_r,
			funcDef__funcHead_r,
		))

	intToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intToken)

	// Int -> Int
	intToIntType := types.Apply[token.Token](arrowConst, intType, intType)
	intToInt := TypeNode{Type, intToIntType}

	fToken := makeIdToken_test("f", 1, 1)

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)

	yToken := makeIdToken_test("y", 1, 1)
	yVar := expr.Var(yToken)

	globalContext__.typeMutex.Lock()
	free0 := globalContext__.typeCxt.NewVar()
	free1 := globalContext__.typeCxt.NewVar()
	free2 := globalContext__.typeCxt.NewVar()
	globalContext__.typeMutex.Unlock()

	// a0 -> a1
	freeToFreeType := types.Apply[token.Token](arrowConst, free0, free1)
	// a1 -> a2
	freeToFreeType2 := types.Apply[token.Token](arrowConst, free1, free2)
	// a0 -> a1 -> a2
	freeToFreeToFreeType := types.Apply[token.Token](arrowConst, free0, freeToFreeType2)

	colonToken := makeToken_test(token.Typing, ":", 1, 1)
	colon := ast.TokenNode(colonToken)

	fHead := FunctionHeadNode{fToken, []expr.Expression[token.Token]{xVar}}

	fxyHead := FunctionHeadNode{fToken, []expr.Expression[token.Token]{xVar, yVar}}

	fNameHead := FunctionHeadNode{fToken, nil}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{fHead, colon, intToInt},
			parser.MakeSource("test/parser/function-def", "f x: Int -> Int"),
			ast.AstRoot{
				FunctionDefNode{
					head:   fHead,
					typing: intToInt.Type,
				},
			},
		},
		{
			[]ast.Ast{fNameHead},
			parser.MakeSource("test/parser/function-def", "f"),
			ast.AstRoot{
				FunctionDefNode{
					head:   fNameHead,
					typing: free0,
				},
			},
		},
		{
			[]ast.Ast{fHead},
			parser.MakeSource("test/parser/function-def", "f x"),
			ast.AstRoot{
				FunctionDefNode{
					head:   fHead,
					typing: freeToFreeType,
				},
			},
		},
		{
			[]ast.Ast{fxyHead},
			parser.MakeSource("test/parser/function-decl", "f x y"),
			ast.AstRoot{
				FunctionDefNode{
					head:   fxyHead,
					typing: freeToFreeToFreeType,
				},
			},
		},
	}

	for i, test := range tests {
		reInit()

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
