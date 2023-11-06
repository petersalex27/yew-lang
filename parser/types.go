package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser/ast"
)

const (
	value_class = IntValue
	// value class
	IntValue    ast.Type = ast.Type(token.IntValue)
	StringValue          = ast.Type(token.StringValue)
	CharValue            = ast.Type(token.CharValue)
	FloatValue           = ast.Type(token.FloatValue)

	// spacing
	Indent = ast.Type(token.Indent)

	// other
	Wildcard        = ast.Type(token.Wildcard)
	Comment         = ast.Type(token.Comment)
	AnnotationToken = ast.Type(token.Annotation)
	Colon           = ast.Type(token.Typing)
	Assign          = ast.Type(token.Assign)
	Bar             = ast.Type(token.Bar)
	Arrow           = ast.Type(token.Arrow)
	Backslash       = ast.Type(token.Backslash)
	Dot             = ast.Type(token.Dot)
	DotDot          = ast.Type(token.DotDot)

	// grouping
	LeftParen    = ast.Type(token.LeftParen)
	RightParen   = ast.Type(token.RightParen)
	LeftBracket  = ast.Type(token.LeftBracket)
	RightBracket = ast.Type(token.RightBracket)
	LeftBrace    = ast.Type(token.LeftBrace)
	RightBrace   = ast.Type(token.RightBrace)

	// separators
	SemiColon = ast.Type(token.SemiColon)
	Comma     = ast.Type(token.Comma)

	name_class = Symbol
	// name class
	Symbol  = ast.Type(token.Symbol)
	Id      = ast.Type(token.Id)
	TypeId  = ast.Type(token.TypeId)
	Infixed = ast.Type(token.Infixed)

	shiftKeywords_class = Let
	// shift keywords
	Let    = ast.Type(token.Let)
	Traot  = ast.Type(token.Trait)
	Use    = ast.Type(token.Use)
	Family = ast.Type(token.Family)
	Forall = ast.Type(token.Forall)
	Mapall = ast.Type(token.Mapall)
	Module = ast.Type(token.Module)

	Import               = ast.Type(token.Import)
	Match                = ast.Type(token.Match)
	From                 = ast.Type(token.From)
	In                   = ast.Type(token.In)
	Where                = ast.Type(token.Where)
	Qualified            = ast.Type(token.Qualified)
	Of                   = ast.Type(token.Of)
	Derives              = ast.Type(token.Derives)
	Do                   = ast.Type(token.Do)
	LAST_TERMINAL_TYPE__ = ast.Type(token.LAST_TYPE__)
	// non-terminal types

	// expression, e.g.,
	//		1 + 1
	Expr ast.Type = iota + LAST_TERMINAL_TYPE__
	// literal value or array, e.g.,
	//		1
	Val
	// expr expr
	Application
	// let x = y in z
	LetExpr
	// z where x = y
	WhereExpr
	// --@MyAnnotation args and stuff
	Annotation
	// trait C a
	TraitDeclaration
	// trait C a where f: a -> a
	TraitDefinition
	// C of MyType
	InstanceDeclaration
	// C of MyType where f mt = mt
	InstanceDefinition
	// x { member -> SomeType a }
	StructDefinition
	// Brach (BinaryTree a) (BinaryTree a)
	Constructor
	// definitions and instances of functions and traits and type definitions
	//		trait C a where f: a -> a
	//		MyType a = MyData a Int
	//		C of MyType where f mt = mt
	//		myFunc: a -> a
	//		myFunc x = x
	Definitions
	// module myModule ( myFunc,
	ExportHead
	// top level node
	// 		module myModule ( myFunc ) where
	//			myFunc: a -> a
	//			myFunc x = x
	Source
	// f x: Int -> Int
	FunctionDefinition
	// f x
	FunctionHead
	// f
	FunctionDecl
	// f x: Int -> Int = x
	Function
	// BinaryTree a = Leaf a | Brach (BinaryTree a) (BinaryTree a)
	TypeDef
	// BinaryTree a
	TypeDecl
	// import ( myModule, anotherModule, .. )
	ImportList
	// use import ( myModule, anotherModule, .. )
	UseList
	// qualified import ( myModule, anotherModule )
	QualifiedList
	// use myModule in
	Namespace
	// module myModule
	ModuleDeclaration
	// module myModule ( myFunc, (%>) )
	ModuleDefinition
	// module myModule ( myFunc, (%>),
	ExportList
	// module myModule ( myFunc, MyType )
	ExportDone
	// (\x -> x)
	AnonymousFunction
	// Just x -> x
	Case
	// pather match:
	//		g when
	//			Just x -> f x
	//			Nothing -> Nothing
	PatternMatch
	// Just x
	Pattern
	// Just
	PatternC
	// [1, a+b, b]
	Empty
	Array
	// [1, 2, 3]
	LiteralArray
	// [1, a+b,
	ArrayValHead
	// [1, 2,
	LitArrHead
	// Just 1
	Data
	// myName
	Name
	// myName
	Param
	// (>>=)
	FuncName
	// 1
	Literal
	// x: T
	TypeJudgement
	// forall a b . MyType a b
	Polytype
	// mapall (a: A) (b: B) . (MyTypeFunc; a b)
	Dependtype
	// type node for a type T, dep 𝚷(a: A)B(a) s.t. 𝚷(a: A)B(a) ⊑ T
	Dependtyped
	// Maybe Int
	Monotype
	// most generic "type" node
	Type
	// , Maybe Int, Int)
	TupleType
	// (SomeType a;
	DependIndexHead
	// (SomeType a; (x+1) y)
	DependInstance
	// a: Int
	VarJudgement
	// mapall (a: A) (b: B)
	DependHead
	// mapall (a: A) (b: B)
	DependBinders
	// forall a b
	PolyHead
	// forall a b
	PolyBinders
	//
	IndentExprBlock
	//
	Error
	_last_type_
)
