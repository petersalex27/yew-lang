package lexer

import (
	"testing"

	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestAnalyzeNumber(t *testing.T) {
	tests := []struct{
		source []string
		expect []itoken.Token
	}{
		{
			[]string{"1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("1").SetLineChar(1,1),
			},
		},
		{
			[]string{"0x1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0x1").SetLineChar(1,1),
			},
		},
		{
			[]string{"0xa",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0xa").SetLineChar(1,1),
			},
		},
		{
			[]string{"0Xa",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0Xa").SetLineChar(1,1),
			},
		},
		{
			[]string{"0o1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0o1").SetLineChar(1,1),
			},
		},
		{
			[]string{"0O1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0O1").SetLineChar(1,1),
			},
		},
		{
			[]string{"0b1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0b1").SetLineChar(1,1),
			},
		},
		{
			[]string{"0B1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0B1").SetLineChar(1,1),
			},
		},
		{
			[]string{"1.0",},
			[]itoken.Token{
				token.FloatValue.Make().AddValue("1.0").SetLineChar(1,1),
			},
		},
		{
			[]string{"1e1",},
			[]itoken.Token{
				token.FloatValue.Make().AddValue("1e1").SetLineChar(1,1),
			},
		},
		{
			[]string{"1E1",},
			[]itoken.Token{
				token.FloatValue.Make().AddValue("1E1").SetLineChar(1,1),
			},
		},
		{
			[]string{"1e+1",},
			[]itoken.Token{
				token.FloatValue.Make().AddValue("1e+1").SetLineChar(1,1),
			},
		},
		{
			[]string{"1e-1",},
			[]itoken.Token{
				token.FloatValue.Make().AddValue("1e-1").SetLineChar(1,1),
			},
		},
		{
			[]string{"1.0e1",},
			[]itoken.Token{
				token.FloatValue.Make().AddValue("1.0e1").SetLineChar(1,1),
			},
		},
		{
			[]string{"0_1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("1").SetLength(3).SetLineChar(1,1),
			},
		},
		{
			[]string{"00_1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("1").SetLength(4).SetLineChar(1,1),
			},
		},
		{
			[]string{"11__1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("111").SetLength(5).SetLineChar(1,1),
			},
		},
		{
			[]string{"0x1_1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0x11").SetLength(5).SetLineChar(1,1),
			},
		},
		{
			[]string{"0o1_1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0o11").SetLength(5).SetLineChar(1,1),
			},
		},
		{
			[]string{"0b1_1",},
			[]itoken.Token{
				token.IntValue.Make().AddValue("0b11").SetLength(5).SetLineChar(1,1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-num.yew")
		stat := analyzeNumber(lex)

		if es := lex.GetErrors(); len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}
		
		if stat.NotOk() {
			t.Fatalf("failed test #%d: analyzeNumber(lex).NotOk() == true\n", i+1)
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