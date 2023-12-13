// =============================================================================
// Author-Date: Alex Peters - December 01, 2023
// =============================================================================
package imports

import (
	"testing"

	"github.com/petersalex27/yew-lang/parser/internal"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestShiftImportName(t *testing.T) {
	nameToken := token.Id.Make().AddValue("test")
	notNameToken := token.Alias.Make()

	tests := []struct{
		desc testutil.Description
		tokens []itoken.Token
		expect bool
	}{
		{
			desc: "incorrect look-ahead",
			tokens: []itoken.Token{notNameToken},
			expect: false,
		},
		{
			desc: "correct look-ahead",
			tokens: []itoken.Token{nameToken},
			expect: true,
		},
	}

	for _, test := range tests {
		p := internal.InitInternal(parser.EmptySource{}, test.tokens)
		actual := shiftImportName(p)
		if test.expect != actual {
			t.Fatal(test.desc.FailMessage(test.expect, actual))
		}
	}
}

func TestShiftOptionalFrom(t *testing.T) {
	fromToken := token.From.Make()
	stringValueToken := token.StringValue.Make().AddValue("test")
	badToken := token.Alias.Make()

	tests := []struct{
		desc testutil.Description
		tokens []itoken.Token
		expectReduce bool
		expectOk bool
	}{
		{
			desc: "first look-ahead is not 'from'",
			tokens: []itoken.Token{badToken},
			expectReduce: false,
			expectOk: true,
		},
		{
			desc: "incorrect second look-ahead",
			tokens: []itoken.Token{fromToken, badToken},
			expectReduce: false,
			expectOk: false,
		},
		{
			desc: "not enough look-ahead tokens",
			tokens: []itoken.Token{fromToken},
			expectReduce: false,
			expectOk: false,
		},
		{
			desc: "correct look-ahead tokens",
			tokens: []itoken.Token{fromToken, stringValueToken},
			expectReduce: true,
			expectOk: true,
		},
	}

	for _, test := range tests {
		p := internal.InitInternal(parser.EmptySource{}, test.tokens)
		actualReduce, actualOk := shiftOptionalFrom(p)
		if test.expectReduce != actualReduce {
			testing := testutil.Testing("reduce", string(test.desc))
			t.Fatal(testing.FailMessage(test.expectReduce, actualReduce))
		}

		if test.expectOk != actualOk {
			testing := testutil.Testing("ok", string(test.desc))
			t.Fatal(testing.FailMessage(test.expectOk, actualOk))
		}
	}
}
