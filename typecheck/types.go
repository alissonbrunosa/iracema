package typecheck

var (
	INT     Type
	FLOAT   Type
	BOOL    Type
	NONE    Type
	STRING  Type
	INVALID Type
)

var LIT_TYPES map[string]Type

func define(name string) Type {
	obj := &object{name: name}
	return obj
}

func init() {
	INT = define("Int")
	FLOAT = define("Float")
	BOOL = define("Bool")
	NONE = define("None")
	STRING = define("String")
	INVALID = define("INVALID")

	LIT_TYPES = map[string]Type{
		"Int":    INT,
		"Float":  FLOAT,
		"Bool":   BOOL,
		"String": STRING,
		"None":   NONE,
	}
}
