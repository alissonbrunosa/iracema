package lang

type IrObject interface {
	Class() *Class
	Is(*Class) bool
	LookupMethod(methodName string) *Method
}
