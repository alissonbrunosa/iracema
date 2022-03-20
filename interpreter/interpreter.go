package interpreter

import (
	"fmt"
	"iracema/bytecode"
	"iracema/lang"
	"os"
)

const (
	ADD byte = iota
	SUB
	MUL
	DIV
	EQ
	GT
	GE
	NE
	LT
	LE
)

var names = map[byte]string{
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	EQ:  "==",
	GT:  ">",
	GE:  ">=",
	NE:  "!=",
	LT:  "<",
	LE:  "<=",
}

type Interpreter struct {
	*frame

	frameCount int
	err        *lang.ErrorObject
}

func (i *Interpreter) Init(top *lang.Method) {
	i.PushFrame(lang.NewObject(), top, TOP_FRAME)
}

func (i *Interpreter) Dispatch() (lang.IrObject, error) {
	for {
	start_frame:
		if i.instrPointer >= len(i.instrs) {
			break
		}

	resume_frame:
		instrs := i.instrs
		constants := i.constants
		self := i.self

	next_instr:
		instr := instrs[i.instrPointer]
		op := bytecode.Opcode(instr >> 8)
		operand := byte(instr & 255)
		i.instrPointer++

		switch op {
		case bytecode.Push:
			value := constants[operand]
			i.Push(value)
			goto next_instr

		case bytecode.PushNone:
			i.Push(lang.None)
			goto next_instr

		case bytecode.PushSelf:
			i.Push(self)
			goto next_instr

		case bytecode.Pop:
			i.Pop()
			goto next_instr

		case bytecode.Return:
			ret := i.Pop()
			if i.PopFrame() {
				return ret, nil
			}

			i.Push(ret)
			goto resume_frame

		case bytecode.MatchType:
			err := i.Top(0)
			name := constants[operand]
			class := lang.TypeLookup(name)
			if class == nil {
				i.err = lang.NewNameError(name)
				goto fail
			}

			if err.Is(class) {
				i.Push(lang.True)
			} else {
				i.Push(lang.False)
			}
			goto next_instr

		case bytecode.Throw:
			i.err = i.Pop().(*lang.ErrorObject)
			i.PopFrame()
			goto next_instr

		case bytecode.Jump:
			i.JumpTo(operand)
			goto next_instr

		case bytecode.JumpIfFalse:
			if !lang.IsTruthy(i.Pop()) {
				i.JumpTo(operand)
			}
			goto next_instr

		case bytecode.Binary:
			rhs := i.Pop()
			lhs := i.Pop()
			name := names[operand]

			method := lhs.LookupMethod(name)
			val := method.Body().(lang.Native).Invoke(i, lhs, rhs)
			if val == nil {
				goto fail
			}

			i.Push(val)
			goto next_instr

		case bytecode.BuildArray:
			elements := i.PopN(operand)
			ary := lang.NewArray(elements)
			i.Push(ary)
			goto next_instr

		case bytecode.SetLocal:
			val := i.Pop()
			i.SetLocal(operand, val)
			goto next_instr

		case bytecode.GetLocal:
			val := i.GetLocal(operand)
			i.Push(val)
			goto next_instr

		case bytecode.SetAttr:
			attr := constants[operand]
			value := i.Pop()

			if err := lang.SetAttr(self, attr, value); err != nil {
				i.SetError(err)
				goto fail
			}

			goto next_instr

		case bytecode.GetAttr:
			attr := constants[operand]

			value, err := lang.GetAttr(self, attr)
			if err != nil {
				i.SetError(err)
				goto fail
			}

			i.Push(value)

			goto next_instr

		case bytecode.GetConstant:
			name := constants[operand]
			class := lang.TypeLookup(name)
			if class == nil {
				i.err = lang.NewNameError(name)
				goto fail
			}

			i.Push(class)
			goto next_instr

		case bytecode.DefineObject:
			body := constants[operand].(*lang.Method)
			class := lang.NewClass(body.Name(), lang.ObjectClass)
			lang.DefineType(body.Name(), class)
			i.PushObjectFrame(class, body)
			goto start_frame

		case bytecode.DefineFunction:
			class := self.(*lang.Class)
			meth := constants[operand].(*lang.Method)
			class.AddMethod(meth.Name(), meth)
			goto next_instr

		case bytecode.CallMethod:
			ci := constants[operand].(*lang.CallInfo)
			recv := i.Top(ci.Argc())

			meth := recv.LookupMethod(ci.Name())

			if meth == nil {
				err := lang.NewNoMethodError(recv, ci.Name())
				i.SetError(err)
				goto fail
			}

			switch m := meth.Body().(type) {
			case lang.Native:
				args := i.PopN(ci.Argc() + 1) // +1 recv
				i.PushGoFrame(recv, meth)
				val := m.Invoke(i, recv, args[1:]...)
				if val == nil {
					goto fail
				}

				i.PopFrame()
				i.Push(val)
				goto next_instr

			case []uint16:
				if ci.Argc() != meth.Arity() {
					i.err = lang.NewArityError(int(ci.Argc()), int(meth.Arity()))
					goto fail
				}

				i.PushFrame(recv, meth, IRMETHOD_FRAME)
				goto start_frame
			}

		default:
			panic(op.String())

		}

	fail:
		for i.frame != nil {
			if i.frame.catchOffset > 0 {
				i.Push(i.err)
				i.frame.instrPointer = i.frame.catchOffset
				goto resume_frame
			}
			i.PopFrame()
		}

		fmt.Println(i.err)
		os.Exit(1)
	}

	return nil, nil
}

func (i *Interpreter) PushObjectFrame(self lang.IrObject, fun *lang.Method) {
	i.frame = i.NewObjectFrame(self, fun)
	i.frameCount++
}

func (i *Interpreter) PushGoFrame(self lang.IrObject, meth *lang.Method) {
	i.frame = i.NewGoFrame(self, meth)
	i.frameCount++
}

func (i *Interpreter) PushFrame(self lang.IrObject, fun *lang.Method, flags byte) {
	if i.frame == nil {
		i.frame = TopFrame(self, fun)
	} else {
		i.frame = i.NewFrame(self, fun, flags)
	}
	i.frameCount++
}

func (i *Interpreter) PopFrame() (finished bool) {
	if i.flags&(TOP_FRAME|SINGLE_FRAME) != 0 {
		finished = true
	}

	i.frame = i.frame.previous
	i.frameCount--
	return
}

func (i *Interpreter) SetError(err *lang.ErrorObject) {
	i.err = err
}

func (i *Interpreter) Call(recv lang.IrObject, meth *lang.Method, args ...lang.IrObject) lang.IrObject {
	switch m := meth.Body().(type) {
	case lang.Native:
		i.PushGoFrame(recv, meth)
		val := m.Invoke(i, recv, args...)
		if val == nil {
			panic(i.err)
		}

		i.PopFrame()
		return val

	case []uint16:
		if byte(len(args)) != meth.Arity() {
			i.err = lang.NewArityError(len(args), int(meth.Arity()))
			panic(i.err)
		}

		i.Push(recv)
		for _, arg := range args {
			i.Push(arg)
		}

		i.PushFrame(recv, meth, SINGLE_FRAME|IRMETHOD_FRAME)
		ret, err := i.Dispatch()
		if err != nil {
			panic("ERROR" + err.Error())
		}

		return ret
	}

	return nil
}
