// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content: 
// Util functions for parsing
//
// Notes: -
// =============================================================================
package utils

import (
	"github.com/petersalex27/yew-lang/token"
	itoken "github.com/petersalex27/yew-packages/token"
)

func TokensEquals[T ~[]token.Token](xs, ys T) bool {
	if len(xs) != len(ys) {
		return false
	}

	for i, x := range xs {
		if EquateTokens(x, ys[i]) {
			return false
		}
	}
	return true
}

func EquateTokens(a, b itoken.Token) bool {
	la, ca := a.GetLineChar()
	lb, cb := b.GetLineChar()
	va, vb := a.GetValue(), b.GetValue()
	ta, tb := a.GetType(), b.GetType()
	return la == lb &&
		ca == cb &&
		va == vb &&
		ta == tb
}