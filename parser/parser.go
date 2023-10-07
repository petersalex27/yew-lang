package parser

import (
	//"github.com/petersalex27/yew-packages/parser"
	//"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

func name_maker(s string) token.Token {
	return token.Id.Make().AddValue(s)
}

//var base = types.NewContext[token.Token]().SetNameMaker(name_maker)