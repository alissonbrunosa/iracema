package lang

import (
	"testing"
)

type dummyRuntime struct {
	err *ErrorObject
}

func (rt *dummyRuntime) SetError(err *ErrorObject) { rt.err = err }
func (rt *dummyRuntime) Call(recv IrObject, method *Method, args ...IrObject) IrObject {
	return method.Native().Invoke(rt, recv, args...)
}

var globalTestDummyRuntime = new(dummyRuntime)

func assertEqual(t *testing.T, got IrObject, expected IrObject) {
	t.Helper()

	var match bool
	switch res := expected.(type) {
	case Int:
		match = res == got.(Int)
	case Float:
		match = res == got.(Float)
	case *String:
		match = string(res.Value) == string(got.(*String).Value)
	case Bool:
		match = res == got.(Bool)

	default:
		t.Fatalf("wrong type %T", res)
	}

	if !match {
		t.Errorf("expected value to be %v, got %v", expected, got)
	}
}
