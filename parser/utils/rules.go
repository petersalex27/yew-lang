// =============================================================================
// Author-Date: Alex Peters - November 29, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package utils

import (
	"github.com/petersalex27/yew-packages/parser/ast"
)

// just returns first node
func Id(nodes ...ast.Ast) ast.Ast {
	return nodes[0]
}