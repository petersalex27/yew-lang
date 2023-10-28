package parser

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-lang/token"
)

func getDependIndex(node ast.Ast) types.DependentTypeInstance[token.Token] {
	return GetType(node).Type.(types.DependentTypeInstance[token.Token])
}

func arrEnclose(lbracket ast.Ast) types.EnclosingConst[token.Token] {
	line, char := GetToken(lbracket).GetLineChar()
	tok := token.TypeId.Make().AddValue("[]").SetLineChar(line, char).(token.Token)
	return types.MakeEnclosingConst[token.Token](1, tok)
}

var getDependInstance = getDependIndex

func getDependentTyped(node ast.Ast) types.DependentTyped[token.Token] {
	return GetType(node).Type.(types.DependentTyped[token.Token])
}

func freeJudgement(ex expr.Expression[token.Token]) types.ExpressionJudgement[token.Token, expr.Expression[token.Token]] {
	lockType()
	judge := types.FreeJudge[token.Token, expr.Expression[token.Token]](glb_cxt.typeCxt, ex)
	unlockType()
	return types.ExpressionJudgement[token.Token, expr.Expression[token.Token]](judge)
}

// dependTyped <- dependInstance
var dependTyped__dependInstance_r = parser.
	Get(dependTyped__dependInstance).
	From(DependInstance)

func dependTyped__dependInstance(nodes ...ast.Ast) ast.Ast {
	return TypeNode{Dependtyped, GetType(nodes[0]).Type}
}

// dependTyped <- depend
var dependTyped__depend_r = parser.
	Get(dependTyped__depend).
	From(Dependtype)

func dependTyped__depend(nodes ...ast.Ast) ast.Ast {
	return TypeNode{Dependtyped, GetType(nodes[0]).Type}
}

// dependInstance <- dependIndexHead expr
var dependInstance__dependIndexHead_expr_r = parser.
	Get(dependInstance__dependIndexHead_expr).
	From(DependIndexHead, Expr)

func dependInstance__dependIndexHead_expr(nodes ...ast.Ast) ast.Ast {
	ty := getDependIndex(nodes[0])
	ex := getExpression(nodes[1]).Expression
	// create type judgement
	judge := freeJudgement(ex)

	// make dependent type instance
	inst := types.Index(ty.Application, judge) // == (<ty>; <ex>: newvar)
	return TypeNode{DependInstance, inst}
}

// dependInstance <- dependIndexHead judgement
var dependTyped__dependIndexHead_judge_r = parser.
	Get(dependTyped__dependIndexHead_judge).
	From(DependIndexHead, TypeJudgement)

func dependTyped__dependIndexHead_judge(nodes ...ast.Ast) ast.Ast {
	ty := getDependIndex(nodes[0])
	judge := getExprJudgement_ty(nodes[1])

	// make dependent type instance
	inst := types.Index(ty.Application, judge) // == (<ty>; <judge>)
	return TypeNode{DependInstance, inst}
}

// dependInstance <- arrayHead RightBracket
var dependInstance__arrayHead_RightBracket_r = parser.
	Get(dependInstance__arrayHead_RightBracket).
	From(ArrayHead, RightBracket)

func dependInstance__arrayHead_RightBracket(nodes ...ast.Ast) ast.Ast {
	dep := getDependentTyped(nodes[0])
	return TypeNode{DependInstance, dep}
}

// arrayHead <- LeftBracket typeApp
var arrayHead__LeftBracket_typeApp_r = parser.
	Get(arrayHead__LeftBracket_typeApp).
	From(LeftBracket, TypeApp)

func arrayHead__LeftBracket_typeApp(nodes ...ast.Ast) ast.Ast {
	app := getApplicationType(nodes[1])
	ty := types.Apply[token.Token](arrEnclose(nodes[0]), app)
	// generate arbitrary type judgement for arbitrary kind-variable
	// newvar: Uint
	uintNode := TypeNode{Monotype, getUint()}
	judge := makeFreeJudgementOf(uintNode)

	// make dependent type instance
	inst := types.Index(ty, getExprJudgement_free(judge)) // == (Array <typeApp>; newvar: Uint)
	return TypeNode{ArrayHead, inst}
}

// arrayHead <- LeftBracket TypeId
var arrayHead__LeftBracket_TypeId_r = parser.
	Get(arrayHead__LeftBracket_TypeId).
	From(LeftBracket, TypeId)

func arrayHead__LeftBracket_TypeId(nodes ...ast.Ast) ast.Ast {
	name := types.MakeConst(GetToken(nodes[1]))
	ty := types.Apply[token.Token](arrEnclose(nodes[0]), name)
	// generate arbitrary type judgement for arbitrary kind-variable
	// newvar: Uint
	uintNode := TypeNode{Monotype, getUint()}
	judge := makeFreeJudgementOf(uintNode)

	// make dependent type instance
	inst := types.Index(ty, getExprJudgement_free(judge)) // == (Array <TypeId>; newvar: Uint)
	return TypeNode{ArrayHead, inst}
}

// arrayHead <- LeftBracket var
var arrayHead__LeftBracket_var_r = parser.
	Get(arrayHead__LeftBracket_var).
	From(LeftBracket, FreeVar)

func arrayHead__LeftBracket_var(nodes ...ast.Ast) ast.Ast {
	v := getVariable(nodes[1])
	ty := types.Apply[token.Token](arrEnclose(nodes[0]), v)
	// generate arbitrary type judgement for arbitrary kind-variable
	// newvar: Uint
	uintNode := TypeNode{Monotype, getUint()}
	judge := makeFreeJudgementOf(uintNode)

	// make dependent type instance
	inst := types.Index(ty, getExprJudgement_free(judge)) // == (Array <var>; newvar: Uint)
	return TypeNode{ArrayHead, inst}
}

// arrayHead <- LeftBracket dependIndexHead expr
var arrayHead__LeftBracket_dependIndexHead_expr_r = parser.
	Get(arrayHead__LeftBracket_dependIndexHead_expr).
	From(LeftBracket, DependIndexHead, Expr)

func arrayHead__LeftBracket_dependIndexHead_expr(nodes ...ast.Ast) ast.Ast {
	ty := getDependIndex(nodes[1])
	ex := getExpression(nodes[2]).Expression
	// create type judgement
	judge := freeJudgement(ex)

	// make array
	arr := types.Apply[token.Token](arrEnclose(nodes[0]), ty.Application)
	// make dependent type instance
	inst := types.Index(arr, judge) // == (Array <ty>; <ex>: newvar)
	return TypeNode{ArrayHead, inst}
}

// arrayHead <- LeftBracket dependIndexHead judgement
var arrayHead__LeftBracket_dependIndexHead_judge_r = parser.
	Get(arrayHead__LeftBracket_dependIndexHead_judge).
	From(LeftBracket, DependIndexHead, TypeJudgement)

func arrayHead__LeftBracket_dependIndexHead_judge(nodes ...ast.Ast) ast.Ast {
	ty := getDependIndex(nodes[1])
	judge := getExprJudgement_free(nodes[2])

	// make array
	arr := types.Apply[token.Token](arrEnclose(nodes[0]), ty.Application)
	// make dependent type instance
	inst := types.Index(arr, judge) // == (Array <ty>; <judge>)
	return TypeNode{ArrayHead, inst}
}

// dependIndexHead <- typeApp SemiColon
var dependIndexHead__typeApp_SemiColon_r = parser.
	Get(dependIndexHead__typeApp_SemiColon).
	From(TypeApp, SemiColon)

func dependIndexHead__typeApp_SemiColon(nodes ...ast.Ast) ast.Ast {
	ty := getApplicationType(nodes[0])
	// create empty judgement
	judge_ := types.TypeJudgement[token.Token, expr.Expression[token.Token]]{}
	judge := types.ExpressionJudgement[token.Token, expr.Expression[token.Token]](judge_)

	// make dependent type instance
	inst := types.Index(ty, judge) // == (<ty>; _)
	return TypeNode{DependIndexHead, inst}
}

// dependIndexHead <- TypeId SemiColon
var dependIndexHead__TypeId_SemiColon_r = parser.
	Get(dependIndexHead__TypeId_SemiColon).
	From(TypeId, SemiColon)

func dependIndexHead__TypeId_SemiColon(nodes ...ast.Ast) ast.Ast {
	// get name and make "application" type (really just constant wrapped in application)
	name := types.MakeConst(GetToken(nodes[0]))
	ty := types.Apply[token.Token](name)
	// create empty judgement
	judge_ := types.TypeJudgement[token.Token, expr.Expression[token.Token]]{}
	judge := types.ExpressionJudgement[token.Token, expr.Expression[token.Token]](judge_)

	// make dependent type instance
	inst := types.Index(ty, judge) // == (<ty>; _)
	return TypeNode{DependIndexHead, inst}
}

// dependIndexHead <- var SemiColon
var dependIndexHead__var_SemiColon_r = parser.
	Get(dependIndexHead__var_SemiColon).
	From(FreeVar, SemiColon)

func dependIndexHead__var_SemiColon(nodes ...ast.Ast) ast.Ast {
	v := getVariable(nodes[0])
	ty := types.Apply[token.Token](v)
	// create empty judgement
	judge_ := types.TypeJudgement[token.Token, expr.Expression[token.Token]]{}
	judge := types.ExpressionJudgement[token.Token, expr.Expression[token.Token]](judge_)

	// make dependent type instance
	inst := types.Index(ty, judge) // == (<ty>; _)
	return TypeNode{DependIndexHead, inst}
}

// dependTyped <- typeApp
/*
func dependTyped__typeApp(nodes ...ast.Ast) ast.Ast {
	return TypeNode{Dependtyped, GetType(nodes[0]).Type}
}
*/

// dependTyped <- TypeId
/*
func dependTyped__TypeId(nodes ...ast.Ast) ast.Ast {
	c := types.MakeConst(GetToken(nodes[0]))
	return TypeNode{Dependtyped, c}
}
*/

// dependTyped <- var
/*
func dependTyped__var(nodes ...ast.Ast) ast.Ast {
	return TypeNode{Dependtyped, getVariable(nodes[0])}
}
*/

// depend <- Mapall typing.var Dot monotype

// depend <- Mapall typing.freevars Dot monotype
