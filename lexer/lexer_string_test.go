package lexer

import (
	"testing"
	"alex.peters/yew/lexer"
	itoken "alex.peters/yew/token"
	"yew.lang/main/token"
	"yew.lang/main/errors"
)

func TestAnalyzeString(t *testing.T) {
	tests := []struct{
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`""`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue(``).SetLineChar(1,1),
			},
		},
		{
			[]string{`" "`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue(` `).SetLineChar(1,1),
			},
		},
		{
			[]string{`"--"`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue(`--`).SetLineChar(1,1),
			},
		},
		{
			[]string{`"this is a string"`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue(`this is a string`).SetLineChar(1,1),
			},
		},
		{
			[]string{`"\n\t\a\b\v\f\r\"\\"`,},
			[]itoken.Token{
				token.StringValue.Make().AddValue("\n\t\a\b\v\f\r\"\\").SetLineChar(1,1),
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