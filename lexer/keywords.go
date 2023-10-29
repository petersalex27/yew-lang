package lexer

import "github.com/petersalex27/yew-lang/token"

var keywords = map[string]token.TokenType{
	"derives":   token.Derives,
	"do":        token.Do,
	"family":    token.Family,
	"forall":    token.Forall,
	"from":      token.From,
	"import":    token.Import,
	"in":        token.In,
	"let":       token.Let,
	"mapall":    token.Mapall,
	"module":    token.Module,
	"of":        token.Of,
	"trait":     token.Trait,
	"qualified": token.Qualified,
	"use":       token.Use,
	"when":      token.When,
	"where":     token.Where,
}

var builtinSymbols = map[string]token.TokenType{
	"(":  token.LeftParen,
	")":  token.RightParen,
	"{":  token.LeftBrace,
	"}":  token.RightBrace,
	"[":  token.LeftBracket,
	"]":  token.RightBracket,
	",":  token.Comma,
	";":  token.SemiColon,
	":":  token.Typing,
	"=":  token.Assign,
	"|":  token.Bar,
	`\`:  token.Backslash,
	"->": token.Arrow,
	`..`: token.DotDot,
	`.`:  token.Dot,
}

var keywordTrie = map[byte]map[string]token.TokenType{
	'd': {"derives": token.Derives, "do": token.Do},
	'f': {"family": token.Family, "forall": token.Forall, "from": token.From},
	'i': {"import": token.Import, "in": token.In},
	'l': {"let": token.Let},
	'm': {"mapall": token.Mapall, "module": token.Module},
	'o': {"of": token.Of},
	'q': {"qualified": token.Qualified},
	't': {"trait": token.Trait},
	'u': {"use": token.Use},
	'w': {"when": token.When, "where": token.Where},
}
