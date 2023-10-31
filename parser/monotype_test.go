package parser


import (
	"testing"

	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestMonotype(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.Order(
			monotype__monotype_monotype_r,
			monotype__monotype_Arrow_monotype_r,
			monotype__enclosed_r,
			monotype__TypeId_r,
			monotype__Id_r,
			monotype__LeftParen_monotype_tupleType_r,
		))


	comma := ast.TokenNode(token.Comma.Make())
	commaType := types.MakeInfixConst[token.Token](comma.Token.(token.Token))

	arrow := ast.TokenNode(token.Arrow.Make())
	arrowType := types.MakeInfixConst[token.Token](arrow.Token.(token.Token))
		
	// a
	aToken := makeIdToken_test("a", 1, 1)
	aTokenNode := ast.TokenNode(aToken)
	aType := types.Var(aToken)
	aTypeVar := TypeNode{Monotype, aType}

	// Int
	intNameToken := makeTypeIdToken_test("Int", 1, 1)
	intTokenNode := ast.TokenNode(intNameToken)
	intType := types.MakeConst(intNameToken)
	int_ := TypeNode{Monotype, intType}
	// Int -> Int
	intToIntType := types.Apply[token.Token](arrowType, intType, intType)
	intToInt := TypeNode{Monotype, intToIntType}
	// (Int Int)
	intIntType := types.Apply[token.Token](intType, intType)
	intInt := TypeNode{Monotype, intIntType}

	// .., Int, Int)
	tupleTailWorkingType := types.Apply[token.Token](commaType, intType, intType)
	tupleTail := NodeSequence{
		TupleType,
		[]ast.Ast{
			TypeNode{Monotype, commaType},
			TypeNode{Monotype, tupleTailWorkingType},
		},
	}

	// (Int, Int, Int)
	tupleType := types.Apply[token.Token](commaType, intType, tupleTailWorkingType)
	tuple := TypeNode{Monotype, tupleType}

	rparen := ast.TokenNode(makeToken_test(token.RightParen,")",1,1))
	lparen := ast.TokenNode(makeToken_test(token.LeftParen,"(",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{int_, int_},
			parser.MakeSource("test/parser/monotype", "Int Int"),
			ast.AstRoot{intInt},
		},
		{
			[]ast.Ast{int_, arrow, int_},
			parser.MakeSource("test/parser/monotype", "Int -> Int"),
			ast.AstRoot{intToInt},
		},
		{
			[]ast.Ast{lparen, intToInt, rparen},
			parser.MakeSource("test/parser/monotype", "(Int -> Int)"),
			ast.AstRoot{intToInt},
		},
		{
			[]ast.Ast{intTokenNode},
			parser.MakeSource("test/parser/monotype", "Int"),
			ast.AstRoot{int_},
		},
		{
			[]ast.Ast{aTokenNode},
			parser.MakeSource("test/parser/monotype", "a"),
			ast.AstRoot{aTypeVar},
		},
		{
			[]ast.Ast{lparen, int_, tupleTail},
			parser.MakeSource("test/parser/monotype", "(Int, Int, Int)"),
			ast.AstRoot{tuple},
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
