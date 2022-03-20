package lang

type MethodType byte

const (
	GoMethod MethodType = 1 << iota
	IrMethod
)

type Method struct {
	*base

	methodType MethodType
	name       string
	arity      byte
	body       interface{}
	localCount byte
	constants  []IrObject
}

func (m *Method) Name() string           { return m.name }
func (m *Method) Arity() byte            { return m.arity }
func (m *Method) Body() interface{}      { return m.body }
func (m *Method) Instrs() []uint16       { return m.body.([]uint16) }
func (m *Method) MethodType() MethodType { return m.methodType }
func (m *Method) Constants() []IrObject  { return m.constants }
func (m *Method) LocalCout() byte        { return m.localCount }

func NewGoMethod(name string, body Native, arity byte) *Method {
	return &Method{
		methodType: GoMethod,
		name:       name,
		arity:      arity,
		body:       body,
	}
}

func NewIrMethod(name string, arity byte, body []uint16, localCount byte, consts []IrObject) *Method {
	return &Method{
		methodType: IrMethod,
		name:       name,
		arity:      arity,
		body:       body,
		localCount: localCount,
		constants:  consts,
	}
}
