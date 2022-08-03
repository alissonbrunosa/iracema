package compile

import "iracema/bytecode"

/*
* Basic block represents a maximal-length sequence of branch-free code.
* It ends with jump or return instr, when a basicblock ends with an
* unconditional jump, it will be considered a block with a fallthrough.
* This makes the basicblock referenced by the next pointer the branch
* in case the conditional is not met.
*
*      ┌──basicblock──┐
*      │              │
*      │              │
*      │              │
*    ┌─┤ JUMP_IF_TRUE │
*    │ └──────┬───────┘
*    │        └──────────┐
*    │ ┌──basicblock──┐  │
*    │ │              ◄──┘
*    │ │              │
*    │ │              │
*    │ │              │
*    │ └──────────────┘
*    │
*    │ ┌──basicblock──┐
*    └─►              │
*      │              │
*      │              │
*      │              │
*      └──────────────┘
****/

type basicblock struct {
	instrs    []*instr
	offset    int
	next      *basicblock
	hasReturn bool
	visited   bool
	reachable bool
}

func (b *basicblock) isDone() bool {
	index := len(b.instrs) - 1
	if index < 0 {
		return false
	}

	ins := b.instrs[index]
	return ins.opcode == bytecode.Jump ||
		ins.opcode == bytecode.JumpIfTrue ||
		ins.opcode == bytecode.JumpIfFalse ||
		ins.opcode == bytecode.Return
}

func (b *basicblock) hasFallthrough() bool {
	// no instrs we consider a fallthrough
	if len(b.instrs) == 0 {
		return true
	}

	index := len(b.instrs) - 1
	ins := b.instrs[index]
	return ins.opcode != bytecode.Jump && ins.opcode != bytecode.Return
}
