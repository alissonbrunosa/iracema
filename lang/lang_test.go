package lang

type dummyRuntime struct {
	err *ErrorObject
}

func (rt *dummyRuntime) SetError(err *ErrorObject) { rt.err = err }
func (rt *dummyRuntime) Call(recv IrObject, method *Method, args ...IrObject) IrObject {
	if meth, ok := method.Body().(Native); ok {
		return meth.Invoke(rt, recv, args...)
	}

	panic("can't call user defined method with dummy runtime")
}

var runtime = new(dummyRuntime)
