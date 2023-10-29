package lexer

import (
	"testing"

	"github.com/petersalex27/yew-packages/lexer"
	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-lang/util"
)

func TestAdvanceChar(t *testing.T) {
	tests := []struct{
		line, char int
		expect byte
		stat source.Status
	}{
		{1,2,' ',source.Ok}, {1,3,' ',source.Ok}, {1,4,' ',source.Ok}, {2,1,'\n',source.Ok},
		{2,2,' ',source.Ok}, {2,3,' ',source.Ok}, {2,4,' ',source.Ok}, {2,4,0,source.Eof},
	}
	lex := lexer.NewLexer(lexerWhitespace, 0, 0, 1)
	lex.SetSource([]string{`   `,`   `})
	lex.SetPath("./test-lex-advance-char-position.yew")

	for i, test := range tests {
		actual, stat := lex.AdvanceChar()
		if !stat.Is(test.stat) {
			t.Fatal(util.TestFail2("stat", test.stat, stat, i))
		}
		if actual != test.expect {
			t.Fatalf(util.TestFail2("byte", test.expect, actual, i))
		}
		
		line, char := lex.GetLineChar()
		if test.line != line {
			t.Fatal(util.TestFail2("line", test.line, line, i))
		}
		if test.char != char {
			t.Fatalf(util.TestFail2("char", test.char, char, i))
		}
	}
}

func TestPositioning(t *testing.T) { 

}