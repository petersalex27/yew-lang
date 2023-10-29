package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestSource(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			source__module_Where_definitions_r,
			source_module_r,
		))

	where := ast.TokenNode(token.Where.Make())

	moduleNameToken := token.Id.Make().AddValue("main")
	//moduleName := ast.TokenNode(moduleNameToken)

	myFuncToken := makeIdToken_test("myFunc", 1, 1)

	intToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intToken)

	// Int -> Int
	intToIntType := types.Apply[token.Token](arrowConst, intType, intType)

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)

	myFuncDef := FunctionDefNode{
		head: FunctionHeadNode{myFuncToken, []expr.Expression[token.Token]{xVar}},
		typing: intToIntType,
	}

	myFunc := FunctionNode{
		def:  myFuncDef,
		body: xVar,
	}

	defs := DefinitionsNode{[]FunctionNode{myFunc}, []FunctionDefNode{myFuncDef}}

	//module := ast.TokenNode(token.Module.Make())
	module := ModuleNode{
		ModuleDefinition,
		moduleNameToken,
		[]token.Token{myFuncToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	srcNoDefs := ModuleNode{
		Source,
		moduleNameToken,
		[]token.Token{myFuncToken},
		DefinitionsNode{[]FunctionNode{}, []FunctionDefNode{}},
	}

	src := ModuleNode{
		Source,
		moduleNameToken,
		[]token.Token{myFuncToken},
		defs,
	}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{module, where, defs},
			parser.MakeSource("test/parser/source", 
				"module main ( myFunc ) where",
				"myFunc: Int -> Int",
				"myFunc x = x",
			),
			ast.AstRoot{src},
		},
		{
			[]ast.Ast{module},
			parser.MakeSource("test/parser/source", "module main ( )"),
			ast.AstRoot{srcNoDefs},
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
