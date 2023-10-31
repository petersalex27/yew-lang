package parser

import (
	"sync"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/lexer"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	itoken "github.com/petersalex27/yew-packages/token"
)

type AnnotationNode struct {
	// used to lock annotation from parent thread while its being parsed or
	// tokenized
	sync.WaitGroup
	// name of annotation
	annotationName token.Token
	// body of annotation
	annotationBody []itoken.Token
}

// == annotation reduction rules ==============================================

var annotation__Annotation_r = parser.
	ActionRule().
	Get(annotationReduction).
	From(AnnotationToken)

// == annotation reductions ===================================================

func annotationReduction(call func(name string) func(any), nodes ...ast.Ast) ast.Ast {
	const annotationIndex int = 0
	annotToken := GetToken(nodes[annotationIndex])
	annot := new(AnnotationNode)

	// if main thread tries to access members of annotation, force it to wait
	// until annotation is lexed
	annot.WaitGroup.Add(1)

	// make sure main thread does not end before annotation finishes its work
	parseJobs.Add(1)

	// annotation will have its wait group decremented inside call to
	// lexAnnotation
	go func() {
		lexAnnotation(call, annot, annotToken)
		parseJobs.Done()
	}()

	return annot
}

// == annotation implementation of ast.Ast ====================================

func (annot *AnnotationNode) Equals(a ast.Ast) bool {
	annot.WaitGroup.Wait() // wait until annotation is lexed

	annot2, ok := a.(*AnnotationNode)
	if !ok {
		return false
	}

	annot2.WaitGroup.Wait() // wait until annotation is lexed

	// both annotations have been read, so now they can be compared for equality
	if !EqualsToken(annot.annotationName, annot2.annotationName) {
		return false
	}

	if len(annot.annotationBody) != len(annot2.annotationBody) {
		return false
	}

	for i, tok := range annot.annotationBody {
		if !EqualsToken(tok, annot2.annotationBody[i]) {
			return false
		}
	}

	return true
}

func (annot *AnnotationNode) InOrderTraversal(f func(itoken.Token)) {
	annot.WaitGroup.Wait() // make sure annotation is done being lexed

	f(annot.annotationName)

	for _, tok := range annot.annotationBody {
		f(tok)
	}
}

func (*AnnotationNode) NodeType() ast.Type { return Annotation }

// == annotation utils ========================================================

// TODO: need a way to report errors
func lexAnnotation(call func(name string) func(any), annot *AnnotationNode, annotToken token.Token) {
	// once function is done, wait group can be decremented.
	defer annot.WaitGroup.Done()

	// create a lexer to read the annotation's info
	lex := lexer.NewStringLexer(globalContext__.path, annotToken.GetValue())
	if lex == nil {
		return // TODO: signal to parent thread that an error happend
	}

	// run lexer
	_, es := lexer.RunLexer(lex)
	if len(es) != 0 {
		addErrors(es...)
		return // TODO: signal to parent thread that an error happend
	}

	// apply offsets to tokens' lines and chars based on annotation token's so
	// the tokens have the correct line and char number
	lex.ApplyOffset(annotToken.GetLineChar())
	// now get correctly offset tokens
	itoks := lex.GetTokens()

	if len(itoks) < 1 { // were tokens returned?
		// no tokens recovered
		var src parser.MinSource
		call(getSource)(&src) // request source code
		e := errors.Parser(src, []token.Token{annotToken}, errors.ExpectedNamedAnnotation)
		addErrors(e)
		return
	}

	// first token is annotation name, remaining tokens are annotation body
	const annotNameIndex, annotBodyStart int = 0, 1
	justName := lexer.CastTokens(itoks[annotNameIndex:annotBodyStart])
	(*annot).annotationName = justName[0]
	(*annot).annotationBody = itoks[annotBodyStart:]
}
