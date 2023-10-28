package parser

import (
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-lang/token"
)

func _nodeAsToken(node ast.Ast) token.Token {
	return node.(Node).Token
}

var funcNameAsToken, nameAsToken, paramAsToken = 
	_nodeAsToken, _nodeAsToken, _nodeAsToken

var nameReduction = simpleNodeRule(Name)

var paramReduction = simpleNodeRule(Param)

var funcNameReduction = simpleNodeRule(FuncName)

var anyNameReduction = rewrapNodeRule(AnyName)

// == name rules ============================================

// name <- Id
var name__Id_r = parser.Get(nameReduction).From(Id)

// name <- TypeId
var name__TypeId_r = parser.Get(nameReduction).From(TypeId)

// == param rules ==============================================

// param <- Id
var param__Id_r = parser.Get(paramReduction).From(Id)

// param <- Thunked
var param__Thunked_r = parser.Get(paramReduction).From(Thunked)

// == funcName rules ===========================================

// funcName <- Symbol
var funcName__Symbol_r = parser.Get(funcNameReduction).From(Symbol)

// funcName <- Infixed
var funcName__Infixed_r = parser.Get(funcNameReduction).From(Infixed)

// funcName <- Id
var funcName__Id_r = parser.Get(funcNameReduction).From(Id)

// == anyName rules ============================================

// anyName <- name
var anyName__name_r = parser.Get(anyNameReduction).From(Name)

// anyName <- param
var anyName__param_r = parser.Get(anyNameReduction).From(Param)

// anyName <- funcName
var anyName__funcName_r = parser.Get(anyNameReduction).From(FuncName)