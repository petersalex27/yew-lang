package errors

import (
	"fmt"
	"strconv"

	"github.com/petersalex27/yew-packages/lexer"
	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/errors"
)

const LexErrorHead string = " Error (Lexical Analysis): "

const UnknownSymbol string = "unknown symbol"
const UnexpectedSymbol string = "unexpected symbol"
const IllegalEscape string = "illegal escape sequence"
const IllegalChar string = "illegal character literal"

func Lex(lex *lexer.Lexer, msg string, args ...any) errors.Err {
	line, char := lex.GetLineChar()
	path := lex.GetPath()
	msg = fmt.Sprintf(msg, args...)
	src, stat := lex.SourceLine(line)
	format := "tflcm"
	if stat.Is(source.Ok) {
		format = format + "s"
	}
	return errors.Ferr(format, "Lexical Analysis", path, line, char, msg, src)
}

func LazyLex(msg string, args ...any) LazyErrorFn {
	return func(lex *lexer.Lexer) errors.Err {
		return Lex(lex, msg, args...)
	}
}

type LazyErrorFn func(*lexer.Lexer) errors.Err

func makePadding(n int) string {
	bs := make([]byte, n)
	for i := range bs {
		bs[i] = ' '
	}
	return string(bs)
}

func getLineNumber(line, lastLine int) string {
	if lastLine < line {
		panic("illegal arguments: lastLine < line")
	}
	lastLineLen := len(strconv.Itoa(lastLine))
	bs := make([]byte, 0, lastLineLen)
	lineStr := strconv.Itoa(line)
	for _, r := range lineStr {
		bs = append(bs, byte(r))
	}

	for len(bs) < lastLineLen {
		bs = append(bs, ' ')
	}
	return string(bs)
}

func MessageFromStatus(stat source.Status) string {
	switch stat {
	case source.Eof:
		return "unexpected end of file"
	case source.Eol:
		return "unexpected end of line"
	case source.Bad:
		return "bad input"
	default:
		return "bug in messageFromStatus(status=#" + strconv.Itoa(int(stat)) + ")"
	}
}