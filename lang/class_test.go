package lang

import (
	"testing"
)

func Test_Name(t *testing.T) {
	class := NewClass("Dummy", nil)

	if class.Name() != "Dummy" {
		t.Errorf("expected class name to be Dummy, got %s", class.Name())
	}
}

func Test_String(t *testing.T) {
	class := NewClass("Dummy", nil)

	if class.String() != class.Name() {
		t.Errorf("expected class.String() to be %s, got %s", class.Name(), class.String())
	}
}

func Test_AddGoMethod(t *testing.T) {
	class := NewClass("Object", nil)

	tests := []struct {
		scenario      string
		expectedName  string
		expectedArity byte
		method        Native
	}{
		{
			scenario:      "zeroArgs",
			expectedName:  "zero",
			expectedArity: 0,
			method: zeroArgs(func(self IrObject) IrObject {
				return nil
			}),
		},
		{
			scenario:      "oneArg",
			expectedName:  "one",
			expectedArity: 1,
			method: oneArg(func(self IrObject, value IrObject) IrObject {
				return nil
			}),
		},
		{
			scenario:      "twoArgs",
			expectedName:  "two",
			expectedArity: 2,
			method: twoArgs(func(self IrObject, value1 IrObject, value2 IrObject) IrObject {
				return nil
			}),
		},
		{
			scenario:      "nArgs",
			expectedName:  "many",
			expectedArity: 255,
			method: nArgs(func(self IrObject, values ...IrObject) IrObject {
				return nil
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			class.AddGoMethod(test.expectedName, test.method)

			method := class.methods[test.expectedName]
			if method.Name() != test.expectedName {
				t.Errorf("expected mathod name to be '%s', got '%s'", test.expectedName, method.Name())
			}

			if method.Arity() != test.expectedArity {
				t.Errorf("expected mathod arity to be %d, got '%d'", test.expectedArity, method.Arity())
			}

			if _, ok := method.Body().(Native); !ok {
				t.Errorf("expected body to be Native type, got %T", method.Body())
			}
		})
	}
}

func Test_classNew_CallsInitMethod(t *testing.T) {
	class := NewClass("Dummy", ObjectClass)

	var called = false
	init := zeroArgs(func(rt Runtime, self IrObject) IrObject {
		called = true
		return nil
	})

	class.AddGoMethod("init", init)

	classNew(runtime, class)
	if !called {
		t.Error("expected init to be called")
	}
}
func Test_classNew_ReturnObjectFromTargetClass(t *testing.T) {
	class := NewClass("Dummy", ObjectClass)

	object := classNew(runtime, class)

	if object.Class() != class {
		t.Errorf("expected class to be %s, got %s", class.name, object.Class().name)
	}
}

func Test_Alloc_PanicsIfNotDefiend(t *testing.T) {
	class := NewClass("Dummy", nil)

	defer func() {
		err := recover()
		expectedReason := "undefined method new for Dummy"
		if err != expectedReason {
			t.Errorf("expected panic reason to be %q, got %q", expectedReason, err)
		}
	}()

	class.Alloc()
	t.Error("expected function to panic")
}

func Test_Alloc(t *testing.T) {
	called := false

	class := NewClass("Dummy", nil)
	class.allocator = func(class *Class) IrObject {
		called = true
		return nil
	}

	class.Alloc()

	if !called {
		t.Error("expected class alloc function to be called")
	}
}

func Test_Alloc_WhenSuperClassHasItDefined(t *testing.T) {
	called := false

	super := NewClass("SuperDummy", nil)
	super.allocator = func(class *Class) IrObject {
		called = true
		return nil
	}

	class := NewClass("Dummy", super)
	class.Alloc()

	if !called {
		t.Error("expected super's alloc function to be called")
	}
}
