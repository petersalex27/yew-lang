package types

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser/ast"
)

const (
	// values
	IntValue    ast.Type = ast.Type(token.IntValue)
	StringValue          = ast.Type(token.StringValue)
	CharValue            = ast.Type(token.CharValue)
	FloatValue           = ast.Type(token.FloatValue)
	// spacing
	Newline = ast.Type(token.Newline)
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
	Comma     = ast.Type(token.Comma)
	// names
	Id      = ast.Type(token.Id)
	TypeId  = ast.Type(token.TypeId)
	Infixed = ast.Type(token.Infixed)
	// keywords
	Let       = ast.Type(token.Let)
	Trait     = ast.Type(token.Trait)
	Use       = ast.Type(token.Use)
	Family    = ast.Type(token.Family)
	Forall    = ast.Type(token.Forall)
	Mapval    = ast.Type(token.Mapval)
	Module    = ast.Type(token.Module)
	Import    = ast.Type(token.Import)
	Match     = ast.Type(token.Match)
	From      = ast.Type(token.From)
	In        = ast.Type(token.In)
	Where     = ast.Type(token.Where)
	Qualified = ast.Type(token.Qualified)
	Of        = ast.Type(token.Of)
	Derives   = ast.Type(token.Derives)
	Alias     = ast.Type(token.Alias)
	// end of token types
	LAST_TERMINAL_TYPE__ = ast.Type(token.LAST_TYPE__)

	Application ast.Type = iota + LAST_TERMINAL_TYPE__
	Function
	LetExpr
	RecExpr
	Variable
	WhereDef

	Judgement

	Monotype
	DependentType
	Polytype

	TypeDef

	Constructor

	ImportContext
	ImportElement

	ModuleDef
	ExportList
	TypeExport

	Pattern

	Annotation
)

const END_TOKEN token.TokenType = token.TokenType(LAST_TERMINAL_TYPE__)
