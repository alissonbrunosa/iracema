package typecheck

type object struct {
	name string
}

func (o *object) Name() string { return o.name }

func (o *object) Is(t Type) bool {
	return o == t
}
