package lexer

import "alex.peters/yew/token"

func tokensEqual(a, b token.Token) bool {
	lineA, charA := a.GetLineChar()
	lineB, charB := b.GetLineChar()
	return a.GetValue() == b.GetValue() && 
		a.GetType() == b.GetType() &&
		lineA == lineB && charA == charB
}