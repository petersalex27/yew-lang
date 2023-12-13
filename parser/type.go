package parser

import (
	"github.com/petersalex27/yew-lang/token"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
)

type TypeNodeType ast.Type

const (
	TypeMonotype      TypeNodeType = TypeNodeType(Monotype)
	TypeDependentType TypeNodeType = TypeNodeType(DependentType)
	TypePolytype      TypeNodeType = TypeNodeType(Polytype)
)

type TypeNode struct {
	TypeNodeType
	types.Type[token.Token]
}
