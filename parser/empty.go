package parser

import (
	"github.com/petersalex27/yew-packages/parser"
)

// == empty node reduction rules ==============================================

var empty__LeftParen_RightParen_r = parser.
	Get(giveTypeToTokenProductionGen(Empty)).
	From(LeftParen, RightParen)
