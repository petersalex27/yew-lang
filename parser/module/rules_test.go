package module

import (
	"testing"

	. "github.com/petersalex27/yew-lang/parser/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/util/testutil"
)

func TestProduceInitialHelper(t *testing.T) {
	for _, variable := range []bool{false, true} {
		nameToken := token.Id.Make().AddValue("test")
		nameTokenNode := ast.TokenNode(nameToken)
		var __ ast.Ast = nil
		ty := ast.Type(0)

		expect := &ModuleSourceNode{
			Name: nameToken,
			exportAll: variable,
			Type: ty,
		}

		actual := produceInitialHelper(ty, variable, __, nameTokenNode)

		if !expect.Equals(actual) {
			t.Fatal(testutil.FailMessage(*expect, *actual))
		}
	}
}

// all produceModuleAndExportList should do is change the type from the old type to `ModuleDef`
func TestProduceModuleAndExportList(t *testing.T) {
	var __ ast.Ast = nil
	inputExportList := &ModuleSourceNode{
		Type: ModuleDef+1, // this is just some type that isn't ModuleDef
	}

	expect := &ModuleSourceNode{
		Type: ModuleDef,
	}

	actual := produceModuleAndExportList(inputExportList, __)

	if !expect.Equals(actual) {
		t.Fatal(testutil.FailMessage(expect, actual))
	}
}