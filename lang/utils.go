package lang

func GoString(value IrObject) string {
	return string(unwrapString(value))
}
