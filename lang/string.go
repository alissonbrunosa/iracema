package lang

import (
	"bytes"
	"fmt"
	"strings"
)

func toString(value IrObject) (*String, *ErrorObject) {
	if s, ok := value.(*String); ok {
		return s, nil
	}

	var mesg = new(strings.Builder)
	fmt.Fprintf(mesg, "no implicit conversion of %s into Regexp", value.Class())
	return nil, NewTypeError(mesg.String())
}

func STRING(obj IrObject) *String {
	return obj.(*String)
}

func unwrapString(obj IrObject) []byte {
	return STRING(obj).Value
}

func stringSize(rt Runtime, self IrObject) IrObject {
	str := self.(*String)
	return Int(len(str.Value))
}

func stringEqual(rt Runtime, self IrObject, rhs IrObject) IrObject {
	left := STRING(self)
	right := STRING(rhs)

	return NewBoolean(bytes.Equal(left.Value, right.Value))
}

func stringPlus(rt Runtime, self IrObject, rhs IrObject) IrObject {
	var buf strings.Builder

	left := unwrapString(self)
	right := unwrapString(rhs)

	buf.Grow(len(left) + len(right))
	buf.Write(left)
	buf.Write(right)

	return NewString(buf.String())
}

func stringHash(rt Runtime, self IrObject) IrObject {
	str := unwrapString(self)

	var hash Int
	if len(str) > 0 {
		for i := 0; i < len(str); i++ {
			hash = 31*hash + Int(str[i])
		}
	}

	return hash
}

func stringInspect(rt Runtime, self IrObject) IrObject {
	return self
}

var StringClass *Class

func InitString() {
	if StringClass != nil {
		return
	}

	StringClass = NewClass("String", ObjectClass)

	StringClass.AddGoMethod("==", oneArg(stringEqual))
	StringClass.AddGoMethod("hash", zeroArgs(stringHash))
	StringClass.AddGoMethod("size", zeroArgs(stringSize))
	StringClass.AddGoMethod("+", oneArg(stringPlus))
	StringClass.AddGoMethod("inspect", zeroArgs(stringInspect))
}

/*
Represets strings object
*/

type String struct {
	*base
	Value []byte
}

func (s *String) String() string {
	return string(s.Value)
}

/*
Creates a new string object
*/
func NewString(value string) *String {
	return &String{
		Value: []byte(value),

		base: &base{class: StringClass},
	}
}
