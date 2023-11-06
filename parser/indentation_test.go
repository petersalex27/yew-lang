package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestIndentation0(t *testing.T) {
	indentToken0 := token.Indent.Make().AddValue("")
	indent0 := ExprBlockStart(indentToken0)

	eq := ast.TokenNode(token.Assign.Make())
	where := ast.TokenNode(token.Where.Make())
	let := ast.TokenNode(token.Let.Make())
	in := ast.TokenNode(token.In.Make())
	of := ast.TokenNode(token.Of.Make())
	match := ast.TokenNode(token.Match.Make())

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			exprBlock__Assign_r,
			exprBlock__Where_r,
			exprBlock__Let_r,
			exprBlock__In_r,
			exprBlock__Of_r,
			exprBlock__Match_r,
		))

	tests := []struct{
		description string
		src parser.MinSource
		nodes []ast.Ast
		expect ast.AstRoot
	}{
		{
			"when '='",
			parser.MakeSource(
				"test/parser/indentation", 
				"=",
				"\t",
			),
			[]ast.Ast{eq},
			ast.AstRoot{eq, indent0},
		},
		{
			"when 'where'",
			parser.MakeSource(
				"test/parser/indentation", 
				"where",
				"\t",
			),
			[]ast.Ast{where},
			ast.AstRoot{where, indent0},
		},
		{
			"when 'let'",
			parser.MakeSource(
				"test/parser/indentation", 
				"let",
			),
			[]ast.Ast{let},
			ast.AstRoot{let, indent0},
		},
		{
			"when 'in'",
			parser.MakeSource(
				"test/parser/indentation", 
				"in",
			),
			[]ast.Ast{in},
			ast.AstRoot{in, indent0},
		},
		{
			"when 'of'",
			parser.MakeSource(
				"test/parser/indentation", 
				"of",
			),
			[]ast.Ast{of},
			ast.AstRoot{of, indent0},
		},
		{
			"when 'match'",
			parser.MakeSource(
				"test/parser/indentation", 
				"match",
			),
			[]ast.Ast{match},
			ast.AstRoot{match, indent0},
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
			t.Fatal(
				testutil.Testing("errors", test.description).
				FailMessage(nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(
				testutil.Testing("equality", test.description).
				FailMessage(test.expect, actual, i))
		}
	}
}

func TestIndentationN(t *testing.T) {
	indentToken1 := token.Indent.Make().AddValue("\t")
	indentTokenNode := ast.TokenNode(indentToken1)
	indent1 := ExprBlockStart(indentToken1)

	eq := ast.TokenNode(token.Assign.Make())
	where := ast.TokenNode(token.Where.Make())
	let := ast.TokenNode(token.Let.Make())
	in := ast.TokenNode(token.In.Make())
	of := ast.TokenNode(token.Of.Make())
	match := ast.TokenNode(token.Match.Make())

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			exprBlock__Assign_Indent_r,
			exprBlock__Where_Indent_r,
			exprBlock__Let_Indent_r,
			exprBlock__In_Indent_r,
			exprBlock__Of_Indent_r,
			exprBlock__Match_Indent_r,
		))

	tests := []struct{
		description string
		src parser.MinSource
		nodes []ast.Ast
		expect ast.AstRoot
	}{
		{
			"when '='",
			parser.MakeSource(
				"test/parser/indentation", 
				"=",
				"\t",
			),
			[]ast.Ast{eq, indentTokenNode},
			ast.AstRoot{eq, indent1},
		},
		{
			"when 'where'",
			parser.MakeSource(
				"test/parser/indentation", 
				"where",
				"\t",
			),
			[]ast.Ast{where, indentTokenNode},
			ast.AstRoot{where, indent1},
		},
		{
			"when 'let'",
			parser.MakeSource(
				"test/parser/indentation", 
				"let",
				"\t",
			),
			[]ast.Ast{let, indentTokenNode},
			ast.AstRoot{let, indent1},
		},
		{
			"when 'in'",
			parser.MakeSource(
				"test/parser/indentation", 
				"in",
				"\t",
			),
			[]ast.Ast{in, indentTokenNode},
			ast.AstRoot{in, indent1},
		},
		{
			"when 'of'",
			parser.MakeSource(
				"test/parser/indentation", 
				"of",
				"\t",
			),
			[]ast.Ast{of, indentTokenNode},
			ast.AstRoot{of, indent1},
		},
		{
			"when 'match'",
			parser.MakeSource(
				"test/parser/indentation", 
				"match",
				"\t",
			),
			[]ast.Ast{match, indentTokenNode},
			ast.AstRoot{match, indent1},
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
			t.Fatal(
				testutil.Testing("errors", test.description).
				FailMessage(nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(
				testutil.Testing("equality", test.description).
				FailMessage(test.expect, actual, i))
		}
	}
}

func TestIndentation(t *testing.T) {
	indentToken1 := token.Indent.Make().AddValue("\t")
	indentToken2 := token.Indent.Make().AddValue("\t\t")
	indentToken3 := token.Indent.Make().AddValue("\t\t\t")
	indent1 := ast.TokenNode(indentToken1)
	indent2 := ast.TokenNode(indentToken2)
	indent3 := ast.TokenNode(indentToken3)

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(_Indent__Indent_Indent_r))

	tests := []struct{
		description string
		src parser.MinSource
		nodes []ast.Ast
		expect ast.AstRoot
	}{
		{
			"base case",
			parser.MakeSource(
				"test/parser/indentation", 
				"\t\t", 
				"\t\t\t",
			),
			[]ast.Ast{indent2, indent3},
			ast.AstRoot{indent3},
		},
		{
			"inductive case",
			parser.MakeSource(
				"test/parser/indentation", 
				"\t", 
				"\t\t",
				"\t\t\t",
			),
			[]ast.Ast{indent1, indent2, indent3},
			ast.AstRoot{indent3},
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
			t.Fatal(
				testutil.Testing("errors", test.description).
				FailMessage(nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(
				testutil.Testing("equality", test.description).
				FailMessage(test.expect, actual, i))
		}
	}
}