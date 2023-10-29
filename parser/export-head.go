package parser

import "github.com/petersalex27/yew-packages/parser"

/*
exportHead   ::= moduleHead '('
                  | export ','
*/

// == export reduction rules ==================================================

var exportHead__moduleDec_LeftParen_r = parser. 
	Get(rewriteModuleTypeReduction(ExportHead)). 
	From(ModuleDeclaration, LeftParen)

var exportHead__export_Comma_r = parser. 
	Get(rewriteModuleTypeReduction(ExportHead)). 
	From(ExportList, Comma)