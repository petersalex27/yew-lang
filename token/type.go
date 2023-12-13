package token

type TokenType uint

const (
	// values
	IntValue TokenType = iota
	StringValue
	CharValue
	FloatValue
	// spacing
	Newline
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
	Comma
	// names
	Id
	TypeId
	Infixed
	// alphabetic keywords
	Let
	Match
	Trait
	Import
	Use
	Family
	Forall
	From
	In
	Mapval
	Where
	Module
	Qualified
	Of
	Derives
	Alias
	LAST_TYPE__ // for use with ast node type
)

func (t TokenType) IsConstant() bool {
	_, found := constantTokenMap[t]
	return found
}

// map of constant tokens.
//
// Let, for example, is a constant token--there is only one string that can be tokenized for it.
//
// On the other hand, IntValue is not a constant token--there are many strings that can be tokenized
// for it.
var constantTokenMap = map[TokenType]string{
	Let:          "let",
	Match:        "match",
	Trait:        "trait",
	Import:       "import",
	Use:          "use",
	Forall:       "forall",
	Mapval:       "mapval",
	Where:        "where",
	Module:       "module",
	Derives:      "derives",
	Alias:        "alias",
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
	Comma:        ",",
	Typing:       ":",
	Assign:       "=",
	Bar:          "|",
	Arrow:        "->",
	Backslash:    `\`,
	Dot:          `.`,
	DotDot:       `..`,
}

// If `t` is a constant token, a string representation of it is returned; otherwise, an empty string
// is returned.
func getBuiltinTokenValue(t TokenType) string {
	value, found := constantTokenMap[t]
	if !found {
		value = ""
	}
	return value
}

// Make makes a token with a type from its token type. As a special case when the token type is of a
// keyword, the length and value fields of Token will be assigned correct values
func (t TokenType) Make() Token {
	value := getBuiltinTokenValue(t)
	return Token{
		ty: t,
		// this may not always be correct, but can be set to something else later
		length: len(value),
		value:  value,
	}
}

// only adds value when token is not keyword, else just returns token
func (t Token) MaybeAddValue(value string) Token {
	if t.ty.IsConstant() {
		return t
	}
	return t.AddValue(value)
}

// AddValue gives a token a value even if it already has one. AddValue also sets the length field of
// receiver `t` to length(value).
func (t Token) AddValue(value string) Token {
	t.value = value
	if t.length <= 0 { // length set?
		// no, set it
		t.length = len(value)
	}
	return t
}
