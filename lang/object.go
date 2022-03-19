package lang

import (
	"fmt"
	"reflect"
)

func returnFalse(self IrObject) IrObject {
	return False
}

func objectId(self IrObject) IrObject {
	ptr := reflect.ValueOf(self).Pointer()
	return Int(ptr)
}

func objectEqual(self IrObject, rhs IrObject) IrObject {
	return Bool(self == rhs)
}

func objectInspect(self IrObject) IrObject {
	id := objectId(self)
	str := fmt.Sprintf("#<%s:0x%x>", self.Class(), INT(id))

	return NewString(str)
}

var ObjectClass *Class

func InitObject() {
	if ObjectClass != nil {
		return
	}

	ObjectClass = NewClass("Object", nil)
	ObjectClass.AddGoMethod("==", oneArg(objectEqual))
	ObjectClass.AddGoMethod("hash", zeroArgs(objectId))
	ObjectClass.AddGoMethod("object_id", zeroArgs(objectId))
	ObjectClass.AddGoMethod("inspect", zeroArgs(objectInspect))
	ObjectClass.AddGoMethod("nil?", zeroArgs(returnFalse))
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
