package vals

import (
	"strconv"

	"github.com/llir/llvm/ir/constant"
	llirtypes "github.com/llir/llvm/ir/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/types"
)

var charToken = token.TypeId.Make().AddValue("_Char_")
var charType = types.MakeConst(charToken)

type Char Integer

func NewChar(tok token.Token) *Char {
	c := new(Char)
	c.Token = tok
	c.FromString(tok.GetValue())
	return c
}

// assigns a value to the underlying constant by converting s to a int8
//
// returns error value from strconv.ParseInt
func (c *Char) FromString(s string) error {
	x, e := strconv.ParseInt(s, 0, 8)
	if e != nil {
		return e
	}

	c.Int = constant.NewInt(llirtypes.I8, x)
	return nil
}

// returns true when p is a char and c and p have the same constant value
func (c *Char) Equals(p bridge.PrimInterface[token.Token]) bool {
	c2, ok := p.(*Char)
	if !ok {
		return false
	}

	return c.Int.X.Int64() == c2.Int.X.Int64()
}

// returns char token that this value is from
func (c *Char) Val() token.Token {
	return c.Token
}

// returns _Char_
func (c *Char) GetType() types.Monotyped[token.Token] {
	return charType
}
