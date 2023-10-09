package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"yew.lang/main/token"
)

/*
arrayValHead  ::= '[' expr
                  | arrayValHead ',' expr
									| litArrHead ',' expr

array         ::= arrayValHead ',' ']'
                  | arrayValHead ']'
*/

type ArrayNode struct {
	isLiteral bool
	elems     *expr.List[token.Token]
}

func buildArray(seq NodeSequence, elems *expr.List[token.Token]) {

}

var arrayValHeadMonoReduction = monoSelect(createNodeSequenceRule(ArrayValHead), 1)

var arrayValHeadRevConsReduction = binarySelect(reverseConsRule(ArrayValHead), 0, 2)

var arrayReduction = rewrapNodeSequenceRule(Array)

var litArrHeadMonoReduction = monoSelect(createNodeSequenceRule(LitArrHead), 1)

var litArrHeadRevConsReduction = binarySelect(reverseConsRule(LitArrHead), 0, 2)

var literalArrayReduction = rewrapNodeSequenceRule(LiteralArray)

// arrayValHead <- LeftBracket expr
var arrayValHead__LeftBracket_expr_r = parser.
	Get(arrayValHeadMonoReduction).From(LeftBracket, Expr)

// arrayValHead <- arrayValHead Comma expr
var arrayValHead__arrayValHead_Comma_expr_r = parser.
	Get(arrayValHeadRevConsReduction).From(ArrayValHead, Comma, Expr)

// arrayValHead <- litArrHead Comma expr
var arrayValHead__litArrHead_Comma_expr_r = parser.
	Get(arrayValHeadRevConsReduction).From(LitArrHead, Comma, Expr)

// array <- arrayValHead Comma RightBracket
var array__arrayValHead_Comma_RightBracket_r = parser.
	Get(arrayReduction).From(ArrayValHead, Comma, RightBracket)

// array <- arrayValHead RightBracket
var array__arrayValHead_RightBracket_r = parser.
	Get(arrayReduction).From(ArrayValHead, RightBracket)

/*
litArrHead    ::= '[' literal
                  | litArrHead ',' literal

literalArray  ::= literalArray ',' ']'
                  | literalArray ']'
*/

// litArrHead <- LeftBracket expr
var litArrHead__LeftBracket_expr_r = parser.
	Get(litArrHeadMonoReduction).From(LeftBracket, Expr)

// litArrHead <- litArrHead Comma expr
var litArrHead__litArrHead_Comma_expr_r = parser.
	Get(litArrHeadRevConsReduction).From(LitArrHead, Comma, Expr)

// literalArray <- litArrHead Comma RightBracket
var literalArray__litArrHead_Comma_RightBracket_r = parser.
	Get(literalArrayReduction).From(LitArrHead, Comma, RightBracket)

// literalArray <- litArrHead RightBracket
var literalArray__litArrHead_RightBracket_r = parser.
	Get(literalArrayReduction).From(LitArrHead, RightBracket)
