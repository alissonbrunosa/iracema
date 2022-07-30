package bytecode

//go:generate stringer -type=Opcode  -linecomment
type Opcode byte

const (
	_              Opcode = iota
	Pop                   // POP
	Push                  // PUSH
	Throw                 // THROW
	Return                // RETURN
	PushNone              // PUSH_NONE
	SetAttr               // SET_ATTR
	GetAttr               // GET_ATTR
	PushSelf              // PUSH_SELF
	SetLocal              // SET_LOCAL
	GetLocal              // GET_LOCAL
	MatchType             // MATCH_TYPE
	BuildArray            // BUILD_ARRAY
	CallMethod            // CALL_METHOD
	CallSuper             // CALL_SUPER
	SetConstant           // SET_CONSTANT
	GetConstant           // GET_CONSTANT
	DefineObject          // DEFINE_OBJECT
	DefineFunction        // DEFINE_FUNCTION
	Jump                  // JUMP
	JumpIfFalse           // JUMP_IF_FALSE
	JumpIfTrue            // JUMP_IF_TRUE
	Iterate               // ITERATE
	NewIterator           // NEWITERATOR
)
