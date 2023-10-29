package lexer

import (
	"testing"
	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-lang/errors"
)

func TestAnalyzeChar(t *testing.T) {
	tests := []struct{
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`'a'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue(`a`).SetLineChar(1,1),
			},
		},
		{
			[]string{`' '`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue(` `).SetLineChar(1,1),
			},
		},
		{
			[]string{`'@'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue(`@`).SetLineChar(1,1),
			},
		},
		{
			[]string{`'\n'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\n").SetLineChar(1,1),
			},
		},
		{
			[]string{`'\t'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\t").SetLineChar(1,1),
			},
		},
		{
			[]string{`'\a'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\a").SetLineChar(1,1),
			},
		},
		{
			[]string{`'\b'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\b").SetLineChar(1,1),
			},
		},
		{
			[]string{`'\v'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\v").SetLineChar(1,1),
			},
		},
		{
			[]string{`'\f'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\f").SetLineChar(1,1),
			},
		},
		{
			[]string{`'\r'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\r").SetLineChar(1,1),
			},
		},
		{
			[]string{`'\''`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("'").SetLineChar(1,1),
			},
		},
		{
			[]string{`'\\'`,},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\\").SetLineChar(1,1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-char.yew")
		stat := analyzeChar(lex)

		if es := lex.GetErrors(); len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}
		
		if stat.NotOk() {
			t.Fatalf("failed test #%d: analyzeChar(lex).NotOk() == true\n", i+1)
		}

		actuals := lex.GetTokens()

		if len(test.expect) != len(actuals) {
			t.Fatalf("failed test #%d: expected len(actuals)==%d but got len(actuals)==%d\n", i+1,
				len(test.expect), len(actuals))
		}

		for j, tok := range test.expect {
			if !tokensEqual(tok, actuals[j]) {
				t.Fatalf("failed test #%d.%d: expected:\n%v\nactual:\n%v\n", i+1, j+1,
					tok, actuals[j])
			}
		}
	}
}