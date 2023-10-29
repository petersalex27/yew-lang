package lexer

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/lexer"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
)

// Returns true iff source position is at the start of a line AND source is
// positioned at a non-empty sequence of whitespace. There are two cases in
// which a token is pushed into the lexer's token buffer:
//   - (1) function would return true AND source is positioned at a non-empty
//     sequence of whitespace; in this case, the sequence of whitespace is
//     pushed as an indent token to the lexer's token buffer
//   - (2) function would return false AND source is positioned at the start of
//     a line (this implies an empty sequence of whitespace at position); in
//     this case, an empty indent token is pushed to the lexer's token
//     buffer
//
// NOTE: you can think of the return value of the function as whether or not
// the source position is advanced
func indentation(lex *lexer.Lexer) (isPositionAdvanced bool) {
	line, char := lex.GetLineChar()
	isLineStart := char == 1
	ws, statWs := source.GetLeadingWhitespace(lex)
	if statWs.NotOk() {
		lex.AddError(errors.Lex(lex, errors.MessageFromStatus(statWs)))
		return false
	}

	if ws != "" { // found leading whitespace
		tok := token.Indent.Make().AddValue(ws)
		lex.PushToken(tok.SetLineChar(line, char))
		return true
	} else if isLineStart {
		// add empty indentation
		tok := token.Indent.Make().AddValue("")
		lex.PushToken(tok.SetLineChar(line, char))
		return false // NOTE: notice false return
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
	stringClass
	identifier
	underscore
	comment
	char
)

// ( ) [ ] { } ! @ # $ % ^ & * ~ , < > . ? / ; : | - + = \ `
const symbolRegexClassRaw string = `[\(\)\[\]\{\}!@#\$%\^\&\*~,<>\.\?/;:\|\-\+=\\` + "`]"

var symbolRegex = regexp.MustCompile(symbolRegexClassRaw)

const freeSymbolRegexClassRaw string = `[!#\$%\^\&\*~,<>\.\?/:\|\-\+=\\` + "`]"

//var freeSymbolRegex = regexp.MustCompile(freeSymbolRegexClassRaw)

func isSymbol(c byte) bool {
	return symbolRegex.Match([]byte{c})
}

// Determines the class of some input section based on some byte `c` of the
// input. Unless there's a good reason to do otherwise, `c` is the first
// character of that input section.
func determineClass(lex *lexer.Lexer, c byte) (class symbolClass, e error) {
	r := rune(c)
	e = nil
	if unicode.IsLetter(r) {
		class = identifier
	} else if unicode.IsDigit(r) {
		class = number
	} else if c == '\'' {
		class = char
	} else if c == '"' {
		class = stringClass
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
		e = errors.Lex(lex, errors.UnknownSymbol)
	}
	return
}

var endMultiCommentRegex = regexp.MustCompile(`\*-`)

/*
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
}*/

// reads and pushes respective token for single-line comments and single-line annotations.
//
// NOTE:
//		--@MyAnnot
// is an annotation, but
//		-- @MyAnnot
// is not an annotation because there is whitespace before '@'
func getSingleLineComment(lex *lexer.Lexer, line string, lineNum, charNum int) source.Status {
	var ty token.TokenType = token.Comment // type of token to be pushed

	length := len(line)

	// check first char of `line`; if it's '@' comment is an annotation
	if length > 0 && line[0] == '@' {
		ty = token.Annotation
		line = line[1:] // remove '@' from "comment"
	}

	// remove extra whitespace
	// NOTE: the order of this statement and the previous conditional branch is
	// important. It prevents `--<whitespace>@` from being considered an annotation 
	trimmed := strings.TrimSpace(line)
	tok := ty.Make().AddValue(trimmed)
	lex.PushToken(tok.SetLineChar(lineNum, charNum))
	lex.SetLineChar(lineNum, length+charNum) // set to end of line
	return source.Ok
}

// reads and pushes respective token for multi-line comments and multi-line annotations 
//
// NOTE:
//		-*@MyAnnot .. *-
// is an annotation, but
//		-* @MyAnnot .. *-
// is not an annotation because there is whitespace before '@'
func getMultiLineComment(lex *lexer.Lexer, line string, lineNum, charNum int) source.Status {
	var ty token.TokenType = token.Comment // type of token to be pushed
	stat := source.Ok
	lineNum_0, charNum_0 := lineNum, charNum // initial line and char numbers
	var comment string = "" // full comment
	var next string = line // next line to analyze

	// check first char of `next`; if it's '@' comment is an annotation
	if len(next) > 0 && next[0] == '@' {
		ty = token.Annotation
		next = next[1:] // remove '@' from "comment"
	}

	loc := endMultiCommentRegex.FindStringIndex(line) // check for end of comment
	// append input read to comment until end of comment is reached
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
	tok := ty.Make().AddValue(comment)
	lex.PushToken(tok.SetLineChar(lineNum_0, charNum_0))
	lex.SetLineChar(lineNum, loc[1])
	return stat
}

func analyzeComment(lex *lexer.Lexer) source.Status {
	lineNum, charNum := lex.GetLineChar()

	var c byte
	_, stat := lex.AdvanceChar() // remove initial '-'
	if stat.NotOk() {
		statError(lex, stat)
		return stat
	}
	// remove second thing ('-' or '*'), the value of `c` will determine the
	// branch to take in the condition below the next one
	c, stat = lex.AdvanceChar()
	if stat.NotOk() {
		statError(lex, stat)
		return stat
	}

	// given a line
	//		-- abc ..
	// or
	//		-* abc ..
	// `line` is
	//		line = " abc .."
	line, _ := lex.RemainingLine()
	if c == '-' { // single line comment
		return getSingleLineComment(lex, line, lineNum, charNum)
	} else if c == '*' { // multi line comment
		return getMultiLineComment(lex, line, lineNum, charNum)
	}
	panic("bug in analyzeComment: else branch reached")
}

var intRegex = regexp.MustCompile(`[0-9](_*[0-9]+)*`)
var hexRegex = regexp.MustCompile(`(0x|0X)[0-9a-fA-F](_*[0-9a-fA-F]+)*`)
var octRegex = regexp.MustCompile(`(0o|0O)[0-7](_*[0-7]+)*`)
var binRegex = regexp.MustCompile(`(0b|0B)(0|1)(_*(0|1)+)*`)

func checkNumTail(line string, numEnd int) bool {
	if len(line) <= numEnd {
		return true
	}

	return line[numEnd] != '_' && (line[numEnd] == '_' || line[numEnd] == '\t' || isSymbol(line[numEnd]))
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
		efunc = errors.LazyLex(errors.UnexpectedSymbol)
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
			efunc = errors.LazyLex(errors.MessageFromStatus(source.Eol))
			return
		}

		frac, ok := locateAtStart(line[numChars:], intRegex)
		if !ok { // <integer>.<non-integer>
			efunc = errors.LazyLex(errors.UnexpectedSymbol)
			return
		}

		numChars = numChars + len(frac)
		num = num + "." + frac

		if len(line) > numChars {
			eNum = isE(line, numChars)
			if !eNum && !checkNumTail(line, numChars) { // <integer>.<integer><illegal-char>
				efunc = errors.LazyLex(errors.UnexpectedSymbol)
				return
			}
		}
	}

	if eNum {
		e := line[numChars] // 'e' or 'E'
		numChars = numChars + 1
		if len(line) <= numChars { // <float>eEOL
			efunc = errors.LazyLex(errors.MessageFromStatus(source.Eol))
			return
		}

		signed := isSign(line, numChars)
		sign := ""
		if signed {
			sign = string(line[numChars])
			numChars = numChars + 1
		}

		if len(line) <= numChars { // <float>e<sign>EOL
			efunc = errors.LazyLex(errors.MessageFromStatus(source.Eol))
			return
		}

		frac, ok := locateAtStart(line[numChars:], intRegex)
		if !ok { // <float>e[sign]<illegal-char>
			efunc = errors.LazyLex(errors.UnexpectedSymbol)
			return
		}

		numChars = numChars + len(frac)
		if !checkNumTail(line, numChars) { // <float>e[sign]<integer><illegal-char>
			efunc = errors.LazyLex(errors.UnexpectedSymbol)
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
		efunc = errors.LazyLex(errors.UnexpectedSymbol)
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
	fixed, unfixedLen := fixEnclosedId(id)
	numChars = unfixedLen
	tok = token.Infixed.Make().AddValue(fixed)
	return
}

// does not check for '(' and ')' (or '{' and '}')--assumes the id is enclosed within and w/o spaces
func fixEnclosedId(enclosedId string) (fixed string, unfixedLen int) {
	unfixedLen = len(enclosedId)
	fixed = enclosedId[1 : unfixedLen-1]
	return
}

// \([a-zA-Z][a-zA-Z0-9'_]+\)
var infixIdRegex = regexp.MustCompile(`\(` + idRegexClassRaw + `*\)`)

// \([!@#\$%\^\&\*~,<>\.\?/:\|-\+=`]+\)
var infixSymbolRegex = regexp.MustCompile(`\(` + freeSymbolRegexClassRaw + `+\)`)

// following regex is used after confirming symbol/id is enclosed by parens:
// (\(.*-\*.*?\*-\))|(\(.*--.*?\))
var commentEmbededRegex = regexp.MustCompile(`(\(.*-\*.*?\*-\))|(\(.*--.*?\))`)

//var lineComment = regexp.MustCompile(`--`)
//var multiLineComment = regexp.MustCompile(`-*`)
//var commentRegex = regexp.MustCompile(`(--)|(-*)`)

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
	id, isInfix := maybeInfixId(line)
	if isInfix {
		tok, numChars = analyzeInfix(id)
	} else {
		tok, numChars = token.LeftParen.Make(), 1
	}
	return
}

var freeSymbolFullRegex = regexp.MustCompile(freeSymbolRegexClassRaw + "+")

func tokenizeSymbol(line string) (tok token.Token, numChars int, efunc errors.LazyErrorFn) {
	efunc = nil
	res, ok := locateAtStart(line, freeSymbolFullRegex)
	if !ok {
		efunc = errors.LazyLex(errors.UnexpectedSymbol)
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
	if ty, found := builtinSymbols[res]; found {
		tok = ty.Make()
	} else {
		tok = token.Symbol.Make().AddValue(res)
	}
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
		var efunc errors.LazyErrorFn
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
	lex.AddError(errors.Lex(lex, errors.MessageFromStatus(stat)))
}

func lexError(lex *lexer.Lexer, msg string, args ...any) {
	lex.AddError(errors.Lex(lex, msg, args...))
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

func readEscapable(line string, end byte) (string, int, source.Status) {
	index := 0
	escaped := false
	for _, c := range line {
		if escaped {
			escaped = false
		} else if byte(c) == end {
			return line[:index], index, source.Ok
		} else if byte(c) == '\\' {
			escaped = true
		}
		index = index + 1
	}
	// `end` not found
	return "", index, source.Eol
}

func updateEscape(s string, escapeString bool) (string, bool, int) {
	var builder strings.Builder
	var next bool = false
	out := len(s) - 1
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

func analyzeChar(lex *lexer.Lexer) source.Status {
	line, char := lex.GetLineChar()
	c, stat := lex.AdvanceChar()
	if stat.NotOk() || c != '\'' {
		if c != '\'' {
			stat = source.Bad
			lexError(lex, errors.UnexpectedSymbol)
		} else {
			statError(lex, stat)
		}
		return stat
	}

	remainingLine, eof := lex.RemainingLine()
	if eof {
		statError(lex, source.Eof)
		return source.Eof
	}

	res, length, stat := readEscapable(remainingLine, '\'')
	if stat.NotOk() {
		statError(lex, stat)
		return stat
	}

	lex.SetLineChar(line, char+length)
	if _, stat = lex.AdvanceChar(); stat.NotOk() { // remove closing `'`
		statError(lex, stat)
		return stat
	}

	var ok bool
	var index int
	res, ok, index = updateEscape(res, false)
	if !ok {
		lex.SetLineChar(line, char+index+1)
		lexError(lex, errors.IllegalEscape)
		return source.Bad
	}
	if len(res) != 1 {
		lex.SetLineChar(line, char)
		lexError(lex, errors.IllegalChar)
		return source.Bad
	}
	tok := token.CharValue.
		Make().
		AddValue(res).
		SetLineChar(line, char)
	lex.PushToken(tok)

	return stat
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
	for section := ""; true; {
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
			for i := len(resHead) - 1; i >= 0; i-- {
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
		//_, _ = lex.AdvanceChar() // eat `"`
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
	if unicode.IsUpper(rune(s[0])) {
		return token.TypeId
	}

	if keySelect, found := keywordTrie[s[0]]; found {
		if ty, found := keySelect[s]; found {
			return ty
		}
	}
	return token.Id
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

	// add token
	tok := resolveType(res).
		Make().
		MaybeAddValue(res).
		SetLineChar(line, char)
	lex.PushToken(tok)

	// set lexer's char num
	lex.SetLineChar(line, char+length)
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
	case char:
		return analyzeChar(lex)
	case stringClass:
		return analyzeString(lex)
	case identifier:
		return analyzeIdentifier(lex)
	case underscore:
		return analyzeUnderscore(lex)
	case comment:
		return analyzeComment(lex)
	}

	e := errors.Lex(lex, errors.UnknownSymbol)
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

func RunLexer(lex *lexer.Lexer) (tokens []itoken.Token, es []error) {
	// keep reading tokens until stat is not Ok
	stat := analyze(lex)
	for stat.IsOk() {
		stat = analyze(lex)
	}

	tokens, es = lex.GetTokens(), lex.GetErrors()
	if len(tokens) > 0 {
		// prepend 0 length indent
		// create indent token that exists at start of first line
		tok := token.Indent.Make().AddValue("").SetLineChar(1, 1)
		// add indent token
		tokens = append([]itoken.Token{tok}, tokens...)
	}
	return
}

func CastTokens(ts []itoken.Token) []token.Token {
	out := make([]token.Token, len(ts))
	for i, tok := range ts {
		out[i] = tok.(token.Token)
	}
	return out
}

func NewLexer(path string) *lexer.Lexer {
	lex, e := lexer.Lex(path, lexerWhitespace)
	if e != nil {
		errors.PrintErrors(errors.MkSystemError(e.Error()))
		return nil
	}
	return lex
}

func GetSourceRaw(lex *lexer.Lexer) []string {
	n := lex.NumLines()
	out := make([]string, n)
	for i := range out {
		line, _ := lex.SourceLine(i + 1)
		out[i] = line
	}
	return out
}

func LexPath(path string) ([]itoken.Token, []error) {
	lex, e := lexer.Lex(path, lexerWhitespace)
	if e != nil {
		return nil, []error{errors.MkSystemError(e.Error())}
	}

	return RunLexer(lex)
}
