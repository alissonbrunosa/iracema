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
		result := floatAdd(globalTestDummyRuntime, test.Left, test.Right)
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
		result := floatSub(globalTestDummyRuntime, test.Left, test.Right)
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
		result := floatMultiply(globalTestDummyRuntime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_floatUnaryMinus(t *testing.T) {
	result := floatUnarySub(globalTestDummyRuntime, Float(20.40))
	assertEqual(t, result, Float(-20.40))
}

func Test_floatUnaryPlus(t *testing.T) {
	value := Float(20.40)
	result := floatUnaryAdd(globalTestDummyRuntime, value)
	assertEqual(t, result, value)
}

func Test_floatGreat(t *testing.T) {
	tests := []struct {
		scenario string
		left     IrObject
		right    IrObject
		expected IrObject
	}{
		{
			scenario: "lhs(Float) > rhs(Float)",
			left:     Float(20.0),
			right:    Float(2.0),
			expected: True,
		},
		{
			scenario: "lhs(Float) > rhs(Int)",
			left:     Float(20.0),
			right:    Int(2),
			expected: True,
		},
		{
			scenario: "lhs(Float) > rhs(Float) are equal",
			left:     Float(2.0),
			right:    Float(2.0),
			expected: False,
		},
		{
			scenario: "lhs(Float) > rhs(Int) are equal",
			left:     Float(2.0),
			right:    Int(2),
			expected: False,
		},
		{
			scenario: "lhs(Float) > rhs(Float): when left is smaller",
			left:     Float(1.0),
			right:    Float(2.0),
			expected: False,
		},
		{
			scenario: "lhs(Float) > rhs(Int): when left is smaller",
			left:     Float(1.0),
			right:    Int(2),
			expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := floatGreat(globalTestDummyRuntime, test.left, test.right)
			assertEqual(t, result, test.expected)
		})
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
		result := floatEqual(globalTestDummyRuntime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_floatGreatEqual(t *testing.T) {
	tests := []struct {
		scenario string
		left     IrObject
		right    IrObject
		expected IrObject
	}{
		{
			scenario: "lhs(Float) >= rhs(Float)",
			left:     Float(20.0),
			right:    Float(2.0),
			expected: True,
		},
		{
			scenario: "lhs(Float) >= rhs(Int)",
			left:     Float(20.0),
			right:    Int(2),
			expected: True,
		},
		{
			scenario: "lhs(Float) >= rhs(Float) are equal",
			left:     Float(2.0),
			right:    Float(2.0),
			expected: True,
		},
		{
			scenario: "lhs(Float) >= rhs(Int) are equal",
			left:     Float(2.0),
			right:    Int(2),
			expected: True,
		},
		{
			scenario: "lhs(Float) >= rhs(Float): when left is smaller",
			left:     Float(1.0),
			right:    Float(2.0),
			expected: False,
		},
		{
			scenario: "lhs(Float) >= rhs(Int): when left is smaller",
			left:     Float(1.0),
			right:    Int(2),
			expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := floatGreatEqual(globalTestDummyRuntime, test.left, test.right)
			assertEqual(t, result, test.expected)
		})
	}
}

func Test_floatLess(t *testing.T) {
	tests := []struct {
		scenario string
		left     IrObject
		right    IrObject
		expected IrObject
	}{
		{
			scenario: "lhs(Float) < rhs(Float)",
			left:     Float(2.0),
			right:    Float(20.0),
			expected: True,
		},
		{
			scenario: "lhs(Float) < rhs(Int)",
			left:     Float(2.0),
			right:    Int(20),
			expected: True,
		},
		{
			scenario: "lhs(Float) < rhs(Float) are equal",
			left:     Float(2.0),
			right:    Float(2.0),
			expected: False,
		},
		{
			scenario: "lhs(Float) < rhs(Int) are equal",
			left:     Float(2),
			right:    Int(2),
			expected: False,
		},
		{
			scenario: "lhs(Float) < rhs(Float): when left is bigger",
			left:     Float(2.0),
			right:    Float(1.0),
			expected: False,
		},
		{
			scenario: "lhs(Float) < rhs(Int): when left is bigger",
			left:     Float(2.0),
			right:    Int(1),
			expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := floatLess(globalTestDummyRuntime, test.left, test.right)
			assertEqual(t, result, test.expected)
		})
	}
}

func Test_floatLessEqual(t *testing.T) {
	tests := []struct {
		scenario string
		left     IrObject
		right    IrObject
		expected IrObject
	}{
		{
			scenario: "lhs(Float) <= rhs(Float)",
			left:     Float(2.0),
			right:    Float(20.0),
			expected: True,
		},
		{
			scenario: "lhs(Float) <= rhs(Int)",
			left:     Float(2.0),
			right:    Int(20),
			expected: True,
		},
		{
			scenario: "lhs(Float) <= rhs(Float) are equal",
			left:     Float(2.0),
			right:    Float(2.0),
			expected: True,
		},
		{
			scenario: "lhs(Float) <= rhs(Int) are equal",
			left:     Float(2),
			right:    Int(2),
			expected: True,
		},
		{
			scenario: "lhs(Float) <= rhs(Float): when left is bigger",
			left:     Float(2.0),
			right:    Float(1.0),
			expected: False,
		},
		{
			scenario: "lhs(Float) <= rhs(Int): when left is bigger",
			left:     Float(2.0),
			right:    Int(1),
			expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := floatLessEqual(globalTestDummyRuntime, test.left, test.right)
			assertEqual(t, result, test.expected)
		})
	}
}

func Test_floatInspect(t *testing.T) {
	result := floatInspect(globalTestDummyRuntime, Float(2.9010))
	assertEqual(t, result, NewString("2.901000"))
}

func TestFloatOperationWithInvalidOperand(t *testing.T) {
	tests := []struct {
		Scenario      string
		Left          IrObject
		Right         IrObject
		operation     func(Runtime, IrObject, IrObject) IrObject
		ExpectedMesg  string
		ExpectedError *Class
	}{
		{
			Scenario:      "div by a non numeric",
			Left:          Float(10.5),
			Right:         NewString("1"),
			operation:     floatDivide,
			ExpectedMesg:  "unsupported operand type(s): 'Float' / 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "add with a non numeric",
			Left:          Float(10.5),
			Right:         NewString("1"),
			operation:     floatAdd,
			ExpectedMesg:  "unsupported operand type(s): 'Float' + 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "sub with a non numeric",
			Left:          Float(10.5),
			Right:         NewString("1"),
			operation:     floatSub,
			ExpectedMesg:  "unsupported operand type(s): 'Float' - 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "mult with a non numeric",
			Left:          Float(10.5),
			Right:         NewString("1"),
			operation:     floatMultiply,
			ExpectedMesg:  "unsupported operand type(s): 'Float' * 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(>) with a non numeric",
			Left:          Float(10.5),
			Right:         False,
			operation:     floatGreat,
			ExpectedMesg:  "invalid comparison (>) between 'Float' and 'Bool'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(>=) with a non numeric",
			Left:          Float(10.5),
			Right:         True,
			operation:     floatGreatEqual,
			ExpectedMesg:  "invalid comparison (>=) between 'Float' and 'Bool'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(<) with a non numeric",
			Left:          Float(10.5),
			Right:         NewString("5"),
			operation:     floatLess,
			ExpectedMesg:  "invalid comparison (<) between 'Float' and 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(<=) with a non numeric",
			Left:          Float(10.5),
			Right:         NewString("5"),
			operation:     floatLessEqual,
			ExpectedMesg:  "invalid comparison (<=) between 'Float' and 'String'",
			ExpectedError: TypeError,
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			var rt = new(dummyRuntime)

			if value := test.operation(rt, test.Left, test.Right); value != nil {
				t.Error("expected value to be nil when an error occurs")
			}

			if rt.err == nil {
				t.Error("expected an error to be set in Runtime")
			}

			if rt.err.message != test.ExpectedMesg {
				t.Errorf("expected error message to be %s, got %s", test.ExpectedMesg, rt.err.message)
			}

			if rt.err.Class() != test.ExpectedError {
				t.Errorf("expected error class to be %s, got %s", test.ExpectedError, rt.err.Class())
			}
		})
	}
}
