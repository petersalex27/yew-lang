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

func TestDefinitions(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			definitions__Indent_functions_Indent_definitions_r,
			definitions__Indent_functionDefs_Indent_definitions_r,
			definitions__exprBlock_functions_Indent_definitions_r,
			definitions__exprBlock_functionDefs_Indent_definitions_r,
			definitions__Indent_functions_r,
			definitions__Indent_functionDefs_r,
			definitions__exprBlock_functions_r,
			definitions__exprBlock_functionDefs_r,
		))

	intToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intToken)

	// Int -> Int
	intToIntType := types.Apply[token.Token](arrowConst, intType, intType)

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)

	fToken := makeIdToken_test("f", 1, 1)
	gToken := makeIdToken_test("g", 1, 1)
	hToken := makeIdToken_test("h", 1, 1)

	indentToken := token.Indent.Make().AddValue("").SetLineChar(1, 1).(token.Token)
	indent := ast.TokenNode(indentToken)
	block := ExprBlockStart(indentToken)

	fDef := FunctionDefNode{
		head:   FunctionHeadNode{fToken, []expr.Expression[token.Token]{xVar}},
		typing: intToIntType,
	}
	gDef := FunctionDefNode{
		head:   FunctionHeadNode{gToken, []expr.Expression[token.Token]{xVar}},
		typing: intToIntType,
	}
	hDef := FunctionDefNode{
		head:   FunctionHeadNode{hToken, []expr.Expression[token.Token]{xVar}},
		typing: intToIntType,
	}

	f := FunctionNode{def: fDef, body: xVar}
	g := FunctionNode{def: gDef, body: xVar}
	h := FunctionNode{def: hDef, body: xVar}

	functions_f := FunctionNodeInstances{f}
	functions_g_h := FunctionNodeInstances{g, h}
	functions_f_g_h := FunctionNodeInstances{f, g, h}

	functionDefs_f := FunctionNodeDefs{fDef}
	functionDefs_g_h := FunctionNodeDefs{gDef, hDef}
	functionDefs_f_g_h := FunctionNodeDefs{fDef, gDef, hDef}

	defs_empty := DefinitionsNode{}
	defs_fs1_ds0 := DefinitionsNode{functions_f, FunctionNodeDefs{}}
	defs_fs0_ds1 := DefinitionsNode{FunctionNodeInstances{}, functionDefs_f}
	defs_fs2_ds2 := DefinitionsNode{functions_g_h, functionDefs_g_h}
	defs_fs3_ds2 := DefinitionsNode{functions_f_g_h, functionDefs_g_h}
	defs_fs2_ds3 := DefinitionsNode{functions_g_h, functionDefs_f_g_h}

	tests := []struct {
		desc   string
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		// definitions ::= ( INDENT(n) ) functions INDENT(n) definitions
		{
			"(1) definitions ::= ( INDENT(n) ) functions INDENT(n) definitions",
			[]ast.Ast{indent, functions_f, indent, defs_empty},
			parser.MakeSource(
				"test/parser/definitions",
				"f x = x",
			),
			ast.AstRoot{indent, defs_fs1_ds0},
		},
		{
			"(2) definitions ::= ( INDENT(n) ) functions INDENT(n) definitions",
			[]ast.Ast{indent, functions_f, indent, defs_fs2_ds2},
			parser.MakeSource(
				"test/parser/definitions",
				"f x = x",
				"g: Int -> Int",
				"g x = x",
				"h: Int -> Int",
				"h x = x",
			),
			ast.AstRoot{indent, defs_fs3_ds2},
		},
		// definitions ::= ( INDENT(n) ) functionDefs INDENT(n) definitions
		{
			"(1) definitions ::= ( INDENT(n) ) functionDefs INDENT(n) definitions",
			[]ast.Ast{indent, functionDefs_f, indent, defs_empty},
			parser.MakeSource(
				"test/parser/definitions",
				"f: Int -> Int",
			),
			ast.AstRoot{indent, defs_fs0_ds1},
		},
		{
			"(2) definitions ::= ( INDENT(n) ) functionDefs INDENT(n) definitions",
			[]ast.Ast{indent, functionDefs_f, indent, defs_fs2_ds2},
			parser.MakeSource(
				"test/parser/definitions",
				"f: Int -> Int",
				"g: Int -> Int",
				"g x = x",
				"h: Int -> Int",
				"h x = x",
			),
			ast.AstRoot{indent, defs_fs2_ds3},
		},
		// definitions ::= ( indent(n) ) functions INDENT(n) definitions
		{
			"(1) definitions ::= ( indent(n) ) functions INDENT(n) definitions",
			[]ast.Ast{block, functions_f, indent, defs_empty},
			parser.MakeSource(
				"test/parser/definitions",
				"f x = x",
			),
			ast.AstRoot{block, defs_fs1_ds0},
		},
		{
			"(2) definitions ::= ( indent(n) ) functions INDENT(n) definitions",
			[]ast.Ast{block, functions_f, indent, defs_fs2_ds2},
			parser.MakeSource(
				"test/parser/definitions",
				"f x = x",
				"g: Int -> Int",
				"g x = x",
				"h: Int -> Int",
				"h x = x",
			),
			ast.AstRoot{block, defs_fs3_ds2},
		},
		// definitions ::= ( indent(n) ) functionDefs INDENT(n) definitions
		{
			"(1) definitions ::= ( indent(n) ) functionDefs INDENT(n) definitions",
			[]ast.Ast{block, functionDefs_f, indent, defs_empty},
			parser.MakeSource(
				"test/parser/definitions",
				"f: Int -> Int",
			),
			ast.AstRoot{block, defs_fs0_ds1},
		},
		{
			"(2) definitions ::= ( indent(n) ) functionDefs INDENT(n) definitions",
			[]ast.Ast{block, functionDefs_f, indent, defs_fs2_ds2},
			parser.MakeSource(
				"test/parser/definitions",
				"f: Int -> Int",
				"g: Int -> Int",
				"g x = x",
				"h: Int -> Int",
				"h x = x",
			),
			ast.AstRoot{block, defs_fs2_ds3},
		},
		// definitions ::= ( INDENT(n) ) functions
		{
			"definitions ::= ( INDENT(n) ) functions",
			[]ast.Ast{indent, functions_f},
			parser.MakeSource(
				"test/parser/definitions",
				"f x = x",
			),
			ast.AstRoot{indent, defs_fs1_ds0},
		},
		// definitions ::= ( INDENT(n) ) functionDefs
		{
			"definitions ::= ( INDENT(n) ) functionDefs",
			[]ast.Ast{indent, functionDefs_f},
			parser.MakeSource(
				"test/parser/definitions",
				"f: Int -> Int",
			),
			ast.AstRoot{indent, defs_fs0_ds1},
		},
		// definitions ::= ( indent(n) ) functions
		{
			"(1) definitions ::= ",
			[]ast.Ast{block, functions_f},
			parser.MakeSource(
				"test/parser/definitions",
				"f x = x",
			),
			ast.AstRoot{block, defs_fs1_ds0},
		},
		// definitions ::= ( indent(n) ) functionDefs
		{
			"(1) definitions ::= ( indent(n) ) functionDefs",
			[]ast.Ast{block, functionDefs_f},
			parser.MakeSource(
				"test/parser/definitions",
				"f: Int -> Int",
			),
			ast.AstRoot{block, defs_fs0_ds1},
		},
	}

	for i, test := range tests {
		p := parser.
			NewParser().
			LA(1).
			UsingReductionTable(table).
			Load([]itoken.Token{}, test.src, nil, nil).
			InitialStackPush(test.nodes...)

		actual := p.Parse()
		if p.HasErrors() {
			es := p.GetErrors()
			errors.PrintErrors(es...)
			t.Fatal(testutil.Testing("errors", test.desc).FailMessage(nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(
				testutil.
					Testing("equality", test.desc).
					FailMessage(test.expect, actual, i))
		}
	}
}
