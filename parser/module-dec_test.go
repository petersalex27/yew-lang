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

func TestModuleDec(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			moduleDec__Indent_Module_Id_r,
		))

	indent := ast.TokenNode(token.Indent.Make().AddValue(""))
	moduleNameToken := token.Id.Make().AddValue("main")
	moduleName := ast.TokenNode(moduleNameToken)

	module := ast.TokenNode(token.Module.Make())

	moduleDec := ModuleNode{
		ModuleDeclaration,
		moduleNameToken,
		[]token.Token{},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{indent, module, moduleName},
			parser.MakeSource("test/parser/module", "module main"),
			ast.AstRoot{moduleDec},
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
