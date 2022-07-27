package lang

import "testing"

func Test_stringEqual(t *testing.T) {
	table := []struct {
		scenario   string
		lhs        IrObject
		rhs        IrObject
		wantOutput Bool
	}{
		{scenario: "equal", lhs: NewString("a"), rhs: NewString("a"), wantOutput: True},
		{scenario: "not equal", lhs: NewString("a"), rhs: NewString("aa"), wantOutput: False},
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			result := stringEqual(globalTestDummyRuntime, test.lhs, test.rhs)
			assertEqual(t, result, test.wantOutput)
		})
	}
}

func Test_stringLength(t *testing.T) {
	str := NewString("string")
	length := stringSize(globalTestDummyRuntime, str)
	assertEqual(t, length, Int(6))
}

func Test_stringInspect(t *testing.T) {
	str := NewString("string")
	inspect := stringInspect(globalTestDummyRuntime, str)
	assertEqual(t, inspect, NewString("string"))
}

func Test_stringPlus(t *testing.T) {
	a := NewString("a")
	b := NewString("b")

	result := stringPlus(globalTestDummyRuntime, a, b)
	length := stringSize(globalTestDummyRuntime, result)
	assertEqual(t, result, NewString("ab"))
	assertEqual(t, length, Int(2))
}
