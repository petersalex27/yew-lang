package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser/ast"
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

func GetToken(a ast.Ast) token.Token {
	tmp, _ := a.(ast.Token)
	tok, _ := tmp.Token.(token.Token)
	return tok
}

// a <- LeftParen a RightParen
func parenEnclosedReduction(nodes ...ast.Ast) ast.Ast {
	return nodes[1]
}

func EqualsToken(a, b token.Token) bool {
	lineA, charA := a.GetLineChar()
	lineB, charB := b.GetLineChar()
	tyA, tyB := a.GetType(), b.GetType()
	valA, valB := a.GetValue(), b.GetValue()
	return lineA == lineB &&
		charA == charB &&
		tyA == tyB &&
		valA == valB
}

func monoSelect(rule func(nodes ...ast.Ast) ast.Ast, at int) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		return rule(nodes[at])
	}
}
