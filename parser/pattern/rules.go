// =============================================================================
// Author-Date: Alex Peters - December 11, 2023
//
// Grammar:
//
//	pattern_inner ::= ID
//										| INFIXED
//	                  | '_'
//	                  | INT_VALUE
//	                  | CHAR_VALUE
//	                  | STRING_VALUE
//	                  | FLOAT_VALUE
//	                  | pattern pattern
//	                  | TYPE_ID { pattern }
//	pattern_mid   ::= pattern_inner
//	                  | '(' pattern_inner ')'
//	pattern       ::= pattern_mid
//	                  | '(' pattern ',' pattern { ',' pattern } [ ',' ] ')'
//
// =============================================================================
package pattern

import (
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

func produceFromToken(nodes ...ast.Ast) ast.Ast {
	const tokenIndex int = 0
	tok := nodes[tokenIndex].(ast.Token).Token.(token.Token)
	toks := []token.Token{tok}
	return &PatternNode{toks, nil}
}

func produceAppend(nodes ...ast.Ast) ast.Ast {
	const leftIndex, rightIndex int = 0, 1
	left := nodes[leftIndex].(*PatternNode)
	right := nodes[rightIndex].(*PatternNode)
	if len(*right) == 1 {
		*left = append(*left, *right...)
	}
}

// =============================================================================
// rules
// =============================================================================

var fromIdRule = parser.Get(produceAppendToken).From(Id)

var infixedRule = parser.Get(produceAppendToken).From(Infixed)

var wildcardRule = parser.Get(produceAppendToken).From(Wildcard)

var intRule = parser.Get(produceAppendToken).From(IntValue)

var charRule = parser.Get(produceAppendToken).From(CharValue)

var stringRule = parser.Get(produceAppendToken).From(StringValue)

var floatRule = parser.Get(produceAppendToken).From(FloatValue)

var constantConsRule = parser.Get(produceAppendToken).From(TypeId)

var applicationRule = parser.Get(produceAppend).From(Pattern, Pattern)

var tupleRule = parser.Get(produceConsTuple).From(Pattern, Comma, Pattern)
