package lang

import "testing"

func Test_Is(t *testing.T) {
	class := NewClass("ClassOne", nil)

	b := &base{class: class}

	if !b.Is(class) {
		t.Errorf("expected to be instance of %s", class)
	}
}

func Test_Is_whenTestAgainstSuperClass(t *testing.T) {
	super := NewClass("Super", nil)
	class := NewClass("ClassOne", super)

	b := &base{class: class}

	if !b.Is(super) {
		t.Errorf("expected to be instance of %s", class)
	}
}

func Test_Is_whenNotInstanceClass(t *testing.T) {
	anotherClass := NewClass("ClassOne", nil)
	class := NewClass("ClassTow", nil)

	b := &base{class: class}

	if b.Is(anotherClass) {
		t.Errorf("expected not to be instance of %s", anotherClass)
	}
}

func Test_LookupMethod(t *testing.T) {
	class := NewClass("ClassOne", nil)
	class.AddGoMethod("method", zeroArgs(func(rt Runtime, recv IrObject) IrObject {
		return nil
	}))

	b := &base{class: class}

	method := b.LookupMethod("method")
	if method == nil {
		t.Error("expected not to be nil")
	}
}

func Test_LookupMethod_WhenNotDefined(t *testing.T) {
	class := NewClass("ClassOne", nil)

	b := &base{class: class}

	method := b.LookupMethod("method")
	if method != nil {
		t.Error("expected to be nil")
	}
}

func Test_LookupMethod_WhenDefinedInSuper(t *testing.T) {
	super := NewClass("Super", nil)
	super.AddGoMethod("method", zeroArgs(func(rt Runtime, recv IrObject) IrObject {
		return nil
	}))
	class := NewClass("ClassOne", super)

	b := &base{class: class}

	method := b.LookupMethod("method")
	if method == nil {
		t.Error("expected not to be nil")
	}
}
