package vals

import (
	"strconv"

	"github.com/llir/llvm/ir/constant"
	llirtypes "github.com/llir/llvm/ir/types"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/types"
)

var floatToken = token.TypeId.Make().AddValue("_Float_")
var floatType = types.MakeConst(floatToken)

type Float struct {
	*constant.Float
	token.Token
}

func NewFloat(tok token.Token) *Float {
	f := new(Float)
	f.Token = tok
	f.FromString(tok.GetValue())
	return f
}

// assigns a value to the underlying constant by converting s to a float
//
// returns error value from strconv.ParseFloat
func (f *Float) FromString(s string) error {
	x, e := strconv.ParseFloat(s, 64)
	if e != nil {
		return e
	}

	f.Float = constant.NewFloat(llirtypes.Double, x)
	return nil
}

// returns true when p is a float and f and p have the same constant values
func (f *Float) Equals(p bridge.PrimInterface[token.Token]) bool {
	f2, ok := p.(*Float)
	if !ok {
		return false
	}

	x, _ := f.Float.X.Float64()
	x2, _ := f2.Float.X.Float64()
	return x == x2
}

// returns float token that this value is from
func (f *Float) Val() token.Token {
	return f.Token
}

// returns _Float_
func (f *Float) GetType() types.Monotyped[token.Token] {
	return floatType
}
