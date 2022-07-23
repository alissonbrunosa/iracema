package lang

import (
	"fmt"
	"reflect"
)

func returnFalse(rt Runtime, self IrObject) IrObject {
	return False
}

func objectPuts(rt Runtime, self IrObject, args ...IrObject) IrObject {
	for _, arg := range args {
		inspect := call(rt, arg, "inspect")
		if inspect == nil {
			return nil
		}

		fmt.Println(inspect)
	}

	return None
}

func objectId(rt Runtime, self IrObject) IrObject {
	ptr := reflect.ValueOf(self).Pointer()
	return Int(ptr)
}

func objectEqual(rt Runtime, self IrObject, rhs IrObject) IrObject {
	return Bool(self == rhs)
}

func objectInspect(rt Runtime, self IrObject) IrObject {
	id := objectId(rt, self)
	str := fmt.Sprintf("#<%s:0x%x>", self.Class(), INT(id))

	return NewString(str)
}

func objectUnaryNot(rt Runtime, self IrObject) IrObject {
	return !IsTruthy(self)
}

func objectAlloc(class *Class) IrObject {
	return &Object{
		base:  &base{class: class},
		attrs: make(map[string]IrObject, 3),
	}
}

type Object struct {
	*base

	attrs map[string]IrObject
}

func (o *Object) Set(name string, value IrObject) {
	o.attrs[name] = value
}

func (o *Object) Get(name string) IrObject {
	return o.attrs[name]
}

func (o *Object) String() string {
	return "<Object:" + o.class.Name() + ">"
}

func NewObject() *Object {
	return &Object{
		base: &base{class: ObjectClass},
	}
}

var ObjectClass *Class

func InitObject() {
	if ObjectClass != nil {
		return
	}

	ObjectClass = NewClass("Object", nil)
	ObjectClass.allocator = objectAlloc
	ObjectClass.AddGoMethod("==", oneArg(objectEqual))
	ObjectClass.AddGoMethod("hash", zeroArgs(objectId))
	ObjectClass.AddGoMethod("puts", nArgs(objectPuts))
	ObjectClass.AddGoMethod("object_id", zeroArgs(objectId))
	ObjectClass.AddGoMethod("inspect", zeroArgs(objectInspect))
	ObjectClass.AddGoMethod("nil?", zeroArgs(returnFalse))
	ObjectClass.AddGoMethod("unot", zeroArgs(objectUnaryNot))
}
