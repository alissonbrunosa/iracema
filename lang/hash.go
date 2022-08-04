package lang

import "strings"

func HASH(obj IrObject) *Hash {
	return obj.(*Hash)
}

type entry struct {
	hashCode Int
	key      IrObject
	value    IrObject
	next     *entry
}

func retrieveHashCode(rt Runtime, obj IrObject) Int {
	hash := call(rt, obj, "hash")
	if hash == nil {
		return 0
	}

	return INT(hash)
}

func hashLookup(rt Runtime, this IrObject, key IrObject) IrObject {
	h := HASH(this)

	hashCode := retrieveHashCode(rt, key)
	length := Int(len(h.table))
	index := (hashCode & 0x7FFFFFFF) % length

	for entry := h.table[index]; entry != nil; entry = entry.next {
		equal := call(rt, entry.key, "==", key)
		if equal == nil {
			return nil
		}

		if entry.hashCode == hashCode && BOOL(equal) {
			return entry.value
		}
	}
	return nil
}

func hashInsert(rt Runtime, this IrObject, key IrObject, value IrObject) IrObject {
	h := HASH(this)

	hashCode := retrieveHashCode(rt, key)
	length := Int(len(h.table))
	index := (hashCode & 0x7FFFFFFF) % length

	for entry := h.table[index]; entry != nil; entry = entry.next {
		equal := call(rt, entry.key, "==", key)
		if equal == nil {
			return nil
		}

		if entry.hashCode == hashCode && BOOL(equal) {
			entry.value = value
			return True
		}
	}

	h.addEntry(rt, hashCode, key, value, index)
	return True
}

func hashKeys(rt Runtime, this IrObject) IrObject {
	h := HASH(this)

	var keys []IrObject
	for _, entry := range h.table {
		if entry == nil {
			continue
		}

		for ; entry != nil; entry = entry.next {
			keys = append(keys, entry.key)
		}
	}

	return NewArray(keys)
}

func hashValues(rt Runtime, this IrObject) IrObject {
	h := HASH(this)

	var values []IrObject
	for _, entry := range h.table {
		if entry == nil {
			continue
		}

		for ; entry != nil; entry = entry.next {
			values = append(values, entry.value)
		}
	}

	return NewArray(values)
}

func hashInspect(rt Runtime, this IrObject) IrObject {
	h := HASH(this)

	if h.count == 0 {
		return NewString("{}")
	}

	var out strings.Builder
	out.WriteByte('{')
	for _, entry := range h.table {
		if entry == nil {
			continue
		}

		for ; entry != nil; entry = entry.next {
			if out.Len() > 2 { // 2 is the two chars { and whitespace
				out.WriteString(", ")
			}
			val := call(rt, entry.key, "inspect")
			if val == nil {
				return nil
			}

			out.Write(unwrapString(val))
			out.WriteString(": ")

			val = call(rt, entry.value, "inspect")
			if val == nil {
				return nil
			}

			out.Write(unwrapString(val))
		}
	}

	out.WriteByte('}')

	return NewString(out.String())
}

func hashSize(rt Runtime, this IrObject) IrObject {
	return HASH(this).count
}

func hashValuesAt(rt Runtime, this IrObject, keys ...IrObject) IrObject {
	result := make([]IrObject, len(keys))
	if len(keys) == 0 {
		return NewArray(result)
	}

	h := HASH(this)
	for i, key := range keys {
		result[i] = hashLookup(rt, h, key)
	}

	return NewArray(result)
}

func hashHasKey(rt Runtime, this, key IrObject) IrObject {
	h := HASH(this)

	hashCode := retrieveHashCode(rt, key)
	length := Int(len(h.table))
	index := (hashCode & 0x7FFFFFFF) % length

	for entry := h.table[index]; entry != nil; entry = entry.next {
		equal := call(rt, entry.key, "==", key)
		if equal == nil {
			return nil
		}

		if entry.hashCode == hashCode && BOOL(equal) {
			return True
		}
	}

	return False
}

var HashClass *Class

func InitHash() {
	if HashClass != nil {
		return
	}

	HashClass = NewClass("Hash", ObjectClass)

	HashClass.AddGoMethod("put", twoArgs(hashInsert))
	HashClass.AddGoMethod("insert", twoArgs(hashInsert))
	HashClass.AddGoMethod("get", oneArg(hashLookup))
	HashClass.AddGoMethod("key?", oneArg(hashHasKey))
	HashClass.AddGoMethod("keys", zeroArgs(hashKeys))
	HashClass.AddGoMethod("keys", zeroArgs(hashValues))
	HashClass.AddGoMethod("size", zeroArgs(hashSize))
	HashClass.AddGoMethod("values_at", nArgs(hashValuesAt))
	HashClass.AddGoMethod("inspect", zeroArgs(hashInspect))
}

func NewHash() *Hash {
	return &Hash{
		table:      make([]*entry, 1<<4),
		loadFactor: 0.75,

		base: &base{class: HashClass},
	}
}

type Hash struct {
	*base

	table      []*entry
	threshold  Int
	count      Int
	loadFactor float32
}

func (h *Hash) addEntry(rt Runtime, hashCode Int, key IrObject, value IrObject, index Int) {
	if h.count >= h.threshold {
		h.rehash()
		hashCode := retrieveHashCode(rt, key)
		index = (hashCode & 0x7FFFFFFF) % Int(len(h.table))
	}

	e := h.table[index]
	h.table[index] = &entry{
		hashCode: hashCode,
		key:      key,
		value:    value,
		next:     e,
	}
	h.count++
}

func (h *Hash) BulkInsert(rt Runtime, elements []IrObject) {
	for i := 0; i < len(elements); i += 2 {
		if res := hashInsert(rt, h, elements[i], elements[i+1]); res == nil {
			break
		}
	}
}

func (h *Hash) rehash() {
	oldTable := h.table
	oldCapacity := Int(len(oldTable))
	newCapacity := (oldCapacity << 1)
	h.threshold = Int(float32(newCapacity) * h.loadFactor)
	h.table = make([]*entry, newCapacity)

	for i := oldCapacity - 1; i >= 0; i-- {
		for old := oldTable[i]; old != nil; {
			entry := old
			old = old.next
			index := (entry.hashCode & 0x7FFFFFFF) % newCapacity
			entry.next = h.table[index]
			h.table[index] = entry
		}
	}
}
