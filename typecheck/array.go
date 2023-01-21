package typecheck

type array struct {
	*object

	value Type
}

func (a *array) Is(t Type) bool {
	o, ok := t.(*array)

	return ok && a.value == o.value
}

func arrayFuns(value Type) map[string]*signature {
	return map[string]*signature{
		"insert":  &signature{name: "insert", params: []Type{INT, value}, ret: BOOL},
		"get":     &signature{name: "get", params: []Type{INT}, ret: value},
		"at":      &signature{name: "at", params: []Type{INT}, ret: value},
		"size":    &signature{name: "size", params: nil, ret: INT},
		"to_str":  &signature{name: "to_str", params: nil, ret: STRING},
		"inspect": &signature{name: "inspect", params: nil, ret: STRING},
	}
}

func newArray(value Type) Type {
	return &array{
		value: value,

		object: &object{
			name:      "Array",
			parent:    OBJECT,
			methodSet: arrayFuns(value),
		},
	}
}
