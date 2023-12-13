package typesexport

import (
	internal "github.com/petersalex27/yew-lang/parser/internal"
	"github.com/petersalex27/yew-lang/token"
)

func initializeTypeExport(p *internal.Parser) bool {
	p.DropShiftReduce(typeExportInitProductions)
	return !p.Panicking
}

// Shifts symbolic export tokens for "export all type constructors" and "export
// no type constructors". Returns the number of tokens shifted
func shiftSymbolicExport(p *internal.Parser) uint {
	if p.LookAhead(token.DotDot) {
		// NOTE: new lines are NOT permitted in an "export all constructors" decl
		const numberOfShifts uint = 2
		// shift '..' and ')'
		p.ShiftN(numberOfShifts)
		return numberOfShifts
	} else if p.LookAhead(token.RightParen) {
		// NOTE: new lines are NOT permitted
		const numberOfShifts uint = 1
		p.Shift()
		return numberOfShifts
	}

	return 0
}

// parses type export w/ explicitly exported constructors
func parseExplicitConsExport(p *internal.Parser) bool {
	end := false
	for !end && !p.Panicking {
		// export type constructor
		p.DropShiftReduce(typeExportProductions)
		_, end = p.ParseListComma(typeExportCommaProductions)
	}

	// shift right paren
	p.Reduce(typeExportEndProductions)
	return !p.Panicking
}

// parses a type export that has some kind of information attached:
//   - (1) symbolic export of all type constructors for some type
//   - (2) symbolic export of abstract type (i.e., export of type w/o constructors)
//   - (3) explicit export of one or more type constructors for some type
func parseInformedTypeExport(p *internal.Parser) bool {
	if shiftSymbolicExport(p) != 0 {
		return true
	}

	// initialize type export
	p.DropShiftReduce(typeExportInitProductions)
	// parse member list
	return parseExplicitConsExport(p)
}

// Parses a type export list.
func Parse(p *internal.Parser) bool {
	p.DropNewlines()

	// is export an implicit abstract type export?
	if !p.LookAhead(token.LeftParen) {
		return true // yes, it is
	}

	p.Shift() // shift left paren
	// export type that has more export information attached
	return parseInformedTypeExport(p)
}
