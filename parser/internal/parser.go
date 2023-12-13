// =============================================================================
// Author-Date: Alex Peters - November 26, 2023
//
// Content:
// Parser struct and methods
//
// Notes: -
// =============================================================================
package internal

import (
	token "github.com/petersalex27/yew-lang/token"
	parser "github.com/petersalex27/yew-packages/parser"
	ast "github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/parser/status"
	source "github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/stack"
)

// internal parser data
type Parser struct {
	parser.Parser
	parseStack *stack.SaveStack[ast.Ast]
	shift     func(*Parser)
	reduce    func(*Parser, parser.ProductionOrder, bool, bool)
	Errors    []error
	Panicking bool
}

// removes leading indent tokens from unseen tokens
// func (p *Parser) DropIndents() {
// 	for p.LookAhead(token.Indent) {
// 		parser.DropNext(p.Parser)
// 	}
// }

// removes leading `n` tokens w/ `ty` token type from unseen tokens.
//
// negative integer for `n` drops all leading, unseen tokens w/ type `ty`
//
// returns number of tokens dropped
func (p *Parser) DropN(n int, ty token.TokenType) int {
	dropped := 0
	for (dropped != n) && p.LookAhead(ty) {
		parser.DropNext(p.Parser)
		dropped++
	}
	return dropped
}

// see description for `DropN`; when n >= 0, this does the same thing as DropN but sets p.Panicking
// to true if exactly `n` tokens with type `ty` are not dropped. Additionally, an error message is
// reported saying an unexpected token was found. However, if n < 0, then p.Panicking will not
// be changed
func (p *Parser) RequireDropN(n int, ty token.TokenType) {
	dropped := p.DropN(n, ty)
	if n < 0 || n == dropped {
		return
	}

	// if already panicking, don't report another error
	if p.Panicking {
		return
	}
	// set panic, and report unexpected token
	p.Panicking = true
	p.UnexpectedLookAheadToken()
}

// removes leading newline tokens from unseen tokens
func (p *Parser) DropNewlines() {
	p.DropN(-1, token.Newline)
}

// remove next lookahead and return it; this does NOT shift the lookahead token!
//
// call ignored when p.Panicking == true
func (p *Parser) GetAndRemoveLookAhead() (tok token.Token) {
	if p.Panicking {
		return
	}
	tokenInterface, ok := parser.RemoveNext(p.Parser)
	if !ok {
		// can happen when p.Parser is acted on async w/o locking both lookahead check and call to this
		// and/or lookahead's existence was not confirmed
		panic("failed to retrieve next token")
	}

	return tokenInterface.(token.Token)
}

// shift function for when parser isn't panicking
func parserInternalShift(p *Parser) {
	if stat := p.Parser.Shift(); stat.NotOk() {
		// report error and set internal functions to internal skip functions
		p.ShiftError(stat)
		p.shift, p.reduce = parserInternalSkip, parserInternalReductionSkip
		p.Panicking = true
	}
}

// reduce function for when parser isn't panicking
func parserInternalReduce(p *Parser, productions parser.ProductionOrder, loop bool, baseOkay bool) {
	var stat status.Status
	var applied bool
	countAsApplied := baseOkay
	ok := baseOkay
	for {
		stat, applied = p.Parser.Reduce(productions)
		ok = applied && stat.IsOk()
		if !ok {
			break
		}
		countAsApplied = true
		if loop {
			continue
		}
		break
	}

	if countAsApplied {
		if applied && stat.IsOk() {
			return
		}

		if loop && stat.Is(status.EndAction) {
			return
		}
	}

	if !applied {
		// TODO: report error based on handle
	} else if stat.NotOk() {
		// TODO: report error based on handle

	}

	p.shift, p.reduce = parserInternalSkip, parserInternalReductionSkip
	p.Panicking = true
}

// replaces `parserInternalShift` when parser has an error
func parserInternalSkip(*Parser) {}

// replaces `parserInternalReduce` when parser has an error
func parserInternalReductionSkip(*Parser, parser.ProductionOrder, bool, bool) {}

func (p *Parser) ShiftN(n uint) {
	for i := uint(0); i < n; i++ {
		p.shift(p)
	}
}

// moves first element from unseen token queue onto top of parse stack
func (p *Parser) Shift() { p.shift(p) }

// transforms some number of nodes on the top of the parse stack into a single
// node
func (p *Parser) Reduce(productions parser.ProductionOrder) { p.reduce(p, productions, false, false) }

func (p *Parser) ReductionLoop(productions parser.ProductionOrder) {
	p.reduce(p, productions, true, false)
}

func (p *Parser) ReductionAllowingEmptyLoop(productions parser.ProductionOrder) {
	p.reduce(p, productions, true, true)
}

func initInternalBase(src source.Source, tokenStream []itoken.Token) *Parser {
	out := new(Parser)
	out.shift = parserInternalShift
	out.reduce = parserInternalReduce

	return out
}

// initializes internal parser data
func InitInternal(src source.Source, tokenStream []itoken.Token) *Parser {
	out := initInternalBase(src, tokenStream)
	// out.Parser = parser.
	// 	NewParser().LA(1).
	// 	Load(tokenStream, src, nil, nil)
	const initialCap uint = 32
	out.parseStack = stack.NewSaveStack[ast.Ast](initialCap)
	return out
}

// mostly for mocking rule sets allowing asts to be returned
func InitInternalWithTable(src source.Source, tokenStream []itoken.Token, table parser.ReductionTable) *Parser {
	p := InitInternal(src, tokenStream)
	p.Parser = parser.
		NewParser().
		LA(1).
		UsingReductionTable(table).
		Load(tokenStream, src, nil, nil)
	return p
}

// returns next unseen token
func (p *Parser) Next() token.TokenType {
	tok := parser.LookAheadTokens(p.Parser)
	if len(tok) != 1 {
		return token.TokenType(ast.None)
	}
	return token.TokenType(tok[0].GetType())
}

// returns true iff next unseen token has the token type `tokenType`
func (p *Parser) LookAhead(tokenType token.TokenType) bool {
	if p.Panicking {
		return false
	}

	tok := parser.LookAheadTokens(p.Parser)
	if len(tok) != 1 {
		return false
	}
	return tok[0].GetType() == uint(tokenType)
}

// parseListComma parses comma in some kind of list.
//
// `commaGrammarProductions` is the ordered list of grammar rules to use for
// parsing the comma
//
// The first return value is true iff `p.Panicking` is false; the second return
// value is true iff grammar expects the list to end (i.e., iff the next token
// is not a comma)
func (p *Parser) ParseListComma(commaGrammarProductions parser.ProductionOrder) (ok bool, endOfList bool) {
	p.DropNewlines()
	if !p.LookAhead(token.Comma) {
		return !p.Panicking, true
	}
	p.Shift()
	p.Reduce(commaGrammarProductions)
	return !p.Panicking, false
}

// wrapper for very common sequence of calls:
//   - drop leading newline tokens
//   - shift next token ("drop" action implies this is a non-newline token)
//   - apply reduction action from one of the productions passed as an arg
func (p *Parser) DropShiftReduce(productions parser.ProductionOrder) {
	p.DropNewlines()
	p.Shift()
	p.Reduce(productions)
}
