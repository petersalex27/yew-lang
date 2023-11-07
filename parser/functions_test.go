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

func TestFunctionInstances(t *testing.T) {
	reInit()

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			functions__Indent_function_Indent_functions_r,
			functions__Indent_function_r,
			functions__exprBlock_function_Indent_functions_r,
			functions__exprBlock_function_r,
		))

	intToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intToken)

	// int -> int
	intToIntType := types.Apply[token.Token](arrowConst, intType, intType)

	fToken := makeIdToken_test("f", 1, 1)
	gToken := makeIdToken_test("g", 1, 1)
	hToken := makeIdToken_test("h", 1, 1)

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)
	xExpr := ExpressionNode{xVar}

	fDecl := FunctionHeadNode{fToken, []expr.Expression[token.Token]{xVar}}
	gDecl := FunctionHeadNode{gToken, []expr.Expression[token.Token]{xVar}}
	hDecl := FunctionHeadNode{hToken, []expr.Expression[token.Token]{xVar}}

	fDef := FunctionDefNode{fDecl, intToIntType}
	gDef := FunctionDefNode{gDecl, intToIntType}
	hDef := FunctionDefNode{hDecl, intToIntType}

	f := FunctionNode{def: fDef, body: xExpr.Expression}
	g := FunctionNode{def: gDef, body: xExpr.Expression}
	h := FunctionNode{def: hDef, body: xExpr.Expression}

	functions_h := FunctionNodeInstances{h}
	functions_g_h := FunctionNodeInstances{g, h}
	functions_f_g_h := FunctionNodeInstances{f, g, h}

	indentToken0 := token.Indent.Make().AddValue("")
	indent := ast.TokenNode(indentToken0)
	indentExprBlock := ExprBlockStart(indentToken0)

	tests := []struct {
		desc   string
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			"functions ::= ( INDENT(n) ) function INDENT(n) functions",
			[]ast.Ast{indent, g, indent, functions_h},
			parser.MakeSource("test/parser/functions", 
				"g x: Int -> Int = x",
				"h x: Int -> Int = x",
			),
			ast.AstRoot{indent, functions_g_h},
		},
		{
			"functions ::= ( INDENT(n) ) function",
			[]ast.Ast{indent, h},
			parser.MakeSource("test/parser/functions",
				"h x: Int -> Int = x",
			),
			ast.AstRoot{indent, functions_h},
		},
		{
			"functions ::= ( indent(n) ) function",
			[]ast.Ast{indentExprBlock, h},
			parser.MakeSource("test/parser/functions",
				"h x: Int -> Int = x",
			),
			ast.AstRoot{indentExprBlock, functions_h},
		},
		{
			"functions ::= ( indent(n) ) function INDENT(n) functions",
			[]ast.Ast{indentExprBlock, g, indent, functions_h},
			parser.MakeSource("test/parser/functions", 
				"g x: Int -> Int = x",
				"h x: Int -> Int = x",
			),
			ast.AstRoot{indentExprBlock, functions_g_h},
		},

		{
			"prepend to existing functions with > 1 elements",
			[]ast.Ast{indent, f, indent, functions_g_h},
			parser.MakeSource("test/parser/functions", 
				"f x: Int -> Int = x",
				"g x: Int -> Int = x",
				"h x: Int -> Int = x",
			),
			ast.AstRoot{indent, functions_f_g_h},
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
			t.Fatal(testutil.Testing("errors", test.desc).FailMessage(nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(testutil.Testing("equality", test.desc).FailMessage(test.expect, actual, i))
		}
	}
}
