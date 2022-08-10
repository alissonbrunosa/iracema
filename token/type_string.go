// Code generated by "stringer -type=Type -linecomment"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Illegal-1]
	_ = x[EOF-2]
	_ = x[If-3]
	_ = x[Is-4]
	_ = x[For-5]
	_ = x[Switch-6]
	_ = x[Case-7]
	_ = x[Default-8]
	_ = x[In-9]
	_ = x[Stop-10]
	_ = x[Next-11]
	_ = x[While-12]
	_ = x[Else-13]
	_ = x[Fun-14]
	_ = x[None-15]
	_ = x[Catch-16]
	_ = x[Block-17]
	_ = x[Object-18]
	_ = x[Return-19]
	_ = x[Super-20]
	_ = x[Or-21]
	_ = x[And-22]
	_ = x[This-23]
	_ = x[Int-24]
	_ = x[Float-25]
	_ = x[String-26]
	_ = x[Bool-27]
	_ = x[Minus-28]
	_ = x[Plus-29]
	_ = x[Slash-30]
	_ = x[Star-31]
	_ = x[Dot-32]
	_ = x[Colon-33]
	_ = x[NewLine-34]
	_ = x[Not-35]
	_ = x[Arrow-36]
	_ = x[Comma-37]
	_ = x[Assign-38]
	_ = x[Equal-39]
	_ = x[NotEqual-40]
	_ = x[Less-41]
	_ = x[LessEqual-42]
	_ = x[Great-43]
	_ = x[GreatEqual-44]
	_ = x[Ident-45]
	_ = x[LeftParen-46]
	_ = x[RightParen-47]
	_ = x[LeftBracket-48]
	_ = x[RightBracket-49]
	_ = x[LeftBrace-50]
	_ = x[RightBrace-51]
}

const _Type_name = "IllegalEOFifisforswitchcasedefaultinstopnextwhileelsefunnonecatchblockobjectreturnsuperorandthisIntFloatStringBool-+/*.:\\n!->Comma===!=<<=>>=Ident()[]{}"

var _Type_index = [...]uint8{0, 7, 10, 12, 14, 17, 23, 27, 34, 36, 40, 44, 49, 53, 56, 60, 65, 70, 76, 82, 87, 89, 92, 96, 99, 104, 110, 114, 115, 116, 117, 118, 119, 120, 122, 123, 125, 130, 131, 133, 135, 136, 138, 139, 141, 146, 147, 148, 149, 150, 151, 152}

func (i Type) String() string {
	i -= 1
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
