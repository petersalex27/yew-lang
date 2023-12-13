package lexer

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
)

func TestAnalyzeChar(t *testing.T) {
	tests := []struct {
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`'a'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue(`a`).SetLength(3).SetLineChar(1, 1),
			},
		},
		{
			[]string{`' '`},
			[]itoken.Token{
				token.CharValue.Make().AddValue(` `).SetLength(3).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'@'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue(`@`).SetLength(3).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\n'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\n").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\t'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\t").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\a'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\a").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\b'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\b").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\v'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\v").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\f'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\f").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\r'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\r").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\''`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("'").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			[]string{`'\\'`},
			[]itoken.Token{
				token.CharValue.Make().AddValue("\\").SetLength(4).SetLineChar(1, 1),
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

		actual := lex.GetTokens()

		if len(test.expect) != len(actual) {
			t.Fatalf("failed test #%d: expected len(actual)==%d but got len(actual)==%d\n", i+1,
				len(test.expect), len(actual))
		}

		for j, tok := range test.expect {
			if !tokensEqual(tok, actual[j]) {
				t.Fatalf("failed test #%d.%d: expected:\n%v\nactual:\n%v\n", i+1, j+1,
					tok, actual[j])
			}
		}
	}
}
