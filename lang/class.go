package lang

type Class struct {
	super   *Class
	name    string
	methods map[string]*Method
}

func (c *Class) Name() string {
	return c.name
}

func (c *Class) String() string {
	return c.name
}
