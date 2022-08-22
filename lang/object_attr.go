package lang

import "fmt"

func checkObject(obj IrObject) (*Object, *ErrorObject) {
	o, ok := obj.(*Object)

	if !ok {
		err := fmt.Sprintf("can't get attribute for instance of %s", obj.Class())
		return nil, NewError(err, RuntimeError)
	}

	return o, nil
}

func SetAttr(obj IrObject, attr IrObject, value IrObject) *ErrorObject {
	object, err := checkObject(obj)
	if err != nil {
		return err
	}

	class := obj.Class()
	name := GoString(attr)
	pos, ok := class.fields[name]
	if !ok {
		return NewError("'%s' object has no field '%s'", RuntimeError, class, name)
	}

	object.Set(pos, value)
	return nil
}

func GetAttr(obj IrObject, attr IrObject) (IrObject, *ErrorObject) {
	var object *Object
	var err *ErrorObject

	if object, err = checkObject(obj); err != nil {
		return nil, err
	}

	class := obj.Class()
	name := GoString(attr)
	pos, ok := class.fields[name]
	if !ok {
		return nil, NewError("'%s' object has no field '%s'", RuntimeError, class, name)
	}

	return object.Get(pos), nil
}
