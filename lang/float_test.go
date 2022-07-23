package lang

import "testing"

func Test_floatPlus(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Float(1.0),
			Right:    Int(2),
			Expected: Float(3.0),
		},
		{
			Left:     Float(1.0),
			Right:    Float(2.5),
			Expected: Float(3.5),
		},
	}

	for _, test := range tests {
		result := floatAdd(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_floatMinus(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Float(1.0),
			Right:    Int(2),
			Expected: Float(-1.0),
		},
		{
			Left:     Float(1.0),
			Right:    Float(2.5),
			Expected: Float(-1.5),
		},
	}

	for _, test := range tests {
		result := floatSub(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_floatMultiply(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Float(2.0),
			Right:    Int(2),
			Expected: Float(4.0),
		},
		{
			Left:     Float(2.0),
			Right:    Float(2.5),
			Expected: Float(5.0),
		},
	}

	for _, test := range tests {
		result := floatMultiply(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_floatEqual(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Float(2.0),
			Right:    Int(2),
			Expected: True,
		},
		{
			Left:     Float(2.0),
			Right:    Float(3.0),
			Expected: False,
		},
		{
			Left:     Float(2.0),
			Right:    Float(2.0),
			Expected: True,
		},
	}

	for _, test := range tests {
		result := floatEqual(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_floatUnaryMinus(t *testing.T) {
	result := floatUnarySub(runtime, Float(20.40))
	assertEqual(t, result, Float(-20.40))
}

func Test_floatUnaryPlus(t *testing.T) {
	value := Float(20.40)
	result := floatUnaryAdd(runtime, value)
	assertEqual(t, result, value)
}

func Test_floatGreatThan(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Float(2.0),
			Right:    Int(2.0),
			Expected: False,
		},
		{
			Left:     Float(3.0),
			Right:    Float(2.0),
			Expected: True,
		},
		{
			Left:     Float(2.0),
			Right:    Float(1.0),
			Expected: True,
		},
	}

	for _, test := range tests {
		result := floatGreatThan(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_floatLessThan(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Float(2.0),
			Right:    Int(3.0),
			Expected: True,
		},
		{
			Left:     Float(3.0),
			Right:    Float(2.0),
			Expected: False,
		},
		{
			Left:     Float(2.0),
			Right:    Float(1.0),
			Expected: False,
		},
	}

	for _, test := range tests {
		result := floatLessThan(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_floatInspect(t *testing.T) {
	result := floatInspect(runtime, Float(2.9010))
	assertEqual(t, result, NewString("2.901000"))
}
