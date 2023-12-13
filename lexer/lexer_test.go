package lexer

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`Maybe a = Just a | Nothing`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`Maybe`).SetLineChar(1, 1),
				token.Id.Make().AddValue(`a`).SetLineChar(1, 7),
				token.Assign.Make().AddValue(`=`).SetLineChar(1, 9),
				token.TypeId.Make().AddValue(`Just`).SetLineChar(1, 11),
				token.Id.Make().AddValue(`a`).SetLineChar(1, 16),
				token.Bar.Make().AddValue(`|`).SetLineChar(1, 18),
				token.TypeId.Make().AddValue(`Nothing`).SetLineChar(1, 20),
			},
		},
		{
			[]string{`[`},
			[]itoken.Token{
				token.LeftBracket.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`{`},
			[]itoken.Token{
				token.LeftBrace.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`,`},
			[]itoken.Token{
				token.Comma.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`}`},
			[]itoken.Token{
				token.RightBrace.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`]`},
			[]itoken.Token{
				token.RightBracket.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`)`},
			[]itoken.Token{
				token.RightParen.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`:`},
			[]itoken.Token{
				token.Typing.Make().SetLineChar(1, 1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-run.yew")
		actual, es := RunLexer(lex)

		if len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}

		if len(test.expect) != len(actual) {
			t.Fatalf("failed test #%d: expected len(actual)==%d but got len(actual)==%d\n", i+1,
				len(test.expect), len(actual))
		}

		for j, tok := range test.expect {
			if !tokensEqual(tok, actual[j]) {
				t.Fatalf("failed test #%d.%d:\nexpected:\n%v\nactual:\n%v\n", i+1, j+1,
					tok, actual[j])
			}
		}
	}
}
