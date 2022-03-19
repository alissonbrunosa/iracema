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

func (c *Class) AddGoMethod(name string, body Native) {
	var arity byte
	switch body.(type) {
	case nArgs:
		arity = 255
	case zeroArgs:
		arity = 0
	case oneArg:
		arity = 1
	case twoArgs:
		arity = 2
	}

	c.methods[name] = NewGoMethod(name, body, arity)
}

func NewClass(name string, super *Class) *Class {
	return &Class{
		name:    name,
		super:   super,
		methods: make(map[string]*Method),

		base: &base{class: irClass},
	}
}
