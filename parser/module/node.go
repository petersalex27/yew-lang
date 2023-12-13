package module

import (
	"github.com/petersalex27/yew-lang/parser/function/bound"
	"github.com/petersalex27/yew-lang/parser/imports"
	nodes "github.com/petersalex27/yew-lang/parser/node"
	typedef "github.com/petersalex27/yew-lang/parser/type-def"

	//. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/parser/utils"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"

	//"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type ModuleSourceNode struct {
	ast.Type
	// module Name ( functionNames, typeNames (constructors), ... ) where definitions
	Name token.Token

	exportAll        bool
	FunctionNames    []token.Token
	TypeNames        []token.Token
	ConstructorNames [][]token.Token

	Imports []*imports.ImportNode
	Uses []*imports.ImportNode

	TypeDefinitions     []*typedef.TypeDefNode
	FunctionDefinitions []*bound.BoundFunctionNode
}

// module ::= export ')' 'where'
//var explicitModuleExportRule = parser.Get(produceModuleAndExportList).From(ExportList, RightParen, Where)

func (msn *ModuleSourceNode) ExportContext() *inf.ExportableContext[token.Token] {
	cxt, export := inf.Export(msn.Name, makeToken, msn.FunctionNames, msn.TypeNames, msn.ConstructorNames)

	// do imports
	for _, imprt := range msn.Imports {
		imprt.Visit(cxt)
	}

	// do uses
	for _, use := range msn.Uses {
		use.Visit(cxt)
	}

	// define types
	for _, typeDef := range msn.TypeDefinitions {
		typeDef.Visit(cxt)
	}

	// define functions
	for _, function := range msn.FunctionDefinitions {
		function.Visit(cxt)
	}

	return export()
}

func makeToken(s string) token.Token {
	return token.Id.Make().AddValue(s)
}

func (msn *ModuleSourceNode) Equals(node ast.Ast) bool {
	msn2, ok := node.(*ModuleSourceNode)
	if !ok {
		return false
	}

	ok = utils.EquateTokens(msn.Name, msn2.Name) &&
		utils.TokensEquals(msn.FunctionNames, msn2.FunctionNames) &&
		utils.TokensEquals(msn.TypeNames, msn2.TypeNames)
	if !ok {
		return false
	}

	if len(msn.ConstructorNames) != len(msn2.ConstructorNames) {
		return false
	}

	for i, cons := range msn.ConstructorNames {
		for j, con := range cons {
			if !utils.EquateTokens(con, msn2.ConstructorNames[i][j]) {
				return false
			}
		}
	}

	return nodes.NodesEquals(msn.TypeDefinitions, msn2.TypeDefinitions) &&
		nodes.NodesEquals(msn.FunctionDefinitions, msn2.FunctionDefinitions)
}

func (msn *ModuleSourceNode) NodeType() ast.Type { return msn.Type }

func (msn *ModuleSourceNode) InOrderTraversal(action func(itoken.Token)) {
	action(msn.Name)

	for _, token := range msn.FunctionNames {
		action(token)
	}

	for _, token := range msn.TypeNames {
		action(token)
	}

	for _, tokens := range msn.ConstructorNames {
		for _, token := range tokens {
			action(token)
		}
	}

	for _, imprt := range msn.Imports {
		imprt.InOrderTraversal(action)
	}

	for _, typeDef := range msn.TypeDefinitions {
		typeDef.InOrderTraversal(action)
	}

	for _, function := range msn.FunctionDefinitions {
		function.InOrderTraversal(action)
	}
}
