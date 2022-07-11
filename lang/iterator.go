package lang

type Iterator struct {
	index int
	list  *Array

	*base
}

func (i *Iterator) HasNext() bool {
	return i.index != i.list.Length()
}

func (i *Iterator) Next() IrObject {
	item := i.list.At(i.index)
	i.index++
	return item
}

func NewIterator(obj IrObject) *Iterator {
	return &Iterator{
		index: 0,
		list:  ARRAY(obj),
		base:  &base{class: NewClass("Iterator", ObjectClass)},
	}
}
