package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
)

/*
functionDef ::= functionHead ':' type
								| functionHead							# only when L.A. is '='
*/

// FunctionDefNode is a function declaration with a type
type FunctionDefNode struct {
	head   FunctionHeadNode
	typing types.Type[token.Token]
}

// funcDef <- funcHead ':' type
var funcDef__funcHead_Colon_type_r = parser.
	Get(typedFunctionReduction).
	From(FunctionHead, Colon, Type)

// funcDef <- funcHead
var funcDef__funcHead_r = parser.
	Get(untypedFunctionReduction).
	From(FunctionHead)

// creates a function defintion node from
//
//	nodes[0].(FunctionHeadNode)
//
// and
//
//	nodes[2].(TypeNode)
func typedFunctionReduction(nodes ...ast.Ast) ast.Ast {
	const funcHeadIndex, _, typeIndex int = 0, 1, 2
	funcHead := nodes[funcHeadIndex].(FunctionHeadNode)
	typing := nodes[typeIndex].(TypeNode).Type
	return FunctionDefNode{
		head:   funcHead,
		typing: typing,
	}
}

// creates a function definition node with a function type made completely from
// new free type variables
func untypedFunctionReduction(nodes ...ast.Ast) ast.Ast {
	const funcHeadIndex int = 0
	funcHead := nodes[funcHeadIndex].(FunctionHeadNode)

	// number of params and number of type variables (excluding return type)
	length := len(funcHead.params)

	// generated type for function declaration
	var typing types.Monotyped[token.Token]

	if length < 1 {
		// just name
		globalContext__.typeMutex.Lock()
		typing = globalContext__.typeCxt.NewVar()
		globalContext__.typeMutex.Unlock()
	} else {
		// name and an application pattern

		// create an array for `length` type variables
		vars := make([]types.Variable[token.Token], length)

		globalContext__.typeMutex.Lock()

		// generate `length` new type variables
		//
		// The type variables are generated in reverse order so that the next
		// for-loop that generates the function type, from left to right, in
		// oldest to newest type variable
		for i := range vars {
			vars[length-1-i] = globalContext__.typeCxt.NewVar()
		}

		// return type of function
		retType := globalContext__.typeCxt.NewVar()

		globalContext__.typeMutex.Unlock()

		// aN -> retType
		functionType := types.Apply[token.Token](arrowConst, vars[0], retType)

		// creates the following given vars={aN, .., a2, a1}:
		//  a1 -> (a2 -> (.. -> (aN -> retType)))
		for _, variable := range vars[1:] {
			functionType = types.Apply[token.Token](arrowConst, variable, functionType)
		}

		typing = functionType
	}

	return FunctionDefNode{head: funcHead, typing: typing}
}

func (fd FunctionDefNode) Equals(a ast.Ast) bool {
	fd2, ok := a.(FunctionDefNode)
	if !ok {
		return false
	}

	return fd.head.Equals(fd2.head) && fd.typing.Equals(fd2.typing)
}

// returns FunctionDefinition
func (fd FunctionDefNode) NodeType() ast.Type { return FunctionDefinition }

func (fd FunctionDefNode) InOrderTraversal(f func(itoken.Token)) {
	fd.head.InOrderTraversal(f)
	tokens := fd.typing.Collect()
	for _, token := range tokens {
		f(token)
	}
}
