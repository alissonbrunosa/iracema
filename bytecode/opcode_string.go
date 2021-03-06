// Code generated by "stringer -type=Opcode -linecomment"; DO NOT EDIT.

package bytecode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Pop-1]
	_ = x[Push-2]
	_ = x[Throw-3]
	_ = x[Return-4]
	_ = x[PushNone-5]
	_ = x[SetAttr-6]
	_ = x[GetAttr-7]
	_ = x[PushSelf-8]
	_ = x[SetLocal-9]
	_ = x[GetLocal-10]
	_ = x[MatchType-11]
	_ = x[BuildArray-12]
	_ = x[CallMethod-13]
	_ = x[CallSuper-14]
	_ = x[SetConstant-15]
	_ = x[GetConstant-16]
	_ = x[DefineObject-17]
	_ = x[DefineFunction-18]
	_ = x[Jump-19]
	_ = x[JumpIfFalse-20]
	_ = x[JumpIfTrue-21]
	_ = x[Iterate-22]
	_ = x[NewIterator-23]
}

const _Opcode_name = "POPPUSHTHROWRETURNPUSH_NONESET_ATTRGET_ATTRPUSH_SELFSET_LOCALGET_LOCALMATCH_TYPEBUILD_ARRAYCALL_METHODCALL_SUPERSET_CONSTANTGET_CONSTANTDEFINE_OBJECTDEFINE_FUNCTIONJUMPJUMP_IF_FALSEJUMP_IF_TRUEITERATENEWITERATOR"

var _Opcode_index = [...]uint8{0, 3, 7, 12, 18, 27, 35, 43, 52, 61, 70, 80, 91, 102, 112, 124, 136, 149, 164, 168, 181, 193, 200, 211}

func (i Opcode) String() string {
	i -= 1
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
