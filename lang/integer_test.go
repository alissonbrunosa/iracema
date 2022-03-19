package lang

import (
	"testing"
)

func eq(t *testing.T, got IrObject, expected IrObject) {
	t.Helper()

	var result bool
	var value IrObject
	switch res := expected.(type) {
	case Int:
		result = res == got.(Int)
	case Float:
		result = res == got.(Float)
	case *String:
		result = string(res.Value) == string(got.(*String).Value)
	case Bool:
		result = res == got.(Bool)

	default:
		t.Fatalf("wrong type %T", res)
	}

	if !result {
		t.Errorf("expected value to be %v, got %v", expected, value)
	}
}

func Test_intPlus(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Int(1),
			Right:    Int(2),
			Expected: Int(3),
		},
		{
			Left:     Int(1),
			Right:    Float(2.5),
			Expected: Float(3.5),
		},
	}

	for _, test := range tests {
		result := intPlus(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_intMinus(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Int(1),
			Right:    Int(2),
			Expected: Int(-1),
		},
		{
			Left:     Int(1),
			Right:    Float(2.5),
			Expected: Float(-1.5),
		},
	}

	for _, test := range tests {
		result := intMinus(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_intMultiply(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Int(2),
			Right:    Int(2),
			Expected: Int(4),
		},
		{
			Left:     Int(2),
			Right:    Float(2.5),
			Expected: Float(5.0),
		},
	}

	for _, test := range tests {
		result := intMultiply(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_intEqual(t *testing.T) {
	tests := []struct {
		Scenario string
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Scenario: "when both are equal",
			Left:     Int(2),
			Right:    Int(2),
			Expected: True,
		},
		{
			Scenario: "when both are different",
			Left:     Int(2),
			Right:    Int(3),
			Expected: False,
		},
		{
			Scenario: `when rhs is a float but has the "same" value`,
			Left:     Int(2),
			Right:    Float(2.0),
			Expected: True,
		},
		{
			Scenario: `when rhs is non numeric`,
			Left:     Int(2),
			Right:    True,
			Expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			result := intEqual(test.Left, test.Right)
			eq(t, result, test.Expected)
		})
	}
}

func Test_intNegate(t *testing.T) {
	tests := []struct {
		Obj      IrObject
		Expected IrObject
	}{
		{
			Obj:      Int(2),
			Expected: Int(-2),
		},
		{
			Obj:      Int(-2),
			Expected: Int(2),
		},
	}

	for _, test := range tests {
		result := intNegate(test.Obj)
		eq(t, result, test.Expected)
	}
}

func Test_intGreatThan(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Int(2),
			Right:    Int(2),
			Expected: False,
		},
		{
			Left:     Int(3),
			Right:    Int(2),
			Expected: True,
		},
		{
			Left:     Int(2),
			Right:    Float(1.0),
			Expected: True,
		},
	}

	for _, test := range tests {
		result := intGreatThan(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_intGreatThanOrEqual(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Int(2),
			Right:    Int(2),
			Expected: True,
		},
		{
			Left:     Int(3),
			Right:    Int(2),
			Expected: True,
		},
		{
			Left:     Int(2),
			Right:    Float(2.0),
			Expected: True,
		},
	}

	for _, test := range tests {
		result := intGreaterThanOrEqual(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_intLessThan(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Int(2),
			Right:    Int(3),
			Expected: True,
		},
		{
			Left:     Int(3),
			Right:    Int(2),
			Expected: False,
		},
		{
			Left:     Int(2),
			Right:    Float(1.0),
			Expected: False,
		},
	}

	for _, test := range tests {
		result := intLessThan(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_intLessThanOrEqual(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Int(3),
			Right:    Int(3),
			Expected: True,
		},
		{
			Left:     Int(5),
			Right:    Int(2),
			Expected: False,
		},
		{
			Left:     Int(2),
			Right:    Float(2.0),
			Expected: True,
		},
	}

	for _, test := range tests {
		result := intLessThanOrEqual(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}
func Test_intDivide(t *testing.T) {
	tests := []struct {
		Left     IrObject
		Right    IrObject
		Expected IrObject
	}{
		{
			Left:     Int(8),
			Right:    Int(2),
			Expected: Int(4),
		},
		{
			Left:     Int(5),
			Right:    Float(2.0),
			Expected: Float(2.5),
		},
	}

	for _, test := range tests {
		result := intDivide(test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_intInspect(t *testing.T) {
	result := intInspect(Int(2))
	eq(t, result, NewString("2"))
}
