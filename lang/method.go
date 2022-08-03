package lang

type MethodType byte

const (
	GoFunction MethodType = 1 << iota
	IrMethod
)

type Method struct {
	*base

	methodType  MethodType
	name        string
	arity       byte
	body        interface{}
	localCount  byte
	constants   []IrObject
	catchOffset int
}

func (m *Method) Name() string           { return m.name }
func (m *Method) Arity() byte            { return m.arity }
func (m *Method) Native() Native         { return m.body.(Native) }
func (m *Method) Instrs() []uint16       { return m.body.([]uint16) }
func (m *Method) MethodType() MethodType { return m.methodType }
func (m *Method) Constants() []IrObject  { return m.constants }
func (m *Method) LocalCount() byte       { return m.localCount }
func (m *Method) CatchOffset() int       { return m.catchOffset }

func NewGoMethod(name string, body Native, arity byte) *Method {
	return &Method{
		methodType: GoFunction,
		name:       name,
		arity:      arity,
		body:       body,
	}
}

func NewIrMethod(name string, arity byte, body []uint16, localCount byte, consts []IrObject, catchOffset int) *Method {
	return &Method{
		methodType:  IrMethod,
		name:        name,
		arity:       arity,
		body:        body,
		localCount:  localCount,
		constants:   consts,
		catchOffset: catchOffset,
	}
}
