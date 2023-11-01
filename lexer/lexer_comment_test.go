package lexer

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestSingleLineComment(t *testing.T) {
	tests := []struct {
		description testutil.Description
		source      []string
		expect      []itoken.Token
	}{
		{
			"empty",
			[]string{"--"},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLength(2).SetLineChar(1, 1),
			},
		},
		{
			"empty with trailing space",
			[]string{"-- "},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLength(3).SetLineChar(1, 1),
			},
		},
		{
			"single line comment with just non-whitespace content",
			[]string{"--comment"},
			[]itoken.Token{
				token.Comment.Make().AddValue("comment").SetLength(9).SetLineChar(1, 1),
			},
		},
		{
			"single line comment with leading whitespace and content",
			[]string{"-- comment"},
			[]itoken.Token{
				token.Comment.Make().AddValue(" comment").SetLength(10).SetLineChar(1, 1),
			},
		},
		{
			"single line comment with leading whitespace and content (2)",
			[]string{"--\tcomment"},
			[]itoken.Token{
				token.Comment.Make().AddValue("\tcomment").SetLength(10).SetLineChar(1, 1),
			},
		},
		{
			"single line comment with trailing whitespace and content",
			[]string{"--comment "},
			[]itoken.Token{
				token.Comment.Make().AddValue("comment").SetLength(10).SetLineChar(1, 1),
			},
		},
		{
			"single line comment with trailing whitespace and content (2)",
			[]string{"--comment\t"},
			[]itoken.Token{
				token.Comment.Make().AddValue("comment").SetLength(10).SetLineChar(1, 1),
			},
		},
		{
			"single line comment with whitespace and non-whitespace content",
			[]string{"--a comment"},
			[]itoken.Token{
				token.Comment.Make().AddValue("a comment").SetLength(11).SetLineChar(1, 1),
			},
		},
		{
			"single line comment with leading whitespace and non-whitespace content",
			[]string{"-- a comment"},
			[]itoken.Token{
				token.Comment.Make().AddValue(" a comment").SetLength(12).SetLineChar(1, 1),
			},
		},
		{
			"single line comment with extra dashes",
			[]string{"-------"},
			[]itoken.Token{
				token.Comment.Make().AddValue("-----").SetLength(7).SetLineChar(1, 1),
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
			t.Fatal(test.description.FailMessagef("see above errors")(i))
		}
		if stat.NotOk() {
			msg := test.description.FailMessagef("analyzeComment(lex).NotOk() == true")(i)
			t.Fatal(msg)
		}

		actuals := lex.GetTokens()

		if len(test.expect) != len(actuals) {
			t.Fatal(test.description.FailMessage(len(test.expect), len(actuals), i))
		}

		for j, tok := range test.expect {
			if !tokensEqual(tok, actuals[j]) {
				t.Fatalf(test.description.FailMessage(tok, actuals[j], i, j))
			}
		}
	}
}

func TestMultiLineComment(t *testing.T) {
	tests := []struct {
		description testutil.Description
		source      []string
		expect      []itoken.Token
	}{
		{
			"empty multi-line comment",
			[]string{"-**-"},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLength(4).SetLineChar(1, 1),
			},
		},
		{
			"multi-line comment with extra *'s",
			[]string{"-*****-"},
			[]itoken.Token{
				token.Comment.Make().AddValue("***").SetLength(7).SetLineChar(1, 1),
			},
		},
		{
			"single dash immediately enclosed by a multi-line comment",
			[]string{"-*-*-"},
			[]itoken.Token{
				token.Comment.Make().AddValue("-").SetLength(5).SetLineChar(1, 1),
			},
		},
		{
			"double dash immediately enclosed by a multi-line comment",
			[]string{"-*--*-"},
			[]itoken.Token{
				token.Comment.Make().AddValue("--").SetLength(6).SetLineChar(1, 1),
			},
		},
		{
			"open dash enclosed by a multi-line comment",
			[]string{"-*-**-"},
			[]itoken.Token{
				token.Comment.Make().AddValue("-*").SetLength(6).SetLineChar(1, 1),
			},
		},
		{
			"open dash enclosed by a multi-line comment next to immediate id",
			[]string{"-*-**-a"},
			[]itoken.Token{
				token.Comment.Make().AddValue("-*").SetLength(6).SetLineChar(1, 1),
			},
		},
		{
			"multi-line comment with whitespace and non-whitespace content",
			[]string{"-*a comment*-"},
			[]itoken.Token{
				token.Comment.Make().AddValue("a comment").SetLength(13).SetLineChar(1, 1),
			},
		},
		{
			"multi-line comment with whitespace and non-whitespace content enclosed by whitespace",
			[]string{"-* a comment *-"},
			[]itoken.Token{
				token.Comment.Make().AddValue(" a comment").SetLength(15).SetLineChar(1, 1),
			},
		},
		{
			"empty multi-line comment spanning two lines",
			[]string{
				"-*",
				"*-",
			},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLength(2).SetLineChar(1, 1),
			},
		},
		{
			"empty multi-line comment spanning three lines with only whitespace",
			[]string{
				"-* ",
				"\t\t",
				"   \t*-",
			},
			[]itoken.Token{
				token.Comment.Make().AddValue("").SetLength(3).SetLineChar(1, 1),
			},
		},
		{
			"multi-line comment spanning three line with non-whitespace content on each line",
			[]string{
				"-* this",
				"is a multi-",
				"line comment *-",
			},
			[]itoken.Token{
				token.Comment.Make().AddValue(" this is a multi- line comment").SetLength(7).SetLineChar(1, 1),
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
			t.Fatal(test.description.FailMessagef("see above errors")(i))
		}
		if stat.NotOk() {
			msg := test.description.FailMessagef("analyzeComment(lex).NotOk() == true")(i)
			t.Fatal(msg)
		}

		actuals := lex.GetTokens()

		if len(test.expect) != len(actuals) {
			t.Fatal(test.description.FailMessage(len(test.expect), len(actuals), i))
		}

		for j, tok := range test.expect {
			if !tokensEqual(tok, actuals[j]) {
				t.Fatalf(test.description.FailMessage(tok, actuals[j], i, j))
			}
		}
	}
}
