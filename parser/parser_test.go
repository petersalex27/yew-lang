package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/petersalex27/yew-lang/lexer"
	"github.com/petersalex27/yew-lang/parser/internal"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func Test(t *testing.T) {
	const nTests int = 6
	for i := 1; i <= nTests; i++ {
		path := fmt.Sprintf("./test-files/test%d.yew", i)
		lex := lexer.NewLexer(path)
		tokens, es := lexer.RunLexer(lex)
		if len(es) != 0 {
			for _, e := range es {
				fmt.Fprintf(os.Stderr, "%v\n", e.Error())
			}
			t.Fatalf("failed test #%d\n", i)
		}

		p := internal.InitInternal(lex, tokens)
		run(p)
		if len(p.Errors) != 0 {
			for _, e := range p.Errors {
				fmt.Fprintf(os.Stderr, "%v\n", e.Error())
			}
			t.Fatalf("failed test #%d\n", i)
		}
	}
}

func TestShiftX(t *testing.T) {
	const x token.TokenType = token.Alias
	tok := x.Make()

	// should return true and parser shouldn't be panicking
	{
		p := internal.InitInternal(parser.EmptySource{}, []itoken.Token{tok})
		expect := true
		actual := shiftX(p, x)
		if actual != expect {
			t.Fatal(testutil.Description("tok is not next").FailMessage(expect, actual))
		}
		if p.Panicking {
			t.Fatal(testutil.Description("panicking").FailMessage(false, p.Panicking))
		}
	}

	// should return false and parser should be panicking
	{
		p := internal.InitInternal(parser.EmptySource{}, []itoken.Token{tok})
		expect := false
		// +1 makes types not equal
		actual := shiftX(p, x+1)
		if actual != expect {
			t.Fatal(testutil.Description("tok is next").FailMessage(expect, actual))
		}
		if !p.Panicking {
			t.Fatal(testutil.Description("not panicking").FailMessage(true, p.Panicking))
		}
	}
}
