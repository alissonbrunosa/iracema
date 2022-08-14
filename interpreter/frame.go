package interpreter

import (
	"iracema/lang"
)

const (
	TOP_FRAME      = 0x01
	OBJECT_FRAME   = 0x02
	IRMETHOD_FRAME = 0x04
	FLAG_DONE      = 0x08
)

type frame struct {
	flags        byte
	name         string
	method       *lang.Method
	this         lang.IrObject
	class        *lang.Class
	instrs       []uint16
	constants    []lang.IrObject
	instrPointer int
	stack        []lang.IrObject
	stackPointer byte
	previous     *frame
	catchOffset  int
}

const STACK_SIZE = 1024

func TopFrame(this lang.IrObject, fun *lang.Method) *frame {
	stack := make([]lang.IrObject, STACK_SIZE)

	for i := byte(0); i < fun.LocalCount(); i++ {
		stack[i] = lang.None
	}

	return &frame{
		flags:        TOP_FRAME | FLAG_DONE,
		method:       fun,
		this:         this,
		class:        this.Class(),
		stack:        stack,
		name:         fun.Name(),
		instrs:       fun.Instrs(),
		constants:    fun.Constants(),
		stackPointer: fun.LocalCount(),
		catchOffset:  -1,
	}
}

func (f *frame) NewObjectFrame(this lang.IrObject, meth *lang.Method) *frame {
	return &frame{
		flags:        OBJECT_FRAME,
		name:         meth.Name(),
		method:       meth,
		this:         this,
		class:        this.(*lang.Class),
		stack:        f.stack[f.stackPointer:],
		instrs:       meth.Instrs(),
		constants:    meth.Constants(),
		stackPointer: meth.LocalCount(),
		catchOffset:  -1,
		previous:     f,
	}
}

func (f *frame) NewFrame(this lang.IrObject, argc byte, meth *lang.Method, flags byte) *frame {
	locals := meth.LocalCount() - meth.Arity()
	for i := f.stackPointer; i <= f.stackPointer+locals; i++ {
		f.stack[i] = lang.None
	}

	f.stackPointer -= argc
	frame := &frame{
		flags:        flags,
		method:       meth,
		this:         this,
		class:        this.Class(),
		stack:        f.stack[f.stackPointer:],
		name:         meth.Name(),
		instrs:       meth.Instrs(),
		constants:    meth.Constants(),
		stackPointer: meth.LocalCount(),
		catchOffset:  meth.CatchOffset(),
		previous:     f,
	}

	f.Pop() //recv

	return frame
}

func (f *frame) SetLocal(index byte, value lang.IrObject) {
	f.stack[index] = value
}

func (f *frame) GetLocal(index byte) lang.IrObject {
	return f.stack[index]
}

func (f *frame) Push(object lang.IrObject) {
	f.stack[f.stackPointer] = object
	f.stackPointer++
}

func (f *frame) Pop() lang.IrObject {
	f.stackPointer--
	var popped lang.IrObject
	popped, f.stack[f.stackPointer] = f.stack[f.stackPointer], nil
	return popped
}

func (f *frame) Top(nth byte) lang.IrObject {
	return f.stack[f.stackPointer-nth-1]
}

func (f *frame) PeekAt(index byte) lang.IrObject {
	return f.stack[index]
}

func (f *frame) PopN(n byte) []lang.IrObject {
	values := f.stack[f.stackPointer-n : f.stackPointer]
	copied := make([]lang.IrObject, n)
	copy(copied, values)

	for i := range values {
		values[i] = nil
	}

	f.stackPointer -= n
	return copied
}

func (f *frame) JumpTo(offset byte) {
	f.instrPointer = int(offset)
}

func (f *frame) Clean() {
	for i := byte(0); i < f.stackPointer; i++ {
		f.stack[i] = nil
	}
}
