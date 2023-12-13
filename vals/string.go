package vals

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/types"
)

var stringToken = token.TypeId.Make().AddValue("_String_")
var stringType = types.MakeConst(stringToken)

type String struct {
	*constant.CharArray
	token.Token
}

func NewString(tok token.Token) *String {
	s := new(String)
	s.Token = tok
	s.FromString(tok.GetValue())
	return s
}

// assigns a value to the underlying constant by converting z to a byte slice
//
// returns nil
func (s *String) FromString(z string) error {
	in := []byte(z)
	s.CharArray = constant.NewCharArray(in)
	return nil
}

// returns true when p is a string and s and p have the same constant value
func (s *String) Equals(p bridge.PrimInterface[token.Token]) bool {
	s2, ok := p.(*String)
	if !ok {
		return false
	}

	return string(s.CharArray.X) == string(s2.CharArray.X)
}

// returns string token that this value is from
func (s *String) Val() token.Token {
	return s.Token
}

// returns _String_
func (*String) GetType() types.Monotyped[token.Token] {
	return stringType
}
