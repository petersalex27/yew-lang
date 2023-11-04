package parser

import "github.com/petersalex27/yew-packages/parser"

// ============================================================================
// exportDone production rules
// ============================================================================

var exportDone__export_RightParen_r = parser.
	Get(rewriteModuleTypeReduction(ExportDone)).
	From(ExportList, RightParen)