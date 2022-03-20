package lang

func checkBoundaries(index int, size int) int {
	if index < 0 {
		index += size
	}

	if index < 0 || index > size-1 {
		panic("IndexOutOfBoundsException")
	}

	return index
}

func ARRAY(obj IrObject) *Array {
	return obj.(*Array)
}

func arrayInsert(self, index, element IrObject) IrObject {
	array := ARRAY(self)

	idx := INT(index)
	size := len(array.Elements)

	pos := checkBoundaries(int(idx), size)

	array.Elements[pos] = element
	return element
}

func arrayAt(self IrObject, index IrObject) IrObject {
	array := ARRAY(self)
	idx := INT(index)
	size := len(array.Elements)
	pos := checkBoundaries(int(idx), size)
	return array.Elements[pos]
}

func arrayPush(self IrObject, elements ...IrObject) IrObject {
	array := ARRAY(self)
	array.Elements = append(array.Elements, elements...)
	return array
}

func arrayFlatten(self IrObject) IrObject {
	array := ARRAY(self)

	var result []IrObject
	for _, element := range array.Elements {
		switch element.(type) {
		case *Array:
			ary := ARRAY(arrayFlatten(element))
			result = append(result, ary.Elements...)
		default:
			result = append(result, element)
		}
	}

	return NewArray(result)
}

func arrayReverse(self IrObject) IrObject {
	array := ARRAY(self)
	length := len(array.Elements)
	elements := make([]IrObject, length)

	for i := length - 1; i >= 0; i-- {
		elements[length-i-1] = array.Elements[i]
	}

	return NewArray(elements)
}

func arrayLength(self IrObject) IrObject {
	array := ARRAY(self)
	length := len(array.Elements)
	return Int(length)
}

func arrayValuesAt(self IrObject, indices ...IrObject) IrObject {
	array := ARRAY(self)
	elements := make([]IrObject, len(indices))

	for i, index := range indices {
		elements[i] = array.Elements[INT(index)]
	}

	return NewArray(elements)
}

func arrayShift(self IrObject, size IrObject) IrObject {
	ary := ARRAY(self)
	n := INT(size)

	if n < 0 {
		panic("negative array size")
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

func arrayPlus(self, other IrObject) IrObject {
	x := ARRAY(self)
	y := ARRAY(other)
	result := make([]IrObject, len(x.Elements)+len(y.Elements))
	i := copy(result, x.Elements)
	copy(result[i:], y.Elements)

	return NewArray(result)
}

var ArrayClass *Class

func InitArray() {
	if ArrayClass != nil {
		return
	}

	ArrayClass = NewClass("Array", ObjectClass)
	ArrayClass.AddGoMethod("at", oneArg(arrayAt))
	ArrayClass.AddGoMethod("get", oneArg(arrayAt))
	ArrayClass.AddGoMethod("+", oneArg(arrayPlus))
	ArrayClass.AddGoMethod("values_at", nArgs(arrayValuesAt))
	ArrayClass.AddGoMethod("push", nArgs(arrayPush))
	ArrayClass.AddGoMethod("insert", twoArgs(arrayInsert))
	ArrayClass.AddGoMethod("reverse", zeroArgs(arrayReverse))
	ArrayClass.AddGoMethod("length", zeroArgs(arrayLength))
	ArrayClass.AddGoMethod("size", zeroArgs(arrayLength))
	ArrayClass.AddGoMethod("flatten", zeroArgs(arrayFlatten))
	ArrayClass.AddGoMethod("shift", oneArg(arrayShift))
}

type Array struct {
	*base

	Elements []IrObject
}

func NewArray(elements []IrObject) *Array {
	return &Array{
		Elements: elements,
		base:     &base{class: ArrayClass},
	}
}
