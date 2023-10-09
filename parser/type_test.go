package parser

import (
	"fmt"
	"testing"

	"github.com/petersalex27/yew-packages/parser"
	"github.com/petersalex27/yew-packages/parser/ast"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
	"yew.lang/main/errors"
	"yew.lang/main/lexer"
	"yew.lang/main/token"
)

func tytok(name string, line, char int) token.Token {
	return token.TypeId.Make().AddValue(name).SetLineChar(line, char).(token.Token)
}

func arrowtok(line, char int) token.Token {
	return token.Arrow.Make().AddValue("->").SetLineChar(line, char).(token.Token)
}

func idtok(name string, line, char int) token.Token {
	return token.Id.Make().AddValue(name).SetLineChar(line, char).(token.Token)
}

func commatok(line, char int) token.Token {
	return token.Arrow.Make().AddValue(",").SetLineChar(line, char).(token.Token)
}

func arrtok(line, char int) token.Token {
	return token.TypeId.Make().AddValue("[]").SetLineChar(line, char).(token.Token)
}

func getTestPath(name string, num int) string { 
	return fmt.Sprintf("./token-streams/%s/test-%s-%d.yew", name, name, num)
}

func TestTypeRules(t *testing.T) {
	tests := []struct{
		name string; num int
		expect ast.AstRoot
	}{
		{
			"type-id", 1,
			ast.AstRoot{TypeNode{Monotype, types.MakeConst(tytok("Int", 1, 1))}},
		},
		{
			"var", 1,
			ast.AstRoot{TypeNode{Monotype, types.Var(idtok("a", 1, 1))}},
		},
		{
			"frees", 1,
			ast.AstRoot{
				TypeNode{
					Monotype, 
					types.Apply(
						types.ReferableType[token.Token](types.Var(idtok("a", 1, 1))), 
						types.Monotyped[token.Token](types.Var(idtok("b", 1, 3))),
					),
				},
			},
		},
		{
			"function", 1,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeInfixConst(arrowtok(1, 5))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 1))),
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 8))),
					),
				},
			},
		},
		{
			"function", 2,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeInfixConst(arrowtok(1, 5))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 1))),
						types.Monotyped[token.Token](
							types.Apply[token.Token](
								types.ReferableType[token.Token](types.MakeInfixConst(arrowtok(1, 12))), 
								types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 8))),
								types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 15))),
							),
						),
					),
				},
			},
		},
		{
			"app", 1,
			ast.AstRoot{
				TypeNode{
					Monotype, 
					types.Apply(
						types.ReferableType[token.Token](types.MakeConst(tytok("Maybe", 1, 1))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 7))),
					),
				},
			},
		},
		{
			"app", 2,
			ast.AstRoot{
				TypeNode{
					Monotype, 
					types.Apply(
						types.ReferableType[token.Token](types.MakeConst(tytok("Maybe", 1, 1))), 
						types.Monotyped[token.Token](types.Var(idtok("a", 1, 7))),
					),
				},
			},
		},
		{
			"app", 3,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeConst(tytok("Either", 1, 1))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 8))),
					).Merge(
						types.Monotyped[token.Token](types.MakeConst(tytok("Bool", 1, 12))),
					),
				},
			},
		},
		{
			"app", 4,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeConst(tytok("Either", 1, 1))), 
						types.Monotyped[token.Token](types.Var(idtok("a", 1, 12))),
					).Merge(
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 8))),
					),
				},
			},
		},
		{
			"app", 5,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeConst(tytok("Either", 1, 1))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 8))),
					).Merge(
						types.Monotyped[token.Token](types.Var(idtok("a", 1, 12))),
					),
				},
			},
		},
		{
			"list", 1,
			ast.AstRoot{
				TypeNode{
					Monotype, 
					types.MakeConst(tytok("Int", 1, 1)),
				},
			},
		},
		{
			"list", 2,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeInfixConst(commatok(1, 5))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 2))),
					).Merge(
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 6))),
					),
				},
			},
		},
		{
			"list", 3,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeInfixConst(commatok(1, 5))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 2))),
					).Merge(
						types.Apply(
							types.ReferableType[token.Token](types.MakeInfixConst(commatok(1, 10))), 
							types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 7))),
						).Merge(
							types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 12))),
						),
					),
				},
			},
		},
		{
			"list", 4,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeInfixConst(commatok(1, 5))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 2))),
					).Merge(
						types.Apply(
							types.ReferableType[token.Token](types.MakeInfixConst(commatok(1, 10))), 
							types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 7))),
						).Merge(
							types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 12))),
						),
					),
				},
			},
		},
		{
			"paren", 1,
			ast.AstRoot{
				TypeNode{
					Monotype, 
					types.MakeConst(tytok("Int", 1, 1)),
				},
			},
		},
		{
			"paren", 2,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeInfixConst(arrowtok(1, 6))), 
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 2))),
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 9))),
					),
				},
			},
		},
		{
			"paren", 3,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](types.MakeInfixConst(arrowtok(1, 14))), 
						types.Monotyped[token.Token](
							types.Apply(
								types.ReferableType[token.Token](types.MakeInfixConst(arrowtok(1, 6))), 
								types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 2))),
								types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 9))),
							),
						),
						types.Monotyped[token.Token](types.MakeConst(tytok("Int", 1, 17))),
					),
				},
			},
		},
		{
			"paren", 4,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Apply(
						types.ReferableType[token.Token](
							types.MakeConst[token.Token](tytok("Maybe", 1, 1)),
						),
						types.Monotyped[token.Token](
							types.Apply(
								types.ReferableType[token.Token](types.Var(idtok("a", 1, 8))), 
								types.Monotyped[token.Token](types.Var(idtok("b", 1, 10))),
							),
						),
					),
				},
			},
		},
		{
			"forall", 1,
			ast.AstRoot{
				TypeNode{
					Polytype,
					types.Forall[token.Token](
						types.Var(idtok("a", 1, 8)),
					).Bind(
						types.Var(idtok("a", 1, 12)),
					),
				},
			},
		},
		{
			"forall", 2,
			ast.AstRoot{
				TypeNode{
					Polytype,
					types.Forall[token.Token](
						types.Var(idtok("a", 1, 8)),
					).Bind(
						types.Apply[token.Token](
							types.MakeConst(tytok("Maybe", 1, 12)),
							types.Var(idtok("a", 1, 18)),
						),
					),
				},
			},
		},
		{
			"forall", 3,
			ast.AstRoot{
				TypeNode{
					Polytype,
					types.Forall[token.Token](
						types.Var(idtok("a", 1, 8)),
					).Bind(
						types.Apply[token.Token](
							types.MakeInfixConst(arrowtok(1, 14)),
							types.Var(idtok("a", 1, 12)),
							types.Var(idtok("a", 1, 17)),
						),
					),
				},
			},
		},
		{
			"forall", 3,
			ast.AstRoot{
				TypeNode{
					Polytype,
					types.Forall[token.Token](
						types.Var(idtok("a", 1, 8)),
					).Bind(
						types.Apply[token.Token](
							types.MakeInfixConst(arrowtok(1, 14)),
							types.Var(idtok("a", 1, 12)),
							types.Var(idtok("a", 1, 17)),
						),
					),
				},
			},
		},
		{
			"forall", 4,
			ast.AstRoot{
				TypeNode{
					Polytype,
					types.Forall[token.Token](
						types.Var(idtok("a", 1, 8)),
					).Bind(
						types.Apply[token.Token](
							types.MakeInfixConst(commatok(1, 14)),
							types.Var(idtok("a", 1, 13)),
							types.Var(idtok("a", 1, 16)),
						),
					),
				},
			},
		},
		{
			"forall", 5,
			ast.AstRoot{
				TypeNode{
					Polytype,
					types.Forall[token.Token](
						types.Var(idtok("a", 1, 8)),
						types.Var(idtok("b", 1, 10)),
					).Bind(
						types.MakeConst(tytok("Int", 1, 14)),
					),
				},
			},
		},
		{
			"array", 1,
			ast.AstRoot{
				TypeNode{
					Monotype,
					types.Index[token.Token](
						types.Apply[token.Token](
							types.MakeEnclosingConst(1, arrtok(1,1)),
							types.MakeConst(tytok("A", 1, 2)),
						), 
						mkFreeJudge(glb_cxt.exprCxt.NewVar(), types.MakeConst(tytok("Uint",0,0))),
					),
				},
			},
		},
	}

	reInit()

	for i, test := range tests {
		pth := getTestPath(test.name, test.num)
		lex := lexer.NewLexer(pth)
		toks, es := lexer.RunLexer(lex)
		if len(es) != 0 {
			t.Fatal(testutil.TestFail2("lex", nil, es, i))
		}
		
		//toks := lexer.CastTokens(tmp)
		src := parser.MakeSource(pth, lexer.GetSourceRaw(lex)...)
		p := parser.
			NewParser().
			LA(1).
			UsingReductionTable(typeReduceTable).
			Load(toks, src, nil, nil)//.LogActions()

		root := p.Parse()
		//println(p.FlushLog())

		if p.HasErrors() {
			es := p.GetErrors()
			errors.PrintErrors(es...)
			t.Fatal(testutil.TestFail2("errors", nil, es, i))
		}

		if !root.Equals(test.expect) {
			rootStr := ast.GetOrderedString(root)
			expStr := ast.GetOrderedString(test.expect)
			t.Fatal(testutil.TestFail(expStr, rootStr, i))
		}
	}
}