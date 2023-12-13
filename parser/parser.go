package parser

import (
	"github.com/petersalex27/yew-lang/parser/annotation"
	"github.com/petersalex27/yew-lang/parser/declaration"
	"github.com/petersalex27/yew-lang/parser/export"
	"github.com/petersalex27/yew-lang/parser/imports"
	"github.com/petersalex27/yew-lang/parser/internal"
	"github.com/petersalex27/yew-lang/parser/module"
	typedef "github.com/petersalex27/yew-lang/parser/type-def"
	"github.com/petersalex27/yew-lang/token"
	source "github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
)

// shiftX shifts a token with the token type passed as an argument for `x`. If a token cannot be
// shifted, p.Panicking is set to true
//
// returns true iff a token is shifted
func shiftX(p *internal.Parser, x token.TokenType) bool {
	if !p.LookAhead(x) {
		// TODO: report error
		p.Panicking = true
		return false
	}
	p.Shift()
	return true
}

// TODO: if module is created for `where`, remove this function and make a 
// `where.Shift(*internal.Parser)` function
//
// shiftWhere shifts 'where' token. If it cannot, p.Panicking is set to true
//
// returns true iff 'where' token is shifted
func shiftWhere(p *internal.Parser) bool {
	const x token.TokenType = token.Where
	return shiftX(p, x)
}

// TODO: if module is created for `in`, remove this function and make an
// `in.Shift(*internal.Parser)` function
//
// shiftIn shifts 'in' token. If it cannot, p.Panicking is set to true
//
// returns true iff 'in' token is shifted
func shiftIn(p *internal.Parser) bool {
	const x token.TokenType = token.In
	return shiftX(p, x)
}

func parseTopLevelBinder(p *internal.Parser) {
	if p.Panicking {
		return
	}

	if p.LookAhead(token.Typing) {
		declaration.Parse(p)
	} else {
		definition.Parse(p)
	}
}

// parses source body
func parseModuleBody(p *internal.Parser) {
	for !p.Panicking {
		p.DropNewlines()
		switch p.Next() {
		case token.Id:
			fallthrough
		case token.Infixed:
			p.Shift()
			parseTopLevelBinder(p)
		case token.TypeId:
			p.Shift()
			typedef.Parse(p)
		case token.Trait:
			// TODO
		case token.Comment:
			// TODO
		case token.Alias:
			p.Shift()
			alias.Parse(p)
		case token.Annotation:
			annotation.Parse(p)
		default:
			return
		}
	}
}

// parses (optional) imports, sets p.Panicking to true on failure
func parseImports(p *internal.Parser) {
	if p.Panicking || !p.LookAhead(token.Import) {
		return
	}

	ok := imports.Parse(p)
	if !ok {
		p.Panicking = true
		return
	}

	// TODO: if module is created for `in`, make an `in.Shift(*internal.Parser)` function
	ok = shiftIn(p)
	if !ok {
		p.UnexpectedLookAheadToken()
		p.Panicking = true
	}
}

// parses preamble annotations, sets p.Panicking to true on failure
func parsePreambleAnnotation(p *internal.Parser, annotationParse func(*internal.Parser) bool) {
	if p.Panicking {
		return
	}

	success := annotationParse(p)
	if !success {
		p.UnexpectedLookAheadToken()
		p.Panicking = true
	}
}

// shifts 'where' token, sets p.Panicking to true on failure 
func finishModuleDeclaration(p *internal.Parser) {
	if p.Panicking {
		return
	}

	// TODO: if module is created for `where`, make a `where.Shift(*internal.Parser)` function
	ok := shiftWhere(p)
	if !ok {
		p.UnexpectedLookAheadToken()
		p.Panicking = true
	}
}

// parses export list, sets p.Panicking to true on failure
func parseExportList(p *internal.Parser, doExport bool) {
	if p.Panicking || !doExport {
		return
	}

	success := export.Parse(p) && module.AttachExportList(p)
	p.Panicking = !success
}

// parses module declaration (module name and its export list), sets p.Panicking to true on failure
func parseModuleDeclaration(p *internal.Parser) {
	if p.Panicking {
		return 
	}

	// parse module name
	doExport, ok := module.Parse(p)
	if !ok {
		p.Panicking = true
	}

	parseExportList(p, doExport)	
	finishModuleDeclaration(p)
}

// parses preamble of source file (import annotations, imports/uses, module annotations, and top-
// level module).
func parseSourcePreamble(p *internal.Parser) bool {
	parsePreambleAnnotation(p, annotation.ImportAnnotationParse)
	parseImports(p)
	parsePreambleAnnotation(p, annotation.ModuleAnnotationParse)
	parseModuleDeclaration(p)
	return !p.Panicking
}

func run(p *internal.Parser) {
	if !parseSourcePreamble(p) {
		return
	}
	parseModuleBody(p)
}

// begins parsing
func Parse(source source.Source, tokens []itoken.Token) {
	p := internal.InitInternal(source, tokens)
	run(p)
}
