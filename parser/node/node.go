// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content: 
// Yew AST node
// =============================================================================
package node

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
)

// Yew AST node
type Node interface {
	Visit(*inf.Context[token.Token])
	ast.Ast
}

// true iff:
//	- len(xs) == len(ys)
//	- for each x in xs, y in ys: x.Equals(y)
func NodesEquals[N Node](xs, ys []N) bool {
	if len(xs) != len(ys) {
		return false
	}

	for i, x := range xs {
		if !x.Equals(ys[i]) {
			return false
		}
	}
	return true
}