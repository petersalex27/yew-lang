# Keywords
The following words are reserved by the compiler and serve specific syntatic roles described below.
```
┌─────────┬─────────┬─────────┬─────────┐
│  class  │  from   │ module  │  where  │
├─────────┼─────────┼─────────┼─────────┤
│ derives │ import  │   of    │         │
├─────────├─────────┼─────────┼─────────┤
│   do    │   in    │qualified│         │
├─────────┼─────────┼─────────┼─────────┤
│ family  │   let   │   use   │         │
├─────────┼─────────┼─────────┼─────────┤
│ forall  │ mapall  │ struct  │         │
└─────────┴─────────┴─────────┴─────────┘
```
### `class`
denotes the start of a type class definition
```
class Monad f => 
  (>>=): f a -> (a -> f b) -> f b
  unit: a -> f a
```
### `derives`
1. creates a default type class implementation for the preceding type
```
FileMode derives (Str, Eq)
```
2. signifies that the implementations of the type class preceding `derives` can derive the type class following `derives`
```
ApplicativeFunctor derives Monad =>
  (>>=) g f = (pure f) <*> g
  unit x = pure x
```
### `family`
denotes a type family definition
```
family LlvmType t = 
  LlvmType Int 
  | LlvmType Byte
  | LlvmType Float
  | LlvmType Pointer
  | LlvmType Bool
  | LlvmType [LlvmType t]
```
### `forall`
binds type variables in an explicit polytype
```
forall a b . a -> b
```
### `from`
denotes from which implemented type class another type class is derived from
```
List derives (Monad from Applicative, Eq, Ord)
```
### `import`
imports a module whose symbols (except for infix and suffix symbols) are qualified with the module's name
```
import sys
```
### `in`
gives context for the expression to the right
```
let x = 3 in x + 1
```
### `let`
binds an identifier to a value
```
let false = 
  let id = (\x -> x) in 
  (\_ -> id)
```
### `mapall`
binds kind variables in an explicit dependent type
```
forall a . mapall (n: Uint) . (Array a; n)
```
### `module`
declares a file as part of a given module
```
module prelude
```
### `of`
matches an expression via cases
```
(>>=) g f = g of
  Just x -> f x
  Nothing -> Nothing
```
### `qualified`
denotes that an import statment should qualify all of a modules's imported symbols (including infix and suffix symbols)
```
import qualified prelude
```
### `use`
`use` denotes a name-space where qualifiers are not required for the given module(s). It can be used in front of an import statement to import all symbol's exported by the given module as unqualified (i.e., the name-space exists for the whole file),
```
use import io
```
or it can be used as context for `in`
```
use (prelude, io) in
  file << write (cast result: [Byte])
```
### `where`
contextualizes the preceding expression with the following expression
```
(.) f g = compose f g where
  compose = (\f g x -> f (g x))
```