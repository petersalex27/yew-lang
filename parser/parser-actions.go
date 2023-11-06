package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/types"
)

const (
	// Creates a parser.MinSource for the current source code being parsed, and
	// then sets it to the value pointed to by srcDest.
	//		getSource(srcDest *parser.MinSource)
	getSource string = "getSource"

	// Sets value pointed to by pathDest to current path being parsed.
	// 		getPath(pathDest *string)
	getPath string = "getPath"

	// Creates a new type variable and sets it to the value pointed to by 
	// typeVariableDest
	//		newTypeVariable(typeVariableDest *types.Variable[token.Token])
	newTypeVariable string = "newTypeVariable"

	// Creates a new kind variable and sets it to the value pointed to by 
	// kindVariableDest
	//		newKindVariable(kindVariableDest *expr.Variable[token.Token])
	newKindVariable string = "newKindVariable"
)

// =============================================================================
// error carrier type
// =============================================================================

// allows error to be sent back to caller and caller to send badItem
type errorCarrier struct {
	error
	badItem any
}

// generates all available parser actions:
//		getPath(pathDest *string)
//		newTypeVariable(typeVariableDest *types.Variable[token.Token])
//		newKindVariable(kindVariableDest *expr.Variable[token.Token])
//
// NOTE: I will try to keep documentation up-to-date here, but see
// github.com/petersalex27/yew-lang/parser/<branch-being-used>/parser-actions.go
// for most up to date list of available functions. The most up-to-date list of 
// functions can also be found in the same file this function is defined in
func (p *Context) generateActions() []parser.Action {
	return []parser.Action{
		{
			Name: getSource, 
			Does: func(srcDest any) {
				srcPtr, _ := srcDest.(*parser.MinSource)
				*srcPtr = p.src
			},
		},
		{
			Name: getPath, 
			Does: func(pathDest any) {
				pathPtr, _ := pathDest.(*string)
				*pathPtr = p.path
			},
		},
		{
			Name: newTypeVariable, 
			Does: func(typeVariableDest any) {
				p.typeMutex.Lock()
				defer p.typeMutex.Unlock()

				typeVariablePtr, _ := typeVariableDest.(*types.Variable[token.Token])
				*typeVariablePtr = p.typeCxt.NewVar()
			},
		},
		{
			Name: newKindVariable,
			Does: func(kindVariableDest any) {
				p.exprMutex.Lock()
				defer p.exprMutex.Unlock()

				kindVariable, _ := kindVariableDest.(*expr.Variable[token.Token])
				*kindVariable = p.exprCxt.NewVar()
			},
		},
	}
}