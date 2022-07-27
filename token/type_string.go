// Code generated by "stringer -type=Type -linecomment"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Illegal-0]
	_ = x[Eof-1]
	_ = x[If-2]
	_ = x[Is-3]
	_ = x[For-4]
	_ = x[Switch-5]
	_ = x[Case-6]
	_ = x[Default-7]
	_ = x[In-8]
	_ = x[Stop-9]
	_ = x[Next-10]
	_ = x[While-11]
	_ = x[Else-12]
	_ = x[Fun-13]
	_ = x[None-14]
	_ = x[Catch-15]
	_ = x[Block-16]
	_ = x[Object-17]
	_ = x[Return-18]
	_ = x[Super-19]
	_ = x[Int-20]
	_ = x[Float-21]
	_ = x[String-22]
	_ = x[Bool-23]
	_ = x[Minus-24]
	_ = x[Plus-25]
	_ = x[Slash-26]
	_ = x[Star-27]
	_ = x[Dot-28]
	_ = x[Colon-29]
	_ = x[NewLine-30]
	_ = x[Not-31]
	_ = x[Arrow-32]
	_ = x[Comma-33]
	_ = x[Assign-34]
	_ = x[Equal-35]
	_ = x[NotEqual-36]
	_ = x[Less-37]
	_ = x[LessEqual-38]
	_ = x[Great-39]
	_ = x[GreatEqual-40]
	_ = x[Ident-41]
	_ = x[LeftParenthesis-42]
	_ = x[RightParenthesis-43]
	_ = x[LeftBracket-44]
	_ = x[RightBracket-45]
	_ = x[LeftBrace-46]
	_ = x[RightBrace-47]
}

const _Type_name = "IllegalEOFifisforswitchcasedefaultinstopnextwhileelsefunnonecatchblockobjectreturnsuperIntFloatStringBool-+/*.:\\n!->Comma===!=<<=>>=Ident()[]{}"

var _Type_index = [...]uint8{0, 7, 10, 12, 14, 17, 23, 27, 34, 36, 40, 44, 49, 53, 56, 60, 65, 70, 76, 82, 87, 90, 95, 101, 105, 106, 107, 108, 109, 110, 111, 113, 114, 116, 121, 122, 124, 126, 127, 129, 130, 132, 137, 138, 139, 140, 141, 142, 143}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
