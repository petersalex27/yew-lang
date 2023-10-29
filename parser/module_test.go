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

func TestModule(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			module__moduleDec_r,
			module__export_RightParen_r,
		))

	rparen := ast.TokenNode(token.RightParen.Make())

	moduleNameToken := token.Id.Make().AddValue("main")
	//moduleName := ast.TokenNode(moduleNameToken)

	myFuncToken := makeIdToken_test("myFunc", 1, 1)

	//module := ast.TokenNode(token.Module.Make())
	module := ModuleNode{
		ModuleDefinition,
		moduleNameToken,
		[]token.Token{},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	moduleWithExport := ModuleNode{
		ModuleDefinition,
		moduleNameToken,
		[]token.Token{myFuncToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	moduleDec := ModuleNode{
		ModuleDeclaration,
		moduleNameToken,
		[]token.Token{},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	export := ModuleNode{
		ExportList,
		moduleNameToken,
		[]token.Token{ myFuncToken },
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	/*
	intToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intToken)

	// Int -> Int
	intToIntType := types.Apply[token.Token](arrowConst, intType, intType)

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)

	indent := ast.TokenNode(makeToken_test(token.Indent, "", 1, 1))

	myFuncDef := FunctionDefNode{
		head: FunctionHeadNode{myFuncToken, []expr.Expression[token.Token]{xVar}},
		typing: intToIntType,
	}

	myFunc := FunctionNode{
		def:  myFuncDef,
		body: xVar,
	}*/

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{moduleDec},
			parser.MakeSource("test/parser/module", "module main"),
			ast.AstRoot{module},
		},
		{
			[]ast.Ast{export, rparen},
			parser.MakeSource("test/parser/module", "module main ( myFunc )"),
			ast.AstRoot{moduleWithExport},
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
