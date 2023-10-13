package parser

import (
	"sync"

	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/stack"
	"yew.lang/main/token"
)

/*
typeDecl      ::= TYPE_ID
                  | typeDecl ID
typeDef       ::= typeDecl '=' constructor
									| typeDef '|' constructor
*/

var typeDecl__TypeId_r = parser.
	Get(simpleNodeRule(TypeDecl)).
	From(TypeId)
	
// creates a tree structure like:
//          /\
//         /\ Id
//        /\ Id
//  TypeId  Id
// typeDecls should not be stored as NodeSequences because they will most likely 
// be converted to constructor nodes; and constructor nodes are stored in a tree 
// structure
var typeDecl__typeDecl_Id_r = parser.
	Get(func(nodes ...ast.Ast) ast.Ast {
		return simpleBinaryNodeRule(TypeDecl)(nodes[0], simpleNodeRule(Name)(nodes[1]))
	}).From(TypeDecl, Id)

// unrolls tree structue and splits the contant head from the free type variables
func typeDeclType(declNode ast.Ast) (name types.Constant[token.Token], vars []types.Variable[token.Token]) {
	// typeDecls always have the following structure:
	//          /\
	//         /\ var
	//        /\ var
	//      ... var
	//      /\
	//  Name  var
	var left, right BinaryRecursiveNode
	left = declNode.(BinaryRecursiveNode)
	vars = []types.Variable[token.Token]{}

	head, childless := left.HasValue() // check if root has children
	hasFreeVars := !childless
	// stack stores free variables in reverse order
	var s *stack.Stack[types.Variable[token.Token]] = nil
	if hasFreeVars { 
		s = stack.NewStack[types.Variable[token.Token]](8) // init stack
	}

	// keep pushing right childrens' tokens (these represent free variables) 
	// until left-most node is reached; the left-most (on loop end, this will 
	// be `head`) node holds the type constant
	for !childless {
		// grab children
		left, right = left.SplitNode()
		// push right child's token (has type Node{ty: Id, Token: _})
		varToken := right.(Node).Token
		freeVar := types.Var(varToken)
		s.Push(freeVar)
		// prepare next iteration; last iteration puts type decl name into head
		head, childless = left.HasValue()
	}

	if hasFreeVars { // typeDecl has free variables
		vars = make([]types.Variable[token.Token], s.GetCount())
		// free vars were pushed in reverse order; thus, they are popped in order
		for i := range vars {
			v, _ := s.Pop()
			vars[i] = v
		}
	}

	// set head as name
	name = types.MakeConst(head.Token)
	return name, vars
}

func binaryNodeAsAppOrConst(bnode BinaryRecursiveNode) types.ReferableType[token.Token] {
	if node, ok := bnode.HasValue(); ok {
		return asType(node)
	}
	left, right := bnode.SplitNode()
	return types.Apply[token.Token](binaryNodeAsAppOrConst(left), binaryNodeAsAppOrConst(right))
}

func exprVar(name string) expr.Variable[token.Token] {
	return expr.Var[token.Token](token.Id.Make().AddValue(name))
}

var constrGenMemoLock sync.Mutex

var constructorGenMemo = make(map[int]expr.Function[token.Token], 10)

func abstractForConstructor(depth int) expr.Function[token.Token] {
	if depth < 0 {
		panic("illegal argument: depth < 0")
	}

	constrGenMemoLock.Lock()
	defer constrGenMemoLock.Unlock()

	f, found := constructorGenMemo[depth] // in table?
	if found {
		return f
	} 

	vars := make([]expr.Variable[token.Token], depth)
	var bound expr.Expression[token.Token] = exprVar("tag")
	// expression context must be locked!!!
	glb_cxt.exprMutex.Lock()
	for i := range vars {
		vars[i] = glb_cxt.exprCxt.NewVar()
		bound = expr.Apply[token.Token](bound, vars[i])
	}
	glb_cxt.exprMutex.Unlock() // matching unlock :)

	// create function
	f = expr.Bind[token.Token](exprVar("tag"), vars...).In(bound)
	// now "memoize" and return function
	constructorGenMemo[depth] = f
	return f
}

var arrowConst = types.InfixConst[token.Token](types.MakeConst(token.Arrow.Make().AddValue("->")))

func binaryNodeInFunction(constr BinaryRecursiveNode, rightTy types.Monotyped[token.Token]) (res types.Monotyped[token.Token], tag token.Token, depth int) {
	if node, ok := constr.HasValue(); ok {
		// bottom of left edge of tree, `node` holds the name of the type constructor/tag
		tag = node.Token
		return rightTy, tag, 1 
	}

	left, right := constr.SplitNode()
	rightType := binaryNodeAsAppOrConst(right)
	rightRes := types.Apply[token.Token](arrowConst, rightType, rightTy)
	res, tag, depth = binaryNodeInFunction(left, rightRes)
	depth = depth + 1
	return res, tag, depth
}

// given 
//  unboundType = (Type _),
//  constr = (Con a1 a2 .. aN)
// return
//  res = a1 -> a2 -> .. -> aN -> (Type _), 
//  tag = Con, 
//  depth = N
func constructorToType(unboundType types.Monotyped[token.Token], constr BinaryRecursiveNode) (res types.Monotyped[token.Token], tag token.Token, depth int) {
	if node, ok := constr.HasValue(); ok {
		// at root of tree, this is also the bottom of the left edge of the tree; thus,
		// this is the name of the type constructor/tag
		tag = node.Token
		return unboundType, tag, 0 // return the type itself; constructor is just a tag, no other data
	}
	// constructor has data members; all the data together (plus the constructor tag) makes
	// the `unboundType` type

	left, right := constr.SplitNode()
	// given constructorToType(`Type a`, `Con a1 .. aN`), 
	// endDomain = `aN`, endFunction = `aN -> (Type a)`
	endDomain := binaryNodeAsAppOrConst(right)
	endFunction := types.Apply[token.Token](arrowConst, endDomain, unboundType)
	return binaryNodeInFunction(left, endFunction)
}

func constructorToJudgement(ty types.Monotyped[token.Token], binders []types.Variable[token.Token], constr BinaryRecursiveNode) types.TypeJudgement[token.Token, expr.Expression[token.Token]] {
	// Type a = Con1 a Int | Con2 (Type a) a
	// Type = forall a . Type a
	// Con1_constr = ((\c x y -> c x y) Con1): forall a . a -> Int -> Type a 
	// Con2_constr = ((\c x y -> c x y) Con2): forall a . Type a -> a -> Type a
	//     /\         /\
	//    /\ Int     /\ a
	// Con1 a    Con2 /\
	//            Type  a
	res, tag, depth := constructorToType(ty, constr)
	constructorConstructor := abstractForConstructor(depth)
	constructorTag := Const(tag)
	constructor := constructorConstructor.Apply(constructorTag)
	var judgedType types.Type[token.Token]
	if len(binders) < 1 {
		judgedType = res
	} else {
		judgedType = types.Forall(binders...).Bind(res)
	}
	return types.Judgement(constructor, judgedType)
}

type TypeDefNode struct {
	constType types.Constant[token.Token]
	closedType types.Type[token.Token]
	constructors []types.TypeJudgement[token.Token, expr.Expression[token.Token]]
}

func (n TypeDefNode) SplitType() ([]types.Variable[token.Token], types.ReferableType[token.Token]) {
	// assertion is not guarenteed to work if `n` is initialized outside of 
	// function `initialTypeDef`; but it is assumed here that it is
	
	if poly, ok := n.closedType.(types.Polytype[token.Token]); ok {
		return poly.GetBinders_shallow(), poly.GetBound().(types.ReferableType[token.Token])
	}

	return []types.Variable[token.Token]{}, n.closedType.(types.ReferableType[token.Token])
}

func (n TypeDefNode) Equals(a ast.Ast) bool {
	n2, ok := a.(TypeDefNode)
	if !ok {
		return false
	}
	if !n.constType.Equals(n2.constType) || !n.closedType.Equals(n2.closedType) {
		return false 
	}
	if len(n.constructors) != len(n2.constructors) {
		return false
	}
	for i, con := range n.constructors {
		exp, ty := con.GetExpression(), con.GetType()
		if !exp.Equals(glb_cxt.exprCxt, n2.constructors[i].GetExpression()) {
			return false 
		}
		if !ty.Equals(n2.constructors[i].GetType()) {
			return false
		}
	}
	return true
}

func (n TypeDefNode) NodeType() ast.Type { return TypeDef }

func (n TypeDefNode) InOrderTraversal(f func(itoken.Token)) {
	for _, tok := range n.constType.Collect() {
		f(tok)
	}

	for _, tok := range n.closedType.Collect() {
		f(tok)
	}

	for _, judge := range n.constructors {
		for _, tok := range judge.Collect() {
			f(tok)
		}
	}
}

func typeDefCast(node ast.Ast) TypeDefNode {
	return node.(TypeDefNode)
}

func initialTypeDef(nodes ...ast.Ast) ast.Ast {
	// nodes[0]: TypeDecl, nodes[1]: _, nodes[2]: Constructor
	const typeDeclIndex, _, constructorIndex int = 0, 1, 2

	// split type declaration into constant name and free type binders
	head, binders := typeDeclType(nodes[typeDeclIndex])
	def := TypeDefNode{
		constType: head,
		constructors: make([]types.TypeJudgement[token.Token, expr.Expression[token.Token]], 0, 1),
	}

	// create constant part of monotype part of type for type def
	var open types.ReferableType[token.Token] = head

	// finish create type for defined type
	if len(binders) != 0 {
		// finish creating monotype part of type for defined type
		for _, v := range binders {
			open = types.Apply[token.Token](open, v)
		}
		// close type by binding all free variables in monotype
		def.closedType = types.Forall[token.Token](binders...).Bind(open)
	} else {
		// no free variables; type is closed as is
		def.closedType = head
	}

	return appendConstructorToTypeDef(def, nil, nodes[constructorIndex])
}

func appendConstructorToTypeDef(nodes ...ast.Ast) ast.Ast {
	// nodes[0]: TypeDef, nodes[1]: _, nodes[2]: Constructor
	const typeDefIndex, _, constructorIndex int = 0, 1, 2
	
	def := typeDefCast(nodes[typeDefIndex])
	constructorNode := constructorCast(nodes[constructorIndex])
	// create left-most type constructor and corresponding type
	binders, mono := def.SplitType()
	judgement := constructorToJudgement(mono, binders, constructorNode)
	// add constructor to definition
	def.constructors = append(def.constructors, judgement)
	return def
	
}

var typeDef__typeDecl_Assign_constructor_r = parser.
	Get(initialTypeDef).From(TypeDecl, Assign, Constructor)

var typeDef__typeDef_Bar_constructor_r = parser.
	Get(appendConstructorToTypeDef).From(TypeDef, Bar, Constructor)