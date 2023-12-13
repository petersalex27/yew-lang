package errors

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/errors"
	"github.com/petersalex27/yew-packages/source"
)

const ParserErrorHead string = " Error (Syntax): "

const ParserErrorType string = "Syntax"

const ExpectedNamedAnnotation string = "expected named annotation"

const UnexpectedToken string = "unexpected token"

const UnexpectedError string = "unexpected error"

const UnexpectedEndOfTokens string = "unexpected end of tokens"

// pretty prints source code
func PrettyPrint(tokenStream []token.Token) string {
	// TODO: finish implementing

	panic("TODO: implement")
}

func appendLine(lines []string, line int, src source.StaticSource) ([]string, source.Status) {
	origin, stat := src.SourceLine(line)
	if stat.NotOk() {
		// cannot reconstruct source lines
		return nil, stat
	}
	lines = append(lines, origin)
	return lines, stat
}

func buildMultiLineSource(src source.StaticSource, fromLine, toLine int) ([]string, source.Status) {
	// fromLine > toLine: order does not follow order of line numbers
	// fromLine < 1: not a line number
	// toLine < 1: not a line number
	if fromLine > toLine || fromLine < 1 || toLine < 1 {
		// cannot build source's requested lines
		return nil, source.BadLineNumber
	}

	// `out` holds slice w/ length equal to nubmer of lines in subsection of source
	out := make([]string, 0, toLine-(fromLine-1))

	// loop until current line is last line that should be included
	currentLine := fromLine
	var stat source.Status = source.Ok
	for ; currentLine < toLine && stat.IsOk(); currentLine++ {
		out, stat = appendLine(out, currentLine, src)
	}
	if stat.NotOk() {
		return nil, stat
	}

	// include last line (this is here to handle case when toLine == math.MaxInt,
	// else loop counter could overflow)
	out, stat = appendLine(out, currentLine, src)
	if stat.NotOk() {
		return nil, stat
	}
	return out, source.Ok
}

// replace source code from fromToken to (inclusive) toToken
func MakeSuggestion(src source.StaticSource, fromToken, toToken token.Token, replacement string) string {
	fromLine, fromChar := fromToken.GetLineChar()
	toLine, toCharStart := toToken.GetLineChar()
	// last char of toToken
	toChar := toCharStart + toToken.GetLength() - 1

	// create string from lines (or line if fromLine == toLine)
	original, stat := buildMultiLineSource(src, fromLine, toLine)

	// cannot make suggestion
	if stat.NotOk() {
		return ""
	}

	lastLine := original[len(original)-1]

	firstPart := strings.Join(original[:len(original)-1], "\n")
	toChar = len(firstPart) + toChar
	together := firstPart + "\n" + lastLine

	// insert replacement
	return together[:fromChar-1] + replacement + together[toChar-1:]
}

//func ParserSuggest(src source.StaticSource, errorRange [2]token.Token, )

// assumes tokens in `errorRange` are in the order they appear
func Parser(src source.StaticSource, errorRange [2]token.Token, msg string, msgArgs ...any) errors.Err {
	type ranged struct{ line, char int }
	var first, last ranged
	var errorRangeBreadth int

	msg = fmt.Sprintf(msg, msgArgs...)

	format := "tfm" // <error-Type><File-path><Message>

	// there will be at least three args and seven at most
	var arguments []any = make([]any, 3, 7)
	// three guarenteed arguments: type, path, and message
	arguments[0], arguments[1], arguments[2] = ParserErrorType, src.GetPath(), msg

	if len(errorRange) > 0 {
		format = format + "lc" // add line and char format spec.
		first.line, first.char = errorRange[0].GetLineChar()
		// add args
		arguments = append(arguments, first.line, first.char)
		if len(errorRange) == 1 {
			tokenLength := errorRange[0].GetLength()
			format = format + "r"
			arguments = append(arguments, tokenLength)
		}
	}

	if len(errorRange) > 1 {
		last.line, last.char = errorRange[1].GetLineChar()
		ln := errorRange[1].GetLength()
		if last.line == first.line {
			format = format + "r"
			//   _  4    <- lengths
			// x ab cdef	 : rng len = 7 = (6 + 4) - 3
			//   3  6    <- char nums
			errorRangeBreadth = (last.char + ln) - first.char
			// add arg
			arguments = append(arguments, errorRangeBreadth)
		} // else too many lines to display range on one line
	}

	srcLine, stat := src.SourceLine(first.line)

	if stat.Is(source.Ok) {
		format = format + "s"
		// add arg
		arguments = append(arguments, srcLine)
	}

	return errors.Ferr(format, arguments...)
}

// expected -> [actual1, actual2, .., actualN]
var didYouMeanMap = map[token.TokenType][]token.TokenType{
	token.LeftParen:    {token.LeftBrace, token.LeftBracket},
	token.RightParen:   {token.RightBrace, token.RightBracket},
	token.LeftBrace:    {token.LeftParen, token.LeftBracket},
	token.RightBrace:   {token.RightParen, token.RightBracket},
	token.LeftBracket:  {token.LeftParen, token.LeftBrace},
	token.RightBracket: {token.RightParen, token.RightBrace},
}

func suggestion(expected token.TokenType, actual token.TokenType) {
	panic("TODO: implement")
}
