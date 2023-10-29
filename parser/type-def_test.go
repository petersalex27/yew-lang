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

func TestTypeDef(t *testing.T) {
	reInit() // reset global context

	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			typeDef__typeDecl_Assign_constructor_r,
			typeDef__typeDef_Bar_constructor_r,
		))

	aToken := makeToken_test(token.Id, "a", 1, 1)
	typeToken := makeTypeIdToken_test("Type", 1, 1)
	conToken := makeTypeIdToken_test("Con", 1, 1)
	konToken := makeTypeIdToken_test("Kon", 1, 1)

	// Type a
	decl := BinaryNode{TypeDecl, Node{Name, typeToken}, Node{Name, aToken}}
	// Con a
	con := BinaryNode{Constructor, Node{Name, conToken}, Node{Name, aToken}}
	kon := Node{Constructor, konToken}
	// Con a (Type a)
	con2 := BinaryNode{Constructor,
		// Con a
		BinaryNode{Constructor,
			Node{Name, conToken},
			Node{Name, aToken},
		},
		// Type a
		BinaryNode{Constructor, Node{Name, typeToken}, Node{Name, aToken}},
	}

	assign := ast.TokenNode(makeToken_test(token.Assign, "=", 1, 1))
	bar := ast.TokenNode(makeToken_test(token.Bar, "|", 1, 1))

	constType := types.ReferableType[token.Token](types.MakeConst(typeToken))
	//conType := types.MakeConst[token.Token]("Con")
	//konType := types.MakeConst[token.Token]("Kon")
	aVar := types.Var(aToken)
	typeA := types.Apply(constType, types.Monotyped[token.Token](aVar))
	closedType := types.Forall(aVar).Bind(typeA)
	newVar := globalContext__.exprCxt.NewVar()
	newVar2 := globalContext__.exprCxt.NewVar()

	// (\$0 -> Con $0): forall a . a -> Type a
	conJudge := types.Judgement[token.Token, expr.Expression[token.Token]](
		// (\$0 -> Con $0)
		expr.Bind[token.Token](newVar).In(expr.Apply[token.Token](Const(conToken), newVar)),
		// forall a . a -> Type a
		types.Forall(aVar).Bind(
			types.Apply(
				types.ReferableType[token.Token](arrowConst),
				types.Monotyped[token.Token](aVar),
				types.Monotyped[token.Token](typeA),
			),
		),
	)

	// Kon: forall a . Type a
	konJudge := types.Judgement[token.Token, expr.Expression[token.Token]](
		Const(konToken),
		types.Forall(aVar).Bind(typeA),
	)

	// (\$0 $1 -> Con $0 $1): forall a . a -> Type a -> Type a
	con2Judge := types.Judgement[token.Token, expr.Expression[token.Token]](
		// (\$0 $1 -> Con $0 $1)
		expr.Bind[token.Token](newVar, newVar2).In(expr.Apply[token.Token](Const(conToken), newVar, newVar2)),
		// forall a . a -> Type a -> Type a
		types.Forall(aVar).Bind(
			types.Apply( // a -> Type a -> Type a
				types.ReferableType[token.Token](arrowConst),
				types.Monotyped[token.Token](aVar),
				types.Monotyped[token.Token](
					types.Apply( // Type a -> Type a
						types.ReferableType[token.Token](arrowConst),
						types.Monotyped[token.Token](typeA),
						types.Monotyped[token.Token](typeA),
					),
				),
			),
		),
	)

	test1Expect := TypeDefNode{
		constType.(types.Constant[token.Token]),
		closedType,
		[]types.TypeJudgement[token.Token, expr.Expression[token.Token]]{
			conJudge,
		},
	}

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{decl, assign, con},
			parser.MakeSource("test/parser/literal", "Type a = Con a"),
			ast.AstRoot{test1Expect},
		},
		{
			[]ast.Ast{decl, assign, con2},
			parser.MakeSource("test/parser/literal", "Type a = Con a (Type a)"),
			ast.AstRoot{
				TypeDefNode{
					constType.(types.Constant[token.Token]),
					closedType,
					[]types.TypeJudgement[token.Token, expr.Expression[token.Token]]{
						con2Judge,
					},
				},
			},
		},
		{
			[]ast.Ast{test1Expect, bar, kon},
			parser.MakeSource("test/parser/literal", "Type a = Con a | Kon"),
			ast.AstRoot{
				TypeDefNode{
					constType.(types.Constant[token.Token]),
					closedType,
					[]types.TypeJudgement[token.Token, expr.Expression[token.Token]]{
						conJudge,
						konJudge,
					},
				},
			},
		},
	}

	for i, test := range tests {
		// reset global context
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
