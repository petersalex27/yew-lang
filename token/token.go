package token

import (
	"alex.peters/yew/token"
)

type Token struct {
	line_char
	ty byte
	value string
}

func (t Token) SetLineChar(line, char int) token.Token {
	return Token{
		line_char: line_char{line, char},
		value: t.value,
	}
}

func (t Token) GetType() uint { return uint(t.ty) }

func (t Token) GetLineChar() (line, char int) { return t.line, t.char }

func (t Token) SetType(ty uint) token.Token {
	t.ty = byte(ty)
	return t
}

func (t Token) SetValue(value string) (token.Token, error) {
	t.value = value
	return t, nil
}

func (t Token) GetValue() string {
	return t.value
}

