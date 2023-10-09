package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

type applicant struct {
	isApplication bool
	node          ast.Ast
}

type ApplicationNode struct {
	ty    ast.Type
	left  applicant
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

func (a ApplicationNode) toTypeApp() types.Application[token.Token] {
	panic("TODO: implement")
	if a.ty != Constructor {
		panic("not a constructor")
	}

	prev := a.right
	curr := a.left
	//var left, right types.Type[token.Token]
	if a.right.isApplication {
		//right = a.right.node.(ApplicationNode).toTypeApp()
	} else {

	}

	if a.left.isApplication {
		//left = a.left.node.(ApplicationNode).toTypeApp()
	}

	var tok *token.Token
	/*grabToken */_ = func(t itoken.Token) {
		*tok = t.(token.Token)
	}

	for curr.isApplication {
		//var left, right types.Application[token.Token]
		if prev.isApplication {
		//	right = prev.node.(ApplicationNode).toTypeApp()
		}
		//left
	}
	return types.Application[token.Token]{}
}

// constr <- TypeId name
var constr__TypeId_name_r = parser.
	Get(constr__false_false).
	From(TypeId, Name)

// constr <- constr name
var constr__constr_name_r = parser.
	Get(constr__true_false).
	From(Constructor, Name)

// constr <- constr constr
var constr__constr_constr_r = parser.
	Get(constr__true_true).
	From(Constructor, Constructor)

// constr <- LeftParen constr RightParen
var constr__LeftParen_constr_RightParen_r = parser.
	Get(grab_enclosed).
	From(LeftParen, Constructor, RightParen)

func constr__false_false(nodes ...ast.Ast) ast.Ast {
	return ApplicationNode{
		Constructor,
		applicant{false, nodes[0]},
		applicant{false, nodes[1]},
	}
}

func constr__true_true(nodes ...ast.Ast) ast.Ast {
	return ApplicationNode{
		Constructor,
		applicant{true, nodes[0]},
		applicant{true, nodes[1]},
	}
}

func constr__true_false(nodes ...ast.Ast) ast.Ast {
	return ApplicationNode{
		Constructor,
		applicant{true, nodes[0]},
		applicant{false, nodes[1]},
	}
}

// app <- expr expr
var app__expr_expr_r = parser.
	Get(app__expr_expr).
	From(Expr, Expr)

func app__expr_expr(nodes ...ast.Ast) ast.Ast {
	left, right := getExpression(nodes[0]).Expression, getExpression(nodes[1]).Expression
	return SomeExpression{
		Application,
		expr.Apply(left, right),
	}
}

// app <- LeftParen app RightParen
var app__enclosed_r = parser.Get(grab_enclosed).From(LeftParen, Application, RightParen)
