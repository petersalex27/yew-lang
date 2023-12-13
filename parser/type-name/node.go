// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
// =============================================================================
package typename

import (
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

// =============================================================================
// type name node interface members
// =============================================================================

type referableMonotype interface {
	types.Monotyped[token.Token]
	GetReferred() token.Token
}

type castableToReferableMonotype interface {
	AsReferableMonotype() referableMonotype
}

// =============================================================================
// type name node
// =============================================================================

// node representing the name of one of two parts of a type: a constant or a 
// variable
type TypeNameNode struct {
	Name castableToReferableMonotype
}

// =============================================================================
// referable monotype pseudo nodes
// =============================================================================

type VariableTypeName token.Token

type ConstantTypeName token.Token

// =============================================================================
// referable monotype pseudo nodes castable impl.
// =============================================================================

func (name VariableTypeName) AsReferableMonotype() referableMonotype {
	return types.Var(token.Token(name))
}

func (name ConstantTypeName) AsReferableMonotype() referableMonotype {
	return types.MakeConst(token.Token(name))
}

// =============================================================================
// helper function for castable interface
// =============================================================================

// gets token from `m`
func castToToken(m castableToReferableMonotype) token.Token {
	return m.AsReferableMonotype().GetReferred()
}

// =============================================================================
// *TypeNameNode ast node impl.
// =============================================================================

// does nothing
func (*TypeNameNode) Visit(*inf.Context[token.Token]) {}

// true iff `node` is a *TypeNameNode and receiver has identical token to `node`
func (typeName *TypeNameNode) Equals(node ast.Ast) bool {
	typeName2, ok := node.(*TypeNameNode)
	if !ok {
		return false
	}

	a, b := castToToken(typeName.Name), castToToken(typeName2.Name)
	return utils.EquateTokens(a, b)
}

// returns type of token (as an ast.Type)
func (typeName *TypeNameNode) NodeType() ast.Type {
	tokenType := castToToken(typeName.Name).GetType()
	return ast.Type(tokenType) 
}

// calls action on underlying token of typeName.Name
func (typeName *TypeNameNode) InOrderTraversal(action func(itoken.Token)) { 
	token := castToToken(typeName.Name)
	action(token)
}
