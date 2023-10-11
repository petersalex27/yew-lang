package lexer

import (
	"testing"
	"github.com/petersalex27/yew-packages/lexer"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
	"yew.lang/main/errors"
)

func TestKey(t *testing.T) {
	tests := []struct{
		source []string
		expect []itoken.Token
	}{
		{
			[]string{`class`},
			[]itoken.Token{token.Class.Make().SetLineChar(1,1),},
		},
		{
			[]string{`derives`},
			[]itoken.Token{token.Derives.Make().SetLineChar(1,1),},
		},
		{
			[]string{`do`},
			[]itoken.Token{token.Do.Make().SetLineChar(1,1),},
		},
		{
			[]string{`family`},
			[]itoken.Token{token.Family.Make().SetLineChar(1,1),},
		},
		{
			[]string{`forall`},
			[]itoken.Token{token.Forall.Make().SetLineChar(1,1),},
		},
		{
			[]string{`from`},
			[]itoken.Token{token.From.Make().SetLineChar(1,1),},
		},
		{
			[]string{`import`},
			[]itoken.Token{token.Import.Make().SetLineChar(1,1),},
		},
		{
			[]string{`in`},
			[]itoken.Token{token.In.Make().SetLineChar(1,1),},
		},
		{
			[]string{`let`},
			[]itoken.Token{token.Let.Make().SetLineChar(1,1),},
		},
		{
			[]string{`mapall`},
			[]itoken.Token{token.Mapall.Make().SetLineChar(1,1),},
		},
		{
			[]string{`module`},
			[]itoken.Token{token.Module.Make().SetLineChar(1,1),},
		},
		{
			[]string{`when`},
			[]itoken.Token{token.When.Make().SetLineChar(1,1),},
		},
		{
			[]string{`qualified`},
			[]itoken.Token{token.Qualified.Make().SetLineChar(1,1),},
		},
		{
			[]string{`struct`},
			[]itoken.Token{token.Struct.Make().SetLineChar(1,1),},
		},
		{
			[]string{`use`},
			[]itoken.Token{token.Use.Make().SetLineChar(1,1),},
		},
		{
			[]string{`where`},
			[]itoken.Token{token.Where.Make().SetLineChar(1,1),},
		},
		
	}

	// validate that all keywords are present 
	// (if keyword doesn't match, then test will fail later)
	if len(keywords) != len(tests) {
		t.Fatalf("failed pre-test: not all keywords present\n")
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