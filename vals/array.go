package vals

// import (
// 	"errors"

// 	"github.com/llir/llvm/ir/constant"
// 	llirtypes "github.com/llir/llvm/ir/types"
// 	"github.com/petersalex27/yew-lang/token"
// 	"github.com/petersalex27/yew-packages/bridge"
// 	"github.com/petersalex27/yew-packages/expr"
// 	"github.com/petersalex27/yew-packages/types"
// )

// var arrayToken = token.TypeId.Make().AddValue("[]")
// var _bind_me_ = expr.Var(token.Id.Make().AddValue("n"))

// func arrayType(elemType types.Monotyped[token.Token]) types.DependentTypeInstance[token.Token] {
// 	arr := types.Apply[token.Token](types.MakeEnclosingConst(1, arrayToken), elemType)
// 	j := types.TypedJudge[token.Token, expr.Referable[token.Token], types.Type[token.Token]](_bind_me_, uintegerType)
// 	return types.Index[token.Token](arr, j)
// }

// type Array struct {
// 	*constant.Array
// 	headToken token.Token
// }

// func NewArray(headToken token.Token, elems ...expr.Expression[token.Token]) (*Array, error) {
// 	a := new(Array)
// 	a.headToken = headToken

// 	var lltype llirtypes.Type = llirtypes.I8Ptr // just a place holder

// 	if len(elems) > 0 {
// 		p, ok := elems[0].(bridge.Prim[token.Token])
// 		if !ok {
// 			return nil, errors.New("unexpected expression")
// 		}
		
// 		ty := p.Val.GetType()
// 		if ty.Equals(integerType) {
// 			lltype = llirtypes.I64
// 		} else if ty.Equals(charType) {
// 			lltype = llirtypes.I8
// 		} else if ty.Equals(stringType) {

// 		}
// 	}

// 	arr := constant.NewArray(
// 		llirtypes.NewArray(uint64(len(elems)), lltype),
// 		([]constant.Constant{})...,
// 	)

// 	for _, elem := range elems {
// 		p, ok := elem.(bridge.Prim[token.Token])
// 		if !ok {
// 			return nil, errors.New("unexpected expression")
// 		}

// 	//	if p.Val.
// 	}

// 	s.FromString(tok.GetValue())
// 	return s
// }

// // assigns a value to the underlying constant by converting z to a byte slice
// //
// // returns nil
// func (s *String) FromString(z string) error {
// 	in := []byte(z)
// 	s.CharArray = constant.NewCharArray(in)
// 	return nil
// }

// // returns true when p is a string and s and p have the same constant value
// func (s *String) Equals(p bridge.PrimInterface[token.Token]) bool {
// 	s2, ok := p.(*String)
// 	if !ok {
// 		return false
// 	}

// 	return string(s.CharArray.X) == string(s2.CharArray.X)
// }

// // returns string token that this value is from
// func (s *String) Val() token.Token {
// 	return s.Token
// }

// // returns _String_
// func (*String) GetType() types.Monotyped[token.Token] {
// 	return stringType
// }
