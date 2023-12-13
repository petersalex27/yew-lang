// =============================================================================
// Author-Date: Alex Peters - November 29, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package assign

import (
	//"github.com/petersalex27/yew-lang/parser/indent"
	"github.com/petersalex27/yew-lang/parser/internal"
	"github.com/petersalex27/yew-lang/token"
)

func Parse(p *internal.Parser) bool {
	if !p.LookAhead(token.Assign) {
		// TODO: report error
		return false
	}
	
	p.Shift()
	//indent.Parse(p)
	return true
}