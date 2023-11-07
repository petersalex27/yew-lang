// =============================================================================
// Author-Date: Alex Peters - 2023
//
// Content: 
// utility functions that are common to multiple files in this package
//
// Notes: -
// =============================================================================

package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

// prefix for compiler generated parameter names
const parameterStringPrefix string = "$p"

// prefix for compiler generated data member variables
//const memberStringPrefix string = "$m"

// prefix for compiler generated type names
//const typeStringPrefix string = "$t"

// indexes must be sorted from low to high
// func sliceRule(rule func(...ast.Ast) ast.Ast, indexes ...int) func(nodes ...ast.Ast) ast.Ast {
// 	n := len(indexes)
// 	return func(nodes ...ast.Ast) ast.Ast {
// 		buff := make([]ast.Ast, n)
// 		for i, index := range indexes {
// 			buff[i] = nodes[index]
// 		}
// 		return rule(buff...)
// 	}
// }

// type assertion/cast from ast.Ast to token
func GetToken(a ast.Ast) token.Token {
	tmp, _ := a.(ast.Token)
	tok, _ := tmp.Token.(token.Token)
	return tok
}

// returns first element of nodes, ignoring remaining nodes
func grabInitialProduction(nodes ...ast.Ast) ast.Ast {
	const initialIndex int = 0
	return nodes[initialIndex]
}

// a ::= LeftParen a RightParen
func parenEnclosedProduction(nodes ...ast.Ast) ast.Ast {
	return nodes[1]
}

// true iff two tokens a, b have the same line, char, type, and value
func EqualsToken[T itoken.Token](a, b T) bool {
	lineA, charA := a.GetLineChar()
	lineB, charB := b.GetLineChar()
	tyA, tyB := a.GetType(), b.GetType()
	valA, valB := a.GetValue(), b.GetValue()
	return lineA == lineB &&
		charA == charB &&
		tyA == tyB &&
		valA == valB
}

// generates a function that takes a single node of the passed handle at index 
// `at` and passes it to the production function `rule`
func monoSelect(rule func(nodes ...ast.Ast) ast.Ast, at int) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return rule(nodes[at])
	}
}
