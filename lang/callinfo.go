package lang

type CallInfo struct {
	*base

	name string
	argc byte
}

func (c *CallInfo) Name() string { return c.name }
func (c *CallInfo) Argc() byte   { return c.argc }

func NewCallInfo(name string, argc byte) *CallInfo {
	return &CallInfo{name: name, argc: argc}
}
