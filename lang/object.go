package lang

import (
	"fmt"
	"reflect"
)

func returnFalse(rt Runtime, this IrObject) IrObject {
	return False
}

func objectPuts(rt Runtime, this IrObject, args ...IrObject) IrObject {
	for _, arg := range args {
		v := call(rt, arg, "to_str")
		if v == nil {
			return nil
		}

		fmt.Println(v)
	}

	return None
}

func objectInit(rt Runtime, this IrObject) IrObject {
	return None
}

func objectId(rt Runtime, this IrObject) IrObject {
	ptr := reflect.ValueOf(this).Pointer()
	return Int(ptr)
}

func objectEqual(rt Runtime, this IrObject, rhs IrObject) IrObject {
	return Bool(this == rhs)
}

func objectNotEqual(rt Runtime, this IrObject, other IrObject) IrObject {
	result := call(rt, this, "==", other)
	return !IsTruthy(result)
}

func objectInspect(rt Runtime, this IrObject) IrObject {
	id := objectId(rt, this)
	str := fmt.Sprintf("#<%s:0x%x>", this.Class(), INT(id))

	return NewString(str)
}

func objectUnaryNot(rt Runtime, this IrObject) IrObject {
	return !IsTruthy(this)
}

type Object struct {
	*base

	values []IrObject
}

func (o *Object) Set(pos byte, value IrObject) {
	o.values[pos] = value
}

func (o *Object) Get(pos byte) IrObject {
	return o.values[pos]
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

	ObjectClass.allocator = func(class *Class) IrObject {
		return &Object{
			base:   &base{class: class},
			values: make([]IrObject, len(class.fields)),
		}
	}

	ObjectClass.AddGoMethod("init", zeroArgs(objectInit))
	ObjectClass.AddGoMethod("==", oneArg(objectEqual))
	ObjectClass.AddGoMethod("!=", oneArg(objectNotEqual))
	ObjectClass.AddGoMethod("hash", zeroArgs(objectId))
	ObjectClass.AddGoMethod("puts", nArgs(objectPuts))
	ObjectClass.AddGoMethod("object_id", zeroArgs(objectId))
	ObjectClass.AddGoMethod("inspect", zeroArgs(objectInspect))
	ObjectClass.AddGoMethod("to_str", zeroArgs(objectInspect))
	ObjectClass.AddGoMethod("nil?", zeroArgs(returnFalse))
	ObjectClass.AddGoMethod("unot", zeroArgs(objectUnaryNot))
}
