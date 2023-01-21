package typecheck

// type hash struct {
// 	*object
//
// 	key   Type
// 	value Type
// }
//
// func (h *hash) isUntyped() bool {
// 	return h.key == nil && h.value == nil
// }
//
// func (h *hash) setType(key, value Type) {
// 	h.key = key
// 	h.value = value
// 	h.methodSet = hashFuns(key, value)
// }
//
// func untypedHash() *hash {
// 	return &hash{
// 		object: &object{name: "Hash"},
// 	}
// }
//
// func hashFuns(key Type, value Type) map[string]*signature {
// 	return map[string]*signature{
// 		"insert":  &signature{name: "insert", params: []Type{key, value}, ret: BOOL},
// 		"put":     &signature{name: "put", params: []Type{key, value}, ret: BOOL},
// 		"get":     &signature{name: "get", params: []Type{key}, ret: value},
// 		"key?":    &signature{name: "key?", params: []Type{key}, ret: BOOL},
// 		"size":    &signature{name: "size", params: nil, ret: INT},
// 		"to_str":  &signature{name: "to_str", params: nil, ret: STRING},
// 		"inspect": &signature{name: "inspect", params: nil, ret: STRING},
// 	}
// }
//
// func newHash(key Type, value Type) Type {
// 	return &hash{
// 		key:   key,
// 		value: value,
//
// 		object: &object{
// 			name:      "Hash",
// 			methodSet: hashFuns(key, value),
// 		},
// 	}
// }
