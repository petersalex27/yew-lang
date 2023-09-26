package lexer

import "github.com/petersalex27/yew-packages/token"

func tokensEqual(a, b token.Token) bool {
	lineA, charA := a.GetLineChar()
	lineB, charB := b.GetLineChar()
	return a.GetValue() == b.GetValue() && 
		a.GetType() == b.GetType() &&
		lineA == lineB && charA == charB
}