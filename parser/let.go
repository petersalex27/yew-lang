package parser

import (
	"strconv"

	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
)

// letIn ::= 'let' function 'in' expr

var letIn__Let_function_In_expr_r = parser.
	Get(letReduction).From(Let, Function, In, Expr)

func toApplication[T expr.Expression[token.Token]](params []T) expr.Expression[token.Token] {
	if len(params) == 0 {
		return nil
	}

	var out expr.Expression[token.Token] = params[0]
	for _, ex := range params[1:] {
		out = expr.Apply[token.Token](out, ex)
	}
	return out
}

// let
//
//	  -- function
//		f				-- name
//		x y 		-- params
//		=
//		x + y 	-- body
//
// in
//
//	f 1 1						-- expr
func letReduction(nodes ...ast.Ast) ast.Ast {
	const _, functionIndex, _, exprIndex int = 0, 1, 2, 3

	function := nodes[functionIndex].(FunctionNode)
	expression := nodes[exprIndex].(ExpressionNode)

	// break apart function into component parts
	name := function.def.head.name
	params := function.def.head.params
	body := function.body

	numParams := len(params)

	// create number of dummy params equal to the length of params
	binders := make(expr.BindersOnly[token.Token], numParams)
	dummyBinders := make(expr.BindersOnly[token.Token], numParams)
	dummy := expr.Var(token.Id.Make().AddValue("_"))
	for i := range binders {
		paramStr := "$p" + strconv.Itoa(i)
		paramToken := token.Id.Make().AddValue(paramStr)
		paramVar := expr.Var(paramToken)
		binders[i] = paramVar
		dummyBinders[i] = dummy
	}

	case_ := expr.
		Bind(dummyBinders[0], dummyBinders[1:]...). // _ _ ..
		InCase(toApplication(params), body)         // p1 p2 .. -> body
	patternMatch := expr.Select(toApplication(binders), case_)
	anon := expr.Bind(binders[0], binders[1:]...).In(patternMatch)
	letExpr := expr.Let[token.Token](Const(name), anon, expression.Expression)

	return SomeExpression{LetExpr, letExpr}

	// create data deconstruction expansion (sequence of "let" expressions)
	// for each data parameter in the `params` pattern

	// embed expanded function inside of pattern match expression which is inside
	// the wrapper function
	// 		(\$p1 $p2 .. -> ($p1 $p2 ..) when
	//			(a1 a2 ..) ->
	//					let memb1 = (getMember $p1 (a: Int)) in
	//					.. in
	//					let memb2 = (getMember $pN (b: Int)) in
	//					.. in body
	//		)

	// Let params = p1 p2 ..
	// Then, abstract
	//		(\p1 p2 .. -> body)
	// Then, create let node
	//		let = {name, abstraction, expression}
	// Or, in Yew
	//		let name = (\p1 p2 .. -> body) in expression
}

// exprWhere ::= expr 'where' function
