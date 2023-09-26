package lexer

import (
	"testing"
	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
	"yew.lang/main/errors"
)

func TestLexer(t *testing.T) {
	tests := []struct{
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`Maybe a = Just a | Nothing`,},
			[]itoken.Token{
				token.IdType.Make().AddValue(`Maybe`).SetLineChar(1,1),
				token.IdType.Make().AddValue(`a`).SetLineChar(1,7),
				token.Assign.Make().AddValue(`=`).SetLineChar(1,9),
				token.IdType.Make().AddValue(`Maybe`).SetLineChar(1,1),
			},
		},
		{
			[]string{`[`,},
			[]itoken.Token{
				token.LeftBracket.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`{`,},
			[]itoken.Token{
				token.LeftBrace.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`;`,},
			[]itoken.Token{
				token.SemiColon.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`,`,},
			[]itoken.Token{
				token.Comma.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`}`,},
			[]itoken.Token{
				token.RightBrace.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`]`,},
			[]itoken.Token{
				token.RightBracket.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`)`,},
			[]itoken.Token{
				token.RightParen.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`@`,},
			[]itoken.Token{
				token.At.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`()`,},
			[]itoken.Token{
				token.Empty.Make().SetLineChar(1,1),
			},
		},
		{
			[]string{`:`,},
			[]itoken.Token{
				token.Typing.Make().SetLineChar(1,1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-run.yew")
		actuals, es := runLexer(lex)

		if len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}

		if len(test.expect) != len(actuals) {
			t.Fatalf("failed test #%d: expected len(actuals)==%d but got len(actuals)==%d\n", i+1,
				len(test.expect), len(actuals))
		}

		for j, tok := range test.expect {
			if !tokensEqual(tok, actuals[j]) {
				t.Fatalf("failed test #%d.%d:\nexpected:\n%v\nactual:\n%v\n", i+1, j+1,
					tok, actuals[j])
			}
		}
	}
}