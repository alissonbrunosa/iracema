// Code generated by "stringer -type=Opcode -linecomment"; DO NOT EDIT.

package bytecode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Nop-0]
	_ = x[Pop-1]
	_ = x[Push-2]
	_ = x[Throw-3]
	_ = x[Return-4]
	_ = x[PushNone-5]
	_ = x[SetField-6]
	_ = x[GetField-7]
	_ = x[PushThis-8]
	_ = x[SetLocal-9]
	_ = x[GetLocal-10]
	_ = x[MatchType-11]
	_ = x[BuildArray-12]
	_ = x[BuildHash-13]
	_ = x[CallMethod-14]
	_ = x[CallSuper-15]
	_ = x[SetConstant-16]
	_ = x[GetConstant-17]
	_ = x[DefineObject-18]
	_ = x[DefineField-19]
	_ = x[DefineFunction-20]
	_ = x[Jump-21]
	_ = x[JumpIfFalse-22]
	_ = x[JumpIfTrue-23]
	_ = x[Iterate-24]
	_ = x[NewIterator-25]
	_ = x[LoadFile-26]
	_ = x[WithCatch-27]
}

const _Opcode_name = "NOPPOPPUSHTHROWRETURNPUSH_NONESET_FIELDGET_FIELDPUSH_THISSET_LOCALGET_LOCALMATCH_TYPEBUILD_ARRAYBUILD_HASHCALL_METHODCALL_SUPERSET_CONSTANTGET_CONSTANTDEFINE_OBJECTDEFINE_FIELDDEFINE_FUNCTIONJUMPJUMP_IF_FALSEJUMP_IF_TRUEITERATENEWITERATORLOAD_FILEWITH_CATCH"

var _Opcode_index = [...]uint16{0, 3, 6, 10, 15, 21, 30, 39, 48, 57, 66, 75, 85, 96, 106, 117, 127, 139, 151, 164, 176, 191, 195, 208, 220, 227, 238, 247, 257}

func (i Opcode) String() string {
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
