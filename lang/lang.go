package lang

var classes map[string]*Class

func init() {
	InitClass()
	InitObject()
	InitError()
	InitString()
	InitFloat()
	InitInt()
	InitNone()
	InitBool()
	InitHash()
	InitArray()

	classes = map[string]*Class{
		"Object":  ObjectClass,
		"Int":     IntClass,
		"Float":   FloatClass,
		"String":  StringClass,
		"None":    NoneClass,
		"Boolean": BoolClass,
		"Hash":    HashClass,
		"Array":   ArrayClass,

		"Error":             Error,
		"NameError":         NameError,
		"RuntimeError":      RuntimeError,
		"ArgumentError":     ArgumentError,
		"NoMethodError":     NoMethodError,
		"ZeroDivisionError": ZeroDivisionError,
	}
}

func IsTruthy(obj IrObject) Bool {
	if obj.Class() == NoneClass {
		return false
	}

	if obj.Class() != BoolClass {
		return true
	}

	return BOOL(obj)
}

func TypeLookup(name IrObject) *Class {
	n := unwrapString(name)
	return classes[string(n)]
}

func DefineType(name string, class *Class) {
	classes[name] = class
}

type IrObject interface {
	Class() *Class
	Is(*Class) bool
}

type Runtime interface {
	SetError(*ErrorObject)
	Call(IrObject, *Method, ...IrObject) IrObject
}
