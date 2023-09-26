package lexer

import "yew.lang/main/token"

var keywords = map[string]token.TokenType{
	"class": token.Class,
	"derives": token.Derives,
	"family": token.Family,
	"forall": token.Forall,
	"from": token.From,
	"import": token.Import,
	"in": token.In,
	"let": token.Let,
	"mapall": token.Mapall,
	"module": token.Module,
	"of": token.Of,
	"qualified": token.Qualified,
	"use": token.Use,
	"where": token.Where,
}

var builtinSymbols = map[string]token.TokenType{
	"()": token.Empty,
	"(": token.LeftParen,
	")": token.RightParen,
	"{": token.LeftBrace,
	"}": token.RightBrace,
	"[": token.LeftBracket,
	"]": token.RightBracket,
	",": token.Comma,
	";": token.SemiColon,
	"@": token.At,
	":": token.Typing,
	"=": token.Assign,
	"|": token.Bar,
	`\`: token.Backslash,
	"->": token.Arrow,
	`..`: token.DotDot,
}

var keywordTrie = map[byte]map[string]token.TokenType{
	'c': {"class": token.Class},
	'd': {"derives": token.Derives},
	'f': {"family": token.Family, "forall": token.Forall, "from": token.From},
	'i': {"import": token.Import, "in": token.In},
	'l': {"let": token.Let},
	'm': {"mapall": token.Mapall, "module": token.Module},
	'o': {"of": token.Of},
	'q': {"qualified": token.Qualified},
	'u': {"use": token.Use},
	'w': {"where": token.Where},
}