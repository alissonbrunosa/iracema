package lang

func noneInspect(rt Runtime, self IrObject) IrObject {
	return NewString("none")
}

func noneCheck(rt Runtime, self IrObject) IrObject {
	return True
}

var None none = 0x01
var NoneClass *Class

func InitNone() {
	if NoneClass != nil {
		return
	}

	NoneClass = NewClass("None", ObjectClass)
	NoneClass.AddGoMethod("inspect", zeroArgs(noneInspect))
	NoneClass.AddGoMethod("none?", zeroArgs(noneCheck))
}

/*
Represets nil object
*/

type none byte

func (none) String() string {
	return "none"
}

func (none) LookupMethod(name string) *Method {
	for class := NoneClass; class != nil; class = class.super {
		if method, ok := class.methods[name]; ok {
			return method
		}
	}

	return nil
}

func (none) Is(parent *Class) bool {
	for class := NoneClass; class != nil; class = class.super {
		if class == parent {
			return true
		}
	}

	return false
}

func (none) Class() *Class { return NoneClass }
