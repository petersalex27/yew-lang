package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestFunctionInstance(t *testing.T) {
	reInit()

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			function__functionDef_Assign_exprBlock_expr_r,
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

	indentToken0 := token.Indent.Make().AddValue("")
	indent := ExprBlockStart(indentToken0)

	tests := []struct {
		desc   string
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			"function ::= functionDef '=' indent(n) expr",
			[]ast.Ast{fDef, assign, indent, xExpr},
			parser.MakeSource("test/parser/function", "f x: Int -> Int = x"),
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
			t.Fatal(testutil.Testing("errors", test.desc).FailMessage(nil, es, i))
		}

		if !actual.Equals(test.expect) {
			t.Fatal(testutil.TestFail(test.expect, actual, i))
		}
	}
}
