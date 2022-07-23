package lang

import "testing"

func Test_objectUnaryNot(t *testing.T) {
	table := []struct {
		input      IrObject
		wantOutput Bool
	}{
		{input: Int(20), wantOutput: False},
		{input: Float(20.20), wantOutput: False},
		{input: NewString("s"), wantOutput: False},
		{input: True, wantOutput: False},
		{input: False, wantOutput: True},
		{input: None, wantOutput: True},
	}

	for _, test := range table {
		result := objectUnaryNot(runtime, test.input)
		assertEqual(t, result, test.wantOutput)
	}
}
