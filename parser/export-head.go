package parser

import (
	"github.com/petersalex27/yew-packages/parser"
)

/*
exportHead    ::= moduleHead '('
                  | export ','
                  | export TYPE_ID ','
                  | export TYPE_ID '..' ','
*/

// == export reduction rules ==================================================

var exportHead__moduleDec_LeftParen_r = parser. 
	Get(rewriteModuleTypeReduction(ExportHead)). 
	From(ModuleDeclaration, LeftParen)

var exportHead__export_Comma_r = parser. 
	Get(rewriteModuleTypeReduction(ExportHead)). 
	From(ExportList, Comma)

var exportHead__export_TypeId_Comma_r = parser. 
	// when reading the export list, finding '..' informs the reader to take 
	// previous element (which must exist and must be a type id), and place all 
	// that types constructors into the export list. Effectively, 
	//		`Type ..` = `Type (Constructor1, Constructor2, Constructor3)`
	Get(someExportTypeAppendProductionGen(ExportHead, false)). 
	From(ExportList, TypeId, Comma)

var exportHead__export_TypeId_DotDot_Comma_r = parser. 
	// when reading the export list, finding '..' informs the reader to take 
	// previous element (which must exist and must be a type id), and place all 
	// that types constructors into the export list. Effectively, 
	//		`Type ..` = `Type (Constructor1, Constructor2, Constructor3)`
	Get(someExportTypeAppendProductionGen(ExportHead, true)). 
	From(ExportList, TypeId, DotDot, Comma)