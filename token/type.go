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
	None    // none
	Catch   // catch
	Block   // block
	Object  // object
	Return  // return
	Super   // super
	Or      // or
	And     // and
	This    // this

	Int    // Int
	Float  // Float
	String // String
	Bool   // Bool

	Minus // -
	Plus  // +
	Slash // /
	Star  // *

	Dot     // .
	Colon   // :
	NewLine // \n
	Not     // !
	Arrow   // ->
	Comma
	Assign           // =
	Equal            // ==
	NotEqual         // !=
	Less             // <
	LessEqual        // <=
	Great            // >
	GreatEqual       // >=
	Ident            // Ident
	LeftParenthesis  // (
	RightParenthesis // )
	LeftBracket      // [
	RightBracket     // ]
	LeftBrace        // {
	RightBrace       // }
)
