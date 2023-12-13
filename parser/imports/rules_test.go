// =============================================================================
// Author-Date: Alex Peters - December 01, 2023
//
// Content: tests import rules
// =============================================================================
package imports

import (
	"testing"

	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/inf"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestProduceInitialImportElemHelper(t *testing.T) {
	// create mock input
	qualifier := inf.NameQualified
	nameToken := token.Id.Make().AddValue("_")
	nameNode := ast.TokenNode(nameToken)

	expected := &ImportElementNode{
		QualificationType: qualifier,
		Name:              nameToken,
		As:                nameToken,
		From:              "",
	}

	// test
	actual := produceInitialImportElemHelper(qualifier, nameNode)
	if !expected.Equals(actual) {
		t.Fatal(testutil.FailMessage(expected, actual))
	}
}

func TestProduceImportElemFrom(t *testing.T) {
	const fromValueString string = "test"

	// = create mock input =======================================================
	// modified and returned node
	importElem := &ImportElementNode{}
	// setting this to nil at least makes sure that this value is not deref. in
	// produceImportElemFrom
	var __ ast.Ast = nil
	fromValue := token.StringValue.Make().AddValue(fromValueString)
	// node used for from's value
	fromValueNode := ast.TokenNode(fromValue)

	expected := &ImportElementNode{
		From: fromValueString,
	}

	// test
	actual := produceImportElemFrom(importElem, __, fromValueNode)
	if !expected.Equals(actual) {
		t.Fatal(testutil.FailMessage(expected, actual))
	}
}
