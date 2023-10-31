package parser

import (
	"testing"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestAnnotation(t *testing.T) {
	table := parser.
		ForTypesThrough(_last_type_).
		UseReductions().
		Finally(
			parser.Order(annotation__Annotation_r),
		)

	// 4,12,17,20,24
	annotName := token.TypeId.Make().AddValue("MyAnnot").SetLineChar(1, 4)
	thisToken := token.Id.Make().AddValue("this").SetLineChar(1, 12)
	isToken := token.Id.Make().AddValue("is").SetLineChar(1, 17)
	theToken := token.Id.Make().AddValue("the").SetLineChar(1, 3)
	bodyToken := token.Id.Make().AddValue("body").SetLineChar(1, 24)

	tests := []struct {
		nodes    []ast.Ast
		src      source.StaticSource
		expected ast.AstRoot
	}{
		{
			[]ast.Ast{
				ast.TokenNode(token.Annotation.Make().AddValue(
					"MyAnnot this is the body",
				)),
			},
			parser.MakeSource(
				"test/parser/annotation",
				"--@MyAnnot this is the body",
			),
			ast.AstRoot{
				&AnnotationNode{
					annotationName: annotName.(token.Token),
					annotationBody: []itoken.Token{
						thisToken, isToken, theToken, bodyToken,
					},
				},
			},
		},
	}

	for i, test := range tests {
		// make context
		cxt := newContext("test/parser/annotation")
		// generate `cxt` captured actions
		actionClosures := cxt.generateActions()

		p := parser.
			NewParser().
			LA(1).
			UsingReductionTable(table).
			Load(nil, test.src, nil, nil).
			Attach(actionClosures...).
			InitialStackPush(test.nodes...)

		actual := p.Parse()

		// must wait for parser to finish
		parseJobs.Wait()

		// errors?
		if p.HasErrors() {
			errors.PrintErrors(p.GetErrors()...)
			t.Fatal(testutil.TestFail2("errors", nil, p.GetErrors(), i))
		}

		// result?
		if !actual.Equals(test.expected) {
			t.Fatal(testutil.TestFail2("equality", test.expected, actual, i))
		}
	}
}
