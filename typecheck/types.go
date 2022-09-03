package typecheck

type Type interface {
	Name() string
	Parent() Type
	Is(Type) bool
	LookupMethod(string) *signature
	Field(string) Type
}

var (
	INT     *object
	FLOAT   *object
	BOOL    *object
	NONE    *object
	STRING  *object
	OBJECT  *object
	SCRIPT  *object
	INVALID *object
)

var LIT_TYPES map[string]Type

func init() {
	OBJECT = &object{name: "Object", parent: nil}
	INT = &object{name: "Int", parent: OBJECT}
	FLOAT = &object{name: "Float", parent: OBJECT}
	BOOL = &object{name: "Bool", parent: OBJECT}
	NONE = &object{name: "None", parent: OBJECT}
	STRING = &object{name: "String", parent: OBJECT}
	INVALID = &object{name: "INVALID", parent: nil}

	OBJECT.defineMethodSet(
		[]*signature{
			{name: "init", params: nil, ret: NONE},
			{name: "==", params: []Type{OBJECT}, ret: BOOL},
			{name: "!=", params: []Type{OBJECT}, ret: BOOL},
			{name: "hash", params: nil, ret: INT},
			{name: "puts", params: []Type{OBJECT}, ret: NONE},
			{name: "object_id", params: nil, ret: INT},
			{name: "inspect", params: nil, ret: STRING},
			{name: "to_str", params: nil, ret: STRING},
			{name: "nil?", params: nil, ret: BOOL},
			{name: "unot", params: nil, ret: OBJECT},
		},
	)

	SCRIPT = &object{
		name:      "Script",
		parent:    OBJECT,
		methodSet: make(map[string]*signature),
	}

	STRING.defineMethodSet(
		[]*signature{
			{name: "==", params: []Type{STRING}, ret: BOOL},
		},
	)

	LIT_TYPES = map[string]Type{
		"Int":    INT,
		"Float":  FLOAT,
		"Bool":   BOOL,
		"String": STRING,
		"Object": OBJECT,
		"None":   NONE,
	}
}
