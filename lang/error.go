package lang

import (
	"fmt"
	"strings"
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

func errMessage(rt Runtime, self IrObject) IrObject {
	err := ERROR(self)
	return NewString(err.message)
}

type ErrorObject struct {
	*base

	message string
	stack   []string
}

func (err *ErrorObject) String() string {
	return err.message
}

func InitError() {
	Error = NewClass("Error", ObjectClass)
	Error.AddGoMethod("message", zeroArgs(errMessage))
	Error.AddGoMethod("inspect", zeroArgs(errMessage))

	NameError = NewClass("NameError", Error)
	RegexpError = NewClass("RegexpError", Error)
	ArgumentError = NewClass("ArgumentError", Error)
	NoMethodError = NewClass("NoMethodError", NameError)
	RuntimeError = NewClass("RuntimeError", Error)

	TypeError = NewClass("TypeError", RuntimeError)
	ZeroDivisionError = NewClass("ZeroDivisionError", RuntimeError)
}

func NewNoMethodError(recv IrObject, name string) *ErrorObject {
	var buf = new(strings.Builder)
	fmt.Fprintf(buf, "undefined method '%s' for %s", name, recv.Class())

	return &ErrorObject{
		message: buf.String(),
		base:    &base{class: NoMethodError},
	}
}

func NewArityError(given, expected int) *ErrorObject {
	var buf = new(strings.Builder)
	fmt.Fprintf(buf, "wrong number of arguments (given %d, expected %d)", given, expected)

	return &ErrorObject{
		message: buf.String(),
		base:    &base{class: ArgumentError},
	}
}

func NewNameError(name IrObject) *ErrorObject {
	var buf = new(strings.Builder)
	fmt.Fprintf(buf, "uninitialized constant %s", name)

	return &ErrorObject{
		message: buf.String(),
		base:    &base{class: NameError},
	}
}

func NewRegexpError(mesg string) *ErrorObject {
	return &ErrorObject{
		message: mesg,
		base:    &base{class: RegexpError},
	}
}

func NewTypeError(mesg string) *ErrorObject {
	return &ErrorObject{
		message: mesg,
		base:    &base{class: TypeError},
	}
}

func NewError(msg string, class *Class) *ErrorObject {
	return &ErrorObject{
		message: msg,
		base:    &base{class: class},
	}
}
