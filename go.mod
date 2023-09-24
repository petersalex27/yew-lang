module yew.lang/main

go 1.20

replace (
	alex.peters/yew => ../yew-packages
	alex.peters/yew/expr => ../yew-packages/expression
	alex.peters/yew/lexer => ../yew-packages/lexer
	alex.peters/yew/source => ../yew-packages/source
	alex.peters/yew/token => ../yew-packages/token
	alex.peters/yew/types => ../yew-packages/types2
	alex.peters/yew/util => ../yew-packages/util
	alex.peters/yew/str => ../yew-packages/strings
	alex.peters/yew/errors => ../yew-packages/errors
)

require (
	alex.peters/yew/lexer v0.0.0-00010101000000-000000000000
	alex.peters/yew/source v0.0.0-00010101000000-000000000000
	alex.peters/yew/token v0.0.0-00010101000000-000000000000
	alex.peters/yew/str v0.0.0-00010101000000-000000000000
	alex.peters/yew/util v0.0.0-00010101000000-000000000000
	alex.peters/yew/errors v0.0.0-00010101000000-000000000000
)
