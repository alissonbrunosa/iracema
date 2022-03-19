package lang

type MethodType byte

const (
	GoMethod MethodType = 1 << iota
	IrMethod
)

type Method struct {
	methodType MethodType
	name       string
	arity      byte
	body       interface{}
}

func (m *Method) Name() string           { return m.name }
func (m *Method) Arity() byte            { return m.arity }
func (m *Method) Body() interface{}      { return m.body }
func (m *Method) MethodType() MethodType { return m.methodType }

func NewGoMethod(name string, body Native, arity byte) *Method {
	return &Method{
		methodType: GoMethod,
		name:       name,
		arity:      arity,
		body:       body,
	}
}
