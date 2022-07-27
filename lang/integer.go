package lang

import (
	"fmt"
	"strings"
)

func INT(obj IrObject) Int {
	return obj.(Int)
}

func toInt(value IrObject) (Int, *ErrorObject) {
	if i, ok := value.(Int); ok {
		return i, nil
	}

	var mesg = new(strings.Builder)
	fmt.Fprintf(mesg, "no implicit conversion of %s into Int", value.Class())
	return 0, NewError(mesg.String(), TypeError)
}

func intAdd(rt Runtime, self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return left + right
	case Float:
		return Float(left) + right
	default:
		err := NewTypeError("unsupported operand type(s): '%s' + '%s'", IntClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func intSub(rt Runtime, self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return left - right
	case Float:
		return Float(left) - right
	default:
		err := NewTypeError("unsupported operand type(s): '%s' - '%s'", IntClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func intMultiply(rt Runtime, lhs, rhs IrObject) IrObject {
	left := INT(lhs)
	switch right := rhs.(type) {
	case Int:
		return left * right
	case Float:
		return Float(left) * right
	default:
		err := NewTypeError("unsupported operand type(s): '%s' * '%s'", IntClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func intDivide(rt Runtime, lhs, rhs IrObject) IrObject {
	left := INT(lhs)
	switch right := rhs.(type) {
	case Int:
		if right == 0 {
			err := NewError("divided by 0", ZeroDivisionError)
			rt.SetError(err)
			return nil
		}
		return left / right
	case Float:
		return Float(left) / right
	default:
		err := NewTypeError("unsupported operand type(s): '%s' / '%s'", IntClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func intEqual(rt Runtime, lhs, rhs IrObject) IrObject {
	left := INT(lhs)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left == right)
	case Float:
		return NewBoolean(left == Int(right))
	default:
		return False
	}
}

func intGreat(rt Runtime, self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left > right)
	case Float:
		return NewBoolean(left > Int(right))
	default:
		err := NewTypeError("invalid comparison (>) between '%s' and '%s'", IntClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func intGreatEqual(rt Runtime, self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left >= right)
	case Float:
		return NewBoolean(left >= Int(right))
	default:
		err := NewTypeError("invalid comparison (>=) between '%s' and '%s'", IntClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func intLess(rt Runtime, self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left < right)
	case Float:
		return NewBoolean(left < Int(right))
	default:
		err := NewTypeError("invalid comparison (<) between '%s' and '%s'", IntClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func intLessEqual(rt Runtime, self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left <= right)
	case Float:
		return NewBoolean(left <= Int(right))
	default:
		err := NewTypeError("invalid comparison (<=) between '%s' and '%s'", IntClass, right.Class())
		rt.SetError(err)
		return nil
	}
}

func intUnaryAdd(rt Runtime, self IrObject) IrObject {
	return self
}

func intUnarySub(rt Runtime, self IrObject) IrObject {
	return -INT(self)
}

func intInspect(rt Runtime, self IrObject) IrObject {
	inspect := fmt.Sprintf("%d", INT(self))
	return NewString(inspect)
}

func intHash(rt Runtime, self IrObject) IrObject {
	return self
}

var IntClass *Class

func InitInt() {
	if IntClass != nil {
		return
	}

	IntClass = NewClass("Int", ObjectClass)
	IntClass.AddGoMethod("==", oneArg(intEqual))
	IntClass.AddGoMethod("hash", zeroArgs(intHash))
	IntClass.AddGoMethod("+", oneArg(intAdd))
	IntClass.AddGoMethod("add", oneArg(intAdd))
	IntClass.AddGoMethod("-", oneArg(intSub))
	IntClass.AddGoMethod("sub", oneArg(intSub))
	IntClass.AddGoMethod("*", oneArg(intMultiply))
	IntClass.AddGoMethod("multiply", oneArg(intMultiply))
	IntClass.AddGoMethod("/", oneArg(intDivide))
	IntClass.AddGoMethod("divide", oneArg(intDivide))
	IntClass.AddGoMethod(">", oneArg(intGreat))
	IntClass.AddGoMethod(">=", oneArg(intGreatEqual))
	IntClass.AddGoMethod("<", oneArg(intLess))
	IntClass.AddGoMethod("<=", oneArg(intLessEqual))
	IntClass.AddGoMethod("inspect", zeroArgs(intInspect))
	IntClass.AddGoMethod("uadd", zeroArgs(intUnaryAdd))
	IntClass.AddGoMethod("usub", zeroArgs(intUnarySub))
}

/*
Represents integer numbers
*/
type Int int

func (Int) LookupMethod(name string) *Method {
	for class := IntClass; class != nil; class = class.super {
		if method, ok := class.methods[name]; ok {
			return method
		}
	}

	return nil
}

func (Int) Is(class *Class) bool {
	for cls := IntClass; cls != nil; cls = cls.super {
		if cls == class {
			return true
		}
	}

	return false
}

func (Int) Class() *Class {
	return IntClass
}

func (i Int) String() string {
	return fmt.Sprintf("%d", i)
}
