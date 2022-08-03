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
	case Float:
		return left + right
	case Int:
		return left + Float(right)
	default:
		err := NewTypeError("unsupported operand type(s): '%s' + '%s'", FloatClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func floatSub(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Float:
		return left - right
	case Int:
		return left - Float(right)
	default:
		err := NewTypeError("unsupported operand type(s): '%s' - '%s'", FloatClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func floatMultiply(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Float:
		return left * right
	case Int:
		return left * Float(right)
	default:
		err := NewTypeError("unsupported operand type(s): '%s' * '%s'", FloatClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func floatDivide(rt Runtime, lhs, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Float:
		return left / right
	case Int:
		return left / Float(right)
	default:
		err := NewTypeError("unsupported operand type(s): '%s' / '%s'", FloatClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func floatEqual(rt Runtime, this IrObject, rhs IrObject) IrObject {
	left := FLOAT(this)

	switch right := rhs.(type) {
	case Float:
		return Bool(left == right)
	case Int:
		return Bool(left == Float(right))
	default:
		return False
	}
}

func floatGreat(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Float:
		return Bool(left > right)
	case Int:
		return Bool(left > Float(right))
	default:
		err := NewTypeError("invalid comparison (>) between '%s' and '%s'", FloatClass, right.Class())
		rt.SetError(err)
		return nil
	}
}
func floatGreatEqual(rt Runtime, this IrObject, rhs IrObject) IrObject {
	left := FLOAT(this)

	switch right := rhs.(type) {
	case Float:
		return Bool(left >= right)
	case Int:
		return Bool(left >= Float(right))
	default:
		err := NewTypeError("invalid comparison (>=) between '%s' and '%s'", FloatClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func floatLess(rt Runtime, lhs IrObject, rhs IrObject) IrObject {
	left := FLOAT(lhs)

	switch right := rhs.(type) {
	case Float:
		return Bool(left < right)
	case Int:
		return Bool(left < Float(right))
	default:
		err := NewTypeError("invalid comparison (<) between '%s' and '%s'", FloatClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func floatLessEqual(rt Runtime, this IrObject, rhs IrObject) IrObject {
	left := FLOAT(this)

	switch right := rhs.(type) {
	case Float:
		return Bool(left <= right)
	case Int:
		return Bool(left <= Float(right))
	default:
		err := NewTypeError("invalid comparison (<=) between '%s' and '%s'", FloatClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func floatUnaryAdd(rt Runtime, this IrObject) IrObject {
	return this
}

func floatUnarySub(rt Runtime, this IrObject) IrObject {
	return -FLOAT(this)
}

func floatInspect(rt Runtime, this IrObject) IrObject {
	value := FLOAT(this)
	inspect := fmt.Sprintf("%f", value)
	return NewString(inspect)
}

func floatHash(rt Runtime, this IrObject) IrObject {
	value := FLOAT(this)
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
	FloatClass.AddGoMethod(">", oneArg(floatGreat))
	FloatClass.AddGoMethod(">=", oneArg(floatGreatEqual))
	FloatClass.AddGoMethod("<", oneArg(floatLess))
	FloatClass.AddGoMethod("<=", oneArg(floatLessEqual))
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
