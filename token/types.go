package token

type Type byte

//go:generate stringer -type=Type -linecomment
const (
	Illegal Type = iota // Illegal
	Eof                 // EOF

	If      // if
	Is      // is
	For     // for
	Switch  // switch
	Case    // case
	Default // default
	In      // in
	Stop    // stop
	Next    // next
	While   // while
	Else    // else
	Fun     // fun
	Nil     // nil
	Catch   // catch
	Block   // block
	Object  // object
	Return  // return

	Int    // Int
	Float  // Float
	String // String
	Bool   // Bool

	Minus // -
	Plus  // +
	Slash // /
	Star  // *

	Dot   // .
	Colon // :
	Not   // !
	Arrow // ->
	Comma
	Assign             // =
	Equal              // ==
	LessThan           // <
	LessOrEqualThan    // <=
	GreaterThan        // >
	GreaterOrEqualThan // >=
	Ident              // Ident
	LeftParenthesis    // (
	RightParenthesis   // )
	LeftBracket        // [
	RightBracket       // ]
	LeftBrace          // {
	RightBrace         // }
)
