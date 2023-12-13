package lexer

import (
	"github.com/petersalex27/yew-lang/token"
	itoken "github.com/petersalex27/yew-packages/token"
)

func tokensEqual(a, b itoken.Token) bool {
	a2, aOk := a.(token.Token)
	b2, bOk := b.(token.Token)
	if !(aOk && bOk) {
		return false
	}
	
	return a2.GetLength() == b2.GetLength() && token.TokenEquals(a2, b2)
}