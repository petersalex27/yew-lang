package parser

import (
	"testing"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
	"yew.lang/main/errors"
	"yew.lang/main/token"
)

func TestFunctionInstance(t *testing.T) {
	reInit()

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			function__functionDef_Assign_expr_r,
		))

	intToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intToken)

	// int -> int
	intToIntType := types.Apply[token.Token](arrowConst, intType, intType)

	fToken := makeIdToken_test("f", 1, 1)

	xToken := makeIdToken_test("x", 1, 1)
	xVar := expr.Var(xToken)
	xExpr := ExpressionNode{xVar}

	assignToken := makeToken_test(token.Assign, "=", 1, 1)
	assign := ast.TokenNode(assignToken)

	fDecl := FunctionHeadNode{fToken, []expr.Expression[token.Token]{xVar}}

	fDef := FunctionDefNode{fDecl, intToIntType}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{fDef, assign, xExpr},
			parser.MakeSource("test/parser/function-def", "f x: Int -> Int = x"),
			ast.AstRoot{
				FunctionNode{
					def:  fDef,
					body: xExpr.Expression,
				},
			},
		},
	}

	for i, test := range tests {
		reInit()

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
