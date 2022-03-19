package lang

type Native interface {
	Invoke(IrObject, ...IrObject) IrObject
}

func checkArity(given, expected int) bool {
	if given != expected {
		return false
	}

	return true
}

type nArgs func(IrObject, ...IrObject) IrObject
type zeroArgs func(IrObject) IrObject
type oneArg func(IrObject, IrObject) IrObject
type twoArgs func(IrObject, IrObject, IrObject) IrObject

func (fn nArgs) Invoke(recv IrObject, argv ...IrObject) IrObject {
	return fn(recv, argv...)
}

func (fn zeroArgs) Invoke(recv IrObject, argv ...IrObject) IrObject {
	if ok := checkArity(len(argv), 0); !ok {
		return nil
	}

	return fn(recv)
}

func (fn oneArg) Invoke(recv IrObject, argv ...IrObject) IrObject {
	if ok := checkArity(len(argv), 1); !ok {
		return nil
	}

	return fn(recv, argv[0])
}

func (fn twoArgs) Invoke(recv IrObject, argv ...IrObject) IrObject {
	if ok := checkArity(len(argv), 2); !ok {
		return nil
	}

	return fn(recv, argv[0], argv[1])
}
