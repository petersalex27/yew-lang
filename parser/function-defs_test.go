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

func TestFunctionDefs(t *testing.T) {
	reInit()

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			functionDefs__Indent_functionDef_Indent_functionDefs_r,
			functionDefs__Indent_functionDef_r,
			functionDefs__exprBlock_functionDef_Indent_functionDefs_r,
			functionDefs__exprBlock_functionDef_r,
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

	fDecl := FunctionHeadNode{fToken, []expr.Expression[token.Token]{xVar}}
	gDecl := FunctionHeadNode{gToken, []expr.Expression[token.Token]{xVar}}
	hDecl := FunctionHeadNode{hToken, []expr.Expression[token.Token]{xVar}}

	fDef := FunctionDefNode{fDecl, intToIntType}
	gDef := FunctionDefNode{gDecl, intToIntType}
	hDef := FunctionDefNode{hDecl, intToIntType}

	functionDefs_h := FunctionNodeDefs{hDef}
	functionDefs_g_h := FunctionNodeDefs{gDef, hDef}
	functionDefs_f_g_h := FunctionNodeDefs{fDef, gDef, hDef}

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
			"functionDefs ::= ( INDENT(n) ) functionDef INDENT(n) functionDefs",
			[]ast.Ast{indent, gDef, indent, functionDefs_h},
			parser.MakeSource("test/parser/function-defs", 
				"g: Int -> Int",
				"h: Int -> Int",
			),
			ast.AstRoot{indent, functionDefs_g_h},
		},
		{
			"functionDefs ::= ( INDENT(n) ) functionDef",
			[]ast.Ast{indent, hDef},
			parser.MakeSource("test/parser/function-defs",
				"h: Int -> Int",
			),
			ast.AstRoot{indent, functionDefs_h},
		},
		{
			"functionDefs ::= ( indent(n) ) functionDef",
			[]ast.Ast{indentExprBlock, hDef},
			parser.MakeSource("test/parser/function-defs",
				"h: Int -> Int",
			),
			ast.AstRoot{indentExprBlock, functionDefs_h},
		},
		{
			"functionDefs ::= ( indent(n) ) functionDef INDENT(n) functionDefs",
			[]ast.Ast{indentExprBlock, gDef, indent, functionDefs_h},
			parser.MakeSource("test/parser/function-defs", 
				"g: Int -> Int",
				"h: Int -> Int",
			),
			ast.AstRoot{indentExprBlock, functionDefs_g_h},
		},

		{
			"prepend to existing functions with > 1 elements",
			[]ast.Ast{indent, fDef, indent, functionDefs_g_h},
			parser.MakeSource("test/parser/function-defs", 
				"f: Int -> Int",
				"g: Int -> Int",
				"h: Int -> Int",
			),
			ast.AstRoot{indent, functionDefs_f_g_h},
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
