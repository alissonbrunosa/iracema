package interpreter

import (
	"fmt"
	"iracema/bytecode"
	"iracema/lang"
	"os"
)

const (
	CALL_OK = 1 << iota
	CALL_NEW_FRAME
	CALL_ERROR
)

type Interpreter struct {
	*frame

	frameCount int
	err        *lang.ErrorObject
}

func (i *Interpreter) Exec(top *lang.Method) (lang.IrObject, error) {
	i.PushFrame(lang.NewScript(), top, TOP_FRAME)
	return i.dispatch()
}

func (i *Interpreter) dispatch() (lang.IrObject, error) {
	for {
	start_frame:
		if i.instrPointer >= len(i.instrs) {
			break
		}

	resume_frame:
		instrs := i.instrs
		constants := i.constants
		this := i.this

	next_instr:
		instr := instrs[i.instrPointer]
		opcode := bytecode.Opcode(instr >> 8)
		operand := byte(instr & 255)
		i.instrPointer++

		switch opcode {
		case bytecode.Nop:
			goto next_instr

		case bytecode.Push:
			value := constants[operand]
			i.Push(value)
			goto next_instr

		case bytecode.PushNone:
			i.Push(lang.None)
			goto next_instr

		case bytecode.PushThis:
			i.Push(this)
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

		case bytecode.JumpIfTrue:
			if lang.IsTruthy(i.Pop()) {
				i.JumpTo(operand)
			}

			goto next_instr

		case bytecode.BuildArray:
			elements := i.PopN(operand)
			ary := lang.NewArray(elements)
			i.Push(ary)
			goto next_instr

		case bytecode.BuildHash:
			hash := lang.NewHash()
			hash.BulkInsert(i, i.PopN(operand))
			i.Push(hash)
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

			if err := lang.SetAttr(this, attr, value); err != nil {
				i.SetError(err)
				goto fail
			}

			goto next_instr

		case bytecode.GetAttr:
			attr := constants[operand]

			value, err := lang.GetAttr(this, attr)
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

			parent := lang.ObjectClass
			if p := i.Pop(); p != lang.None {
				parent = p.(*lang.Class)
			}

			class := lang.NewClass(body.Name(), parent)
			lang.DefineType(body.Name(), class)
			i.PushObjectFrame(class, body)
			goto start_frame

		case bytecode.DefineFunction:
			class := i.class
			meth := constants[operand].(*lang.Method)
			class.AddMethod(meth.Name(), meth)
			goto next_instr

		case bytecode.NewIterator:
			if it := i.Pop(); it.Is(lang.ArrayClass) {
				iter := lang.NewIterator(it)
				i.Push(iter)
				goto next_instr
			}

			err := lang.NewTypeError("object is not iterable")
			i.SetError(err)
			goto fail

		case bytecode.Iterate:
			iter := i.Top(0).(*lang.Iterator)

			if iter.HasNext() {
				i.Push(iter.Next())
				i.Push(lang.True)
			} else {
				i.Pop()
				i.Push(lang.False)
			}

			goto next_instr

		case bytecode.CallMethod:
			info := constants[operand].(*lang.CallInfo)
			recv := i.Top(info.Argc())
			class := recv.Class()
			method := class.LookupMethod(info.Name())

			if method == nil {
				i.err = lang.NewNoMethodError(recv, info.Name())
				goto fail
			}

			switch i.call0(recv, method, info) {
			case CALL_OK:
				goto next_instr
			case CALL_NEW_FRAME:
				goto start_frame
			default:
				goto fail
			}

		case bytecode.CallSuper:
			if (i.frame.flags & IRMETHOD_FRAME) == 0 {
				i.err = lang.NewError("Really?! Calling super outside of a method", lang.Error)
				goto fail
			}

			info := constants[operand].(*lang.CallInfo)
			recv := i.Top(info.Argc())
			super := recv.Class().Super()
			method := super.LookupMethod(info.Name())

			if method == nil {
				i.err = lang.NewError(
					"no superclass of '%s' has method '%s'",
					lang.NoMethodError,
					recv.Class(),
					info.Name(),
				)
				goto fail
			}

			switch i.call0(recv, method, info) {
			case CALL_OK:
				goto next_instr
			case CALL_NEW_FRAME:
				goto start_frame
			default:
				goto fail
			}

		default:
			fmt.Printf("Instruction %s is not implemented yet\n", opcode)
			os.Exit(1)
		}

	fail:
		if i.catchError() {
			goto resume_frame
		}
	}

	return nil, nil
}

func (i *Interpreter) catchError() bool {
	for i.frame != nil {
		if i.frame.catchOffset > 0 {
			i.Push(i.err)
			i.frame.instrPointer = i.frame.catchOffset
			return true
		}
		i.PopFrame()
	}

	fmt.Println(i.err)
	os.Exit(1)
	return false
}

func (i *Interpreter) PushObjectFrame(this lang.IrObject, fun *lang.Method) {
	i.frame = i.NewObjectFrame(this, fun)
	i.frameCount++
}

func (i *Interpreter) PushGoFrame(this lang.IrObject, meth *lang.Method) {
	i.frame = i.NewGoFrame(this, meth)
	i.frameCount++
}

func (i *Interpreter) PushFrame(this lang.IrObject, fun *lang.Method, flags byte) {
	if i.frame == nil {
		i.frame = TopFrame(this, fun)
	} else {
		i.frame = i.NewFrame(this, fun, flags)
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

func (i *Interpreter) callGoFunc(recv lang.IrObject, method lang.Native, argc byte) int {
	args := i.PopN(argc + 1) // +1 recv

	if val := method.Invoke(i, recv, args[1:]...); val != nil {
		i.Push(val)
		return CALL_OK
	}

	return CALL_ERROR
}

func (i *Interpreter) call0(recv lang.IrObject, method *lang.Method, info *lang.CallInfo) int {
	switch method.MethodType() {
	case lang.GoFunction:
		return i.callGoFunc(recv, method.Native(), info.Argc())

	case lang.IrMethod:
		if info.Argc() != method.Arity() {
			i.err = lang.NewArityError(int(info.Argc()), int(method.Arity()))
			return CALL_ERROR
		}

		i.PushFrame(recv, method, IRMETHOD_FRAME)
		return CALL_NEW_FRAME
	default:
		lang.Unreachable()
		return CALL_ERROR
	}
}

func (i *Interpreter) Call(recv lang.IrObject, method *lang.Method, args ...lang.IrObject) lang.IrObject {
	if byte(len(args)) != method.Arity() {
		i.err = lang.NewArityError(len(args), int(method.Arity()))
		return nil
	}

	i.Push(recv)
	for _, arg := range args {
		i.Push(arg)
	}

	i.PushFrame(recv, method, SINGLE_FRAME|IRMETHOD_FRAME)
	ret, err := i.dispatch()
	if err != nil {
		i.err = lang.NewError("unkown error:", lang.Error)
		return nil
	}

	return ret
}
