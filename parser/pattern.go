package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/util"
)

/*
patternC ::= constructor
pattern	 ::= patternC
             | literal
             | funcName
             | pattern pattern
						 | '(' pattern ')'
*/

var patternC__constructor_r = parser.
	Get(rewrapReduction(PatternC)).
	From(Constructor)

// == Note About Pattern Nodes ================================================
// pattern nodes should always create a tree like so (where clf stands for an
// arbitrary pattern from a single constructor, literal, or funcName):
//          /\
//         /\ clf
//        /\ clf
//      ... clf
//      /\
//   clf  clf
// this is because no other rule has any of the following reductions:
//	x	::= pattern constructor
//	y ::= pattern literal
//	z ::= pattern funcName
// thus, the reduction
//	pattern ::= pattern pattern
// will always appear in the following way
//	(1.) stack = .., pattern, clf				[premise]
//	(2.) stack = .., pattern, pattern		[pattern ::= constructor | literal | funcName]
//	(3.) stack = .., pattern						[pattern ::= pattern pattern]

// Returns a list from an expression.
//
// Case 1:
//
//	ex is an expr.Application=(a b c ..) => expr.List=[a,b,c,..]
//
// Case 2:
//
//	ex is anything else => ex
func linearizeExpression(ex expr.Expression[token.Token]) expr.Expression[token.Token] {
	app, isApp := ex.(expr.Application[token.Token])
	if !isApp { // this checks for "Case 2"
		return ex
	}

	// everything below is "Case 1"

	// left and right elements of application
	var left, right expr.Expression[token.Token]
	// a sub-application w/in application `app`
	var subApp expr.Application[token.Token] = app
	// true iff there are more sub-applications to linearize
	var ok bool = true

	// create a expression list buffer to store elements of application
	buff := make(expr.List[token.Token], 0, 32)

	// add elements in breadth-first order
	for ok {
		left, right = subApp.Split()
		buff = append(buff, right)
		// if ok == true, then left will be updated in next iteration;
		// else left is added as the final element
		subApp, ok = left.(expr.Application[token.Token])
	}

	// add final element
	buff = append(buff, left)

	// breadth-first order of the application tree orders the elements in reverse
	// order; so, reverse the order of the buffer (call to util.Reverse returns a
	// perfect-fit slice)
	return expr.List[token.Token](util.Reverse(buff))
}

var pattern__patternC_r = parser.
	Get(patternCAsPatternReduction).
	From(PatternC)

var pattern__literal_r = parser.
	Get(literalAsPatternReduction).
	From(Literal)

var pattern__funcName_r = parser.
	Get(funcNameAsPatternReduction).
	From(FuncName)

var pattern__pattern_pattern_r = parser.
	Get(applyPatternsReduction).
	From(Pattern, Pattern)

var pattern__enclosed_r = parser.Get(parenEnclosedProduction).From(LeftParen, Pattern, RightParen)

func applyPatternsReduction(nodes ...ast.Ast) ast.Ast {
	const leftIndex, rightIndex int = 0, 1
	left := nodes[leftIndex].(SomeExpression).Expression.(expr.List[token.Token])
	right := nodes[rightIndex].(SomeExpression).Expression.(expr.List[token.Token])
	return SomeExpression{
		Pattern,
		append(left, right...),
	}
}

func funcNameAsPatternReduction(nodes ...ast.Ast) ast.Ast {
	const nameIndex int = 0
	nameToken := nodes[nameIndex].(Node).Token
	name := expr.Const[token.Token]{Name: nameToken}
	pattern := expr.List[token.Token]{name}
	return SomeExpression{Pattern, pattern}
}

func literalAsPatternReduction(nodes ...ast.Ast) ast.Ast {
	const literalIndex int = 0
	literal := nodes[literalIndex].(LiteralNode).Expression
	pattern := expr.List[token.Token]{literal}
	return SomeExpression{Pattern, pattern}
}

func patternCAsPatternReduction(nodes ...ast.Ast) ast.Ast {
	const patternCIndex int = 0
	// patternC == constructor.(BinaryRecursiveNode).
	//		UpdateType(constructor.(BinaryRecursiveNode).NodeType(), PatternC)
	constructorExpr := constructorToExpression(nodes[patternCIndex])
	pattern := expr.List[token.Token]{constructorExpr}
	return SomeExpression{Pattern, pattern}
}
