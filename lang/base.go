package lang

type base struct {
	class *Class
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
