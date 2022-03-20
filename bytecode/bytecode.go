package bytecode

//go:generate stringer -type=Opcode -output=opcode_string.go -linecomment
type Opcode byte

const (
	Not            Opcode = iota // not
	Pop                          // pop
	Push                         // push
	Throw                        // throw
	Unary                        // unary
	Binary                       // binary
	Return                       // return
	Compare                      // compare
	PushNone                     // pushnone
	SetAttr                      // setattr
	GetAttr                      // getattr
	PushSelf                     // pushself
	SetLocal                     // setlocal
	GetLocal                     // getlocal
	MatchType                    // matchtype
	BuildArray                   // buildarray
	CallMethod                   // callmethod
	SetConstant                  // setconstant
	GetConstant                  // getconstant
	DefineObject                 // defineobject
	DefineFunction               // definefunction
	Jump                         // jump
	JumpIfFalse                  // jumpiffalse
)
