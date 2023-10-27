package parser

import "github.com/petersalex27/yew-packages/parser"

// Note: this rule requires a leading INDENT but leaves the required 
// INDENT on the stack as is!
//	(INDENT) functionDecl <- INDENT funcName
var functionDecl__Indent_funcName_r = parser.
	Get(rewrapNodeRule(FunctionDecl)).
	When(Indent).From(FuncName)