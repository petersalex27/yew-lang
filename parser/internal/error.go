// =============================================================================
// Author-Date: Alex Peters - November 24, 2023
//
// Content:
// reports syntax errors
//
// Notes: -
// =============================================================================
package internal

import (
	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/status"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
)

// add an error to the parser's errors
func (p *Parser) AppendError(e error) {
	p.Errors = append(p.Errors, e)
}

// returns two tokens: one at the start of the current node, one at the end of the current node
func (p *Parser) GetTokensForShiftEndOfTokens() [2]token.Token {
	node, ok := parser.ParseStackPeek(p.Parser)
	var toks [2]token.Token
	index := 0
	if !ok {
		return toks
	}

	// writes first token at index 0 and last token at index 1
	node.InOrderTraversal(
		func(t itoken.Token) {
			toks[index] = t.(token.Token)
			// this will cause all subsequent tokens to be written
			// to index 1; this will result in the last token seen
			// being written to the last index
			index = 1
		},
	)
	return toks
}

// reports an error that happened when using the "shift" action
func (p *Parser) ShiftError(stat status.Status) {
	var e error
	switch stat {
	case status.EndOfTokens:
		toks := p.GetTokensForShiftEndOfTokens()
		e = errors.Parser(p.GetSource(), toks, errors.UnexpectedEndOfTokens)
	default:
		e = errors.Parser(p.GetSource(), [2]token.Token{}, errors.UnexpectedError)
	}
	p.AppendError(e)
}

// default syntax error
func DefaultError(src source.StaticSource, tok itoken.Token) error {
	return errors.Parser(src, [2]token.Token{tok.(token.Token), tok.(token.Token)}, errors.UnexpectedToken)
}

func (p *Parser) UnexpectedLookAheadToken() {
	tokens := parser.LookAheadTokens(p.Parser)
	var e error
	if len(tokens) < 1 {
		e = errors.Parser(p.GetSource(), [2]token.Token{}, errors.UnexpectedToken)
	} else {
		tok := tokens[0]
		e = DefaultError(p.GetSource(), tok)
	}
	p.AppendError(e)
}