// =============================================================================
// Author-Date: Alex Peters - 2023
//
// Content:
// constructor production rules and functions
//
// Grammar:
//
//	constructor ::= TYPE_ID
//									| typeDecl
//									| constructor name
//									| constructor constructor
//									| '(' constructor ')'
//									| ( '(' ) constructor INDENT(_)
//
// Notes: -
// =============================================================================
package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

// =============================================================================
// production rules
// =============================================================================

var constructor__TypeId_r = parser.
	Get(constructorSingleProduction).From(TypeId)

var constructor__typeDecl_r = parser.
	Get(typeDeclToConstructor).From(TypeDecl)

var constructor__constructor_name_r = parser.
	Get(constructorBinaryProduction).From(Constructor, Name)

var constructor__constructor_constructor_r = parser.
	Get(constructorBinaryProduction).From(Constructor, Constructor)

var constructor__enclosed_r = parser.
	Get(parenEnclosedProduction).From(LeftParen, Constructor, RightParen)

var constructor__LeftParen_constructor_Indent_r = parser.
	Get(grabInitialProduction).
	When(LeftParen).
	From(Constructor, Indent)

// =============================================================================
// production functions
// =============================================================================

// turns a token into a constructor
var constructorSingleProduction = giveTypeToTokenProductionGen(Constructor)

// takes two constructors and produces a single constructor
var constructorBinaryProduction = simpleBinaryNodeRule(Constructor)

// =============================================================================
// utils
// =============================================================================

func constructorCast(node ast.Ast) BinaryRecursiveNode {
	return node.(BinaryRecursiveNode)
}

func constructorToExpression(constructorNode ast.Ast) expr.Expression[token.Token] {
	exprPtr := new(expr.Expression[token.Token])
	*exprPtr = nil

	toExpression := func(tok itoken.Token) {
		right := expr.Const[token.Token]{Name: tok.(token.Token)}
		if *exprPtr == nil {
			*exprPtr = right
		} else {
			*exprPtr = expr.Apply[token.Token](*exprPtr, right)
		}
	}

	getConstructor(constructorNode).InOrderTraversal(toExpression)
	return *exprPtr
}

var getConstructor = getBinaryRecursiveNode

func typeDeclToConstructor(nodes ...ast.Ast) ast.Ast {
	return getConstructor(nodes[0]).UpdateType(TypeDecl, Constructor)
}
