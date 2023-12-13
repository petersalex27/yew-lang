// =============================================================================
// Author-Date: Alex Peters - November 27, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package imports

import (
	"github.com/petersalex27/yew-lang/parser/internal"
	"github.com/petersalex27/yew-lang/token"
)

// either shifts 'qualified' or does nothing
func parseOptionalQualifier(p *internal.Parser) {
	// shift 'qualified' keyword
	if p.LookAhead(token.Qualified) {
		p.Shift()
	} // else shift nothing
}

// Id token is shifted onto parse stack iff function would return true
//
// returns true iff next token is an Id.
func shiftImportName(p *internal.Parser) bool {
	shift := p.LookAhead(token.Id)
	if !shift {
		p.UnexpectedLookAheadToken()
	} else {
		p.Shift()
	}

	return shift
}

// If the next token is not 'from', then don't reduce but no errors occurred. Otherwise, require the
// next two tokens as these tokens in this order: 'from', STRING_VAL; if this is found, both return
// values are true. Otherwise, return both return values as false
func shiftOptionalFrom(p *internal.Parser) (reduce, ok bool) {
	ok = true
	reduce = p.LookAhead(token.From)
	if !reduce {
		return
	}
	p.Shift() // shift 'from'
	
	reduce = p.LookAhead(token.StringValue)
	ok = reduce
	if !reduce {
		p.UnexpectedLookAheadToken()
		return
	}
	p.Shift() // shift path (i.e., shift string value)
	return
}

func parseOptionalFrom(p *internal.Parser) bool {
	reduce, noError := shiftOptionalFrom(p)
	if !reduce { 
		// don't do reduction
		return noError
	}

	p.Reduce(importElemFromProductions)
	return !p.Panicking
}

func parseElem(p *internal.Parser) bool {
	// no newline allowed after 'qualified'
	parseOptionalQualifier(p)

	// required name parse
	if !shiftImportName(p) {
		return false
	}
	p.Reduce(initialImportElemProductions)
	
	return parseOptionalFrom(p) 
}

// parses indent defined group of imports
//
// returns number of import elems parsed and whether that number is useful (consequence of parser
// having no errors)
func parseElemGroup(p *internal.Parser) (numImports int, success bool) {
	numImports, success = 0, true
	if p.DropN(1, token.LeftBrace) != 1 {
		p.UnexpectedLookAheadToken()
		return 0, false 
	}

	// drop newlines following '{'
	p.DropNewlines()
	
	// successful exit is when '}' token is found
	for {
		if !parseElem(p) {
			return 0, false
		}

		numImports++
		// drop newlines and save number dropped
		newlinesDropped := p.DropN(-1, token.Newline)
		
		if p.DropN(1, token.RightBrace) == 1 { 
			// number of newlines dropped is irrelevant
			break
		}

		// import elements must be newline separated, test here
		if newlinesDropped < 1 {
			p.UnexpectedLookAheadToken()
			return 0, false
		}
	}

	return
}

func reduceImports(numImportElems int, p *internal.Parser) bool {
	for i := 0; i < numImportElems; i++ {
		p.Reduce(importElemsProductions)
		if p.Panicking {
			return false
		}
	}

	p.Reduce(finishImportsProductions)
	return !p.Panicking
}

// parses an import group
//
// this can come in the form of a single import:
//
//	import packageName from "optional/path" in
//
// this can also come in the form of a proper import group:
//
//	import {
//		qualified packageName from "optional/path"
//		qualified additionalPackage
//		maybeMorePackages
//		finalPackage from "arbitrary/path" 
//	} in
func Parse(p *internal.Parser) bool {
	// parse 'import' keyword
	if !p.LookAhead(token.Import) {
		p.UnexpectedLookAheadToken()
		return false
	}
	p.Shift()

	if !p.LookAhead(token.LeftBrace) {
		return parseElem(p) && reduceImports(1, p)
	}

	numImports, success := parseElemGroup(p)
	return success && reduceImports(numImports, p)
}
