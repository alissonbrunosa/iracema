package lang

import (
	"fmt"
	"strings"
)

func toArray(value IrObject) (*Array, *ErrorObject) {
	if array, ok := value.(*Array); ok {
		return array, nil
	}

	var mesg = new(strings.Builder)
	fmt.Fprintf(mesg, "TypeError: no implicit conversion of %s into Array", value.Class())
	return nil, NewError(mesg.String(), RuntimeError)
}

func checkBoundaries(index int, size int) (int, *ErrorObject) {
	if index < 0 {
		index += size
	}

	if index < 0 || index > size-1 {
		return 0, NewError("IndexOutOfBoundsException", RuntimeError)
	}

	return index, nil
}

func ARRAY(obj IrObject) *Array {
	return obj.(*Array)
}

func arrayInspect(rt Runtime, self IrObject) IrObject {
	array := ARRAY(self)

	if len(array.Elements) == 0 {
		return NewString("[]")
	}

	var inspect IrObject
	var buf strings.Builder

	buf.WriteByte('[')
	for i, el := range array.Elements {
		inspect = call(rt, el, "inspect")
		if i > 0 {
			buf.WriteString(", ")
		}

		buf.Write(unwrapString(inspect))
	}
	buf.WriteByte(']')

	return NewString(buf.String())
}

func arrayInsert(rt Runtime, self, index, element IrObject) IrObject {
	array := ARRAY(self)

	idx := INT(index)
	size := len(array.Elements)

	pos, err := checkBoundaries(int(idx), size)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	array.Elements[pos] = element
	return element
}

func arrayAt(rt Runtime, self IrObject, index IrObject) IrObject {
	array := ARRAY(self)

	idx := INT(index)
	size := len(array.Elements)

	pos, err := checkBoundaries(int(idx), size)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	return array.Elements[pos]
}

func arrayPush(rt Runtime, self IrObject, elements ...IrObject) IrObject {
	array := ARRAY(self)
	array.Elements = append(array.Elements, elements...)
	return array
}

func arrayFlatten(rt Runtime, self IrObject) IrObject {
	array := ARRAY(self)

	var result []IrObject
	for _, element := range array.Elements {
		switch element.(type) {
		case *Array:
			ary := ARRAY(arrayFlatten(rt, element))
			result = append(result, ary.Elements...)
		default:
			result = append(result, element)
		}
	}

	return NewArray(result)
}

func arrayReverse(rt Runtime, self IrObject) IrObject {
	array := ARRAY(self)
	length := len(array.Elements)
	elements := make([]IrObject, length)

	for i := length - 1; i >= 0; i-- {
		elements[length-i-1] = array.Elements[i]
	}

	return NewArray(elements)
}

func arrayLength(rt Runtime, self IrObject) IrObject {
	array := ARRAY(self)
	length := len(array.Elements)
	return Int(length)
}

func arrayValuesAt(rt Runtime, self IrObject, indices ...IrObject) IrObject {
	array := ARRAY(self)
	elements := make([]IrObject, len(indices))

	for i, index := range indices {
		elements[i] = array.Elements[INT(index)]
	}

	return NewArray(elements)
}

func arrayHash(rt Runtime, self IrObject) IrObject {
	array := ARRAY(self)

	hash := Int(1)
	for _, el := range array.Elements {
		code := call(rt, el, "hash")
		hash = hash*31 + INT(code)
	}

	return hash
}
func arrayUniq(rt Runtime, self IrObject) IrObject {
	ary := ARRAY(self)
	if len(ary.Elements) <= 1 {
		newEls := make([]IrObject, len(ary.Elements))
		copy(newEls, ary.Elements)
		return NewArray(newEls)
	}

	hash := NewHash()
	for _, el := range ary.Elements {
		hashInsert(rt, hash, el, el)
	}

	return hashValues(rt, hash)
}

func arrayShift(rt Runtime, self IrObject, size IrObject) IrObject {
	ary := ARRAY(self)
	n, err := toInt(size)
	if err != nil {
		return nil
	}

	if n < 0 {
		err := NewError("negative array size", ArgumentError)
		rt.SetError(err)
		return nil
	}

	if n == 1 {
		first := ary.Elements[0]
		ary.Elements = ary.Elements[1:]
		return first
	}

	if n >= Int(len(ary.Elements)) {
		result := make([]IrObject, len(ary.Elements))
		copy(result, ary.Elements)
		ary.Elements = nil

		return NewArray(result)
	}

	head := ary.Elements[:n]
	result := make([]IrObject, n)
	copy(result, head)

	ary.Elements = ary.Elements[n:]
	return NewArray(result)
}

func arrayEqual(rt Runtime, self IrObject, other IrObject) IrObject {
	if self == other {
		return True
	}

	y, err := toArray(other)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	x := ARRAY(self)
	if len(x.Elements) != len(y.Elements) {
		return False
	}

	for i, el := range x.Elements {
		ret := call(rt, el, "==", y.Elements[i])
		if !BOOL(ret) {
			return False
		}
	}

	return True
}

func arrayPlus(rt Runtime, self, other IrObject) IrObject {
	y, err := toArray(other)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	x := ARRAY(self)
	result := make([]IrObject, len(x.Elements)+len(y.Elements))
	i := copy(result, x.Elements)
	copy(result[i:], y.Elements)

	return NewArray(result)
}

func arrayMinus(rt Runtime, self, other IrObject) IrObject {
	y, err := toArray(other)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	h := NewHash()
	for _, el := range y.Elements {
		hashInsert(rt, h, el, True)
	}

	x := ARRAY(self)
	var result []IrObject
	for _, el := range x.Elements {
		if BOOL(hashHasKey(rt, h, el)) {
			continue
		}
		result = append(result, el)
	}

	return NewArray(result)
}

var ArrayClass *Class

func InitArray() {
	if ArrayClass != nil {
		return
	}

	ArrayClass = NewClass("Array", ObjectClass)
	ArrayClass.AddGoMethod("==", oneArg(arrayEqual))
	ArrayClass.AddGoMethod("hash", zeroArgs(arrayHash))
	ArrayClass.AddGoMethod("at", oneArg(arrayAt))
	ArrayClass.AddGoMethod("get", oneArg(arrayAt))
	ArrayClass.AddGoMethod("+", oneArg(arrayPlus))
	ArrayClass.AddGoMethod("-", oneArg(arrayMinus))
	ArrayClass.AddGoMethod("values_at", nArgs(arrayValuesAt))
	ArrayClass.AddGoMethod("push", nArgs(arrayPush))
	ArrayClass.AddGoMethod("insert", twoArgs(arrayInsert))
	ArrayClass.AddGoMethod("reverse", zeroArgs(arrayReverse))
	ArrayClass.AddGoMethod("length", zeroArgs(arrayLength))
	ArrayClass.AddGoMethod("size", zeroArgs(arrayLength))
	ArrayClass.AddGoMethod("inspect", zeroArgs(arrayInspect))
	ArrayClass.AddGoMethod("flatten", zeroArgs(arrayFlatten))
	ArrayClass.AddGoMethod("uniq", zeroArgs(arrayUniq))
	ArrayClass.AddGoMethod("shift", oneArg(arrayShift))
}

type Array struct {
	*base

	Elements []IrObject
}

func (a *Array) Length() int {
	return len(a.Elements)
}

func (a *Array) At(idx int) IrObject {
	return a.Elements[idx]
}

func NewArray(elements []IrObject) *Array {
	return &Array{
		Elements: elements,
		base:     &base{class: ArrayClass},
	}
}
