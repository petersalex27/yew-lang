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

func TestExportHead(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			exportHead__moduleDec_LeftParen_r,
			exportHead__export_Comma_r,
			exportHead__export_TypeId_Comma_r,
			exportHead__export_TypeId_DotDot_Comma_r,
		))

	lparen := ast.TokenNode(token.LeftParen.Make())
	comma := ast.TokenNode(token.Comma.Make())
	dotDot := ast.TokenNode(token.DotDot.Make())

	moduleNameToken := token.Id.Make().AddValue("main")
	myFuncToken := makeIdToken_test("myFunc", 1, 1)
	myFuncExport := exportToken{false, myFuncToken}

	myTypeToken := makeTypeIdToken_test("MyType", 1, 1)
	// note: `false`, this does NOT export all constructors of `MyType`
	myTypeExport := exportToken{false, myTypeToken}
	// note: `true`, this DOES export all constructors of `MyType`
	myTypeAllConstructorsExport := exportToken{true, myTypeToken}
	myType := ast.TokenNode(myTypeToken)

	export := ModuleNode{
		ExportList,
		moduleNameToken,
		[]exportToken{myFuncExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	moduleDec := ModuleNode{
		ModuleDeclaration,
		moduleNameToken,
		[]exportToken{},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportHead := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]exportToken{},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportHead2 := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]exportToken{ myFuncExport },
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportHead_exportingMyType := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]exportToken{myFuncExport, myTypeExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportHead_exportingMyTypeAllConstructors := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]exportToken{myFuncExport, myTypeAllConstructorsExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	tests := []struct {
		description string
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			"test `exportHead ::= moduleHead '('`",
			[]ast.Ast{moduleDec, lparen},
			parser.MakeSource("test/parser/export-head", "module main ( "),
			ast.AstRoot{exportHead},
		},
		{
			"test `exportHead ::= export ','`",
			[]ast.Ast{export, comma},
			parser.MakeSource("test/parser/export-head", "module main ( myFunc, "),
			ast.AstRoot{exportHead2},
		},
		{
			"test abstract type export",
			[]ast.Ast{export, myType, comma},
			parser.MakeSource("test/parser/export-head", "module main ( myFucnc, MyType, "),
			ast.AstRoot{exportHead_exportingMyType},
		},
		{
			"test complete type export",
			[]ast.Ast{export, myType, dotDot, comma},
			parser.MakeSource("test/parser/export-head", "module main ( myFunc, MyType .., "),
			ast.AstRoot{exportHead_exportingMyTypeAllConstructors},
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
