// =============================================================================
// Author-Date: Alex Peters - November 27, 2023
//
// Content: grammar rules
// =============================================================================
package annotation

import (
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

func produceAnnotation(nodes ...ast.Ast) ast.Ast {
	const annotationIndex int = 0
	annotation := new(AnnotationNode)
	annotation.Token = nodes[annotationIndex].(ast.Token).Token.(token.Token)
	return annotation
}

var annotationRule = parser.Get(produceAnnotation).From(AnnotationToken)

var annotationProductions = parser.Order(annotationRule)