package lang

import (
	"fmt"
)

var (
	Error             *Class
	RuntimeError      *Class
	TypeError         *Class
	NameError         *Class
	RegexpError       *Class
	NoMethodError     *Class
	ArgumentError     *Class
	ZeroDivisionError *Class
)

func ERROR(obj IrObject) *ErrorObject {
	return obj.(*ErrorObject)
}

func errMessage(rt Runtime, this IrObject) IrObject {
	err := ERROR(this)
	return NewString(err.message)
}

func errInit(rt Runtime, this IrObject, message IrObject) IrObject {
	err := ERROR(this)
	err.message = string(unwrapString(message))
	return None
}

func errAlloc(class *Class) IrObject {
	return &ErrorObject{
		base: &base{class: class},
	}
}

type ErrorObject struct {
	*base

	message string
}

func (err *ErrorObject) String() string {
	return err.message
}

func InitError() {
	Error = NewClass("Error", ObjectClass)
	Error.allocator = errAlloc
	Error.AddGoMethod("init", oneArg(errInit))
	Error.AddGoMethod("message", zeroArgs(errMessage))
	Error.AddGoMethod("inspect", zeroArgs(errMessage))
	Error.AddGoMethod("to_str", zeroArgs(errMessage))

	NameError = NewClass("NameError", Error)
	RegexpError = NewClass("RegexpError", Error)
	ArgumentError = NewClass("ArgumentError", Error)
	NoMethodError = NewClass("NoMethodError", NameError)
	RuntimeError = NewClass("RuntimeError", Error)

	TypeError = NewClass("TypeError", RuntimeError)
	ZeroDivisionError = NewClass("ZeroDivisionError", RuntimeError)
}

func NewNoMethodError(recv IrObject, name string) *ErrorObject {
	mesg := fmt.Sprintf("undefined method '%s' for %s", name, recv.Class())

	return &ErrorObject{
		message: mesg,
		base:    &base{class: NoMethodError},
	}
}

func NewArityError(given, expected int) *ErrorObject {
	mesg := fmt.Sprintf("wrong number of arguments (given %d, expected %d)", given, expected)

	return &ErrorObject{
		message: mesg,
		base:    &base{class: ArgumentError},
	}
}

func NewNameError(name IrObject) *ErrorObject {
	mesg := fmt.Sprintf("uninitialized constant %s", name)

	return &ErrorObject{
		message: mesg,
		base:    &base{class: NameError},
	}
}

func NewRegexpError(mesg string) *ErrorObject {
	return &ErrorObject{
		message: mesg,
		base:    &base{class: RegexpError},
	}
}

func NewTypeError(mesg string, args ...any) *ErrorObject {
	return &ErrorObject{
		message: fmt.Sprintf(mesg, args...),
		base:    &base{class: TypeError},
	}
}

func NewError(msg string, class *Class, args ...any) *ErrorObject {
	return &ErrorObject{
		message: fmt.Sprintf(msg, args...),
		base:    &base{class: class},
	}
}
