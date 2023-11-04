# Parser README

### Notes on `grammar.bnf`
The notation in `grammar.bnf` is non-standard. It's based on BNF notation but extended to account for things related specifically to shift-reduce parsing

In its most general form:
```bnf
lhs ::= 
  alternative_1 @ 
    (lookahead_1 | lookahead_2 | ..)
  | alternative_2 @
    (lookahead_2_1 | lookahead_2_2 | ..)
  | ..
```
- where `lhs` is a single non-terminal symbol or one of the following special actions: `ERROR`, `SHIFT`
  - when `SHIFT` or `ERROR` are present, `<==` may be used to make it more obvious that these are not replacement rules
- where `alternative` is a sequence of terminal and/or non-terminal symbols separated by a space and with zero of the initial symbols enclosed in parens
  - e.g., `thing ::= (INDENT) ID`
    - this represents a rule that results in the a handle of `INDENT ID` and a replacement of `INDENT thing` where `INDENT` is left on the parse stack and not part of the actual production `thing ::= ID`
- where `lookahead` represents zero or more ***terminal*** symbols
- the `@` may be dropped if there are no lookahead symbols in which case it is assumed to be a rule for _***all***_ lookaheads

#### Examples

- `trait ::= trait INDENT functionDef`
  - replaces the handle `trait 'where' functionDef` with the non-terminal `trait` when the lookahead is anything

- `SHIFT <== exportHead TYPE_ID @ '..'`
  - shifts the lookahead symbol `'..'` onto the parse stack with a handle `exportHead TYPE_ID` resulting in a new stack `exportHead TYPE_ID '..'`

- `functionDecl  ::= (INDENT) funcName`
  - does the production rule `functionDecl ::= funcName` only when the handle is `INDENT funcName` when the lookahead is anything