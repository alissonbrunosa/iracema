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

func SetAttr(obj IrObject, attr IrObject, value IrObject) (err *ErrorObject) {
	var object *Object

	if object, err = checkObject(obj); err != nil {
		return err
	}

	name := unwrapString(attr)
	object.Set(string(name), value)
	return
}

func GetAttr(obj IrObject, attr IrObject) (IrObject, *ErrorObject) {
	var object *Object
	var err *ErrorObject

	if object, err = checkObject(obj); err != nil {
		return nil, err
	}

	name := unwrapString(attr)
	return object.Get(string(name)), nil
}
