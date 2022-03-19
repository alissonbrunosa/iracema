package lang

import (
	"fmt"
)

func INT(obj IrObject) Int {
	return obj.(Int)
}

func intPlus(self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return left + right
	case Float:
		return Float(left) + right
	default:
		err := fmt.Sprintf("unsupported operand type(s): '%s' + '%s'", left.Class(), right.Class())
		panic(err)
		return nil
	}
}

func intMinus(self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return left - right
	case Float:
		return Float(left) - right
	default:
		err := fmt.Sprintf("unsupported operand type(s): '%s' - '%s'", left.Class(), right.Class())
		panic(err)
		return nil
	}
}

func intMultiply(lhs, rhs IrObject) IrObject {
	left := INT(lhs)
	switch right := rhs.(type) {
	case Int:
		return left * right
	case Float:
		return Float(left) * right
	default:
		err := fmt.Sprintf("unsupported operand type(s): '%s' * '%s'", left.Class(), right.Class())
		panic(err)
		return nil
	}
}

func intDivide(lhs, rhs IrObject) IrObject {
	left := INT(lhs)
	switch right := rhs.(type) {
	case Int:
		if right == 0 {
			panic("divided by 0")
		}

		return left / right
	case Float:
		return Float(left) / right
	default:
		err := fmt.Sprintf("unsupported operand type(s): '%s' / '%s'", left.Class(), right.Class())
		panic(err)
		return nil
	}
}

func intEqual(lhs, rhs IrObject) IrObject {
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

func intGreatThan(self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left > right)
	case Float:
		return NewBoolean(left > Int(right))
	default:
		err := fmt.Sprintf("invalid comparison between '%s' and '%s'", left.Class(), right.Class())
		panic(err)
		return nil
	}
}

func intGreaterThanOrEqual(self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left >= right)
	case Float:
		return NewBoolean(left >= Int(right))
	default:
		err := fmt.Sprintf("invalid comparison between '%s' and '%s'", left.Class(), right.Class())
		panic(err)
		return nil
	}
}

func intLessThanOrEqual(self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left <= right)
	case Float:
		return NewBoolean(left <= Int(right))
	default:
		err := fmt.Sprintf("invalid comparison between '%s' and '%s'", left.Class(), right.Class())
		panic(err)
		return nil
	}
}

func intLessThan(self IrObject, rhs IrObject) IrObject {
	left := INT(self)
	switch right := rhs.(type) {
	case Int:
		return NewBoolean(left < right)
	case Float:
		return NewBoolean(left < Int(right))
	default:
		err := fmt.Sprintf("invalid comparison between '%s' and '%s'", left.Class(), right.Class())
		panic(err)
		return nil
	}
}

func intNegate(self IrObject) IrObject {
	return -INT(self)
}

func intInspect(self IrObject) IrObject {
	inspect := fmt.Sprintf("%d", INT(self))
	return NewString(inspect)
}

func intHash(self IrObject) IrObject {
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
	IntClass.AddGoMethod("+", oneArg(intPlus))
	IntClass.AddGoMethod("plus", oneArg(intPlus))
	IntClass.AddGoMethod("-", oneArg(intMinus))
	IntClass.AddGoMethod("minus", oneArg(intMinus))
	IntClass.AddGoMethod("*", oneArg(intMultiply))
	IntClass.AddGoMethod("multiply", oneArg(intMultiply))
	IntClass.AddGoMethod("/", oneArg(intDivide))
	IntClass.AddGoMethod("divide_by", oneArg(intDivide))
	IntClass.AddGoMethod(">", oneArg(intGreatThan))
	IntClass.AddGoMethod(">=", oneArg(intGreaterThanOrEqual))
	IntClass.AddGoMethod("<", oneArg(intLessThan))
	IntClass.AddGoMethod("<=", oneArg(intLessThanOrEqual))
	IntClass.AddGoMethod("inspect", zeroArgs(intInspect))
	IntClass.AddGoMethod("negate", zeroArgs(intNegate))
}

/*
 Represets integer numbers
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
	for class := IntClass; class != nil; class = class.super {
		if class == class {
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
