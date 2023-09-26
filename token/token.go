package token

import (
	"github.com/petersalex27/yew-packages/token"
)

type Token struct {
	line_char
	ty TokenType
	value string
}

func (t Token) SetLineChar(line, char int) token.Token {
	return Token{
		line_char: line_char{line, char},
		ty: t.ty,
		value: t.value,
	}
}

func (t Token) GetType() uint { return uint(t.ty) }

func (t Token) GetLineChar() (line, char int) { return t.line, t.char }

func (t Token) SetType(ty uint) token.Token {
	return TokenType(ty).
		Make().
		AddValue(t.value).
		SetLineChar(t.line, t.char)
}

func (t Token) SetValue(value string) (token.Token, error) {
	return TokenType(t.ty).
		Make().
		AddValue(value).
		SetLineChar(t.line, t.char), nil
}

func (t Token) GetValue() string {
	return t.value
}

