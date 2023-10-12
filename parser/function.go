package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

// TODO
/*
functionDecl  ::= funcName param
                  | funcName constructor
                  | functionDecl param
                  | functionDecl constructor
functionDef   ::= functionDecl ':' type
function      ::= functionDef '=' expr
                  | functionDecl '=' expr
*/

type paramInterface interface {
	getAsExpression() expr.Expression[token.Token]
}

type FunctionDeclNode struct {
	name   token.Token
	params []ast.Ast
}

type FunctionDefNode types.TypeJudgement[token.Token, expr.Application[token.Token]]

type FunctionNode struct {
	FunctionDefNode
	body expr.Expression[token.Token]
}

func (fd FunctionDeclNode) MakeApplication() expr.Application[token.Token] {
	panic("TODO: implement")
}

func (fd FunctionDeclNode) Equals(a ast.Ast) bool {
	fd2, ok := a.(FunctionDeclNode)
	if !ok {
		return false
	}
	panic("TODO: implement")
	if /*!fd.name.Equals(fd2.name)*/false {
		return false
	}

	if len(fd.params) != len(fd2.params) {
		return false
	}

	for i, param := range fd.params {
		if !param.Equals(fd2.params[i]) {
			return false
		}
	}
	return true
}

func (fd FunctionDeclNode) NodeType() ast.Type { return FunctionDecl }

func (fd FunctionDeclNode) InOrderTraversal(f func(itoken.Token)) {
	panic("TODO: implement")
	//fd.name.InOrderTraversal(f)
	for _, param := range fd.params {
		param.InOrderTraversal(f)
	}
}

func (fd FunctionDefNode) Equals(a ast.Ast) bool {
	fd2, ok := a.(FunctionDefNode)
	if !ok {
		return false
	}

	return types.JudgesEquals[token.Token, expr.Application[token.Token], expr.Application[token.Token]](
			fd.TypeJudgement, fd2.TypeJudgement) &&
			fd.body.Equals(glb_cxt.exprCxt, fd2.body)
}

func (fd FunctionDefNode) NodeType() ast.Type { return FunctionDefinition }

func (fd FunctionDefNode) InOrderTraversal(f func(itoken.Token)) {
	for _, tok := range fd.TypeJudgement.Collect() {
		f(tok)
	}
	for _, tok := range fd.body.Collect() {
		f(tok)
	}
}

// fndecl <- name construct
var fndecl__name_construct_r = parser. 
	Get(fndecl__name_Ast). 
	From(Name, Constructor)

// fndecl <- name Id
var fndecl__name_Id_r = parser. 
	Get(fndecl__name_Ast). 
	From(Name, Id)

func fndecl__name_Ast(nodes ...ast.Ast) ast.Ast {
	panic("TODO: implement")
	return FunctionDeclNode{
		//name: getNameNode(nodes[0]),
		params: []ast.Ast{nodes[1]},
	}
}

// fndecl <- fndecl construct
var fndecl__fndecl_construct_r = parser. 
	Get(fndecl__fndecl_Ast). 
	From(Name, Constructor)

// fndecl <- fndecl Id
var fndecl__fndecl_Id_r = parser. 
	Get(fndecl__fndecl_Ast). 
	From(Name, Id)

func fndecl__fndecl_Ast(nodes ...ast.Ast) ast.Ast {
	panic("TODO: implement")
	return FunctionDeclNode{
		//name: getNameNode(nodes[0]),
		params: []ast.Ast{nodes[1]},
	}
}

func fndef__fndecl_Assign_expr(nodes ...ast.Ast) ast.Ast {
	panic("TODO: implement")
}