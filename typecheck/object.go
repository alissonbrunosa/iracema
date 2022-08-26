package typecheck

type FieldSet map[string]Type
type MethodSet map[string]*signature

type object struct {
	parent Type

	name      string
	fieldSet  FieldSet
	methodSet MethodSet
}

func (o *object) Name() string {
	return o.name
}

func (o *object) LookupMethod(name string) *signature {
	if m, ok := o.methodSet[name]; ok {
		return m
	}

	if o.parent != nil {
		return o.parent.LookupMethod(name)
	}

	return nil
}

func (o *object) addMethod(sig *signature) *signature {
	if method := o.methodSet[sig.name]; method != nil {
		return method
	}

	o.methodSet[sig.name] = sig
	return nil
}

func (o *object) addField(name string, typ Type) Type {
	if field := o.fieldSet[name]; field != nil {
		return field
	}

	o.fieldSet[name] = typ
	return nil
}

func (o *object) Field(name string) Type {
	return o.fieldSet[name]
}

func (o *object) Is(typ Type) bool {
	if o == typ {
		return true
	}

	if o.parent != nil && o.parent.Is(typ) {
		return true
	}

	return false
}

func (o *object) defineMethodSet(sigs []*signature) {
	o.methodSet = make(MethodSet, len(sigs))

	for _, sig := range sigs {
		o.methodSet[sig.name] = sig
	}
}

func (o *object) String() string { return o.name }

func (o *object) complete() {
	if _, ok := o.methodSet["new"]; !ok {
		o.methodSet["new"] = &signature{name: "new", ret: o}
	}
}

func newObject(name string, parent Type) *object {
	return &object{
		name:      name,
		parent:    parent,
		fieldSet:  make(FieldSet),
		methodSet: make(MethodSet),
	}
}
