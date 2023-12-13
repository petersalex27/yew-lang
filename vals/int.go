package vals

import (
	"strconv"

	"github.com/llir/llvm/ir/constant"
	llirtypes "github.com/llir/llvm/ir/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/types"
)

var integerToken = token.TypeId.Make().AddValue("_Int_")
var integerType = types.MakeConst(integerToken)

type Integer struct {
	*constant.Int
	token.Token
}

func NewInteger(tok token.Token) *Integer {
	i := new(Integer)
	i.Token = tok
	i.FromString(tok.GetValue())
	return i
}

// assigns a value to the underlying constant by converting s to an int
//
// returns error value from strconv.ParseInt
func (i *Integer) FromString(s string) error {
	x, e := strconv.ParseInt(s, 0, 64)
	if e != nil {
		return e
	}

	i.Int = constant.NewInt(llirtypes.I64, x)
	return nil
}

// returns true when p is an integer and i and p have the same constant value
func (i *Integer) Equals(p bridge.PrimInterface[token.Token]) bool {
	i2, ok := p.(*Integer)
	if !ok {
		return false
	}

	return i.Int.X.Int64() == i2.Int.X.Int64()
}

// returns integer token that this value is from
func (i *Integer) Val() token.Token {
	return i.Token
}

// returns _Int_
func (i *Integer) GetType() types.Monotyped[token.Token] {
	return integerType
}
