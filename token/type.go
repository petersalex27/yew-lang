package token

type TokenType uint

const (
	// values
	IntValue TokenType = iota
	StringValue
	CharValue
	FloatValue
	// spacing
	IndentType
	// other
	Wildcard
	Empty
	Comment
	// grouping
	LeftParen
	RightParen
	LeftBracket
	RightBracket
	LeftBrace
	RightBrace
	// separators
	SemiColon
	Comma
	// names
	SymbolType
	IdType
	Infixed
	// keywords
	Let
	Of
	Class
	Import
	Use
	Forall
	Mapall
	Where
	Module
	Hide
	Derives
)

func (t TokenType) IsKeyword() bool {
	return t >= Let && t <= Derives
}

var builtinMap = map[TokenType]string{
	Let:          "let",
	Of:           "of",
	Class:        "class",
	Import:       "import",
	Use:          "use",
	Forall:       "forall",
	Mapall:       "mapall",
	Where:        "where",
	Module:       "module",
	Hide:         "hide",
	Derives:      "derives",
	Wildcard:     "_",
	LeftParen:    "(",
	RightParen:   ")",
	LeftBracket:  "[",
	RightBracket: "]",
	LeftBrace:    "{",
	RightBrace:   "}",
	SemiColon:    ";",
	Comma:        ",",
	Empty:	      "()",
}

func (t TokenType) Make() Token {
	value, found := builtinMap[t]
	if !found {
		value = ""
	}
	return Token{
		ty:    byte(t),
		value: value,
	}
}

func (t Token) AddValue(value string) Token {
	t.value = value
	return t
}
