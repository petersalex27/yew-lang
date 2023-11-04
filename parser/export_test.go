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

func TestExportList(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			export__exportHead_Id_r,
			export__exportHead_Symbol_r,
			export__exportHead_Infixed_r,
		))

	//rparen := ast.TokenNode(token.RightParen.Make())

	moduleNameToken := token.Id.Make().AddValue("main")
	//moduleName := ast.TokenNode(moduleNameToken)

	idToken := makeIdToken_test("id", 1, 1)
	idExport := exportToken{false, idToken}
	id := ast.TokenNode(idToken)

	identToken := makeIdToken_test("ident", 1, 1)
	identExport := exportToken{false, identToken}
	ident := ast.TokenNode(identToken)

	notToken := token.Symbol.Make().AddValue("!").SetLineChar(1,1).(token.Token)
	notExport := exportToken{false, notToken}
	not := ast.TokenNode(notToken)

	bindToken := token.Infixed.Make().AddValue("(>>=)").SetLineChar(1,1).(token.Token)
	bindExport := exportToken{false, bindToken}
	bind := ast.TokenNode(bindToken)

	exportHead := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]exportToken{},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportHead2 := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]exportToken{idExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportId := ModuleNode{
		ExportList,
		moduleNameToken,
		[]exportToken{idExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportSymbol := ModuleNode{
		ExportList,
		moduleNameToken,
		[]exportToken{notExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportInfixed := ModuleNode{
		ExportList,
		moduleNameToken,
		[]exportToken{bindExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportIdAndIdent := ModuleNode{
		ExportList,
		moduleNameToken,
		[]exportToken{idExport, identExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportIdAndInfixed := ModuleNode{
		ExportList,
		moduleNameToken,
		[]exportToken{idExport, bindExport},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	tests := []struct {
		description string
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			"test `export ::= exportHead ID`",
			[]ast.Ast{exportHead, id},
			parser.MakeSource("test/parser/export", "module main ( id"),
			ast.AstRoot{exportId},
		},
		{
			"test `export ::= exportHead SYMBOL`",
			[]ast.Ast{exportHead, not},
			parser.MakeSource("test/parser/export", "module main ( !"),
			ast.AstRoot{exportSymbol},
		},
		{
			"test `export ::= exportHead INFIXED`",
			[]ast.Ast{exportHead, bind},
			parser.MakeSource("test/parser/export", "module main ( (>>=)"),
			ast.AstRoot{exportInfixed},
		},
		{
			"test export production w/ ID on module node w/ items already exported",
			[]ast.Ast{exportHead2, ident},
			parser.MakeSource("test/parser/export", "module main ( id, ident"),
			ast.AstRoot{exportIdAndIdent},
		},
		{
			"test export production w/ SYMBOL on module node w/ items already exported",
			[]ast.Ast{exportHead2, bind},
			parser.MakeSource("test/parser/export", "module main ( id, (>>=)"),
			ast.AstRoot{exportIdAndInfixed},
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
