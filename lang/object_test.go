package lang

import (
	"testing"
)

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

func Test_objectEqual(t *testing.T) {
	obj := NewObject()
	result := objectEqual(runtime, obj, obj)
	assertEqual(t, result, True)

	result = objectEqual(runtime, obj, NewObject())
	assertEqual(t, result, False)
}

func Test_objectNotEqual(t *testing.T) {
	table := []struct {
		scenario   string
		lhs        IrObject
		rhs        IrObject
		wantOutput Bool
	}{
		{scenario: "Int/Equal", lhs: Int(20), rhs: Int(20), wantOutput: False},
		{scenario: "Int/NotEqual", lhs: Int(20), rhs: Int(22), wantOutput: True},
		{scenario: "Float/Equal", lhs: Float(20.20), rhs: Float(20.20), wantOutput: False},
		{scenario: "Float/NotEqual", lhs: Float(20.20), rhs: Float(22.20), wantOutput: True},
		{scenario: "String/Equal", lhs: NewString("s"), rhs: NewString("s"), wantOutput: False},
		{scenario: "String/NotEqual", lhs: NewString("s"), rhs: NewString("ss"), wantOutput: True},
		{scenario: "Bool(True)/Equal", lhs: True, rhs: True, wantOutput: False},
		{scenario: "Bool(True)/NotEqual", lhs: True, rhs: False, wantOutput: True},
		{scenario: "Bool(False)/Equal", lhs: False, rhs: False, wantOutput: False},
		{scenario: "Bool(False)/NotEqual", lhs: False, rhs: True, wantOutput: True},
		{scenario: "None/Equal", lhs: None, rhs: None, wantOutput: False},
		{scenario: "None/NotEqual", lhs: None, rhs: Int(29), wantOutput: True},
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			result := objectNotEqual(runtime, test.lhs, test.rhs)
			assertEqual(t, result, test.wantOutput)
		})
	}
}
