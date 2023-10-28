package parser

import (
	"strconv"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-lang/token"
)

// letIn ::= 'let' function 'in' expr

var letIn__Let_function_In_expr_r = parser. 
	Get(letReduction).From(Let, Function, In, Expr)	

// 
func generateDataIndexer(data DataStruct) {
	
}

func generateDataIndexerCall(enclosingName string, offset int) {

}

// Deconstructs an application into a NameContext expression.
// `data` is a left-proper-sub-tree of initial argument passed to data.
// return values: 
//  - `cxt` is the deconstruction result
//  - `ok` is true iff `cxt` has a value
//  - `enclosingName` is the name of the constant of the enclosing data 
//		type; this is used to build the correct member access function
//	- `offset` is member index of the the right-most member of the current
//		call frame's `data`
func deconstructHelper(data expr.Application[token.Token]) (cxt expr.NameContext[token.Token], ok bool, enclosingName string, offset int) {
	left, right := data.Split()

	// deconstruct left side
	if dataLeft, ok2 := left.(expr.Application[token.Token]); ok2 {
		// data has more than one member remaining
		cxt, ok, enclosingName, offset = deconstructHelper(dataLeft)
	} else if dataName, okConst := left.(expr.Const[token.Token]); okConst {
		// at data head; e.g., given the following input to the initial call to 
		// this function `(MyData a b c)`, dataName.String() == "MyData"
		ok = false
		enclosingName = dataName.String()
	}

	if dataRight, okR := right.(expr.Application[token.Token]); okR {
		cxt2, ok2, _ := deconstructHelper(dataRight)
		if ok2 {
			if ok {
				cxt = cxt.SetContextualized(cxt2)
			} else {
				cxt = cxt2
			}
		}
	} else {
		if v, isVar := right.(expr.Variable[token.Token]); isVar {
			expr.Let[token.Token](Const(v.Collect()[0]), , nil)
		}
	}
}

func deconstruct(data expr.Application[token.Token]) (cxt expr.NameContext[token.Token], ok bool) {
	
}

func recursiveDeconstruction(params []expr.Expression[token.Token], body expr.Expression[token.Token]) *expr.NameContext[token.Token] {
	if len(params) == 0 {
		return nil
	}

	// attempt type assertion on param[0], asserting that it's an application 
	// (which is basically a struct/data instance in this context)
	data, ok := params[0].(expr.Application[token.Token])
	if !ok {
		// try to deconstruct remaining params
		return recursiveDeconstruction(params[1:], body)
	}

	// deconstruct current parameter
	cxt, success := deconstruct(data)

	// deconstruct remaining parameters
	res := recursiveDeconstruction(params[1:], body)
	if res != nil && success {
		// Reassign res to cxt contexualized with res.
		// This is done so when !success, whatever res is is returned and when
		// sucesss, then cxt contextualized with res is returned
		*res = cxt.SetContextualized(*res)
	} else if success {
		res = new(expr.NameContext[token.Token])
		*res = cxt
	}

	return res
} 

// let 
//	  -- function
//		f				-- name
//		x y 		-- params
//		= 
//		x + y 	-- body
// in 
// 		f 1 1						-- expr
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
	for i := range binders {
		paramToken := token.Id.Make().AddValue("$p"+strconv.Itoa(i))
		paramVar := expr.Var(paramToken)
		binders[i] = paramVar
	}

	// data member deconstructions
	var deconstructionRoot expr.NameContext[token.Token]
	var deconstruction *expr.NameContext[token.Token]
	deconstructedSomeParam := false

	
	// create data deconstruction expansion (sequence of "let" expressions)
	// for each data parameter in the `params` pattern


	// embed expanded function inside of pattern match expression which is inside
	// the wrapper function
	// 		(\$p1 $p2 .. -> ($p1 $p2 ..) when 
	//			(p1 p2 ..) -> 
	//					let memb1 = (getMember p1 (a: Int)) in
	//					.. in
	//					let memb2 = (getMember pN (b: Int)) in 
	//					.. in body
	//		)

	// Let params = p1 p2 ..
	// Then, abstract 
	//		(\p1 p2 .. -> body)
	// Then, create let node
	//		let = {name, abstraction, expression}
	// Or, in Yew 
	//		let name = (\p1 p2 .. -> body) in expression
	return SomeExpression{
		LetDeclaration,
		expr.Select[token.Token](),
	}
}


// exprWhere ::= expr 'where' function