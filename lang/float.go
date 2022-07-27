package lang

import (
	"fmt"
	"unsafe"
)

func FLOAT(obj IrObject) Float {
	return obj.(Float)
}

func floatAdd(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Int:
		return left + Float(right)
	case Float:
		return left + right
	default:
		err := fmt.Sprintf("can't operate on %T", right)
		panic(err)
	}
}

func floatSub(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Int:
		return left - Float(right)
	case Float:
		return left - right
	default:
		err := fmt.Sprintf("can't operate on %T", right)
		panic(err)
	}
}

func floatMultiply(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Int:
		return left * Float(right)
	case Float:
		return left * right
	default:
		err := fmt.Sprintf("can't operate on %T", right)
		panic(err)
	}
}

func floatDivide(rt Runtime, lhs, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Int:
		return left / Float(right)
	case Float:
		return left / right
	default:
		err := fmt.Sprintf("can't operate on %T", right)
		panic(err)
	}
}

func floatEqual(rt Runtime, self IrObject, rhs IrObject) IrObject {
	left := FLOAT(self)

	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left == Float(right))
	case Float:
		return NewBoolean(left == right)
	default:
		return False
	}
}

func floatUnaryAdd(rt Runtime, self IrObject) IrObject {
	return self
}

func floatUnarySub(rt Runtime, self IrObject) IrObject {
	return -FLOAT(self)
}

func floatGreatThan(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left > Float(right))
	case Float:
		return NewBoolean(left > right)
	default:
		err := fmt.Sprintf("can't operate on %T", right)
		panic(err)
	}
}

func floatLessThan(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left < Float(right))
	case Float:
		return NewBoolean(left < right)
	default:
		err := fmt.Sprintf("can't operate on %T", right)
		panic(err)
	}
}

func floatInspect(rt Runtime, self IrObject) IrObject {
	value := FLOAT(self)
	inspect := fmt.Sprintf("%f", value)
	return NewString(inspect)
}

func floatHash(rt Runtime, self IrObject) IrObject {
	value := FLOAT(self)
	bits := *(*uint64)(unsafe.Pointer(&value))
	bits = bits ^ (bits >> 32)
	return Int(bits)
}

var FloatClass *Class

func InitFloat() {
	if FloatClass != nil {
		return
	}

	FloatClass = NewClass("Float", ObjectClass)
	FloatClass.AddGoMethod("==", oneArg(floatEqual))
	FloatClass.AddGoMethod("hash", zeroArgs(floatHash))
	FloatClass.AddGoMethod("+", oneArg(floatAdd))
	FloatClass.AddGoMethod("add", oneArg(floatAdd))
	FloatClass.AddGoMethod("-", oneArg(floatSub))
	FloatClass.AddGoMethod("sub", oneArg(floatSub))
	FloatClass.AddGoMethod("*", oneArg(floatMultiply))
	FloatClass.AddGoMethod("multiply", oneArg(floatMultiply))
	FloatClass.AddGoMethod("/", oneArg(floatDivide))
	FloatClass.AddGoMethod("divide", oneArg(floatDivide))
	FloatClass.AddGoMethod(">", oneArg(floatGreatThan))
	FloatClass.AddGoMethod("<", oneArg(floatLessThan))
	FloatClass.AddGoMethod("inspect", zeroArgs(floatInspect))
	FloatClass.AddGoMethod("uadd", zeroArgs(floatUnaryAdd))
	FloatClass.AddGoMethod("usub", zeroArgs(floatUnarySub))
}

/*
Represents float numbers
*/
type Float float64

func (Float) LookupMethod(name string) *Method {
	for class := FloatClass; class != nil; class = class.super {
		if method, ok := class.methods[name]; ok {
			return method
		}
	}

	return nil
}

func (Float) Is(class *Class) bool {
	for cls := FloatClass; cls != nil; cls = cls.super {
		if cls == class {
			return true
		}
	}

	return false
}

func (Float) Class() *Class {
	return FloatClass
}

func (f Float) String() string {
	return fmt.Sprintf("%.2f", f)
}
