// =============================================================================
// Author-Date: Alex Peters - December 09, 2023
// =============================================================================
package typing

import (
	"github.com/petersalex27/yew-lang/parser/internal"
	//. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
)

func beforeCommaReduction(p *internal.Parser) {
	
}

func parseApplication(p *internal.Parser) {
	apply := false 
	for !p.Panicking {
		switch p.Next() {
		case token.Id:
			fallthrough
		case token.TypeId:
			fallthrough
		case token.LeftBrace:
			
		}
		if nextType == token.Id || nextType == token.TypeId {
			p.Shift()
			p.Reduce(singletonProductions)
		} else {
			return
		}

		if apply {
			p.Reduce(applicationProductions)
		}

		apply = true // more than one type on stack
	}
}

func getBindingPower(ty token.TokenType) byte {
	switch ty {
	case token.Typing:
		return 1
	case token.Comma:
		return 2
	case token.Arrow:
		return 3
	}
	return 0
}

const (
	commaBindingPower byte = iota + 2
	arrowBindingPower
)

func parseArrayType(p *internal.Parser, leftBracket token.Token) {
	for !p.Panicking {
		switch p.Next() {
			//case 
		}
	}
}

// tests: 
//		grouping >= 0
// reports error if not, else does nothing
func validateGrouping(p *internal.Parser, grouping int) {
	if grouping < 0 {
		p.Panicking = true
		p.UnexpectedLookAheadToken() // unexpected ')'
		return
	}
} 

// parses tuple type
func parseTupleType(p *internal.Parser) {
	if p.Panicking {
		return
	}
	p.RequireDropN(1, token.RightParen)
	p.ReductionLoop(loopedProductions)
	p.Reduce(closeTupleProductions)
}

// called when look-ahead is `,`
func parseTupleMember(p *internal.Parser, bindingPower byte) {
	if bindingPower > commaBindingPower {
		p.ReductionLoop(beforeCommaProductions)
	}
	p.Shift()
}

func parseDependentTyped(p *internal.Parser, bindingPower byte) {
	// number of expected left-right paren pairs
	grouping := 0

	for !p.Panicking {
		if grouping > 0 {
			p.DropNewlines()
		}
		
		parseApplication(p)
		switch p.Next() {
		case token.Arrow:
			p.Shift()
		case token.Comma:
			parseTupleMember(p, bindingPower)
		case token.LeftParen:
			p.Shift()
			grouping++
		case token.LeftBracket:
			leftBracket := p.GetAndRemoveLookAhead()
			parseArrayType(p, leftBracket)
		case token.RightParen:
			grouping--
			validateGrouping(p, grouping)
			parseTupleType(p)
		default:
			if grouping > 0 {
				// TODO: make this "unclosed paren"
				p.Panicking = true
				p.UnexpectedLookAheadToken()
				return
			}
			// parse any remaining function types
			p.ReductionAllowingEmptyLoop(beforeCommaProductions)

			// if the input up to this point has been valid, then the entire type should be a single node
			return
		}
	}
}

func Parse(p *internal.Parser) {
	if p.Panicking {
		return
	}

	reduce := false
	var productions parser.ProductionOrder

	for {
		switch p.Next() {
		case token.Forall:
			p.Shift()
			shiftTypeVariables(p)
		case token.Id:
			
		case token.TypeId:
			p.Shift()
		case token.LeftParen:

		case token.LeftBracket:
		}
	}
}