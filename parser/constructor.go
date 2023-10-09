package parser

import (
	"github.com/petersalex27/yew-packages/parser"
)

var constructorSingleReduction = simpleNodeRule(Constructor)

var constructorBinaryReduction = simpleBinaryNodeRule(Constructor)

var getConstructor = getBinaryRecursiveNode

/*
constructor   ::= TYPE_ID
                  | constructor name
                  | constructor constructor
                  | '(' constructor ')'
*/

var constructor__TypeId_r = parser. 
	Get(constructorSingleReduction).From(TypeId)

var constructor__constructor_name_r = parser. 
	Get(constructorBinaryReduction).From(Constructor, Name)

var constructor__constructor_constructor_r = parser. 
	Get(constructorBinaryReduction).From(Constructor, Constructor)

var constructor__enclosed_r = parser.
	Get(grab_enclosed).From(LeftParen, Constructor, RightParen)