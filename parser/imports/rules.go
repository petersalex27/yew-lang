// =============================================================================
// Author-Date: Alex Peters - November 27, 2023
//
// Content: import grammar rules
// =============================================================================
package imports

import (
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/util"
)

// =============================================================================
// production functions
// =============================================================================

// creates an import for the package referred to by `nameName` and qualified by `qualifier`
func produceInitialImportElemHelper(qualifier inf.QualificationType, nameNode ast.Ast) *ImportElementNode {
	name := nameNode.(ast.Token).Token.(token.Token)
	elem := new(ImportElementNode)
	elem.Name = name
	elem.As = name // tentatively set "as name" to name
	elem.From = ""
	elem.QualificationType = qualifier
	return elem
}

// creates an import element with the default qualifier
func produceInitialImportElem(nodes ...ast.Ast) ast.Ast {
	const nameIndex int = 0
	return produceInitialImportElemHelper(inf.NameQualified, nodes[nameIndex])
}

// creates an import element with the qualified qualifier
func produceInitialQualifiedImportElem(nodes ...ast.Ast) ast.Ast {
	const _, nameIndex int = 0, 1
	return produceInitialImportElemHelper(inf.FullyQualified, nodes[nameIndex])
}

// attaches "from" information to an import element
func produceImportElemFrom(nodes ...ast.Ast) ast.Ast {
	const importElemIndex, _, stringIndex int = 0, 1, 2
	importElem := nodes[importElemIndex].(*ImportElementNode)
	from := nodes[stringIndex].(ast.Token).Token.GetValue()
	importElem.From = from
	return importElem
}

// reverse imports to match order of source code
func produceImportImports(nodes ...ast.Ast) ast.Ast {
	const _, importsIndex int = 0, 1
	imports := nodes[importsIndex].(*ImportNode)
	util.ReverseInPlace(*imports)
	return imports
}

// this production appends import elements in reverse order
func produceAttachElement(nodes ...ast.Ast) ast.Ast {
	const elemIndex, importsIndex int = 0, 1
	elem := nodes[elemIndex].(*ImportElementNode)
	imports := nodes[importsIndex].(*ImportNode)
	*imports = append(*imports, elem)
	return imports
}

func produceInitialRightImportElem(nodes ...ast.Ast) ast.Ast {
	const rightMostElemIndex int = 0
	elem := nodes[rightMostElemIndex].(*ImportElementNode)
	// initialize imports
	imports := new(ImportNode)
	*imports = make(ImportNode, 0, 1)
	*imports = append(*imports, elem)
	return imports
}

// =============================================================================
// production rules
// =============================================================================

var qualifiedInitialImportElemRule = parser.Get(produceInitialQualifiedImportElem).From(Qualified, Id)

var initialImportElemRule = parser.Get(produceInitialImportElem).From(Id)

var attachFromRule = parser.Get(produceImportElemFrom).From(ImportElement, From, StringValue)

var attachElementRule = parser.Get(produceAttachElement).From(ImportElement, ImportContext)

var initialRightElementRule = parser.Get(produceInitialRightImportElem).From(ImportElement)

var importsRule = parser.Get(produceImportImports).From(Import, ImportContext)

// =============================================================================
// production groups
// =============================================================================

var initialImportElemProductions = parser.Order(
	qualifiedInitialImportElemRule,
	initialImportElemRule,
)

var importElemFromProductions = parser.Order(attachFromRule)

var importElemsProductions = parser.Order(
	attachElementRule,
	initialRightElementRule,
)

var finishImportsProductions = parser.Order(importsRule)
