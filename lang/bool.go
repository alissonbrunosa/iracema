package lang

import "fmt"

func BOOL(obj IrObject) Bool {
	return obj.(Bool)
}

var BoolClass *Class

var True Bool = true
var False Bool = false

func boolInspect(rt Runtime, self IrObject) IrObject {
	inspect := BOOL(self).String()
	return NewString(inspect)
}

func InitBool() {
	if BoolClass != nil {
		return
	}

	BoolClass = NewClass("Bool", ObjectClass)
	BoolClass.AddGoMethod("inspect", zeroArgs(boolInspect))
}

/*
Represents boolean object
*/
type Bool bool

func (b Bool) String() string {
	return fmt.Sprintf("%t", b)
}

func (Bool) LookupMethod(name string) *Method {
	for class := BoolClass; class != nil; class = class.super {
		if method, ok := class.methods[name]; ok {
			return method
		}
	}

	return nil
}

func (Bool) Is(class *Class) bool {
	for cls := BoolClass; cls != nil; cls = cls.super {
		if cls == class {
			return true
		}
	}

	return false
}

func (Bool) Class() *Class { return BoolClass }

func NewBoolean(value bool) Bool {
	return Bool(value)
}
