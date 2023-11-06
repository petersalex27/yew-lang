package parser

import (
	"runtime"
	"sync"

	"github.com/petersalex27/yew-lang/errors"
	"github.com/petersalex27/yew-lang/lexer"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	modLexer "github.com/petersalex27/yew-packages/lexer"
	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/source"
	"github.com/petersalex27/yew-packages/util/stack"

	//itoken "github.com/petersalex27/yew-packages/token"
	itoken "github.com/petersalex27/yew-packages/token"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util"
)

// TODO: define Yew's grammar here
var grammar = parser.
	ForTypesThrough(_last_type_).
	UseReductions(
	// TODO
	).
	Finally(
		// TODO
		parser.Order(),
	)

var programForest *forest = new(forest)

var createLexer func(string) *modLexer.Lexer

var ticketer *ticketVendor

var __ticket_vendor_lock__ sync.Mutex

var __parse_jobs_lock__ sync.Mutex
var parseJobs sync.WaitGroup

func init() {
	// initialize ticketer
	ticketer = new(ticketVendor)
	ticketer.ticketCap = uint(util.Max(runtime.NumCPU(), 4))
	ticketer.remainingTickets = ticketer.ticketCap
	ticketer.ticketAvailability = sync.NewCond(&__ticket_vendor_lock__)
}

// test path determines display name in errors. return value should be defered
// once its received by caller
func initForTest(testPath string, testGrammar parser.ReductionTable) (restoreAfterTest func()) {
	// save globals
	forestRestore := programForest
	grammarRestore := grammar
	createLexerRestore := createLexer

	restoreAfterTest = func() {
		// restore globals
		createLexer = createLexerRestore
		programForest = forestRestore
		grammar = grammarRestore
	}

	// set globals
	programForest = new(forest)
	createLexer = func(raw string) *modLexer.Lexer {
		return lexer.NewStringLexer(testPath, raw)
	}
	grammar = testGrammar

	return restoreAfterTest
}

// resets globals that need to be reset before running another test
func resetForTest() {
	programForest = new(forest)
}

type ticketVendor struct {
	ticketCap          uint
	mu                 sync.Mutex
	remainingTickets   uint
	ticketAvailability *sync.Cond
}

func newContext(path string) *Context {
	p := new(Context)
	// mutex as pointers so they can be set to nil to guard
	// against accidental escaping of them
	p.typeMutex, p.exprMutex = new(sync.Mutex), new(sync.Mutex)
	p.exprCxt = expr.NewContext[token.Token]()
	p.typeCxt = types.NewContext[token.Token]()
	p.path = path
	p.indentStack = *stack.NewStack[string](8)
	//p.src = parser.MakeSource(path)
	return p
}

// give a ticket to `p` allowing it to parse it's given path or wait until a
// ticket is available and then try to let `p` parse again
//
// IMPORTANT: this should only be called from the main thread!
func (p *Context) giveTicket() {
	ticketer.mu.Lock()

	// the ticket waiting area :)
	for {
		// `p` wants to parse, but no tickets available
		if ticketer.remainingTickets == 0 {
			// unlock so other threads can give tickets back
			ticketer.mu.Unlock()
			// wait for ticket availability
			ticketer.ticketAvailability.Wait()
			// lock to prevent race condition at ticket avail. check
			ticketer.mu.Lock()
			continue // check for avail. again
		}

		ticketer.remainingTickets-- // claim ticket
		ticketer.mu.Unlock()        // unlock so other threads can get/return tickets
		break                       // exit waiting area
	}

	go p.parse() // run parser
}

type Context struct {
	path      string
	src       parser.MinSource            // source being parsed
	typeMutex *sync.Mutex                 // type mutex lock
	exprMutex *sync.Mutex                 // expression mutex lock
	typeCxt   *types.Context[token.Token] // type context
	exprCxt   *expr.Context[token.Token]  // expression context
	indentStack stack.Stack[string]				// indentation stack
}

// Sets context to nil and returns ticket to ticket vendor.
//
// Context `cxt` shouldn't have been shared with any variables that are still
// "alive", but it is set to nil as a guard to hopefully catch any escaped
// living variables if they exist
func (cxt *Context) returnTicket() {
	// guards, see returnTicket documentation
	cxt.typeCxt, cxt.exprCxt = nil, nil
	cxt.exprMutex, cxt.typeCxt = nil, nil

	ticketer.mu.Lock()

	// reclaim parser's ticket
	ticketer.remainingTickets++

	// broadcast (to main thread) that there is a ticket available (it waits when
	// no tickets are available and there are still paths left to parse)
	if ticketer.remainingTickets == 1 {
		// let the main thread know a parser can claim a ticket
		ticketer.ticketAvailability.Broadcast()
	}

	// unclock ticket lock
	ticketer.mu.Unlock()

	// mark parse job as complete
	parseJobs.Done()
}

// parses source file given to parser receiver
//
// IMPORTANT: this should never be entered in the main thread!
func (cxt *Context) parse() {
	// return ticket to vendor at end of parse. This is important!!
	defer cxt.returnTicket()

	path := cxt.path

	// create lexer and check that its creation was successful
	lex := createLexer(path)
	if lex == nil {
		// TODO: signal to main thread that parser failed
		return // error
	}

	itoks, es := lexer.RunLexer(lex)
	if len(es) != 0 {
		// TODO: signal to main thread that parser failed
		addErrors(es...)
		return
	}

	// raw source split at each line
	rawSource := lexer.GetSourceRaw(lex)
	// make a copy of the source for the context to use for errors
	cxt.src = parser.MakeSource(path, rawSource...)

	// default error generating function
	defaultError := func(src source.StaticSource, tok itoken.Token) error {
		return errors.Parser(src, []token.Token{tok.(token.Token)}, errors.UnexpectedToken)
	}
	// default error for syntax error when something unknown goes wrong
	// 	TODO: make non-default value (i.e., make non-nil)
	var couldNotParseError error = nil

	// capture `p` inside of action functions; this enables communication between
	// github.com/petersalex27/yew-packages/parser module (the module that runs
	// the parsing logic) and this package's context.
	actionClosures := cxt.generateActions()

	// build parser
	p := parser.
		NewParser().
		// declare parser as lookahead-1 type
		LA(1).
		// load Yew's grammar
		UsingReductionTable(grammar).
		// load: tokens, source, and default errors
		Load(itoks, cxt.src, defaultError, couldNotParseError).
		// load parser production actions
		Attach(actionClosures...)

	// run parser, sourceRoot is root node of AST (if parse is successful)
	sourceRoot := p.Parse()

	if p.HasErrors() {
		// TODO: signal parser has errors
		addErrors(p.GetErrors()...)
		return
	}

	// all paths are unique, and all paths have already been added to the program
	// forest; thus, each parser may write to the program forest w/o worrying
	// about data races
	programForest.trees[path] = sourceRoot

	// parse of source file is successful :)
}

// map of paths to their respective AST roots
type forest struct {
	sync.Mutex
	trees map[string]ast.AstRoot
}

// For each path in paths, Parse(...string) parses the source file at `path`
func Parse(paths ...string) {
	// allocate program ast forest
	programForest.trees = make(map[string]ast.AstRoot, len(paths))

	// only parse unique paths by adding all the paths first and then
	// using the forest to interate through the paths
	for _, path := range paths {
		programForest.trees[path] = nil
	}

	// set number of parse jobs to complete
	parseJobs.Add(len(programForest.trees))

	// launch jobs
	for path := range programForest.trees {
		newContext(path).giveTicket() // run parser
	}

	// force main thread to wait for all parse jobs to complete
	parseJobs.Wait()

	// TODO: check for errors and actually do something with parse result
}
