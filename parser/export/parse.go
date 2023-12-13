// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content:
// Parses export lists and their children
// =============================================================================
package export

import (
	typesexport "github.com/petersalex27/yew-lang/parser/export/types-export"
	"github.com/petersalex27/yew-lang/parser/internal"
	"github.com/petersalex27/yew-lang/token"
)

// parses next exported thing in export list
func parseExportMember(p *internal.Parser) bool {
	switch p.Next() {
	case token.Id:
		fallthrough
	case token.Infixed:
		p.Shift()
	case token.TypeId:
		p.Shift()
		return typesexport.Parse(p)
	default:
		p.UnexpectedLookAheadToken()
		return false
	}
	return true
}

// Top level export list parsing method.
func Parse(p *internal.Parser) bool {
	ok := true
	p.DropNewlines()

	for ok && !p.Panicking {
		ok = parseExportMember(p)

		p.Reduce(exportProductions)

		_, endExportParsing := p.ParseListComma(exportCommaProductions)
		if endExportParsing {
			break
		}
	}

	return ok && !p.Panicking
}
