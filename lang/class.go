package lang

var irClass *Class

type Class struct {
	*base

	name    string
	super   *Class
	methods map[string]*Method
}

func (c *Class) Name() string {
	return c.name
}

func (c *Class) String() string {
	return c.name
}

func NewClass(name string, super *Class) *Class {
	return &Class{
		name:    name,
		super:   super,
		methods: make(map[string]*Method),

		base: &base{class: irClass},
	}
}
