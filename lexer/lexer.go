package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/lexer"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
)

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

// reads and pushes respective token for single-line comments and single-line annotations.
//
// NOTE:
//
//	--@MyAnnot
//
// is an annotation, but
//
//	-- @MyAnnot
//
// is not an annotation because there is whitespace before '@'
func getSingleLineComment(lex *lexer.Lexer, lineAfterDashes string, lineNum, charNum int) source.Status {
	var ty token.TokenType = token.Comment // type of token to be pushed

	// length of comment/annotation excluding leading `--`
	length := len(lineAfterDashes)
	// length of entire comment/annotation including leading `--`
	totalLength := length + 2 // +2 for "--"

	var trimmed string

	// check first char of `line`; if it's '@' comment is an annotation
	if length > 0 && lineAfterDashes[0] == '@' {
		ty = token.Annotation
		lineAfterDashes = lineAfterDashes[1:] // remove '@' from "comment"
		// identity, don't let annotations get their content modified
		trimmed = lineAfterDashes
	} else {
		// this must be done separate from annotations so their values aren't
		// modified
		trimmed = lex.Documentor.Run(lineAfterDashes)
	}

	tok := ty.Make().AddValue(trimmed).SetLength(totalLength)
	lex.PushToken(tok.SetLineChar(lineNum, charNum))

	lex.SetLineChar(lineNum, totalLength+charNum) // set to end of line

	return source.Ok
}

// just returns comment; use when annotations are found
func identityProcessing(_ *lexer.Documentor, comment string) string { return comment }

// reads and pushes respective token for multi-line comments and multi-line annotations
//
// NOTE:
//
//	-*@MyAnnot .. *-
//
// is an annotation, but
//
//	-* @MyAnnot .. *-
//
// is not an annotation because there is whitespace before '@'
func getMultiLineComment(lex *lexer.Lexer, line string, lineNum, charNum int) source.Status {
	var ty token.TokenType = token.Comment // type of token to be pushed
	stat := source.Ok
	lineNum_0, charNum_0 := lineNum, charNum // initial line and char numbers
	var comment string = ""                  // full comment
	var next string = line                   // next line to analyze
	// +2 to account for '-*'
	firstLineOfCommentLength := +2

	// check first char of `next`; if it's '@' comment is an annotation
	if len(next) > 0 && next[0] == '@' {
		ty = token.Annotation
		next = next[1:] // remove '@' from "comment"

		// THIS IS A VITAL IMPORTANT STEP!!
		restore := lex.Documentor
		defer func() { lex.Documentor = restore }()
		lex.Documentor = lexer.MakeDocumentor(identityProcessing)
	}

	// check for end of comment
	loc := endMultiCommentRegex.FindStringIndex(line)
	// check if comment fits on one line
	if loc != nil {
		// get end location of comment
		//
		// length of comment is: end - (start - 1)
		//  _________________length=15________
		//	______________vvvvvvvvvvvvvvv_____
		//	3 | blah blah -* blah blah *- blah
		//	____^_________^_____________^_____
		//  _______start 3:11______end 3:25___
		//  => length = end - (start - 1)
		//            = 25 - (11 - 1)
		//            = 15
		end, start := loc[1], charNum_0
		// add two to account for leading '-*' not in `line`
		firstLineOfCommentLength = 2 + (end - (start - 1))
	} else {
		// set length of token to length of first line from '-*' to end of line
		firstLineOfCommentLength = firstLineOfCommentLength + len(line)
	}

	// append input read to comment until end of comment is reached
	for loc == nil {
		lineNum = lineNum + 1

		lex.SetLineChar(lineNum, 1)
		stat = lex.PositionStatus()
		if stat.Is(source.BadLineNumber) { // at eof?
			statError(lex, source.Eof)
			return source.Eof
		}

		// remove extra white space when appending comment
		next = lex.Documentor.Run(next)
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

	startIndexOfCommentClose := loc[0]
	endCharNumOfCommentClose := loc[1]

	// grab last line of comment up to (but not including) '*-'
	commentFinalContent := next[:startIndexOfCommentClose]
	finalContent := lex.Documentor.Run(commentFinalContent)

	// create an push comment token
	comment = comment + finalContent
	tok := ty.Make().AddValue(comment).SetLength(firstLineOfCommentLength)
	lex.PushToken(tok.SetLineChar(lineNum_0, charNum_0))

	lex.SetLineChar(lineNum, endCharNumOfCommentClose)

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
	var line string
	line, _ = lex.RemainingLine()
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

func isNumEndCharValid(line string, numEnd int) bool {
	if len(line) <= numEnd {
		return true
	}

	return line[numEnd] != '_' && (line[numEnd] == '_' || line[numEnd] == '\t' || isSymbol(line[numEnd]))
}

// removes `strip` from `s` and returns result
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

// read non-base ten integer token: hexadecimal, octal, or binary
func analyzeNonBase10(num, line string) (tok token.Token, numChars int, errorFunc errors.LazyErrorFn) {
	numChars, errorFunc = len(num), nil
	if !isNumEndCharValid(line, len(num)) {
		errorFunc = errors.LazyLex(errors.UnexpectedSymbol)
	} else {
		tokenLength := len(num)
		num = stripChar(num, '_')
		tok = token.IntValue.Make().AddValue(num).SetLength(tokenLength)
	}
	return
}

func isE(s string, i int) bool {
	return s[i] == 'e' || s[i] == 'E'
}

func isSign(s string, i int) bool {
	return s[i] == '+' || s[i] == '-'
}

// assumes numChars is the correct value and that it corresponds to the length
// of the token
func returnInt(num string, numChars int) (token.Token, int, errors.LazyErrorFn) {
	num = stripChar(num, '_')
	return token.IntValue.Make().AddValue(num).SetLength(numChars), numChars, nil
}

func analyzeDotNum(numOrigin, line string, numCharsOrigin int, hasEOrigin bool) (num string, numChars int, hasE bool, errorFunc errors.LazyErrorFn) {
	num, numChars, hasE = numOrigin, numCharsOrigin, hasEOrigin // init

	numChars = numChars + 1
	if len(line) <= numChars { // <integer>.EOL
		errorFunc = errors.LazyLex(errors.MessageFromStatus(source.Eol))
		return
	}

	frac, ok := locateAtStart(line[numChars:], intRegex)
	if !ok { // <integer>.<non-integer>
		errorFunc = errors.LazyLex(errors.UnexpectedSymbol)
		return
	}

	numChars = numChars + len(frac)
	num = num + "." + frac

	if len(line) <= numChars {
		return
	}

	hasE = isE(line, numChars)
	if !hasE && !isNumEndCharValid(line, numChars) { // <integer>.<integer><illegal-char>
		errorFunc = errors.LazyLex(errors.UnexpectedSymbol)
	}
	return
}

// read sign if one exists and return (possibly) empty sign and new total
// number of chars read
func analyzePossibleSign(line string, numCharsOrigin int) (sign string, numChars int) {
	// init
	numChars = numCharsOrigin

	// if signed, read sign and return it
	signed := isSign(line, numChars)
	sign = ""
	if signed {
		sign = string(line[numChars])
		numChars = numChars + 1
	}

	return sign, numChars
}

// ASSUMPTION: line[numChars] == 'e' or 'E'
//
// reads number from input at exponent marker (i.e., 'e' or 'E') to end of number
func analyzeExponentNum(numOrigin, line string, numCharsOrigin int) (num string, numChars int, errorFunc errors.LazyErrorFn) {
	num, numChars = numOrigin, numCharsOrigin // init

	e := line[numChars] // 'e' or 'E'
	numChars = numChars + 1
	if len(line) <= numChars { // <float>eEOL
		errorFunc = errors.LazyLex(errors.MessageFromStatus(source.Eol))
		return
	}

	var sign string
	sign, numChars = analyzePossibleSign(line, numChars)

	if len(line) <= numChars { // <float>e<sign>EOL
		errorFunc = errors.LazyLex(errors.MessageFromStatus(source.Eol))
		return
	}

	// read integer value that follows 'e'/'E'
	frac, ok := locateAtStart(line[numChars:], intRegex)
	if !ok { // <float>e[sign]<illegal-char>
		errorFunc = errors.LazyLex(errors.UnexpectedSymbol)
		return
	}

	numChars = numChars + len(frac)
	if !isNumEndCharValid(line, numChars) { // <float>e[sign]<integer><illegal-char>
		errorFunc = errors.LazyLex(errors.UnexpectedSymbol)
		return
	}

	// build value as string
	num = num + string(e) + sign + frac
	return
}

// return a number token. Could be either floating point number or integer
func maybeFractional(num, line string) (tok token.Token, numChars int, errorFunc errors.LazyErrorFn) {
	tokenLength := 0
	numChars, errorFunc = len(num), nil
	// remove leading zeros (so 0[integer] isn't mistaken as an octal number by llvm or go)
	for numChars != 0 && num[0] == '0' {
		num = num[1:]
		tokenLength++ // increment token length to account for '0's stripped
	}

	if len(line) <= numChars { // just an integer at the end of the line
		return returnInt(num, numChars)
	}

	// because of above branch, line[numChars] must exist
	hasE := isE(line, numChars)
	hasDot := line[numChars] == '.'

	if !hasDot && !hasE {
		return returnInt(num, numChars)
	}

	// dotNum must be handled first to account for numbers like '123.123e123'
	if hasDot {
		num, numChars, hasE, errorFunc = analyzeDotNum(num, line, numChars, hasE)
	}

	// read 'e' or 'E' and exponent
	if hasE {
		num, numChars, errorFunc = analyzeExponentNum(num, line, numChars)
	}

	tokenLength = tokenLength + len(num) // leading zeros + remaining num
	num = stripChar(num, '_')
	tok = token.FloatValue.Make().AddValue(num).SetLength(tokenLength)
	return
}

// read number from input
func analyzeNumber(lex *lexer.Lexer) source.Status {
	lineNum, charNum := lex.GetLineChar()
	line, eol := lex.RemainingLine()
	if eol { // end of line?
		statError(lex, source.Eol)
		return source.Eol
	}

	var tok token.Token                    // token result
	var numChars int                       // total number of chars read
	var errorFunc errors.LazyErrorFn = nil // error function

	// 0x, 0b, and 0o must be checked first, else the lexer might falsely think
	// '0' is the number
	if num, ok := locateAtStart(line, hexRegex); ok { // hex
		tok, numChars, errorFunc = analyzeNonBase10(num, line)
	} else if num, ok := locateAtStart(line, octRegex); ok { // oct
		tok, numChars, errorFunc = analyzeNonBase10(num, line)
	} else if num, ok := locateAtStart(line, binRegex); ok { // bin
		tok, numChars, errorFunc = analyzeNonBase10(num, line)
	} else if num, ok := locateAtStart(line, intRegex); ok { // int or float
		tok, numChars, errorFunc = maybeFractional(num, line)
	} else {
		numChars = 0
		errorFunc = errors.LazyLex(errors.UnexpectedSymbol)
	}

	lex.SetLineChar(lineNum, charNum+numChars)
	if errorFunc != nil {
		errorFunc(lex)
		return source.Bad
	}

	// finish creating number token, token length already set, just set
	// line number and char number
	token := tok.SetLineChar(lineNum, charNum)
	lex.PushToken(token)
	return source.Ok
}

func fixedRegexGen(element string) string {
	return fmt.Sprintf(`(%s)?_?(((%s_)+(%s)?)|(%s))`, element, element, element, element)
}

// ([a-z][a-zA-Z0-9']*)?_?((([a-z][a-zA-Z0-9']*_)+([a-z][a-zA-Z0-9']*)?)|([a-z][a-zA-Z0-9']*))
var infixIdRegex = regexp.MustCompile(fixedRegexGen(idRegexClassRaw + `*`))

// \([!@#\$%\^\&\*~,<>\.\?/:\|-\+=`]+\)
var infixSymbolRegex = regexp.MustCompile(fixedRegexGen(freeSymbolRegexClassRaw + `+`))

// following regex is used after confirming symbol/id is infix:
// (.*-\*.*?\*-)|(.*--.*?)
//var commentEmbeddedRegex = regexp.MustCompile(`(.*-\*.*?\*-)|(.*--.*?)`)

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
		tok = token.LeftParen.Make()
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
	case ',':
		tok = token.Comma.Make()
	default:
		return analyzeIdentifier(lex)
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
	c, stat := lex.AdvanceChar() // should be leading '
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
		SetLength(length+2). // +2 for enclosing single quotes
		SetLineChar(line, char)
	lex.PushToken(tok)

	return stat
}

// determines what error caused the first step of string tokenizing to fail
func determineInitialFailForAnalyzeString(lex *lexer.Lexer, stat source.Status, c byte) source.Status {
	if c != '"' {
		// failed because not a quote char
		stat = source.Bad
		lexError(lex, errors.UnexpectedSymbol)
	} else {
		// failed b/c of status error
		statError(lex, stat)
	}
	return stat
}

// This counts number of contiguous `c`s at the end of `s`.
//
// Examples:
//
//	countTrailing("employee", 'e') = 2
//	countTrailing("employee", 'y') = 0
//	countTrailing("", 'w') = 0
func countTrailing(s string, c byte) uint {
	trailing := uint(0)
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != c {
			break
		}

		trailing++
	}
	return trailing
}

// returns true if `s` ends in an escaped quote
func hasFinalQuoteEscape(s string) bool {
	length := len(s)
	possibleEscapedQuote := length >= 2 && s[length-2:] == `\"`
	if !possibleEscapedQuote {
		return false
	}

	// remove final '"' so trailing backslashes can be counted
	unquoted := s[:length-2]
	// if number of escapes is 2n for some n, then there are n escaped '\\'; if there are 2n+1
	// '\\', then there are n escaped '\\' and a final escaped '"'
	isEscapedQuote := (countTrailing(unquoted, '\\') % 2) != 0

	return isEscapedQuote
}

func getStringContent(lex *lexer.Lexer) (stat source.Status, content string, charsRead int) {
	var section string
	content = ""

	// reads string (and accounts for escaped '"')
	again := true
	for again {
		section, stat = source.ReadThrough(lex, '"')
		if stat.NotOk() {
			return
		}

		content = content + section
		again = hasFinalQuoteEscape(section)
	}

	charsRead = len(content)
	if charsRead > 0 {
		content = content[:charsRead-1] // remove trailing '"'
	}
	return
}

func analyzeString(lex *lexer.Lexer) source.Status {
	c, stat := lex.AdvanceChar() // should be first quotation mark
	line, char := lex.GetLineChar()
	openQuoteCharNum := char - 1 // char number of '"'
	// check for leading quotation mark
	if stat.NotOk() || c != '"' {
		return determineInitialFailForAnalyzeString(lex, stat, c)
	}

	stat, content, charsRead := getStringContent(lex)
	if stat.IsOk() {
		updatedContent, ok, index := updateEscape(content, true)
		if !ok {
			lex.SetLineChar(line, char+index+1)
			lexError(lex, errors.IllegalEscape)
			return source.Bad
		}

		token := token.StringValue.
			Make().
			AddValue(updatedContent).
			SetLength(1+charsRead). // +1 to account for leading quotation mark
			SetLineChar(line, openQuoteCharNum)
		lex.PushToken(token)
	} else {
		statError(lex, stat)
	}

	lex.SetLineChar(line, char+charsRead)

	return stat
}

// resolveType returns the corresponding keyword type if the argument for `s` represents a
// keyword's string; otherwise, it returns token.Id
func resolveType(s string) token.TokenType {
	if keyType, found := keywords[s]; found {
		return keyType
	}
	return token.Id
}

var idRegexClassRaw = `[a-z][a-zA-Z0-9']`
var typeIdRegexClassRaw = `[A-Z][a-zA-Z0-9']`
var typeIdRegex = regexp.MustCompile(typeIdRegexClassRaw + `*`)

func matchId(lex *lexer.Lexer, src string) (res string, length int, stat source.Status) {
	stat = source.Ok
	res, length = lexer.RegexMatch(infixIdRegex, src)
	if length > 0 {
		return
	}

	res, length = lexer.RegexMatch(infixSymbolRegex, src)
	if length < 1 {
		stat = source.Bad
		statError(lex, stat)
		return
	}
	return
}

func analyzeIdentifier(lex *lexer.Lexer) source.Status {
	line, char := lex.GetLineChar()
	src, stat := source.GetSourceSlice(lex, line, char, -1)
	if stat.NotOk() {
		statError(lex, stat)
		return stat
	}

	var ty token.TokenType

	res, length := lexer.RegexMatch(typeIdRegex, src)
	if length > 0 {
		ty = token.TypeId
	} else {
		res, length, stat = matchId(lex, src)
		if strings.ContainsRune(res, '_') {
			ty = token.Infixed
		} else {
			ty = resolveType(res) // id or some keyword
		}
	}

	// add token
	tok := ty.Make().
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
	case identifier:
		fallthrough
	case symbol:
		return analyzeSymbol(lex)
	case char:
		return analyzeChar(lex)
	case stringClass:
		return analyzeString(lex)
	case underscore:
		return analyzeUnderscore(lex)
	case comment:
		return analyzeComment(lex)
	}

	e := errors.Lex(lex, errors.UnknownSymbol)
	lex.AddError(e)
	return source.Bad
}

// adjusts line and char to start of line if char > number of chars in current line. If lexer is at
// the EOF, then (-1, -1) is returned instead of the line and char number
func syncLineChar(lex *lexer.Lexer) (line, char int) {
	line, char = lex.GetLineChar()
	if char < 1 {
		panic("illegal char value")
	}

	sourceLine, stat := lex.SourceLine(line)
	if !stat.IsOk() {
		return -1, -1
	}

	if len(sourceLine) == char && sourceLine[char-1] == '\n' {
		lex.AdvanceLine()
		line, char = lex.GetLineChar()
	}
	return
}

func newline(lex *lexer.Lexer) {
	line, char := syncLineChar(lex)
	isNewline := line != 1 && char == 1
	if !isNewline {
		return
	}

	newlineToken := token.Newline.Make().AddValue("\n").SetLength(0).SetLineChar(line, char)
	lex.PushToken(newlineToken)
}

func analyze(lex *lexer.Lexer) source.Status {
	newline(lex)

	source.SkipWhitespace(lex)

	// get next char
	c, eof := lex.Peek()
	if eof {
		return source.Eof
	} else if c == '\n' {
		return lex.AdvanceLine()
	}

	// use char to determine what class new token will belong to
	class, e := determineClass(lex, c)
	if e != nil {
		lex.AddError(e)
		return source.Bad
	}

	// use class information to get token
	return class.analyze(lex)
}

var lexerWhitespace = regexp.MustCompile(`(\t| )+`)

func RunLexer(lex *lexer.Lexer) (tokens []itoken.Token, es []error) {
	// keep reading tokens until stat is not Ok
	stat := analyze(lex)
	for stat.IsOk() {
		stat = analyze(lex)
	}

	return lex.GetTokens(), lex.GetErrors()
}

func CastTokens(ts []itoken.Token) []token.Token {
	out := make([]token.Token, len(ts))
	for i, tok := range ts {
		out[i] = tok.(token.Token)
	}
	return out
}

func defaultTrimRight(_ *lexer.Documentor, comment string) string {
	return strings.TrimRight(comment, " \t")
}

func NewLexer(path string) *lexer.Lexer {
	lex, e := lexer.Lex(path, lexerWhitespace)
	if e != nil {
		errors.PrintErrors(errors.MkSystemError(e.Error()))
		return nil
	}

	lex.Documentor = lexer.MakeDocumentor(defaultTrimRight)
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

func NewStringLexer(stringFromPath, input string) *lexer.Lexer {
	lex, e := lexer.Initialize(stringFromPath, []byte(input), lexerWhitespace)
	if e != nil {
		errors.PrintErrors(errors.MkSystemError(e.Error()))
		return nil
	}

	lex.Documentor = lexer.MakeDocumentor(defaultTrimRight)

	return lex
}

func LexPath(path string) ([]itoken.Token, []error) {
	lex, e := lexer.Lex(path, lexerWhitespace)
	if e != nil {
		return nil, []error{errors.MkSystemError(e.Error())}
	}

	return RunLexer(lex)
}
