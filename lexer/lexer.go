package lexer

import (
	"regexp"
	"strings"
	"unicode"

	"alex.peters/yew/lexer"
	"alex.peters/yew/source"
	itoken "alex.peters/yew/token"
	"yew.lang/main/errors"
	"yew.lang/main/token"
)

func indentation(lex *lexer.Lexer) bool {
	line, char := lex.GetLineChar()
	ws, statWs := source.GetLeadingWhitespace(lex)
	if statWs.NotOk() {
		lex.AddError(errors.MkLexError(lex, errors.MessageFromStatus(statWs)))
	}

	if ws != "" { // found leading whitespace
		tok := token.IndentType.Make().AddValue(ws)
		lex.PushToken(tok.SetLineChar(line, char))
		return true
	}
	return false
}

type symbolClass byte

const (
	symbol symbolClass = iota
	number
	hex
	octal
	binary
	string_
	identifier
	underscore
	comment
)

// ( ) [ ] { } ! @ # $ % ^ & * ~ , < > . ? / ; : | - + = `
const symbolRegexClassRaw string = `[\(\)\[\]\{\}!@#\$%\^\&\*~,<>\.\?/;:\|\-\+=` + "`]"

var symbolRegex = regexp.MustCompile(symbolRegexClassRaw)

const freeSymbolRegexClassRaw string = `[!@#\$%\^\&\*~,<>\.\?/:\|\-\+=` + "`]"

var freeSymbolRegex = regexp.MustCompile(freeSymbolRegexClassRaw)

func isSymbol(c byte) bool {
	return symbolRegex.Match([]byte{c})
}

func determineClass(lex *lexer.Lexer, c byte) (class symbolClass, e error) {
	r := rune(c)
	e = nil
	if unicode.IsLetter(r) {
		class = identifier
	} else if unicode.IsDigit(r) {
		class = number
	} else if c == '"' {
		class = string_
	} else if c == '_' {
		class = underscore
	} else if c == '-' {
		_, stat := lex.AdvanceChar() // this is just '-'
		if stat.NotOk() {
			panic("bug in determineClass in branch `c == '-'`")
		}
		var eof bool
		c, eof = lex.Peek()

		exitAsSymbol := eof || !(c == '-' || c == '*')
		lex.UnadvanceChar()
		if exitAsSymbol {
			class = symbol
		} else {
			class = comment
		}
		return
	} else if isSymbol(c) {
		class = symbol
	} else {
		e = errors.MkLexError(lex, errors.UnknownSymbol)
	}
	return
}

var endMultiCommentRegex = regexp.MustCompile(`\*-`)

func trimSpaceRight(s string) string {
	if len(s) == 0 {
		return s
	}

	i := len(s) - 1
	for ; i >= 0; i-- {
		if !unicode.IsSpace(rune(s[i])) {
			break
		}
	}
	return s[:i]
}

func analyzeComment(lex *lexer.Lexer) source.Status {
	lineNum, charNum := lex.GetLineChar()
	lineNum_0, charNum_0 := lineNum, charNum

	var c byte
	_, stat := lex.AdvanceChar() // remove initial '-'
	if stat.NotOk() {
		statError(lex, stat)
		return stat
	}
	c, stat = lex.AdvanceChar()
	if stat.NotOk() {
		statError(lex, stat)
		return stat
	}

	line, _ := lex.RemainingLine()
	if c == '-' { // single line comment
		length := len(line)
		tok := token.Comment.Make().AddValue(strings.TrimSpace(line))
		lex.PushToken(tok.SetLineChar(lineNum_0, charNum_0))
		lex.SetLineChar(lineNum, length+charNum) // set to end of line
	} else if c == '*' { // multi line comment
		var comment string = ""
		var next string = line
		loc := endMultiCommentRegex.FindStringIndex(line)
		for loc == nil {
			lineNum = lineNum + 1

			lex.SetLineChar(lineNum, 1)
			stat = lex.PositionStatus()
			if stat.Is(source.BadLineNumber) { // at eof?
				statError(lex, source.Eof)
				return source.Eof
			}

			next = strings.TrimSpace(next)
			comment = comment + next
			if len(next) > 0 {
				comment = comment + " "
			}

			// get next line
			var eol bool
			line, eol = lex.RemainingLine()
			if eol {
				statError(lex, source.Eol)
				return source.Eol
			}
			next = line
			
			// check for '*-'
			loc = endMultiCommentRegex.FindStringIndex(line)
		}
		comment = comment + strings.TrimSpace(next[:loc[0]])
		tok := token.Comment.Make().AddValue(comment)
		lex.PushToken(tok.SetLineChar(lineNum_0, charNum_0))
		lex.SetLineChar(lineNum, loc[1])
	} else {
		panic("bug in analyzeComment: else branch reached")
	}
	return stat
}

var intRegex = regexp.MustCompile(`[0-9](_*[0-9]+)*`)
var hexRegex = regexp.MustCompile(`(0x|0X)[0-9a-fA-F](_*[0-9a-fA-F]+)*`)
var octRegex = regexp.MustCompile(`(0o|0O)[0-7](_*[0-7]+)*`)
var binRegex = regexp.MustCompile(`(0b|0B)(0|1)(_*(0|1)+)*`)

func checkNumTail(line string, numEnd int) bool {
	if len(line) <= numEnd {
		return true
	}

	return '_' != line[numEnd] && (
		' ' == line[numEnd] || '\t' == line[numEnd] || isSymbol(line[numEnd]))
}

func stripChar(s string, strip byte) string {
	var builder strings.Builder
	c := rune(strip)
	for _, r := range s {
		if r != c {
			builder.WriteByte(byte(r))
		}
	}
	return builder.String()
}

func analyzeNon10(num, line string) (tok token.Token, numChars int, efunc errors.LazyErrorFn) {
	numChars, efunc = len(num), nil
	if !checkNumTail(line, len(num)) {
		efunc = errors.MkLazyLexError(errors.UnexpectedSymbol)
	} else {
		num = stripChar(num, '_')
		tok = token.IntValue.Make().AddValue(num)
	}
	return
}

func isE(s string, i int) bool {
	return s[i] == 'e' || s[i] == 'E'
}

func isSign(s string, i int) bool {
	return s[i] == '+' || s[i] == '-'
}

func returnInt(num string, numChars int) (token.Token, int, errors.LazyErrorFn) {
	return token.IntValue.Make().AddValue(stripChar(num, '_')), numChars, nil
}

func maybeFractional(num, line string) (tok token.Token, numChars int, efunc errors.LazyErrorFn) {
	numChars, efunc = len(num), nil
	// remove leading zeros (so 0[integer] isn't mistaken as an octal number by llvm or go)
	for numChars != 0 && num[0] == '0' {
		num = num[1:]
	}

	if len(line) <= numChars { // just an integer at the end of the line
		return returnInt(num, numChars)
	}

	// because of above branch, line[numChars] must exist
	eNum := isE(line, numChars)
	dotNum := line[numChars] == '.'

	if !(dotNum || eNum) {
		return returnInt(num, numChars)
	}

	// dotNum must be handled first to account for numbers like '123.123e123'
	if dotNum {
		numChars = numChars + 1
		if len(line) <= numChars { // <integer>.EOL
			efunc = errors.MkLazyLexError(errors.MessageFromStatus(source.Eol))
			return
		}

		frac, ok := locateAtStart(line[numChars:], intRegex)
		if !ok { // <integer>.<non-integer>
			efunc = errors.MkLazyLexError(errors.UnexpectedSymbol)
			return
		}

		numChars = numChars + len(frac)
		num = num + "." + frac

		if len(line) > numChars {
			eNum = isE(line, numChars)
			if !eNum && !checkNumTail(line, numChars) { // <integer>.<integer><illegal-char>
				efunc = errors.MkLazyLexError(errors.UnexpectedSymbol)
				return
			}
		}
	}

	if eNum {
		e := line[numChars] // 'e' or 'E'
		numChars = numChars + 1
		if len(line) <= numChars { // <float>eEOL
			efunc = errors.MkLazyLexError(errors.MessageFromStatus(source.Eol))
			return
		}

		signed := isSign(line, numChars)
		sign := ""
		if signed {
			sign = string(line[numChars])
			numChars = numChars + 1
		}

		if len(line) <= numChars { // <float>e<sign>EOL
			efunc = errors.MkLazyLexError(errors.MessageFromStatus(source.Eol))
			return
		}
		
		frac, ok := locateAtStart(line[numChars:], intRegex)
		if !ok { // <float>e[sign]<illegal-char>
			efunc = errors.MkLazyLexError(errors.UnexpectedSymbol)
			return
		}

		numChars = numChars + len(frac)
		if !checkNumTail(line, numChars) { // <float>e[sign]<integer><illegal-char>
			efunc = errors.MkLazyLexError(errors.UnexpectedSymbol)
			return
		}

		num = num + string(e) + sign + frac
	}

	num = stripChar(num, '_')
	tok = token.FloatValue.Make().AddValue(num)
	return
}

func analyzeNumber(lex *lexer.Lexer) source.Status {
	lineNum, charNum := lex.GetLineChar()
	line, eol := lex.RemainingLine()
	if eol {
		statError(lex, source.Eol)
		return source.Eol
	}

	var tok token.Token
	var numChars int
	var efunc errors.LazyErrorFn = nil

	// 0x, 0b, and 0o must be checked first, else the lexer might falsely think 
	// '0' is the number 
	if num, ok := locateAtStart(line, hexRegex); ok {
		tok, numChars, efunc = analyzeNon10(num, line)
	} else if num, ok := locateAtStart(line, octRegex); ok {
		tok, numChars, efunc = analyzeNon10(num, line)
		if efunc == nil {
			v := tok.GetValue()
			v = "0" + v[2:] // 0o<octal> -> 0<octal>
			tok = tok.AddValue(v)
		}
	} else if num, ok := locateAtStart(line, binRegex); ok {
		tok, numChars, efunc = analyzeNon10(num, line)
	} else if num, ok := locateAtStart(line, intRegex); ok {
		tok, numChars, efunc = maybeFractional(num, line)
	} else {
		numChars = 0
		efunc = errors.MkLazyLexError(errors.UnexpectedSymbol)
	}

	lex.SetLineChar(lineNum, charNum+numChars)
	if efunc != nil {
		efunc(lex)
		return source.Bad
	}

	lex.PushToken(tok.SetLineChar(lineNum, charNum))
	return source.Ok
}

func analyzeInfix(id string) (tok token.Token, numChars int) {
	fixed, unfixedLen := fixInfixId(id)
	numChars = unfixedLen
	tok = token.Infixed.Make().AddValue(fixed)
	return
}

// does not check for '(' and ')'--assumes the id is enclosed within and w/o spaces
func fixInfixId(infixId string) (fixed string, unfixedLen int) {
	unfixedLen = len(infixId)
	fixed = infixId[1 : unfixedLen-1]
	return
}

// \([a-zA-Z][a-zA-Z0-9'_]+\)
var infixIdRegex = regexp.MustCompile(`\(` + idRegexClassRaw + `*\)`)

// \([!@#\$%\^\&\*~,<>\.\?/:\|-\+=`]+\)
var infixSymbolRegex = regexp.MustCompile(`\(` + freeSymbolRegexClassRaw + `+\)`)

// following regex is used after confirming symbol/id is enclosed by parens:
// (\(.*-\*.*?\*-\))|(\(.*--.*?\))
var commentEmbededRegex = regexp.MustCompile(`(\(.*-\*.*?\*-\))|(\(.*--.*?\))`)
var lineComment = regexp.MustCompile(`--`)
var multiLineComment = regexp.MustCompile(`-*`)
var commentRegex = regexp.MustCompile(`(--)|(-*)`)

func locateAtStart(s string, regex *regexp.Regexp) (string, bool) {
	loc := regex.FindStringIndex(s)
	if loc != nil && loc[0] == 0 {
		return s[:loc[1]], true
	}
	return "", false
}

// id will include surrounding parens; this is so the char num can be easily calculated
func maybeInfixId(s string) (id string, isInfix bool) {
	if id, isInfix = locateAtStart(s, infixIdRegex); isInfix {
		// ignore
	} else if id, isInfix = locateAtStart(s, infixSymbolRegex); isInfix {
		isInfix = !enclosesComment(id) // make sure no comment inside
	}
	return
}

func enclosesComment(infixedId string) bool {
	return commentEmbededRegex.MatchString(infixedId)
}

func handleLParen(line string) (tok token.Token, numChars int) {
	if len(line) > 1 && ')' == line[1] {
		tok, numChars = token.Empty.Make(), 2
	} else {
		id, isInfix := maybeInfixId(line)
		if isInfix {
			tok, numChars = analyzeInfix(id)
		} else {
			tok, numChars = token.LeftParen.Make(), 1
		}
	}
	return
}

var freeSymbolFullRegex = regexp.MustCompile(freeSymbolRegexClassRaw + "+")
func tokenizeSymbol(line string) (tok token.Token, numChars int, efunc func(*lexer.Lexer)errors.LexError) {
	efunc = nil
	res, ok := locateAtStart(line, freeSymbolFullRegex)
	if !ok {
		efunc = errors.MkLazyLexError(errors.UnexpectedSymbol)
		return
	}
	// check for comment
	loc := commentEmbededRegex.FindStringIndex(res)
	if loc != nil {
		// sanity check--if comments are checked for before here, then 
		// this should be impossible
		if loc[0] == 0 {
			panic("bug: comment not scanned for before scanning symbol")
		}
		res = res[:loc[0]] // remove comment
	}
	numChars = len(res)
	tok = token.SymbolType.Make().AddValue(res)
	return
}

func analyzeSymbol(lex *lexer.Lexer) source.Status {
	line, char := lex.GetLineChar()
	remainingLine, eol := lex.RemainingLine()
	if eol {
		statError(lex, source.Eol)
		return source.Eol
	}

	c := remainingLine[0]

	var tok token.Token
	var numChars int = 1
	switch c {
	case '(':
		tok, numChars = handleLParen(remainingLine)
	case ')':
		tok = token.RightParen.Make()
	case '[':
		tok = token.LeftBracket.Make()
	case ']':
		tok = token.RightBracket.Make()
	case '{':
		tok = token.LeftBrace.Make()
	case '}':
		tok = token.RightBrace.Make()
	case ';':
		tok = token.SemiColon.Make()
	case ',':
		tok = token.Comma.Make()
	default:
		var efunc func(*lexer.Lexer)errors.LexError
		tok, numChars, efunc = tokenizeSymbol(remainingLine)
		if efunc != nil {
			efunc(lex)
			return source.Bad
		}
	}
	lex.PushToken(tok.SetLineChar(line, char))
	lex.SetLineChar(line, numChars+char)
	return source.Ok
}

func statError(lex *lexer.Lexer, stat source.Status) {
	lex.AddError(errors.MkLexError(lex, errors.MessageFromStatus(stat)))
}

func lexError(lex *lexer.Lexer, msg string, args ...any) {
	lex.AddError(errors.MkLexError(lex, msg, args...))
}

func getEscape(r rune, escapeString bool) (c byte, ok bool) {
	ok = true
	switch r {
	case 'n':
		c = '\n'
	case 't':
		c = '\t'
	case 'r':
		c = '\r'
	case 'v':
		c = '\v'
	case 'b':
		c = '\b'
	case 'a':
		c = '\a'
	case 'f':
		c = '\f'
	case '\\':
		c = '\\'
	case '"':
		if escapeString {
			c = '"'
		} else {
			ok = false
		}
	case '\'':
		if !escapeString {
			c = '\''
		} else {
			ok = false
		}
	default:
		ok = false
	}
	return
}

func updateEscape(s string, escapeString bool) (string, bool, int) {
	var builder strings.Builder
	var next bool = false
	out := len(s)-1
	for i, r := range s {
		if next {
			next = false
			c, ok := getEscape(r, escapeString)
			if !ok {
				return "", false, i
			}
			builder.WriteByte(c)
		} else if r == '\\' {
			next = true
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String(), true, out
}

func analyzeString(lex *lexer.Lexer) source.Status {
	line, char := lex.GetLineChar()
	c, stat := lex.AdvanceChar()
	if stat.NotOk() || c != '"' {
		if c != '"' {
			stat = source.Bad
			lexError(lex, errors.UnexpectedSymbol)
		} else {
			statError(lex, stat)
		}
		return stat
	}

	res := ""
	for section := ""; true ; {
		section, stat = source.ReadThrough(lex, '"')
		if stat.NotOk() {
			break
		}

		res = res + section
		secLen := len(section)
		if secLen >= 2 && section[secLen-2:] == `\"` {
			resHead := section[:secLen-2]
			isQuoteEsc := true // \"=t \\"=f \\\"=t \\\\"=f ... flips b/w t and f
			// search from right to left to see if it's an escaped quote or end of string
			for i := len(resHead)-1; i >= 0; i-- {
				if resHead[i] != '\\' {
					break
				}
		
				isQuoteEsc = !isQuoteEsc
			}

			if isQuoteEsc {
				continue
			}
		}

		break
	}

	tot := len(res)
	if tot > 0 {
		res = res[:len(res)-1] // remove trailing '"'
	}

	if stat.IsOk() {
		var ok bool
		var index int
		res, ok, index = updateEscape(res, true)
		if !ok {
			lex.SetLineChar(line, char+index+1)
			lexError(lex, errors.IllegalEscape)
			return source.Bad
		}
		_, _ = lex.AdvanceChar() // eat `"`
		tok := token.StringValue.
			Make().
			AddValue(res).
			SetLineChar(line, char)
		lex.PushToken(tok)
	} else {
		statError(lex, stat)
	}

	lex.SetLineChar(line, char+tot)

	return stat
}

// assumes len(s) >= 1
func resolveType(s string) token.TokenType {
	switch s[0] {
	case 'l':
		if s == "let" {
			return token.Let
		}
	case 'c':
		if s == "class" {
			return token.Class
		}
	case 'o':
		if s == "of" {
			return token.Of
		}
	case 'm':
		if s == "module" {
			return token.Module
		}
	case 'w':
		if s == "where" {
			return token.Where
		}
	case 'u':
		if s == "use" {
			return token.Use
		}
	case 'i':
		if s == "import" {
			return token.Import
		}
	case 'h':
		if s == "hide" {
			return token.Hide
		}
	case 'd':
		if s == "derives" {
			return token.Derives
		}
	case 'f':
		if s == "forall" {
			return token.Forall
		}
	}
	return token.IdType
}

var idRegexClassRaw = `[a-zA-Z][a-zA-Z0-9'_]`
var idRegex = regexp.MustCompile(idRegexClassRaw + `*`)

func analyzeIdentifier(lex *lexer.Lexer) source.Status {
	line, char := lex.GetLineChar()
	src, stat := source.GetSourceSlice(lex, line, char, -1)
	if stat.NotOk() {
		statError(lex, stat)
		return stat
	}

	res, length := lexer.RegexMatch(idRegex, src)
	if length < 1 {
		stat = source.Bad
		statError(lex, stat)
		return stat
	}

	ty := resolveType(res)

	// add token
	tok, _ := (token.Token{}).
		SetLineChar(line, char).
		SetType(uint(ty)).
		SetValue(res)
	lex.PushToken(tok)

	// set lexer's char num
	lex.SetLineChar(line, char+length+1)
	return stat
}

func analyzeUnderscore(lex *lexer.Lexer) source.Status {
	line, char := lex.GetLineChar()
	_, stat := lex.AdvanceChar()
	if stat.NotOk() {
		lexError(lex, errors.UnexpectedSymbol)
		return stat
	}

	tok := token.Wildcard.Make().SetLineChar(line, char)
	lex.PushToken(tok)
	return stat
}

func (class symbolClass) analyze(lex *lexer.Lexer) source.Status {
	switch class {
	case number:
		return analyzeNumber(lex)
	case symbol:
		return analyzeSymbol(lex)
	case string_:
		return analyzeString(lex)
	case identifier:
		return analyzeIdentifier(lex)
	case underscore:
		return analyzeUnderscore(lex)
	case comment:
		return analyzeComment(lex)
	}

	e := errors.MkLexError(lex, errors.UnknownSymbol)
	lex.AddError(e)
	return source.Bad
}

func analyze(lex *lexer.Lexer) source.Status {
	if indentation(lex) {
		return source.Ok
	}

	source.SkipWhitespace(lex)

	c, eof := lex.Peek()
	if eof {
		return source.Eof
	}
	class, e := determineClass(lex, c)
	if e != nil {
		lex.AddError(e)
		return source.Bad
	}
	return class.analyze(lex)
}

var lexerWhitespace = regexp.MustCompile(`\t| `)

func runLexer(path string) ([]itoken.Token, []error) {
	lex, e := lexer.Lex(path, lexerWhitespace)
	if e != nil {
		return nil, []error{errors.MkSystemError(e.Error())}
	}

	stat := analyze(lex)
	for stat.IsOk() {
		stat = analyze(lex)
	}
	return lex.GetTokens(), lex.GetErrors()
}
