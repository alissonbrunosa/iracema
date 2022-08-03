package lang

var ScriptClass *Class

func scriptInspect(_rt Runtime, _this IrObject) IrObject {
	return NewString("script")
}

func InitScript() {
	if ScriptClass != nil {
		return
	}

	ScriptClass = NewClass("Script", ObjectClass)
	ScriptClass.AddGoMethod("inspect", zeroArgs(scriptInspect))
}
