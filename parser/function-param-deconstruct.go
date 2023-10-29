package parser

import (
	"github.com/petersalex27/yew-lang/token"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/util/stack"
)

func deconstructHelper(
	exp expr.Expression[token.Token], 
	stk *stack.SaveStack[deconstructionInstruction], 
	extract, move deconstructionInstruction,
) (restore bool) {
	if exp == nil {
		return true
	}

	data, ok := exp.(expr.Application[token.Token])
	if !ok {
		if _, ok := exp.(expr.Variable[token.Token]); ok {
			stk.Push(extract) // write extract instruction
		}
		return false
	}

	stk.Save()
	stk.Push(move) // write movement instruction
	if restore = deconstruct(data, stk); restore {
		stk.Return() // clear any written instructions
	} else {
		stk.Rebase() // incorp. written instructions
	}
	return
}

func deconstruct(data expr.Application[token.Token], stk *stack.SaveStack[deconstructionInstruction]) (restore bool) {
	left, right := data.Split()

	// restore == true iff nothing is written
	leftRestore, rightRestore := 
		deconstructHelper(left, stk, extractLeft, moveLeft), 
		deconstructHelper(right, stk, extractRight, moveRight)
	return leftRestore && rightRestore // Do NOT restore if either side writes something
}

// for each param do instructions: 
//		([][]deconstructionInstruction{
//			{inst1, instr2, ...}, /* these are the instructions to deconstruct the first param */
//			{inst1, instr2, ...}, /* these are the instructions to deconstruct the second param */
//			..., /* etc., etc. */
//		})[paramIndex][instructionIndex]
func recursiveDeconstruction(params []expr.Expression[token.Token]) [][]deconstructionInstruction {
	if len(params) == 0 {
		return nil // nothing to deconstruct
	}

	// return value
	out := make([][]deconstructionInstruction, len(params))

	// stack keeping track of deconstruction
	stk := stack.NewSaveStack[deconstructionInstruction](8, 8) // 8 is an arbitrary choice

	// for each param, record how to deconstruct it
	for i, param := range params {
		// ok is false iff no instructions are written
		data, ok := param.(expr.Application[token.Token])
		if ok {
			// attemp to deconstruct application (applications represent data 
			// patterns in this context).
			ok = !deconstruct(data, stk) // wrote instructions?
		}

		if !ok { // no instructions written
			out[i] = []deconstructionInstruction{skipParam} // skip param
			continue // end interation
		} 
		
		// copy instructions to return value
		size := int(stk.GetCount()) // number of instructions written
		out[i] = make([]deconstructionInstruction, size) // allocate slice for instructions
		tmp, _ := stk.MultiCheck(int(stk.GetCount())) // return instructions
		copy(out[i], tmp) // copy instructions to output
		stk.Clear(stk.GetCount()) // clear stack
	}

	return out
}