// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package module

import (
	internal "github.com/petersalex27/yew-lang/parser/internal"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
)

func parseKeyword(p *internal.Parser) bool {
	if !p.LookAhead(token.Module) {
		p.UnexpectedLookAheadToken()
		return false
	}

	p.Shift()
	return true
}

func parseModuleName(p *internal.Parser) bool {
	if !p.LookAhead(token.Id) {
		p.UnexpectedLookAheadToken()
		return false
	}

	p.Shift()
	return true
}

func parseDeclaration(p *internal.Parser) bool {
	if !parseKeyword(p) {
		return false
	}

	if !parseModuleName(p) {
		return false
	}

	return true
}

func closeExportList(p *internal.Parser) bool {
	if p.Panicking {
		return false
	}

	if !p.LookAhead(token.RightParen) {
		return false
	}

	p.Shift()
	return true
}

// shifts tokens for either optional grammar rule `module <name> ( .. )` or `module <name> ( )`
func maybeParseExportEverythingOrNothing(p *internal.Parser) {
	if p.Panicking {
		return
	}

	if p.LookAhead(token.DotDot) {
		p.Shift()
		p.Panicking = !closeExportList(p)
		return
	}

	_ = closeExportList(p)
}

// parses module definition
func Parse(p *internal.Parser) (doExport bool, success bool) {
	if !parseDeclaration(p) {
		return false, false
	}

	if p.LookAhead(token.LeftParen) {
		p.Shift()
		maybeParseExportEverythingOrNothing(p)
	} else {
		p.DropNewlines()
		if !p.LookAhead(token.Where) {
			p.UnexpectedLookAheadToken()
			return false, false
		}
	}

	p.Reduce(parseModuleProductions)

	// result of reduction
	node, ok := parser.ParseStackPeek(p.Parser)
	if !ok {
		success = false
		return
	}

	// check type of result
	nodeType := node.NodeType()
	doExport = nodeType == ExportList
	success = doExport || nodeType == ModuleDef
	return
}

func AttachExportList(p *internal.Parser) bool {
	p.DropNewlines()
	if !p.LookAhead(token.RightParen) {
		p.UnexpectedLookAheadToken()
		return false
	}
	p.Shift()

	p.DropNewlines()
	if !p.LookAhead(token.Where) {
		p.UnexpectedLookAheadToken()
		return false
	}

	p.Reduce(attachExportListProductions)
	return !p.Panicking
}
