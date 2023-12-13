package export

import (
	"github.com/petersalex27/yew-packages/parser/ast"
)

// return first arg, drop all others
func produceExportDrop(nodes ...ast.Ast) ast.Ast {
	const moduleIndex int = 0
	return nodes[moduleIndex]
}
