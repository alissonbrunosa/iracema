package lang

import (
	"bytes"
	"strings"
)

func STRING(obj IrObject) *String {
	return obj.(*String)
}

func unwrapString(obj IrObject) []byte {
	return STRING(obj).Value
}

func stringSize(rt Runtime, this IrObject) IrObject {
	str := this.(*String)
	return Int(len(str.Value))
}

func stringEqual(rt Runtime, this IrObject, rhs IrObject) IrObject {
	left := STRING(this)
	right := STRING(rhs)

	return Bool(bytes.Equal(left.Value, right.Value))
}

func stringPlus(rt Runtime, this IrObject, rhs IrObject) IrObject {
	var buf strings.Builder

	left := unwrapString(this)
	right := unwrapString(rhs)

	buf.Grow(len(left) + len(right))
	buf.Write(left)
	buf.Write(right)

	return NewString(buf.String())
}

func stringHash(rt Runtime, this IrObject) IrObject {
	str := unwrapString(this)

	var hash Int
	if len(str) > 0 {
		for i := 0; i < len(str); i++ {
			hash = 31*hash + Int(str[i])
		}
	}

	return hash
}

func stringInspect(rt Runtime, this IrObject) IrObject {
	return this
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
Represents strings object
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
