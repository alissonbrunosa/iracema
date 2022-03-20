package lang

import "testing"

func setupObject() *Object {
	object := new(Object)
	object.attrs = make(map[string]IrObject, 1)
	return object
}

func TestObject_SetAttrReturnsNoError(t *testing.T) {
	object := setupObject()

	err := SetAttr(object, NewString("value"), Int(10))
	if err != nil {
		t.Errorf("expected to not return an error: %s", err)
	}
}

func TestObject_SetAttrSetsCorrectValue(t *testing.T) {
	object := setupObject()

	SetAttr(object, NewString("value"), Int(10))
	value := object.attrs["value"].(Int)
	if value != 10 {
		t.Errorf("expected attribu value to be set to %d, got %d", 10, object.attrs["value"])
	}
}
