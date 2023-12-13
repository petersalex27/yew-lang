package lexer

import (
	"testing"
	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-lang/errors"
)

func TestAnalyzeString(t *testing.T) {
	tests := []struct{
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`""`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue(``).SetLength(2).SetLineChar(1,1),
			},
		},
		{
			[]string{`" "`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue(` `).SetLength(3).SetLineChar(1,1),
			},
		},
		{
			[]string{`"--"`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue(`--`).SetLength(4).SetLineChar(1,1),
			},
		},
		{
			[]string{`"this is a string"`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue(`this is a string`).SetLength(18).SetLineChar(1,1),
			},
		},
		{
			[]string{`"\n\t\a\b\v\f\r\"\\"`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue("\n\t\a\b\v\f\r\"\\").SetLength(20).SetLineChar(1,1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-str.yew")
		stat := analyzeString(lex)

		if es := lex.GetErrors(); len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}
		
		if stat.NotOk() {
			t.Fatalf("failed test #%d: analyzeString(lex).NotOk() == true\n", i+1)
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