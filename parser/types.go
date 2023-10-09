package parser

import (
	"github.com/petersalex27/yew-packages/parser/ast"
	"yew.lang/main/token"
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
	Wildcard  = ast.Type(token.Wildcard)
	Empty     = ast.Type(token.Empty)
	Comment   = ast.Type(token.Comment)
	At        = ast.Type(token.At)
	Colon     = ast.Type(token.Typing)
	Assign    = ast.Type(token.Assign)
	Bar       = ast.Type(token.Bar)
	Arrow     = ast.Type(token.Arrow)
	Backslash = ast.Type(token.Backslash)
	Dot       = ast.Type(token.Dot)
	DotDot    = ast.Type(token.DotDot)

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
	Thunked = ast.Type(token.Thunked)

	shiftKeywords_class = Let
	// shift keywords
	Let    = ast.Type(token.Let)
	Class  = ast.Type(token.Class)
	Use    = ast.Type(token.Use)
	Family = ast.Type(token.Family)
	Forall = ast.Type(token.Forall)
	Mapall = ast.Type(token.Mapall)
	Module = ast.Type(token.Module)

	Import               = ast.Type(token.Import)
	Of                   = ast.Type(token.Of)
	From                 = ast.Type(token.From)
	In                   = ast.Type(token.In)
	Where                = ast.Type(token.Where)
	Qualified            = ast.Type(token.Qualified)
	Struct               = ast.Type(token.Struct)
	Derives              = ast.Type(token.Derives)
	Do                   = ast.Type(token.Do)
	LAST_TERMINAL_TYPE__ = ast.Type(token.LAST_TYPE__)
	// non-terminal types
	Expr ast.Type = iota + LAST_TERMINAL_TYPE__
	Application
	ApplicationId
	FreeApplication
	LetDeclaration
	Context
	ClassDeclaration
	ClassDefinition
	InstanceDeclaration
	InstanceDefinition
	StructDefinition
	ConstructorDefinition
	Constructor
	UnionPair
	FunctionDefinition
	FunctionDecl
	InfixedDefinition
	TypeDefinition
	Annotation
	ImportList
	UseList
	QualifiedList
	Namespace
	ModuleDeclaration
	ModuleDefinition
	ExportList
	AnonymousFunction
	Case
	Pattern
	Array
	LiteralArray
	ArrayValHead
	LitArrHead
	Data
	Block
	Group
	Name
	Param
	FuncName
	AnyName
	Literal
	TypeJudgement // e: (T ..)
	FreeTyping
	Polytype    // (forall a ..) . (T ..)
	Dependtype  // (mapall (a: A) ..) . (B ..; ..)
	Dependtyped // type node for a type T, dep 𝚷(a: A)B(a) s.t. 𝚷(a: A)B(a) ⊑ T
	Monotype    // type node for monotypes
	Type        // most generic "type" node
	MonoList    // (T .., U .., ..)
	MonoTail
	FreeVar         // var
	TypeApp         // T U ..
	DependIndexHead // T ..;
	DependInstance  // T ..; e ..
	VarJudgement
	AppJudgement
	TupleJudgement
	AnonJudgement
	MapHead
	ForallHead
	ArrayHead
	_last_type_
)
