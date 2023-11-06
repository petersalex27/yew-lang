package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

// =============================================================================
// indent shift __reduction__
// =============================================================================

// full reduction rule for L.A. INDENT(_)
var reductionSHIFT_Indent = parser.LookAhead(Indent).Shift()

// =============================================================================
// indent production rules
// =============================================================================

// rule for grabbing right-most indent
var _Indent__Indent_Indent_r = parser.
	Get(indentPushSecondProduction).
	From(Indent, Indent)

// rule for opening expression blocks
var exprBlock__Assign_Indent_r = parser.
	Get(openExprBlockProduction).
	When(Assign).From(Indent)

// rule for opening expression blocks
var exprBlock__Where_Indent_r = parser.
	Get(openExprBlockProduction).
	When(Where).From(Indent)

// rule for opening expression blocks
var exprBlock__Let_Indent_r = parser.
	Get(openExprBlockProduction).
	When(Let).From(Indent)

// rule for opening expression blocks
var exprBlock__In_Indent_r = parser.
	Get(openExprBlockProduction).
	When(In).From(Indent)

// rule for opening expression blocks
var exprBlock__Of_Indent_r = parser.
	Get(openExprBlockProduction).
	When(Of).From(Indent)

// rule for opening expression blocks
var exprBlock__Match_Indent_r = parser.
	Get(openExprBlockProduction).
	When(Match).From(Indent)

// =============================================================================
// 0 length indent production rules
// =============================================================================

// rule for opening expression blocks
var exprBlock__Assign_r = parser.
	Get(openExprBlock0Production).
	When(Assign).From()

// rule for opening expression blocks
var exprBlock__Where_r = parser.
	Get(openExprBlock0Production).
	When(Where).From()

// rule for opening expression blocks
var exprBlock__Let_r = parser.
	Get(openExprBlock0Production).
	When(Let).From()

// rule for opening expression blocks
var exprBlock__In_r = parser.
	Get(openExprBlock0Production).
	When(In).From()

// rule for opening expression blocks
var exprBlock__Of_r = parser.
	Get(openExprBlock0Production).
	When(Of).From()

// rule for opening expression blocks
var exprBlock__Match_r = parser.
	Get(openExprBlock0Production).
	When(Match).From()

// =============================================================================
// expression block start node
// =============================================================================

// represents the start of an expression block, an indentation
type ExprBlockStart token.Token

func (indent ExprBlockStart) Equals(a ast.Ast) bool {
	indent2, ok := a.(ExprBlockStart)
	if !ok {
		return false
	}
	return EqualsToken[token.Token](token.Token(indent), token.Token(indent2))
}

func (ExprBlockStart) NodeType() ast.Type { return IndentExprBlock }

// does nothing
func (indent ExprBlockStart) InOrderTraversal(f func(itoken.Token)) {
	f(token.Token(indent))
}

// =============================================================================
// indentation production functions
// =============================================================================

// just returns second node
func indentPushSecondProduction(nodes ...ast.Ast) ast.Ast {
	const _, indentIndex int = 0, 1
	return nodes[indentIndex] 
}

// just wraps first node (token) in ExprBlockStart
func openExprBlockProduction(nodes ...ast.Ast) ast.Ast {
	const indentIndex int = 0
	indent := GetToken(nodes[indentIndex])
	return ExprBlockStart(indent)
}

// returns 0 length indent
//
// TODO: need a way to put the next token's line-char-len info into this node?
func openExprBlock0Production(nodes ...ast.Ast) ast.Ast {
	indent := token.Indent.Make().AddValue("")
	return ExprBlockStart(indent)
}