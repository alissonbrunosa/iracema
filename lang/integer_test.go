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
		result := intPlus(runtime, test.Left, test.Right)
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
		result := intMinus(runtime, test.Left, test.Right)
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
		result := intMultiply(runtime, test.Left, test.Right)
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
			result := intEqual(runtime, test.Left, test.Right)
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
		result := intNegate(runtime, test.Obj)
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
		result := intGreatThan(runtime, test.Left, test.Right)
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
		result := intGreaterThanOrEqual(runtime, test.Left, test.Right)
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
		result := intLessThan(runtime, test.Left, test.Right)
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
		result := intLessThanOrEqual(runtime, test.Left, test.Right)
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
		result := intDivide(runtime, test.Left, test.Right)
		eq(t, result, test.Expected)
	}
}

func Test_intInspect(t *testing.T) {
	result := intInspect(runtime, Int(2))
	eq(t, result, NewString("2"))
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
			operation:     intPlus,
			ExpectedMesg:  "unsupported operand type(s): 'Int' + 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "sub with a non numeric",
			Left:          Int(5),
			Right:         NewString("1"),
			operation:     intMinus,
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
			operation:     intGreatThan,
			ExpectedMesg:  "invalid comparison between 'Int' and 'Bool'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(>=) with a non numeric",
			Left:          Int(5),
			Right:         True,
			operation:     intGreaterThanOrEqual,
			ExpectedMesg:  "invalid comparison between 'Int' and 'Bool'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(<) with a non numeric",
			Left:          Int(5),
			Right:         NewString("5"),
			operation:     intLessThan,
			ExpectedMesg:  "invalid comparison between 'Int' and 'String'",
			ExpectedError: TypeError,
		},
		{
			Scenario:      "compare(<=) with a non numeric",
			Left:          Int(5),
			Right:         NewString("5"),
			operation:     intLessThanOrEqual,
			ExpectedMesg:  "invalid comparison between 'Int' and 'String'",
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
