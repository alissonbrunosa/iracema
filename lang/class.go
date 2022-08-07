package lang

func CLASS(obj IrObject) *Class {
	return obj.(*Class)
}

func classNew(rt Runtime, this IrObject, args ...IrObject) IrObject {
	c := CLASS(this)
	object := c.Alloc()
	if val := call(rt, object, "init", args...); val == nil {
		return nil
	}

	return object
}

var irClass *Class

func InitClass() {
	if irClass != nil {
		return
	}

	irClass = NewClass("Class", ObjectClass)
	irClass.AddGoMethod("new", nArgs(classNew))
}

type Allocator func(*Class) IrObject

type Class struct {
	*base

	name      string
	super     *Class
	methods   map[string]*Method
	allocator func(*Class) IrObject
}

func (c *Class) Name() string {
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

func (c *Class) LookupMethod(name string) *Method {
	for class := c; class != nil; class = class.super {
		if method, ok := class.methods[name]; ok {
			return method
		}
	}

	return nil
}

func (c *Class) AddMethod(name string, fun *Method) {
	c.methods[name] = fun
}

func (c *Class) Alloc() IrObject {
	return allocator(c)(c)
}

func (c *Class) String() string {
	return c.name
}

func (c *Class) Super() *Class {
	return c.super
}

func NewClass(name string, super *Class) *Class {
	return &Class{
		name:    name,
		super:   super,
		methods: make(map[string]*Method),

		base: &base{class: irClass},
	}
}

func allocator(class *Class) Allocator {
	for cls := class; cls != nil; cls = cls.super {
		if cls.allocator != nil {
			return cls.allocator
		}
	}

	panic("undefined method new for " + class.name)
}
