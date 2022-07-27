package lang

import (
	"testing"
)

func Test_toInt(t *testing.T) {
	var value IrObject = Int(1)

	result, err := toInt(value)

	if err != nil {
		t.Fatal("no error is expected")
	}

	if result != Int(1) {
		t.Errorf("expected result to be %d, got %d", Int(1), result)
	}
}

func Test_toInt_WhenArgumentIsNotInt(t *testing.T) {
	value := NewString("1")
	expectedMesg := "no implicit conversion of String into Int"
	expectedError := TypeError

	_, err := toInt(value)

	if err == nil {
		t.Fatal("an error is expectect")
	}

	if err.Class() != expectedError {
		t.Errorf("expected error to be %s, got %s", expectedError, err.Class())
	}

	if err.message != expectedMesg {
		t.Errorf("expected message to be %s, got %s", expectedMesg, err.message)
	}
}

func Test_intAdd(t *testing.T) {
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
		result := intAdd(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_intSub(t *testing.T) {
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
		result := intSub(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
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
		result := intMultiply(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
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
			result := intEqual(runtime, test.Left, test.Right)
			assertEqual(t, result, test.Expected)
		})
	}
}

func Test_intUnarySub(t *testing.T) {
	result := intUnarySub(runtime, Int(20))
	assertEqual(t, result, Int(-20))
}

func Test_intUnaryAdd(t *testing.T) {
	value := Int(20)
	result := intUnaryAdd(runtime, value)
	assertEqual(t, result, value)
}

func Test_intGreat(t *testing.T) {
	tests := []struct {
		scenario string
		left     IrObject
		right    IrObject
		expected IrObject
	}{
		{
			scenario: "lhs(Int) > rhs(Int)",
			left:     Int(20),
			right:    Int(2),
			expected: True,
		},
		{
			scenario: "lhs(Int) > rhs(Float)",
			left:     Int(20),
			right:    Float(2.0),
			expected: True,
		},
		{
			scenario: "lhs(Int) > rhs(Int) are equal",
			left:     Int(2),
			right:    Int(2),
			expected: False,
		},
		{
			scenario: "lhs(Int) > rhs(Float) are equal",
			left:     Int(2),
			right:    Float(2.0),
			expected: False,
		},
		{
			scenario: "lhs(Int) > rhs(Int): when left is smaller",
			left:     Int(1),
			right:    Int(2),
			expected: False,
		},
		{
			scenario: "lhs(Int) > rhs(Float): when left is smaller",
			left:     Int(1),
			right:    Float(2.0),
			expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := intGreat(runtime, test.left, test.right)
			assertEqual(t, result, test.expected)
		})
	}
}

func Test_intGreatEqual(t *testing.T) {
	tests := []struct {
		scenario string
		left     IrObject
		right    IrObject
		expected IrObject
	}{
		{
			scenario: "lhs(Int) >= rhs(Int)",
			left:     Int(20),
			right:    Int(2),
			expected: True,
		},
		{
			scenario: "lhs(Int) >= rhs(Float)",
			left:     Int(20),
			right:    Float(2.0),
			expected: True,
		},
		{
			scenario: "lhs(Int) >= rhs(Int) are equal",
			left:     Int(2),
			right:    Int(2),
			expected: True,
		},
		{
			scenario: "lhs(Int) >= rhs(Float) are equal",
			left:     Int(2),
			right:    Float(2.0),
			expected: True,
		},
		{
			scenario: "lhs(Int) >= rhs(Int): when left is smaller",
			left:     Int(1),
			right:    Int(2),
			expected: False,
		},
		{
			scenario: "lhs(Int) >= rhs(Float): when left is smaller",
			left:     Int(1),
			right:    Float(2.0),
			expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := intGreatEqual(runtime, test.left, test.right)
			assertEqual(t, result, test.expected)
		})
	}
}

func Test_intLess(t *testing.T) {
	tests := []struct {
		scenario string
		left     IrObject
		right    IrObject
		expected IrObject
	}{
		{
			scenario: "lhs(Int) < rhs(Int)",
			left:     Int(2),
			right:    Int(20),
			expected: True,
		},
		{
			scenario: "lhs(Int) < rhs(Float)",
			left:     Int(2),
			right:    Float(20.0),
			expected: True,
		},
		{
			scenario: "lhs(Int) < rhs(Int) are equal",
			left:     Int(2),
			right:    Int(2),
			expected: False,
		},
		{
			scenario: "lhs(Int) < rhs(Float) are equal",
			left:     Int(2),
			right:    Float(2.0),
			expected: False,
		},
		{
			scenario: "lhs(Int) < rhs(Int): when left is bigger",
			left:     Int(10),
			right:    Int(2),
			expected: False,
		},
		{
			scenario: "lhs(Int) < rhs(Float): when left is bigger",
			left:     Int(10),
			right:    Float(2.0),
			expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := intLess(runtime, test.left, test.right)
			assertEqual(t, result, test.expected)
		})
	}
}

func Test_intLessEqual(t *testing.T) {
	tests := []struct {
		scenario string
		left     IrObject
		right    IrObject
		expected IrObject
	}{
		{
			scenario: "lhs(Int) <= rhs(Int)",
			left:     Int(2),
			right:    Int(20),
			expected: True,
		},
		{
			scenario: "lhs(Int) <= rhs(Float)",
			left:     Int(2),
			right:    Float(20.0),
			expected: True,
		},
		{
			scenario: "lhs(Int) <= rhs(Int) are equal",
			left:     Int(2),
			right:    Int(2),
			expected: True,
		},
		{
			scenario: "lhs(Int) <= rhs(Float) are equal",
			left:     Int(2),
			right:    Float(2.0),
			expected: True,
		},
		{
			scenario: "lhs(Int) <= rhs(Int): when left is bigger",
			left:     Int(10),
			right:    Int(2),
			expected: False,
		},
		{
			scenario: "lhs(Int) <= rhs(Float): when left is bigger",
			left:     Int(10),
			right:    Float(2.0),
			expected: False,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := intLessEqual(runtime, test.left, test.right)
			assertEqual(t, result, test.expected)
		})
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
		result := intDivide(runtime, test.Left, test.Right)
		assertEqual(t, result, test.Expected)
	}
}

func Test_intInspect(t *testing.T) {
	result := intInspect(runtime, Int(2))
	assertEqual(t, result, NewString("2"))
}

func TestIntOperationWithInvalidOperand(t *testing.T) {
	tests := []struct {
		Scenario      string
		Left          IrObject
		Right         IrObject
		operation     func(Runtime, IrObject, IrObject) IrObject
		ExpectedMesg  string
		ExpectedError *Class
	}{
		{
			Scenario:      "Int divided by 0",
			Left:          Int(8),
			Right:         Int(0),
			operation:     intDivide,
			ExpectedMesg:  "divided by 0",
			ExpectedError: ZeroDivisionError,
		},
		{
			Scenario:      "div by a non numeric",
			Left:          Int(5),
			Right:         NewString("1"),
			operation:     intDivide,
			ExpectedMesg:  "unsupported operand type(s): 'Int' / 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "add with a non numeric",
			Left:          Int(5),
			Right:         NewString("1"),
			operation:     intAdd,
			ExpectedMesg:  "unsupported operand type(s): 'Int' + 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "sub with a non numeric",
			Left:          Int(5),
			Right:         NewString("1"),
			operation:     intSub,
			ExpectedMesg:  "unsupported operand type(s): 'Int' - 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "mult with a non numeric",
			Left:          Int(5),
			Right:         NewString("1"),
			operation:     intMultiply,
			ExpectedMesg:  "unsupported operand type(s): 'Int' * 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(>) with a non numeric",
			Left:          Int(5),
			Right:         False,
			operation:     intGreat,
			ExpectedMesg:  "invalid comparison (>) between 'Int' and 'Bool'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(>=) with a non numeric",
			Left:          Int(5),
			Right:         True,
			operation:     intGreatEqual,
			ExpectedMesg:  "invalid comparison (>=) between 'Int' and 'Bool'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(<) with a non numeric",
			Left:          Int(5),
			Right:         NewString("5"),
			operation:     intLess,
			ExpectedMesg:  "invalid comparison (<) between 'Int' and 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(<=) with a non numeric",
			Left:          Int(5),
			Right:         NewString("5"),
			operation:     intLessEqual,
			ExpectedMesg:  "invalid comparison (<=) between 'Int' and 'String'",
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
