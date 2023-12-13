// =============================================================================
// Author-Date: Alex Peters - December 09, 2023
// =============================================================================
package declaration

import "github.com/petersalex27/yew-lang/parser/internal"

func Parse(p *internal.Parser) {
	p.Shift()
	typing.Parse(p)
	p.Reduce(declarationProductions)
}

