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

func TestExportHead(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			exportHead__moduleDec_LeftParen_r,
			exportHead__export_Comma_r,
		))

	lparen := ast.TokenNode(token.LeftParen.Make())
	comma := ast.TokenNode(token.Comma.Make())

	moduleNameToken := token.Id.Make().AddValue("main")
	myFuncToken := makeIdToken_test("myFunc", 1, 1)

	export := ModuleNode{
		ExportList,
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

	exportHead := ModuleNode{
		ExportHead,
		moduleNameToken,
		[]token.Token{},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	exportHead2 := ModuleNode{
		ExportHead,
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
			[]ast.Ast{moduleDec, lparen},
			parser.MakeSource("test/parser/module", "module main ( )"),
			ast.AstRoot{exportHead},
		},
		{
			[]ast.Ast{export, comma},
			parser.MakeSource("test/parser/module", "module main ( myFunc )"),
			ast.AstRoot{exportHead2},
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
