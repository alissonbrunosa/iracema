package lang

type base struct {
	class *Class
}

func (b *base) LookupMethod(name string) *Method {
	for cls := b.class; cls != nil; cls = cls.super {
		if method, ok := cls.methods[name]; ok {
			return method
		}
	}

	return nil
}

func (b *base) Is(class *Class) bool {
	for cls := b.class; cls != nil; cls = cls.super {
		if cls == class {
			return true
		}
	}

	return false
}

func (b *base) Class() *Class {
	return b.class
}
