package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"yew.lang/main/token"
)

type applicant struct {
	isApplication bool
	node ast.Ast
}

type ApplicationNode struct {
	ty ast.Type
	left applicant
	right applicant
}

func (t ApplicationNode) Equals(a ast.Ast) bool {
	t2, ok := a.(ApplicationNode)
	if !ok {
		return false
	}

	return t.left.node.Equals(t2.left.node) && t.right.node.Equals(t2.right.node)
}

func GetApp(node ast.Ast) ApplicationNode {
	return node.(ApplicationNode)
}

func (t ApplicationNode) NodeType() ast.Type { return t.ty }

func (t ApplicationNode) InOrderTraversal(f func(itoken.Token)) {
	t.left.node.InOrderTraversal(f)
	t.right.node.InOrderTraversal(f)
}

// appId <- Id Id
var appId__Id_Id_r = parser. 
	Get(appId__Id_Id). 
	From(Id, Id)

func appId__Id_Id(nodes ...ast.Ast) ast.Ast {
	return ApplicationNode{ApplicationId, applicant{false, nodes[0]}, applicant{false, nodes[1]}}
}

// addId <- appId Id
var appId__appId_Id_r = parser. 
	Get(appId__appId_Id). 
	From(ApplicationId, Id)

func appId__appId_Id(nodes ...ast.Ast) ast.Ast {
	return ApplicationNode{ApplicationId, applicant{true, nodes[0]}, applicant{false, nodes[1]}}
}

// addId <- appId appId
var appId__appId_appId_r = parser. 
	Get(appId__appId_appId). 
	From(ApplicationId, ApplicationId)

func appId__appId_appId(nodes ...ast.Ast) ast.Ast {
	return ApplicationNode{ApplicationId, applicant{true, nodes[0]}, applicant{true, nodes[1]}}
}

// addId <- Id appId
var appId__Id_appId_r = parser. 
	Get(appId__Id_appId). 
	From(Id, ApplicationId)

func appId__Id_appId(nodes ...ast.Ast) ast.Ast {
	return ApplicationNode{ApplicationId, applicant{false, nodes[0]}, applicant{true, nodes[1]}}
}

// addId <- LeftParen appId RightParen
var appId__enclosed_r = parser.Get(grab_enclosed).From(LeftParen, ApplicationId, RightParen)

// app <- appId
var app__appId_r = parser. 
	Get(app__appId). 
	From(ApplicationId)

func app__appId(nodes ...ast.Ast) ast.Ast {
	ids := nodes[0].(ApplicationNode)
	var left, right expr.Expression[token.Token]

	if ids.left.isApplication {
		left = app__appId(ids.left.node).(SomeExpression).Expression
	} else {
		left = expr.Const[token.Token]{Name: GetToken(ids.left.node)}
	}

	if ids.right.isApplication {
		right = app__appId(ids.right.node).(SomeExpression).Expression
	} else {
		right = expr.Const[token.Token]{Name: GetToken(ids.right.node)}
	}

	return SomeExpression{
		Application,
		expr.Apply[token.Token](left, right),
	}
}

// app <- expr expr
var app__expr_expr_r = parser. 
	Get(app__expr_expr). 
	From(Expression, Expression)

func app__expr_expr(nodes ...ast.Ast) ast.Ast {
	return SomeExpression{
		Application,
		expr.Apply(asExpression(nodes[0]), asExpression(nodes[1])),
	}
}

// app <- LeftParen app RightParen
var app__enclosed_r = parser.Get(grab_enclosed).From(LeftParen, Application, RightParen)