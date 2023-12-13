// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content:
//
// Notes: -
// =============================================================================
package typedef

import (
	"github.com/petersalex27/yew-lang/parser/cons"
	nodes "github.com/petersalex27/yew-lang/parser/node"
	typename "github.com/petersalex27/yew-lang/parser/type-name"
	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

type TypeDefNode struct {
	Name         token.Token
	TypeParams   []typename.VariableTypeName
	Constructors []*cons.TypeConstructorNode
}

// declares type defined by receiver as part of the source file's context
func (typeDef *TypeDefNode) declareContext(cxt *inf.Context[token.Token]) {
	// name of type
	var typeConstant types.Monotyped[token.Token] = types.MakeConst(typeDef.Name)
	// create open version of defined type
	openDefinedType := fun.FoldLeft(
		typeConstant,
		typeDef.TypeParams,
		func(left types.Monotyped[token.Token], right typename.VariableTypeName) types.Monotyped[token.Token] {
			mono := right.AsReferableMonotype()
			return types.Apply[token.Token](left, mono)
		},
	)

	_ = cxt.AddType(typeDef.Name, openDefinedType)
}

// defines declared type defined by receiver as part of the source file's
// context. After `declareContext` has been called on receiver, this function's
// successful execution will allow defined type to be used during type inf.
func (typeDef *TypeDefNode) defineContext(cxt *inf.Context[token.Token]) {
	typeName := typeDef.Name
	for _, constructor := range typeDef.Constructors {
		data := constructor.Data
		_ = cxt.AddConstructorFor(typeName, data)
	}
}

// does nothing
func (*TypeDefNode) Visit(*inf.Context[token.Token]) {}

func (typeDef *TypeDefNode) Equals(node ast.Ast) bool {
	typeDef2, ok := node.(*TypeDefNode)
	if !ok {
		return false
	}

	if !utils.EquateTokens(typeDef.Name, typeDef2.Name) {
		return false
	}

	for i, param := range typeDef.TypeParams {
		param2 := typeDef2.TypeParams[i]
		if !utils.EquateTokens(token.Token(param), token.Token(param2)) {
			return false
		}
	}

	return nodes.NodesEquals(typeDef.Constructors, typeDef2.Constructors)
}

func (*TypeDefNode) NodeType() ast.Type { return TypeDef }

func (typeDef *TypeDefNode) InOrderTraversal(action func(itoken.Token)) {
	action(typeDef.Name)

	for _, param := range typeDef.TypeParams {
		action(token.Token(param))
	}

	for _, constructor := range typeDef.Constructors {
		constructor.InOrderTraversal(action)
	}
}