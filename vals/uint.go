package vals

import (
	"strconv"

	"github.com/llir/llvm/ir/constant"
	llirtypes "github.com/llir/llvm/ir/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/types"
)

var uintegerToken = token.TypeId.Make().AddValue("_Uint_")
var uintegerType = types.MakeConst(integerToken)

type Uinteger Integer

func NewUinteger(tok token.Token) *Uinteger {
	u := new(Uinteger)
	u.Token = tok
	u.FromString(tok.GetValue())
	return u
}

// assigns a value to the underlying constant by converting s to an int
//
// returns error value from strconv.ParseInt
func (u *Uinteger) FromString(s string) error {
	x, e := strconv.ParseUint(s, 0, 64)
	if e != nil {
		return e
	}

	u.Int = constant.NewInt(llirtypes.I64, int64(x))
	return nil
}

// returns true when p is an integer and i and p have the same constant value
func (u *Uinteger) Equals(p bridge.PrimInterface[token.Token]) bool {
	u2, ok := p.(*Uinteger)
	if !ok {
		return false
	}

	return u.Int.X.Int64() == u2.Int.X.Int64()
}

// returns integer token that this value is from
func (u *Uinteger) Val() token.Token {
	return u.Token
}

// returns _Uint_
func (u *Uinteger) GetType() types.Monotyped[token.Token] {
	return uintegerType
}
