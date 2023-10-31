package parser

import (
	"sync"

	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/types"
)

// Holds global contexts for types and expressions. The most important function
// this serves is generating unique names
type glb_cxt_t struct {
	path string // path being parsed
	typeMutex sync.Mutex // type mutex lock
	exprMutex sync.Mutex // expression mutex lock
	errorsLock sync.Mutex // errors mutex lock
	typeCxt   *types.Context[token.Token] // type context
	exprCxt   *expr.Context[token.Token] // expression context
	es []error
}

// globalContext__ holds global contexts for types and expressions. Remember to
// lock the respective mutexs before using either context w/in
var globalContext__ *glb_cxt_t

func addErrors(es ...error) {
	globalContext__.errorsLock.Lock()
	defer globalContext__.errorsLock.Unlock()

	globalContext__.es = append(globalContext__.es, es...)
}

// intended for use in tests where the generators for unique names need to be
// reset
func reInit() {
	globalContext__ = new(glb_cxt_t)

	globalContext__.exprCxt = expr.NewContext[token.Token]().SetNameMaker(
		func(s string) token.Token {
			return token.Id.Make().AddValue(s)
		},
	)

	globalContext__.typeCxt = types.NewContext[token.Token]().SetNameMaker(
		func(s string) token.Token {
			return token.Id.Make().AddValue(s)
		},
	)

	globalContext__.es = []error{}
}

func init() {
	reInit()
}
