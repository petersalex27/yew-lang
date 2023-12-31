package lexer

import (
	"testing"

	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestAnalyzeId(t *testing.T) {
	tests := []struct {
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`a`},
			[]itoken.Token{
				token.Id.Make().AddValue(`a`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`var`},
			[]itoken.Token{
				token.Id.Make().AddValue(`var`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`a1`},
			[]itoken.Token{
				token.Id.Make().AddValue(`a1`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`a_`},
			[]itoken.Token{
				token.Id.Make().AddValue(`a_`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`a__`},
			[]itoken.Token{
				token.Id.Make().AddValue(`a__`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`a_a`},
			[]itoken.Token{
				token.Id.Make().AddValue(`a_a`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`a'`},
			[]itoken.Token{
				token.Id.Make().AddValue(`a'`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`a''`},
			[]itoken.Token{
				token.Id.Make().AddValue(`a''`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`a'a`},
			[]itoken.Token{
				token.Id.Make().AddValue(`a'a`).SetLineChar(1, 1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-id.yew")
		stat := analyzeIdentifier(lex)

		if es := lex.GetErrors(); len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}

		if stat.NotOk() {
			t.Fatalf("failed test #%d: analyzeIdentifier(lex).NotOk() == true\n", i+1)
		}

		actuals := lex.GetTokens()

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

func TestAnalyzeTypeId(t *testing.T) {
	tests := []struct {
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`A`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`A`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`Var`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`Var`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`A1`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`A1`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`A_`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`A_`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`A__`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`A__`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`A_a`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`A_a`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`A'`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`A'`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`A''`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`A''`).SetLineChar(1, 1),
			},
		},
		{
			[]string{`A'a`},
			[]itoken.Token{
				token.TypeId.Make().AddValue(`A'a`).SetLineChar(1, 1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-id.yew")
		stat := analyzeIdentifier(lex)

		if es := lex.GetErrors(); len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}

		if stat.NotOk() {
			t.Fatalf("failed test #%d: analyzeIdentifier(lex).NotOk() == true\n", i+1)
		}

		actuals := lex.GetTokens()

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
