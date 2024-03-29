package bytecode

//go:generate stringer -type=Opcode  -linecomment
type Opcode byte

const (
	Nop            Opcode = iota // NOP
	Pop                          // POP
	Push                         // PUSH
	Throw                        // THROW
	Return                       // RETURN
	PushNone                     // PUSH_NONE
	SetField                     // SET_FIELD
	GetField                     // GET_FIELD
	PushThis                     // PUSH_THIS
	SetLocal                     // SET_LOCAL
	GetLocal                     // GET_LOCAL
	MatchType                    // MATCH_TYPE
	BuildArray                   // BUILD_ARRAY
	BuildHash                    // BUILD_HASH
	CallMethod                   // CALL_METHOD
	CallSuper                    // CALL_SUPER
	SetConstant                  // SET_CONSTANT
	GetConstant                  // GET_CONSTANT
	DefineObject                 // DEFINE_OBJECT
	DefineField                  // DEFINE_FIELD
	DefineFunction               // DEFINE_FUNCTION
	Jump                         // JUMP
	JumpIfFalse                  // JUMP_IF_FALSE
	JumpIfTrue                   // JUMP_IF_TRUE
	Iterate                      // ITERATE
	NewIterator                  // NEWITERATOR
	LoadFile                     // LOAD_FILE

	/*
		┌──────────────────────── INTERNAL OPCODES ────────────────────────────┐
		│                                                                      │
		│  WithCatch is a hack to make the markReachable function to be able   │
		│  of marking the catches body after the return instr in the function  │
		│  block. It will get replaced by a Nop instr.                         │
		└──────────────────────────────────────────────────────────────────────┘
	*/
	WithCatch // WITH_CATCH
)
