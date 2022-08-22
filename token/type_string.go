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
	_ = x[Use-24]
	_ = x[Var-25]
	_ = x[Int-26]
	_ = x[Float-27]
	_ = x[String-28]
	_ = x[Bool-29]
	_ = x[Minus-30]
	_ = x[Plus-31]
	_ = x[Slash-32]
	_ = x[Star-33]
	_ = x[Dot-34]
	_ = x[Colon-35]
	_ = x[NewLine-36]
	_ = x[Not-37]
	_ = x[Arrow-38]
	_ = x[Comma-39]
	_ = x[Assign-40]
	_ = x[Equal-41]
	_ = x[NotEqual-42]
	_ = x[Less-43]
	_ = x[LessEqual-44]
	_ = x[Great-45]
	_ = x[GreatEqual-46]
	_ = x[Ident-47]
	_ = x[LeftParen-48]
	_ = x[RightParen-49]
	_ = x[LeftBracket-50]
	_ = x[RightBracket-51]
	_ = x[LeftBrace-52]
	_ = x[RightBrace-53]
}

const _Type_name = "IllegalEOFifisforswitchcasedefaultinstopnextwhileelsefunnonecatchblockobjectreturnsuperorandthisusevarIntFloatStringBool-+/*.:\\n!->Comma===!=<<=>>=Ident()[]{}"

var _Type_index = [...]uint8{0, 7, 10, 12, 14, 17, 23, 27, 34, 36, 40, 44, 49, 53, 56, 60, 65, 70, 76, 82, 87, 89, 92, 96, 99, 102, 105, 110, 116, 120, 121, 122, 123, 124, 125, 126, 128, 129, 131, 136, 137, 139, 141, 142, 144, 145, 147, 152, 153, 154, 155, 156, 157, 158}

func (i Type) String() string {
	i -= 1
	if i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
