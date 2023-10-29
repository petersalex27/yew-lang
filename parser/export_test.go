package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	//"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	//"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestExportList(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			export__exportHead_Id_r,
			export__exportHead_TypeId_r,
			export__exportHead_Symbol_r,
			export__exportHead_Infixed_r,
			export__exportHead_r,
		))

	//rparen := ast.TokenNode(token.RightParen.Make())

	moduleNameToken := token.Id.Make().AddValue("main")
	//moduleName := ast.TokenNode(moduleNameToken)

	idToken := makeIdToken_test("id", 1, 1)
	id := ast.TokenNode(idToken)

	charToken := makeTypeIdToken_test("Char", 1, 1)
	char := ast.TokenNode(charToken)

	notToken := token.Symbol.Make().AddValue("!").SetLineChar(1,1).(token.Token)
	not := ast.TokenNode(notToken)

	bindToken := token.Infixed.Make().AddValue("(>>=)").SetLineChar(1,1).(token.Token)
	bind := ast.TokenNode(bindToken)

	exportHead := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]token.Token{},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportHead2 := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]token.Token{idToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportId := ModuleNode{
		ExportList,
		moduleNameToken,
		[]token.Token{idToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportTypeId := ModuleNode{
		ExportList,
		moduleNameToken,
		[]token.Token{charToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportSymbol := ModuleNode{
		ExportList,
		moduleNameToken,
		[]token.Token{notToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportInfixed := ModuleNode{
		ExportList,
		moduleNameToken,
		[]token.Token{bindToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportIdAndInfixed := ModuleNode{
		ExportList,
		moduleNameToken,
		[]token.Token{idToken, bindToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{exportHead, id},
			parser.MakeSource("test/parser/export", "module main ( id"),
			ast.AstRoot{exportId},
		},
		{
			[]ast.Ast{exportHead, char},
			parser.MakeSource("test/parser/export", "module main ( Char"),
			ast.AstRoot{exportTypeId},
		},
		{
			[]ast.Ast{exportHead, not},
			parser.MakeSource("test/parser/export", "module main ( !"),
			ast.AstRoot{exportSymbol},
		},
		{
			[]ast.Ast{exportHead, bind},
			parser.MakeSource("test/parser/export", "module main ( (>>=)"),
			ast.AstRoot{exportInfixed},
		},
		{
			[]ast.Ast{exportHead2},
			parser.MakeSource("test/parser/export", "module main ( id, )"),
			ast.AstRoot{exportId},
		},
		{
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
			t.Fatal(testutil.TestFail2("errors", nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(testutil.TestFail(test.expect, actual, i))
		}
	}
}
