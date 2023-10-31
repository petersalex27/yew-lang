package errors

import (
	"fmt"

	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/errors"
	"github.com/petersalex27/yew-packages/source"
)

const ParserErrorHead string = " Error (Syntax): "

const ParserErrorType string = "Syntax"

const ExpectedNamedAnnotation string = "expected named annotation"

const UnexpectedToken string = "unexpected token"

func Parser(src source.StaticSource, errorRange []token.Token, msg string, msgArgs ...any) errors.Err {
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
		last.line, last.char = errorRange[len(errorRange)-1].GetLineChar()
		ln := errorRange[len(errorRange)-1].GetLength()
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
