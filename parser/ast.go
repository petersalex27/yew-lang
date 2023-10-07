package parser

import (
	"sync"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
	"yew.lang/main/token"
)

type glb_cxt_t struct {
	typeMutex sync.Mutex
	exprMutex sync.Mutex
	typeCxt   *types.Context[token.Token]
	exprCxt   *expr.Context[token.Token]
}

var glb_cxt *glb_cxt_t

func reInit() {
	glb_cxt = new(glb_cxt_t)
	glb_cxt.exprCxt = expr.NewContext[token.Token]().SetNameMaker(
		func(s string) token.Token {
			return token.Id.Make().AddValue(s)
		},
	)
	glb_cxt.typeCxt = types.NewContext[token.Token]().SetNameMaker(
		func(s string) token.Token {
			return token.Id.Make().AddValue(s)
		},
	)
}

func init() {
	reInit()
}

// creates a new type variable
func newTypeVar() types.Variable[token.Token] {
	glb_cxt.typeMutex.Lock()
	defer glb_cxt.typeMutex.Unlock()

	return glb_cxt.typeCxt.NewVar()
}

func lockType() {
	glb_cxt.typeMutex.Lock()
}

func unlockType() {
	glb_cxt.typeMutex.Unlock()
}

// creates a new kind variable
func newKindVar() expr.Variable[token.Token] {
	glb_cxt.exprMutex.Lock()
	defer glb_cxt.exprMutex.Unlock()

	return glb_cxt.exprCxt.NewVar()
}

func generateGetConst(constantName string) func() types.Constant[token.Token] {
	return func() types.Constant[token.Token] {
		tok := token.TypeId.Make().AddValue(constantName)
		return types.MakeConst(tok)
	}
}

var getUint = generateGetConst("Uint")
var getInt = generateGetConst("Int")
var getAny = generateGetConst("@any")
var getString = generateGetConst("String")
var getChar = generateGetConst("Char")
var getBool = generateGetConst("Bool")

func makeFreeJudgementOf(ty TypeNode) JudgementNode {
	varNode := ExpressionNode{newKindVar()}
	return makeJudgement(varNode, ty)
}

func GetToken(a ast.Ast) token.Token {
	tmp, _ := a.(ast.Token)
	tok, _ := tmp.Token.(token.Token)
	return tok
}

// a <- LeftParen a RightParen 
func grab_enclosed(nodes ...ast.Ast) ast.Ast {
	return nodes[1]
}