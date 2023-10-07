package parser

import (
	"github.com/petersalex27/yew-packages/parser"
)

var typeReduceTable = parser.
	ForTypesThrough(_last_type_).
	UseReductions(
		parser.
			LA(TypeId).
			ForN(1, Id).
				Then(qualifier_rules.Union(apptype_rules, var_rules, monotype_singlesRules)).
				ElseShift(),

		parser.LA(Forall).
			Shift(),

		parser.LA(Mapall).
			Then(polytype_rules).
			ElseShift(),

		parser.LA(LeftParen).
			ForN(1, LeftBracket,).
				Shift(),

		parser.LA(RightParen).
			Then(all_mono_rules.Union(monoList_rules)),

		parser.LA(RightBracket).
			Then(array_rules),

		parser.LA(Arrow).
			Then(all_mono_rules_no_function),

		parser.LA(Comma).
			Then(all_mono_rules_no_close),

		parser.LA(Dot). 
			Then(var_rules.Union(qualifier_rules)),
	).
	Finally(alltypes_rules)

	// v0.0.0-20231003031741-ea929e3ab360
	// v0.0.0-20231003185600-667c1eb14677

var qualifier_rules = forall_rules.Union(mapall_rules)

var alltypes_rules = 
	qualifier_rules.Union(
	//dependTyped_rules,
	polytype_rules,
	depend_rules,
	all_mono_rules,
)

var all_mono_rules = 
	all_mono_rules_no_function.Union(monotype_function_rules)

var all_mono_rules_no_close = 
	apptype_rules.Union(
	var_rules,
	monotype_singlesRules,
	array_rules,
	dependIndexHead_rules,
	dependInstance_rules,
	monotype_function_rules,
	//monoList_rules,
)

var all_mono_rules_no_function = 
	apptype_rules.Union(
	var_rules,
	monotype_singlesRules,
	monotype_close_rules,
	array_rules,
	dependIndexHead_rules,
	dependInstance_rules,
)

var apptype_rules = parser.RuleSet(
	typeApp__TypeId_mono_r,
	typeApp__typeApp_mono_r,
	typeApp__mono_mono_r,
)

var mapall_rules = parser.RuleSet(
	mapallHead__Mapall_varJudgement_r,
	mapallHead__mapallHead_varJudgement_r,
)

var forall_rules = parser.RuleSet(
	forallHead__Forall_var_r,
	forallHead__forallHead_var_r,
)

var polytype_rules = parser.RuleSet(
	polytype__forallHead_Dot_depend_r,
	polytype__forallHead_Dot_mono_r,
)

var depend_rules = parser.RuleSet(
	depend__mapallHead_Dot_TypeId_r,
	depend__mapallHead_Dot_typeApp_r,
)

var dependTyped_rules = parser.RuleSet(
	dependTyped__dependIndexHead_judge_r,
	dependTyped__dependInstance_r,
	dependTyped__depend_r,
)

var dependIndexHead_rules = parser.RuleSet(
	dependIndexHead__TypeId_SemiColon_r,
	dependIndexHead__typeApp_SemiColon_r,
	dependIndexHead__var_SemiColon_r,
)

var dependInstance_rules = parser.RuleSet(
	dependInstance__arrayHead_RightBracket_r,
	dependInstance__dependIndexHead_expr_r,
)

var array_rules = parser.RuleSet(
	arrayHead__LeftBracket_TypeId_r,
	arrayHead__LeftBracket_dependIndexHead_expr_r,
	arrayHead__LeftBracket_dependIndexHead_judge_r,
	arrayHead__LeftBracket_typeApp_r,
	arrayHead__LeftBracket_var_r,
)

var var_rules = parser.RuleSet(var__Id_r)

var monotype_singlesRules = parser.RuleSet(
	monotype__enclosed_r,
	monotype__var_r,
	monotype__typeApp_r,
	monotype__dependInstance_r,
	monotype__TypeId_r,
)

var monoList_rules = parser.RuleSet(
	monoList__mono_Comma_monoList_r,
	monoList__mono_Comma_monoTail_r,
	monoTail__mono_Comma_r,
	monoTail__mono_r,
)

var monotype_close_rules = parser.RuleSet(
	monotype__LeftParen_monoTail_RightParen_r,
	monotype__LeftParen_monoList_RightParen_r,
)

var monotype_function_rules = parser.RuleSet(
	monotype__mono_Arrow_mono_r,
)