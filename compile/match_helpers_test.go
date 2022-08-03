package compile

import (
	"iracema/bytecode"
	"iracema/lang"
	"testing"
)

type Match interface {
	Match(*testing.T, uint16, []lang.IrObject)
}

type singleByteMatch struct {
	opcode bytecode.Opcode
}

func (s *singleByteMatch) Match(t *testing.T, instr uint16, _ []lang.IrObject) {
	t.Helper()

	opcode := bytecode.Opcode(instr >> 8)
	if s.opcode != opcode {
		t.Errorf("expected bytecode.Opcode to be %s, got %s", s.opcode, opcode)
	}
}

func (s *singleByteMatch) withOperand(operand byte) *multiByteMatch {
	return &multiByteMatch{singleByteMatch: s, operand: operand}
}

func (s *singleByteMatch) toHaveOperand(operand byte) *multiByteMatch {
	return &multiByteMatch{singleByteMatch: s, operand: operand}
}

func (s *singleByteMatch) toDefine(name string, matchers []Match) *bodyMatch {
	return &bodyMatch{singleByteMatch: s, name: name, matchers: matchers}
}

type multiByteMatch struct {
	*singleByteMatch
	operand byte
}

func (m *multiByteMatch) toHaveConstant(value interface{}) *constantMatch {
	return &constantMatch{multiByteMatch: m, value: value}
}

func (m *multiByteMatch) toBeMethodCall(name string, argc byte) *methoCallMatch {
	return &methoCallMatch{multiByteMatch: m, name: name, argc: argc}
}

func (m *multiByteMatch) Match(t *testing.T, instr uint16, _ []lang.IrObject) {
	t.Helper()

	opcode := bytecode.Opcode(instr >> 8)
	operand := byte(instr & 255)

	if m.opcode != opcode {
		t.Errorf("expected bytecode.Opcode to be %s, got %s", m.opcode, opcode)
	}

	if m.operand != operand {
		t.Errorf("expected instr(%s)'s operand to be %d, got %d", m.opcode, m.operand, operand)
	}
}

type constantMatch struct {
	*multiByteMatch

	value interface{}
}

func (m *constantMatch) Match(t *testing.T, instr uint16, consts []lang.IrObject) {
	t.Helper()

	opcode := bytecode.Opcode(instr >> 8)
	operand := byte(instr & 255)

	if m.opcode != opcode {
		t.Errorf("expected bytecode.Opcode to be %s, got %s", m.opcode, opcode)
	}

	if m.operand != operand {
		t.Errorf("expected operand for %s to be %d, got %d", opcode, m.operand, operand)
	}

	gotConst := consts[operand]
	var result bool
	switch val := m.value.(type) {
	case int:
		result = lang.Int(val) == gotConst.(lang.Int)
	case float64:
		result = lang.Float(val) == gotConst.(lang.Float)
	case string:
		result = val == string(gotConst.(*lang.String).Value)
	case bool:
		result = lang.Bool(val) == gotConst.(lang.Bool)
	}

	if !result {
		t.Errorf("expected constant for %s to be %v, got %v", opcode, m.value, gotConst)
	}
}

type methoCallMatch struct {
	*multiByteMatch

	name string
	argc byte
}

func (b *bodyMatch) Match(t *testing.T, instr uint16, consts []lang.IrObject) {
	t.Helper()

	opcode := bytecode.Opcode(instr >> 8)
	operand := byte(instr & 255)

	if b.opcode != opcode {
		t.Errorf("expected bytecode.Opcode to be %s, got %s", b.opcode, opcode)
	}

	fun, ok := consts[operand].(*lang.Method)
	if !ok {
		t.Errorf("expected const for %s to be *lang.CompiledFunction, got %T", b.opcode, consts[operand])
	}

	for i, instr := range fun.Instrs() {
		b.matchers[i].Match(t, instr, fun.Constants())
	}
}

func (m *methoCallMatch) Match(t *testing.T, instr uint16, consts []lang.IrObject) {
	t.Helper()

	opcode := bytecode.Opcode(instr >> 8)
	operand := byte(instr & 255)

	if m.opcode != opcode {
		t.Errorf("expected bytecode.Opcode to be %s, got %s", m.opcode, opcode)
	}

	if m.operand != operand {
		t.Errorf("expected operand to be %d, got %d", m.operand, operand)
	}

	ci, ok := consts[operand].(*lang.CallInfo)
	if !ok {
		t.Errorf("expected constant at %d, to be *lang.CallInfo, got %T", operand, consts[operand])
	}

	if ci.Name() != m.name {
		t.Errorf("expected name for %s to be %s, got %s", opcode, m.name, ci.Name())
	}
	if ci.Argc() != m.argc {
		t.Errorf("expected argc for %s to be %d, got %d", opcode, m.argc, ci.Argc())
	}
}

type bodyMatch struct {
	*singleByteMatch

	name     string
	matchers []Match
}

func expect(opcode bytecode.Opcode) *singleByteMatch {
	return &singleByteMatch{opcode: opcode}
}
