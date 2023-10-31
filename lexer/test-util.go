package lexer

import (
	"github.com/petersalex27/yew-lang/token"
	itoken "github.com/petersalex27/yew-packages/token"
)

func tokensEqual(a, b itoken.Token) bool {
	lineA, charA := a.GetLineChar()
	lineB, charB := b.GetLineChar()
	lenA, lenB := 0, 0
	if a2, ok := a.(token.Token); ok {
		if b2, ok := b.(token.Token); ok {
			lenA, lenB = a2.GetLength(), b2.GetLength()
		}
	}
	return a.GetValue() == b.GetValue() && 
		a.GetType() == b.GetType() &&
		lineA == lineB && charA == charB &&
		lenA == lenB
}