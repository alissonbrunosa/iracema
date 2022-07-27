package lang

import "testing"

func Test_hashSize(t *testing.T) {
	h := NewHash()

	key := NewString("int")
	value := Int(10)
	hashInsert(globalTestDummyRuntime, h, key, value)

	if INT(hashSize(globalTestDummyRuntime, h)) != 1 {
		t.Errorf("expected count to be 1, got %d", h.count)
	}
}

func Test_hashInsert(t *testing.T) {
	h := NewHash()

	key := NewString("int")
	value := Int(10)

	if hashInsert(globalTestDummyRuntime, h, key, value) == False {
		t.Fatal("not inserted")
	}
}

func Test_hashLookup(t *testing.T) {
	h := NewHash()

	key := NewString("int")
	value := Int(10)

	if hashInsert(globalTestDummyRuntime, h, key, value) == False {
		t.Fatal("not inserted")
	}

	v := hashLookup(globalTestDummyRuntime, h, key)

	if v != value {
		t.Errorf("expected value to be %+v, got %+v", value, v)
	}
}

func Test_hashKeys(t *testing.T) {
	h := NewHash()

	expectedKeys := []IrObject{
		NewString("a"),
		NewString("b"),
		NewString("c"),
	}

	for i, key := range expectedKeys {
		hashInsert(globalTestDummyRuntime, h, key, Int(i))
	}

	keys := hashKeys(globalTestDummyRuntime, h)

	for i, key := range ARRAY(keys).Elements {
		isSame := stringEqual(globalTestDummyRuntime, key, expectedKeys[i])
		if BOOL(isSame) {
			continue
		}

		t.Errorf("expected key at %d position to be %s, got %s", i, expectedKeys[i], key)
	}
}

func Test_hashValuesAt(t *testing.T) {
	h := NewHash()

	keys := []IrObject{
		NewString("a"),
		NewString("b"),
		NewString("c"),
	}

	for i, key := range keys {
		hashInsert(globalTestDummyRuntime, h, key, Int(i))
	}

	ret := hashValuesAt(globalTestDummyRuntime, h, keys[1], keys[2])

	expected := []IrObject{
		Int(1),
		Int(2),
	}

	if arrayEqual(globalTestDummyRuntime, ret, NewArray(expected)) == False {
		t.Errorf("expected keys to be %+v, got %+v", keys, ret)
	}
}
