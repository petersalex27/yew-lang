package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type ErrorNode struct{ error }

// (ErrorNode) Equals(ast.Ast) returns false if type assertion
//
//	a.(ErrorNode)
//
// fails; Else, function returns
//
//	e1.Error() == a.(ErrorNode).Error()
func (e1 ErrorNode) Equals(a ast.Ast) bool {
	e2, ok := a.(ErrorNode)
	if !ok {
		return false
	}

	return e1.Error() == e2.Error()
}

// returns Error
func (e ErrorNode) NodeType() ast.Type { return Error }

// does nothing
func (e ErrorNode) InOrderTraversal(func(itoken.Token)) {}
