package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

func constructorCast(node ast.Ast) BinaryRecursiveNode {
	return node.(BinaryRecursiveNode)
}

var constructorSingleReduction = simpleNodeRule(Constructor)

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

func setType(node BinaryRecursiveNode, ty ast.Type) BinaryRecursiveNode {
	if n, ok := node.(BinaryNode); ok {
		n.ty = ty
		return n
	} else if n, ok := node.(Node); ok {
		n.ty = ty
		return n
	}
	return node
}

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
	Get(grab_enclosed).From(LeftParen, Constructor, RightParen)