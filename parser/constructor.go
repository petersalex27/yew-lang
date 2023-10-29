package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

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

var constructorSingleReduction = giveTypeToTokenReductionGen(Constructor)

var constructorBinaryReduction = simpleBinaryNodeRule(Constructor)

var getConstructor = getBinaryRecursiveNode

/*
constructor   ::= TYPE_ID
									| typeDecl
                  | constructor name
                  | constructor constructor
                  | '(' constructor ')'
*/

var constructor__TypeId_r = parser.
	Get(constructorSingleReduction).From(TypeId)

func typeDeclToConstructor(nodes ...ast.Ast) ast.Ast {
	return getConstructor(nodes[0]).UpdateType(TypeDecl, Constructor)
}

var constructor__typeDecl_r = parser.
	Get(typeDeclToConstructor).From(TypeDecl)

var constructor__constructor_name_r = parser.
	Get(constructorBinaryReduction).From(Constructor, Name)

var constructor__constructor_constructor_r = parser.
	Get(constructorBinaryReduction).From(Constructor, Constructor)

var constructor__enclosed_r = parser.
	Get(parenEnclosedReduction).From(LeftParen, Constructor, RightParen)
