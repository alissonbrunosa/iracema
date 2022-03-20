package lang

type Native interface {
	Invoke(Runtime, IrObject, ...IrObject) IrObject
}

func checkArity(rt Runtime, given, expected int) bool {
	if given != expected {
		rt.SetError(NewArityError(given, expected))
		return false
	}

	return true
}

type nArgs func(Runtime, IrObject, ...IrObject) IrObject
type zeroArgs func(Runtime, IrObject) IrObject
type oneArg func(Runtime, IrObject, IrObject) IrObject
type twoArgs func(Runtime, IrObject, IrObject, IrObject) IrObject

func (fn nArgs) Invoke(rt Runtime, recv IrObject, argv ...IrObject) IrObject {
	return fn(rt, recv, argv...)
}

func (fn zeroArgs) Invoke(rt Runtime, recv IrObject, argv ...IrObject) IrObject {
	if ok := checkArity(rt, len(argv), 0); !ok {
		return nil
	}

	return fn(rt, recv)
}

func (fn oneArg) Invoke(rt Runtime, recv IrObject, argv ...IrObject) IrObject {
	if ok := checkArity(rt, len(argv), 1); !ok {
		return nil
	}

	return fn(rt, recv, argv[0])
}

func (fn twoArgs) Invoke(rt Runtime, recv IrObject, argv ...IrObject) IrObject {
	if ok := checkArity(rt, len(argv), 2); !ok {
		return nil
	}

	return fn(rt, recv, argv[0], argv[1])
}
