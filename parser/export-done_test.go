package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestExportDone(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			exportDone__export_RightParen_r,
		))

	// create handle[1]
	rparen := ast.TokenNode(token.RightParen.Make())

	// create export element for handle[1]
	moduleNameToken := token.Id.Make().AddValue("main")
	myFuncToken := makeIdToken_test("myFunc", 1, 1)
	myFuncExport := exportToken{false, myFuncToken}

	// create handle[0]
	export := ModuleNode{
		ExportList,
		moduleNameToken,
		[]exportToken{myFuncExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	// expected value
	exportDone := ModuleNode{
		ExportDone,
		moduleNameToken,
		[]exportToken{myFuncExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	tests := []struct {
		description string
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			"test `exportDone ::= export ')'`",
			[]ast.Ast{export, rparen},
			parser.MakeSource("test/parser/export-done", "module main ( myFunc )"),
			ast.AstRoot{exportDone},
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