package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

func _nodeAsToken(node ast.Ast) token.Token {
	return node.(Node).Token
}

var funcNameAsToken, nameAsToken, paramAsToken = _nodeAsToken, _nodeAsToken, _nodeAsToken

var nameReduction = giveTypeToTokenReductionGen(Name)

var paramReduction = giveTypeToTokenReductionGen(Param)

var funcNameReduction = giveTypeToTokenReductionGen(FuncName)

// == name reduction rules ====================================================

// name <- Id
var name__Id_r = parser.Get(nameReduction).From(Id)

// name <- TypeId
var name__TypeId_r = parser.Get(nameReduction).From(TypeId)

// == param reduction rules ===================================================

// param <- Id
var param__Id_r = parser.Get(paramReduction).From(Id)

// == funcName reduction rules ================================================

var funcName__Symbol_r = parser.Get(funcNameReduction).From(Symbol)

var funcName__Infixed_r = parser.Get(funcNameReduction).From(Infixed)

var funcName__Id_r = parser.Get(funcNameReduction).From(Id)
