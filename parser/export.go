package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

/*
export        ::= exportHead ID
                  | exportHead TYPE_ID
                  | exportHead SYMBOL
                  | exportHead INFIXED
                  | exportHead          # only when L.A. is ')'
*/

// == export reduction rules ==================================================

var export__exportHead_Id_r = parser. 
	Get(exportAppendReduction). 
	From(ExportHead, Id)

var export__exportHead_TypeId_r = parser. 
	Get(exportAppendReduction). 
	From(ExportHead, TypeId)

var export__exportHead_Symbol_r = parser. 
	Get(exportAppendReduction). 
	From(ExportHead, Symbol)

var export__exportHead_Infixed_r = parser. 
	Get(exportAppendReduction). 
	From(ExportHead, Infixed)

var export__exportHead_r = parser.
	Get(rewriteModuleTypeReduction(ExportList)). 
	From(ExportHead)

// == export reduction functions ==============================================

func exportAppendReduction(nodes ...ast.Ast) ast.Ast {
	const exportHeadIndex, symIndex int = 0, 1
	exportHead := nodes[exportHeadIndex].(ModuleNode)
	sym := GetToken(nodes[symIndex])
	exportHead.Type = ExportList
	exportHead.exportList = append(exportHead.exportList, sym)
	return exportHead
}