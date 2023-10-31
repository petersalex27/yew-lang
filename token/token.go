package token

import (
	"github.com/petersalex27/yew-packages/token"
)

type Token struct {
	line_char
	length int
	ty    TokenType
	value string
}

func TokenEquals(a, b Token) bool {
	return a.line == b.line && 
		a.char == b.char &&
		a.ty == b.ty &&
		a.value == b.value
}

func (t Token) GetName() string {
	return t.value
}

func (t Token) SetLineChar(line, char int) token.Token {
	return Token{
		line_char: line_char{line, char},
		length: t.length,
		ty:        t.ty,
		value:     t.value,
	}
}

func (t Token) GetLength() int {
	if t.length <= 0 {
		return len(t.value)
	}

	return t.length
}

func (t Token) SetLength(length int) Token {
	return Token{
		line_char: t.line_char,
		length: length,
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
