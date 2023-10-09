package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"yew.lang/main/token"
)

/*
mainly contains functions to help make code more readable by wraping type parameterized
casts inside functions
*/

func Const(tok token.Token) expr.Const[token.Token] {
	return expr.Const[token.Token]{Name: tok}
}

func Expression(e any) expr.Expression[token.Token] {
	return e.(expr.Expression[token.Token])
}