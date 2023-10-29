package parser

import "github.com/petersalex27/yew-lang/token"

// take input from lexer and append INDENT(0) token to each line that does not start with
// one
func prepareInputForParsing(tokens *[][]token.Token) {
	for lineIndex := range (*tokens) {
		toks := (*tokens)[lineIndex]
		// skip if no tokens in line or line starts with indent
		if len(toks) == 0 || toks[lineIndex].GetType() == uint(token.Indent) {
			continue
		}

		// create token that exists at current line number and start of line
		itok := token.Indent.Make().AddValue("").SetLineChar(lineIndex+1,1)
		tok := itok.(token.Token)
		// add indent token
		(*tokens)[lineIndex] = append([]token.Token{tok}, (*tokens)[lineIndex]...)
	}
}

func Parse([][]token.Token) {
	
}