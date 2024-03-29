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

func arrayInspect(rt Runtime, this IrObject) IrObject {
	array := ARRAY(this)

	if len(array.Elements) == 0 {
		return NewString("[]")
	}

	var buf strings.Builder

	buf.WriteByte('[')
	for i, el := range array.Elements {
		if val := call(rt, el, "inspect"); val != nil {
			if i > 0 {
				buf.WriteString(", ")
			}

			buf.Write(unwrapString(val))
			continue
		}

		return nil
	}
	buf.WriteByte(']')

	return NewString(buf.String())
}

func arrayInsert(rt Runtime, this, index, element IrObject) IrObject {
	array := ARRAY(this)

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

func arrayAt(rt Runtime, this IrObject, index IrObject) IrObject {
	array := ARRAY(this)

	idx := INT(index)
	size := len(array.Elements)

	pos, err := checkBoundaries(int(idx), size)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	return array.Elements[pos]
}

func arrayPush(rt Runtime, this IrObject, elements ...IrObject) IrObject {
	array := ARRAY(this)
	array.Elements = append(array.Elements, elements...)
	return array
}

func arrayFlatten(rt Runtime, this IrObject) IrObject {
	array := ARRAY(this)

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

func arrayReverse(rt Runtime, this IrObject) IrObject {
	array := ARRAY(this)
	length := len(array.Elements)
	elements := make([]IrObject, length)

	for i := length - 1; i >= 0; i-- {
		elements[length-i-1] = array.Elements[i]
	}

	return NewArray(elements)
}

func arrayLength(rt Runtime, this IrObject) IrObject {
	array := ARRAY(this)
	length := len(array.Elements)
	return Int(length)
}

func arrayValuesAt(rt Runtime, this IrObject, indices ...IrObject) IrObject {
	array := ARRAY(this)
	elements := make([]IrObject, len(indices))

	for i, index := range indices {
		elements[i] = array.Elements[INT(index)]
	}

	return NewArray(elements)
}

func arrayHash(rt Runtime, this IrObject) IrObject {
	array := ARRAY(this)

	hash := Int(1)
	for _, el := range array.Elements {
		if code := call(rt, el, "hash"); code != nil {
			hash = hash*31 + INT(code)
		}

		return nil
	}

	return hash
}
func arrayUniq(rt Runtime, this IrObject) IrObject {
	ary := ARRAY(this)
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

func arrayShift(rt Runtime, this IrObject, size IrObject) IrObject {
	ary := ARRAY(this)
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

func arrayEqual(rt Runtime, this IrObject, other IrObject) IrObject {
	if this == other {
		return True
	}

	y, err := toArray(other)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	x := ARRAY(this)
	if len(x.Elements) != len(y.Elements) {
		return False
	}

	for i, el := range x.Elements {
		if ret := call(rt, el, "==", y.Elements[i]); ret != nil {
			if !BOOL(ret) {
				return False
			}

			continue
		}

		return nil
	}

	return True
}

func arrayPlus(rt Runtime, this, other IrObject) IrObject {
	y, err := toArray(other)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	x := ARRAY(this)
	result := make([]IrObject, len(x.Elements)+len(y.Elements))
	i := copy(result, x.Elements)
	copy(result[i:], y.Elements)

	return NewArray(result)
}

func arrayMinus(rt Runtime, this, other IrObject) IrObject {
	y, err := toArray(other)
	if err != nil {
		rt.SetError(err)
		return nil
	}

	h := NewHash()
	for _, el := range y.Elements {
		hashInsert(rt, h, el, True)
	}

	x := ARRAY(this)
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
	ArrayClass.AddGoMethod("to_str", zeroArgs(arrayInspect))
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
