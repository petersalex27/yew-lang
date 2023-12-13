// =============================================================================
// Author-Date: Alex Peters - November 27, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package annotation

import (
	internal "github.com/petersalex27/yew-lang/parser/internal"
	"github.com/petersalex27/yew-lang/token"
)

// parses all leading annotations
func parseAnnotations(p *internal.Parser) {
	p.DropNewlines()

	// keep parsing annotation until no more annotations are found
	for p.LookAhead(token.Annotation) {
		p.Shift() // shift annotation
		p.Reduce(annotationProductions)
		p.DropNewlines()
	}
}

func ImportAnnotationParse(p *internal.Parser) bool {
	parseAnnotations(p)
	
	return true
}

// parses top level annotations until module token is found--leaves module token
// as next unseen token
//
// returns true iff zero or more annotation tokens are followed by a module
// token
func ModuleAnnotationParse(p *internal.Parser) bool {
	parseAnnotations(p)

	return p.LookAhead(token.Module)
} 