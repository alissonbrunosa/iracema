package interpreter

import (
	"iracema/lang"
)

const (
	TOP_FRAME      = 0x01
	OBJECT_FRAME   = 0x02
	IRMETHOD_FRAME = 0x04
	SINGLE_FRAME   = 0x08
	GOMETHOD_FRAME = 0x10
)

type frame struct {
	flags        byte
	name         string
	method       *lang.Method
	self         lang.IrObject
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

func TopFrame(self lang.IrObject, meth *lang.Method) *frame {
	return &frame{
		flags:        TOP_FRAME,
		name:         "top",
		method:       meth,
		self:         self,
		class:        self.Class(),
		stack:        make([]lang.IrObject, STACK_SIZE),
		instrs:       meth.Instrs(),
		constants:    meth.Constants(),
		stackPointer: meth.LocalCount(),
		catchOffset:  -1,
	}
}

func (f *frame) NewGoFrame(self lang.IrObject, meth *lang.Method) *frame {
	return &frame{
		flags:        GOMETHOD_FRAME,
		name:         meth.Name(),
		method:       meth,
		self:         self,
		stack:        f.stack[f.stackPointer:],
		instrs:       nil,
		constants:    nil,
		stackPointer: 0,
		catchOffset:  -1,
		previous:     f,
	}
}

func (f *frame) NewObjectFrame(self lang.IrObject, meth *lang.Method) *frame {
	return &frame{
		flags:        OBJECT_FRAME,
		name:         meth.Name(),
		method:       meth,
		self:         self,
		class:        self.(*lang.Class),
		stack:        f.stack[f.stackPointer:],
		instrs:       meth.Instrs(),
		constants:    meth.Constants(),
		stackPointer: meth.LocalCount(),
		catchOffset:  -1,
		previous:     f,
	}
}

func (f *frame) NewFrame(self lang.IrObject, meth *lang.Method, flags byte) *frame {
	f.stackPointer -= meth.Arity()

	frame := &frame{
		flags:        flags,
		method:       meth,
		self:         self,
		class:        self.Class(),
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
