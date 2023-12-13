handle, LA(1), action, rule

𝛆, `module`, shift, -
𝛆, ANNOT, shift, -

ANNOT, ANNOT, shift, -
ANNOT, `module`, setEnv, `𝛆 ::= ANNOT`

`module`, ID, shift, -

`module` ID, `(`, shift, -

`module` ID `(`, ID, reduce, `export ::= 'module' ID '('`
`module` ID `(`, TYPE_ID, reduce, `export ::= 'module' ID '('`
`module` ID `(`, SYMBOL, reduce, `export ::= 'module' ID '('`
`module` ID `(`, INFIXED, reduce, `export ::= 'module' ID '('`

`module` ID, `where`, reduce, `module ::= 'module' ID`

export, ID, shift, -
export, TYPE_ID, shift, -
export, SYMBOL, shift, -
export, INFIXED, shift, -

export ID, `,`, shift, -
export ID, `)`, shift, -

export ID `,`, `)`, shift, -

export ID `,` `)`, `where`, reduce, `module ::= export ID ',' ')'`
export ID `)`, `where`, reduce, `module ::= export ID ')'`

export SYMBOL, `,`, shift, -
export SYMBOL, `)`, shift, -

export SYMBOL `,`, `)`, shift, -

export SYMBOL `,` `)`, `where`, reduce, `module ::= export SYMBOL ',' ')'`
export SYMBOL `)`, `where`, reduce, `module ::= export SYMBOL ')'`

export INFIXED, `,`, shift, -
export INFIXED, `)`, shift, -

export INFIXED `,`, `)`, shift, -

export INFIXED `,` `)`, `where`, reduce, `module ::= export INFIXED ',' ')'`
export INFIXED `)`, `where`, reduce, `module ::= export INFIXED ')'`

export TYPE_ID, `..`, shift, -
export TYPE_ID, `(`, shift, -

export TYPE_ID `(`, `)`, reduce, `typeExport ::= TYPE_ID '('`
export TYPE_ID `(`, TYPE_ID, reduce, `typeExport ::= TYPE_ID '('`
export TYPE_ID `(`, `..`, shift, -

export typeExport, `)`, shift, - 
export typeExport, TYPE_ID, shift, -

export TYPE_ID `(` `..`, `)`, shift, -

export typeExport `)`, `,`, reduce, `export ::= typeExport ')'`
export typeExport `)`, `)`, reduce, `export ::= typeExport ')'`

export typeExport TYPE_ID, `,`, shift, -
export typeExport TYPE_ID, `)`, reduce, `export ::= typeExport TYPE_ID`

export typeExport TYPE_ID `,`, TYPE_ID, reduce, `typeExport ::= typeExport TYPE_ID ','`
export typeExport TYPE_ID `,`, `)`, reduce, `typeExport ::= typeExport TYPE_ID ','`

export TYPE_ID `(` `..` `)`, `,`, reduce, `export ::= export TYPE_ID '(' '..' ')'`
export TYPE_ID `(` `..` `)`, `)`, reduce, `export ::= export TYPE_ID '(' '..' ')'`

export TYPE_ID `..`, `,`, reduce, `export ::= export TYPE_ID '..'`
export TYPE_ID `..`, `)`, reduce, `export ::= export TYPE_ID '..'`

export, `)`, shift, -

export `)`, `where`, reduce, `module ::= export ')'`

module, `where`, shift, -

module `where` imports, `in`, reduce, `moduleDec ::= module 'where' imports`

`import` indent(n) ID, `as`, shift, -

`import` indent(n) ID `as`, ID, shift, -

`import` indent(n) ID `as`, INDENT(_), shift, -

`import` indent(n) ID `as` INDENT(_), ID, shift, -

`import` indent(n) ID `as` INDENT(_) ID, INDENT(\_), dropLA, -

`import` indent(n) ID `as` INDENT(_) ID INDENT(n), INDENT(n), reduce, `mod ::= ID 'as' INDENT(_) ID`

`import` indent(n) ID, `in`, reduce, `mod ::= ID`
`import` indent(n) ID, INDENT(n), reduce, `mod ::= ID`

`import` indent(n) mod, `in`, reduce, `imports ::= 'import' indent(n) mod`

`import` indent(n) ID, INDENT(n), shift, -

`import` indent(n) ID INDENT(n), ID, 

`use` indent(n) ID

`import` indent(n) `qualified` indent(m) ID