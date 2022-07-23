package lang

import "testing"

func Test_stringEqual(t *testing.T) {
	a := NewString("a")
	b := NewString("a")

	if stringEqual(runtime, a, b) != True {
		t.Error("expected to be true")
	}
}

func Test_stringLength(t *testing.T) {
	str := NewString("string")
	length := stringSize(runtime, str)
	assertEqual(t, length, Int(6))
}

func Test_stringInspect(t *testing.T) {
	str := NewString("string")
	inspect := stringInspect(runtime, str)
	assertEqual(t, inspect, NewString("string"))
}

func Test_stringPlus(t *testing.T) {
	a := NewString("a")
	b := NewString("b")

	result := stringPlus(runtime, a, b)
	length := stringSize(runtime, result)
	assertEqual(t, result, NewString("ab"))
	assertEqual(t, length, Int(2))
}
