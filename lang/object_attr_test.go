package lang

import "testing"

func setupObject() IrObject {
	class := &Class{
		fields: map[string]byte{
			"value": 0,
		},
	}

	return ObjectClass.allocator(class)
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

	given := Int(10)
	SetAttr(object, NewString("value"), Int(10))
	got, _ := GetAttr(object, NewString("value"))
	if given != got {
		t.Errorf("expected field value to be set to %d, got %d", given, got)
	}
}
