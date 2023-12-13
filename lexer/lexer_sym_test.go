package lexer

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
)

func TestAnalyzeBuiltinSymbols(t *testing.T) {
	tests := []struct {
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`(`},
			[]itoken.Token{
				token.LeftParen.Make().SetLineChar(1, 1),
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
		{
			[]string{`=`},
			[]itoken.Token{
				token.Assign.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`|`},
			[]itoken.Token{
				token.Bar.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`->`},
			[]itoken.Token{
				token.Arrow.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`\`},
			[]itoken.Token{
				token.Backslash.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`..`},
			[]itoken.Token{
				token.DotDot.Make().SetLineChar(1, 1),
			},
		},
		{
			[]string{`.`},
			[]itoken.Token{
				token.Dot.Make().SetLineChar(1, 1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-sym-builtin.yew")
		stat := analyzeSymbol(lex)

		if es := lex.GetErrors(); len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}

		if stat.NotOk() {
			t.Fatalf("failed test #%d: analyzeSymbol(lex).NotOk() == true\n", i+1)
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

func TestAnalyzeSymbol(t *testing.T) {
	tests := []struct {
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`+`},
			[]itoken.Token{
				token.Id.Make().AddValue("+").SetLineChar(1, 1),
			},
		},
		{
			[]string{`+{`},
			[]itoken.Token{
				token.Id.Make().AddValue("+").SetLineChar(1, 1),
			},
		},
		{
			[]string{`+=`},
			[]itoken.Token{
				token.Id.Make().AddValue("+=").SetLineChar(1, 1),
			},
		},
		{
			[]string{`_+_`},
			[]itoken.Token{
				token.Infixed.Make().AddValue("_+_").SetLength(3).SetLineChar(1, 1),
			},
		},
		{
			[]string{`_>>=_`},
			[]itoken.Token{
				token.Infixed.Make().AddValue("_>>=_").SetLength(5).SetLineChar(1, 1),
			},
		},
		{
			[]string{`_mod_`},
			[]itoken.Token{
				token.Infixed.Make().AddValue("_mod_").SetLength(5).SetLineChar(1, 1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-sym.yew")
		stat := analyzeSymbol(lex)

		if es := lex.GetErrors(); len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}

		if stat.NotOk() {
			t.Fatalf("failed test #%d: analyzeSymbol(lex).NotOk() == true\n", i+1)
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
