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

func TestTupleType(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(parser.RuleSet(
			tupleType__Comma_monotype_RightParen_r,
			tupleType__Comma_monotype_tupleType_r,
		))

	intNameToken := makeTypeIdToken_test("Int", 1, 1)
	intType := types.MakeConst(intNameToken)

	mono := TypeNode{Monotype, intType}

	comma := ast.TokenNode(token.Comma.Make())
	commaType := types.MakeInfixConst[token.Token](comma.Token.(token.Token))

	tupleType := NodeSequence{
		TupleType,
		[]ast.Ast{
			TypeNode{Monotype, commaType},
			mono,
		},
	}

	tupleType2 := NodeSequence{
		TupleType,
		[]ast.Ast{
			TypeNode{Monotype, commaType},
			TypeNode{Monotype, types.Apply[token.Token](commaType, intType, intType)},
		},
	}

	rparen := ast.TokenNode(makeToken_test(token.RightParen,")",1,1))
	//lparen := ast.TokenNode(makeToken_test(token.LeftParen,"(",1,1))

	tests := []struct {
		nodes  []ast.Ast
		src    source.StaticSource
		expect ast.AstRoot
	}{
		{
			[]ast.Ast{comma, mono, rparen},
			parser.MakeSource("test/parser/tupleType", ", Int)"),
			ast.AstRoot{tupleType},
		},
		{
			[]ast.Ast{comma, mono, tupleType},
			parser.MakeSource("test/parser/type", ", Int, Int)"),
			ast.AstRoot{tupleType2},
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
