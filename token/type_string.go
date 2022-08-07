// Code generated by "stringer -type=Type -linecomment"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Illegal-0]
	_ = x[EOF-1]
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
	_ = x[Or-20]
	_ = x[And-21]
	_ = x[This-22]
	_ = x[Int-23]
	_ = x[Float-24]
	_ = x[String-25]
	_ = x[Bool-26]
	_ = x[Minus-27]
	_ = x[Plus-28]
	_ = x[Slash-29]
	_ = x[Star-30]
	_ = x[Dot-31]
	_ = x[Colon-32]
	_ = x[NewLine-33]
	_ = x[Not-34]
	_ = x[Arrow-35]
	_ = x[Comma-36]
	_ = x[Assign-37]
	_ = x[Equal-38]
	_ = x[NotEqual-39]
	_ = x[Less-40]
	_ = x[LessEqual-41]
	_ = x[Great-42]
	_ = x[GreatEqual-43]
	_ = x[Ident-44]
	_ = x[LeftParen-45]
	_ = x[RightParen-46]
	_ = x[LeftBracket-47]
	_ = x[RightBracket-48]
	_ = x[LeftBrace-49]
	_ = x[RightBrace-50]
}

const _Type_name = "IllegalEOFifisforswitchcasedefaultinstopnextwhileelsefunnonecatchblockobjectreturnsuperorandthisIntFloatStringBool-+/*.:\\n!->Comma===!=<<=>>=Ident()[]{}"

var _Type_index = [...]uint8{0, 7, 10, 12, 14, 17, 23, 27, 34, 36, 40, 44, 49, 53, 56, 60, 65, 70, 76, 82, 87, 89, 92, 96, 99, 104, 110, 114, 115, 116, 117, 118, 119, 120, 122, 123, 125, 130, 131, 133, 135, 136, 138, 139, 141, 146, 147, 148, 149, 150, 151, 152}

func (i Type) String() string {
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
