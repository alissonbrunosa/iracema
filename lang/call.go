package lang

import (
	"fmt"
	"os"
)

func call(rt Runtime, recv IrObject, name string, args ...IrObject) IrObject {
	class := recv.Class()

	method := class.LookupMethod(name)
	if method == nil {
		rt.SetError(NewNoMethodError(recv, name))
		return nil
	}

	switch method.methodType {
	case GoFunction:
		return method.Native().Invoke(rt, recv, args...)
	case IrMethod:
		return rt.Call(recv, method, args...)
	default:
		Unreachable()
	}

	return nil
}

func Unreachable() {
	fmt.Println("Damn! That's a bug!")
	os.Exit(1)
}
