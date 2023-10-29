package lexer

import (
	"testing"

	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
)

func TestAnalyzeComment(t *testing.T) {
	tests := []struct{
		source []string
		expect []itoken.Token
	}{
		{
			[]string{"--",},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLineChar(1,1),
			},
		},
		{
			[]string{"-- ",},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLineChar(1,1),
			},
		},
		{
			[]string{"--a comment",},
			[]itoken.Token{
				token.Comment.Make().AddValue("a comment").SetLineChar(1,1),
			},
		},
		{
			[]string{"-- a comment",},
			[]itoken.Token{
				token.Comment.Make().AddValue("a comment").SetLineChar(1,1),
			},
		},
		{
			[]string{"-- a comment \t",},
			[]itoken.Token{
				token.Comment.Make().AddValue("a comment").SetLineChar(1,1),
			},
		},
		{
			[]string{"-------",},
			[]itoken.Token{
				token.Comment.Make().AddValue("-----").SetLineChar(1,1),
			},
		},
		{
			[]string{"-**-",},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLineChar(1,1),
			},
		},
		{
			[]string{"-*****-",},
			[]itoken.Token{
				token.Comment.Make().AddValue("***").SetLineChar(1,1),
			},
		},
		{
			[]string{"-*-*-",},
			[]itoken.Token{
				token.Comment.Make().AddValue("-").SetLineChar(1,1),
			},
		},
		{
			[]string{"-*--*-",},
			[]itoken.Token{
				token.Comment.Make().AddValue("--").SetLineChar(1,1),
			},
		},
		{
			[]string{"-*-**-",},
			[]itoken.Token{
				token.Comment.Make().AddValue("-*").SetLineChar(1,1),
			},
		},
		{
			[]string{"-*a comment*-",},
			[]itoken.Token{
				token.Comment.Make().AddValue("a comment").SetLineChar(1,1),
			},
		},
		{
			[]string{"-* a comment *-",},
			[]itoken.Token{
				token.Comment.Make().AddValue("a comment").SetLineChar(1,1),
			},
		},
		{
			[]string{
				"-*",
				"*-",
			},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLineChar(1,1),
			},
		},
		{
			[]string{
				"-* ",
				"\t\t",
				"   \t*-",
			},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLineChar(1,1),
			},
		},
		{
			[]string{
				"-* this",
				"is a multi-",
				"line comment *-",
			},
			[]itoken.Token{
				token.Comment.Make().AddValue("this is a multi- line comment").SetLineChar(1,1),
			},
		},
		{
			[]string{
				"-* this\t",
				"is a multi-",
				"line comment *-",
			},
			[]itoken.Token{
				token.Comment.Make().AddValue("this is a multi- line comment").SetLineChar(1,1),
			},
		},
	}

	for i, test := range tests {
		lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
		lex.SetSource(test.source)
		lex.SetPath("./test-lex-comment.yew")
		stat := analyzeComment(lex)

		if es := lex.GetErrors(); len(es) != 0 {
			errors.PrintErrors(lex.GetErrors()...)
			t.Fatalf("failed test #%d: see above errors\n", i+1)
		}
		if stat.NotOk() {
			t.Fatalf("failed test #%d: analyzeComment(lex).NotOk() == true\n", i+1)
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