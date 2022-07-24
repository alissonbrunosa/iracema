package lang

func call(rt Runtime, recv IrObject, name string, args ...IrObject) IrObject {
	class := recv.Class()

	method := class.LookupMethod(name)
	if method == nil {
		rt.SetError(NewNoMethodError(recv, name))
		return nil
	}

	return rt.Call(recv, method, args...)
}

func safeCall(rt Runtime, recv IrObject, name string, args ...IrObject) IrObject {
	class := recv.Class()

	method := class.LookupMethod(name)
	if method == nil {
		return nil
	}

	return rt.Call(recv, method, args...)
}
