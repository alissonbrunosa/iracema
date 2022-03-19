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
		result := floatPlus(test.Left, test.Right)
		eq(t, result, test.Expected)
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
		result := floatMinus(test.Left, test.Right)
		eq(t, result, test.Expected)
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
		result := floatMultiply(test.Left, test.Right)
		eq(t, result, test.Expected)
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
		result := floatEqual(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_floatNegate(t *testing.T) {
	tests := []struct {
		Obj      IrObject
		Expected IrObject
	}{
		{
			Obj:      Float(2.0),
			Expected: Float(-2.0),
		},
		{
			Obj:      Float(-2),
			Expected: Float(2.0),
		},
	}

	for _, test := range tests {
		result := floatNegate(test.Obj)
		eq(t, result, test.Expected)
	}
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
		result := floatGreatThan(test.Left, test.Right)
		eq(t, result, test.Expected)
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
		result := floatLessThan(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_floatInspect(t *testing.T) {
	result := floatInspect(Float(2.9010))
	eq(t, result, NewString("2.901000"))
}
