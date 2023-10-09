package parser

import (
	"yew.lang/main/token"
)

func makeToken_test(ty token.TokenType, val string, line, char int) token.Token {
	return ty.Make().AddValue(val).SetLineChar(line, char).(token.Token)
}

func makeIdToken_test(name string, line, char int) token.Token {
	return makeToken_test(token.Id, name, line, char)
}

func makeTypeIdToken_test(name string, line, char int) token.Token {
	return makeToken_test(token.TypeId, name, line, char)
}

func makeSymbolToken_test(name string, line, char int) token.Token {
	return makeToken_test(token.Symbol, name, line, char)
}

func makeInfixed_test(name string, line, char int) token.Token {
	return makeToken_test(token.Infixed, name, line, char)
}

func makeThunkedToken_test(name string, line, char int) token.Token {
	return makeToken_test(token.Thunked, name, line, char)
}