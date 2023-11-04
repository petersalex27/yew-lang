package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
exportDone    ::= export ')'

export        ::= exportHead ID
                  | exportHead SYMBOL
                  | exportHead INFIXED
*/

// wraps a token.Token and flags whether it exports all constructors (when
// token is a type id, that is)
type exportToken struct {
	// includes all constructors for the type referenced by (exportToken).Token
	// if true and (exportToken).Token.GetType() == TypeId. Else, means nothing
	exportAllConstructors bool
	// token refering to some entity: can be an id, symbol, infixed symbol or
	// type/constructor name
	token token.Token
}

// ============================================================================
// export production rules
// ============================================================================

var export__exportHead_Id_r = parser.
	Get(exportAppendProduction).
	From(ExportHead, Id)

var export__exportHead_Symbol_r = parser.
	Get(exportAppendProduction).
	From(ExportHead, Symbol)

var export__exportHead_Infixed_r = parser.
	Get(exportAppendProduction).
	From(ExportHead, Infixed)

// ============================================================================
// export production functions
// ============================================================================

var exportAppendProduction = someExportTypeAppendProductionGen(ExportList, false)

func someExportTypeAppendProductionGen(ty ast.Type, withAllConstructors bool) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		const someExportIndex, symIndex int = 0, 1
		exportHead := nodes[someExportIndex].(ModuleNode)
		tok := GetToken(nodes[symIndex])
		return appendExportToken(exportHead, tok, withAllConstructors, ty)
	}
}

// ============================================================================
// export related utils
// ============================================================================

// appends a token (with the option of exporting `tok` `withAllConstructors` of
// the type referenced by `tok`) and then sets the type of the module node to 
// `typeTransformation`
func appendExportToken(
	modExportList ModuleNode, 
	tok token.Token, 
	withAllConstructors bool, 
	typeTransformtion ast.Type,
) ModuleNode {
	exporting := exportToken{withAllConstructors, tok}
	modExportList.Type = typeTransformtion
	modExportList.exportList = append(modExportList.exportList, exporting)
	return modExportList
} 
