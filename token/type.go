package token

type TokenType uint

const (
	// values
	IntValue TokenType = iota
	StringValue
	CharValue
	FloatValue
	// spacing
	Indent
	// other
	Wildcard
	Comment
	Annotation
	Typing
	Assign
	Bar
	Arrow
	Backslash
	Dot
	DotDot
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
	Symbol
	Id
	TypeId
	Infixed
	/*keywords*/ _keyword_start_ // do not use!
	Let
	Match
	Trait
	Import
	Use
	Family
	Forall
	From
	In
	Mapall
	Where
	Module
	Qualified
	Of
	Derives
	Do
	/*end of keywords*/ _keyword_end_ // do not use!
	LAST_TYPE__                       // for use with ast node type
)

func (t TokenType) IsKeyword() bool {
	return t > _keyword_start_ && t < _keyword_end_
}

var builtinMap = map[TokenType]string{
	Let:          "let",
	Match:        "match",
	Trait:        "trait",
	Import:       "import",
	Use:          "use",
	Forall:       "forall",
	Mapall:       "mapall",
	Where:        "where",
	Module:       "module",
	Derives:      "derives",
	Do:           "do",
	Family:       "family",
	Qualified:    "qualified",
	From:         "from",
	In:           "in",
	Of:           "of",
	Wildcard:     "_",
	LeftParen:    "(",
	RightParen:   ")",
	LeftBracket:  "[",
	RightBracket: "]",
	LeftBrace:    "{",
	RightBrace:   "}",
	SemiColon:    ";",
	Comma:        ",",
	Typing:       ":",
	Assign:       "=",
	Bar:          "|",
	Arrow:        "->",
	Backslash:    `\`,
	Dot:          `.`,
	DotDot:       `..`,
}

func (t TokenType) Make() Token {
	value, found := builtinMap[t]
	if !found {
		value = ""
	}
	return Token{
		ty: t,
		// this may not always be correct, but can be set to something else later
		length: len(value),
		value:  value,
	}
}

// only adds value when token is not keyword, else just returns token
func (t Token) MaybeAddValue(value string) Token {
	if t.ty.IsKeyword() {
		return t
	}
	return t.AddValue(value)
}

func (t Token) AddValue(value string) Token {
	t.value = value
	if t.length <= 0 { // length set?
		// no, set it
		t.length = len(value)
	}
	return t
}
