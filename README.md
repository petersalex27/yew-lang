*NOTE: this file is a work in progress*

## Design Goals
There are three primary design goals for Yew. 

Code written in Yew is ... 
1. Reusable/Flexible
2. Bug resistant
3. Able to use cool, new programming language ideas

## Project Structure

- `./` : main module
  - `update.lua` : useful for when required modules' packages aren't being properly updated
    - `usage: lua update.lua [args ..]`
    - args
      - `t` : gets modules needed to build the tests for the updated packages
      - `n` : prevents script from running `go mod tidy` after getting the updates
    - `lua update.lua` runs:
      ```
      go clean -cache
      go clean -modcache
      go get -u ./...
      go mod tidy
      ```
- `./errors` : module for handling compiler errors
- `./lexer` : module for lexical analysis (text-to-tokens)
- `./lib` : standard library for Yew
- `./parser` : module for parser
- `./token` : module for tokens, the output of the lexical analysis step
- `./util`
- `./ir` : module for translating abstract syntax trees to the intermediate representation (i.e., translation to LLVM-IR)
- `./assets` : module for assets (yew icon image, for example)