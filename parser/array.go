package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
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
	expr.List[token.Token]
}

func (node ArrayNode) getExpression() ExpressionNode {
	return ExpressionNode{node.List}
}

func arrayBuilder(isLiteral bool) func(nodes ...ast.Ast) ast.Ast {
	return func(nodes ...ast.Ast) ast.Ast {
		seq := nodes[0].(NodeSequence)

		arr := ArrayNode{
			isLiteral: isLiteral,
			List:      make(expr.List[token.Token], len(seq.nodes)),
		}

		for i, node := range seq.nodes {
			arr.List[i] = node.(expressionNodeTypes).getExpression().Expression
		}
		return arr
	}
}

func (node ArrayNode) Equals(a ast.Ast) bool {
	arr, ok := a.(ArrayNode)
	if !ok || arr.isLiteral != node.isLiteral {
		return false
	}

	return node.List.Equals(globalContext__.exprCxt, arr.List)
}

func (node ArrayNode) NodeType() ast.Type {
	if node.isLiteral {
		return LiteralArray
	}
	return Array
}

func (node ArrayNode) InOrderTraversal(f func(itoken.Token)) {
	for _, tok := range node.List.Collect() {
		f(tok)
	}
}

var arrayValHeadMonoReduction = monoSelect(createNodeSequenceRule(ArrayValHead), 1)

var arrayValHeadRevConsReduction = binarySelect(reverseConsRule(ArrayValHead), 0, 2)

//var arrayReduction = rewrapNodeSequenceRule(Array)

var litArrHeadMonoReduction = monoSelect(createNodeSequenceRule(LitArrHead), 1)

var litArrHeadRevConsReduction = binarySelect(reverseConsRule(LitArrHead), 0, 2)

//var literalArrayReduction = rewrapNodeSequenceRule(LiteralArray)

// arrayValHead <- LeftBracket expr
var arrayValHead__LeftBracket_expr_r = parser.
	Get(arrayValHeadMonoReduction).From(LeftBracket, Expr)

// arrayValHead <- arrayValHead Comma expr
var arrayValHead__arrayValHead_Comma_expr_r = parser.
	Get(arrayValHeadRevConsReduction).From(ArrayValHead, Comma, Expr)

// arrayValHead <- litArrHead Comma expr
var arrayValHead__litArrHead_Comma_expr_r = parser.
	Get(arrayValHeadRevConsReduction).From(LitArrHead, Comma, Expr)

// array <- arrayValHead RightBracket
var array__arrayValHead_RightBracket_r = parser.
	Get(arrayBuilder(false)).From(ArrayValHead, RightBracket)

// array <- arrayValHead Comma RightBracket
var array__arrayValHead_Comma_RightBracket_r = parser.
	Get(arrayBuilder(false)).From(ArrayValHead, Comma, RightBracket)

/*
litArrHead    ::= '[' literal
                  | litArrHead ',' literal

literalArray  ::= literalArray ',' ']'
                  | literalArray ']'
*/

// litArrHead <- LeftBracket literal
var litArrHead__LeftBracket_literal_r = parser.
	Get(litArrHeadMonoReduction).From(LeftBracket, Literal)

// litArrHead <- litArrHead Comma literal
var litArrHead__litArrHead_Comma_literal_r = parser.
	Get(litArrHeadRevConsReduction).From(LitArrHead, Comma, Literal)

// literalArray <- litArrHead RightBracket
var literalArray__litArrHead_RightBracket_r = parser.
	Get(arrayBuilder(true)).From(LitArrHead, RightBracket)

// literalArray <- litArrHead Comma RightBracket
var literalArray__litArrHead_Comma_RightBracket_r = parser.
	Get(arrayBuilder(true)).From(LitArrHead, Comma, RightBracket)
