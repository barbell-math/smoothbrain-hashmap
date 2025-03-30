package sbmap

import (
	"hash/maphash"
)

type (
	entry[K any, V any] struct {
		key    K
		value  V
		filled bool
	}

	Map[K any, V any] struct {
		data []entry[K, V]
		len  int
		eq   func(l K, r K) bool
		hash func(l K) uint64
	}
)

var (
	_comparableSeed = maphash.MakeSeed()
	// The default initial capacity that will be used if no capacity or zero
	// capacity is supplied
	_defaultInitialCap = 8
	// A value between 0 and 100 that determines how full the map can get before
	// the hash map doubles the underlying slice.
	_growFactor = 75
	// A value between 0 and 100 that determines how empty the map can get
	// before the hash map halves the underlying slice.
	_shrinkFactor = 50
	// Note that growth factor must always be a power of two because of how the
	// hash clamping works!!
	_sliceGrowthFactor = 2
)

// An equality function that can be passed to [NewCustom] when using a
// comparable type. If the key type is comparable then you can simply use [New]
// instead of [NewCustom] and this function will be used by default.
func ComparableEqual[T comparable](l T, r T) bool {
	return l == r
}

// A hash function that can be passed to [NewCustom] when using a comparable
// type. If the key type is comparable then you can simply use [New] instead of
// [NewCustom] and this function will be used by default.
func ComparableHash[T comparable](v T) uint64 {
	return maphash.Comparable(_comparableSeed, v)
}

// Creates a Map where K is the key type and V is the value type.
// ComparableEqual and ComparableHash funcitons will be used by the returned
// Map. For creating a Map with non-comparable types or simply custom
// hash and equality functions refer to [NewCustom].
func New[K comparable, V comparable]() Map[K, V] {
	return Map[K, V]{
		data: make([]entry[K, V], _defaultInitialCap, _defaultInitialCap),
		len:  0,
		eq:   ComparableEqual[K],
		hash: ComparableHash[K],
	}
}

// TODO
// func NewCap[K any, V any](cap int) Map[K, V] {
//
// }

// TODO
// func NewCustom

// Returns the number of elements in the hash map. This is different than the
// maps capacity.
func (h *Map[K, V]) Len() int {
	return h.len
}

func (h *Map[K, V]) clampedHash(hash uint64) uint64 {
	return hash & (uint64(cap(h.data)) - 1)
}

// Gets the value that is related to the supplied key. If the key is found the
// boolean return value will be true and the value will be returned. If the key
// is not found the boolean return value will be false and a zero-initilized
// value of type V will be returned.
func (h *Map[K, V]) Get(k K) (V, bool) {
	hash := h.clampedHash(h.hash(k))
	for h.data[hash].filled {
		if h.eq(h.data[hash].key, k) {
			return h.data[hash].value, true
		}
		hash = h.clampedHash(hash + 1)
	}

	var tmp V
	return tmp, false
}

// Places the supplied key, value pair in the hash map. If the key was already
// present in the map the old value will be overwritten. The map will resize as
// necessary.
func (h *Map[K, V]) Put(k K, v V) {
	// Original equation:
	// 	len/cap *100 >= _growFactor
	// Except dividing ints is bad, we want more precision. So remove the
	// division and we get this:
	if h.len*100 >= _growFactor*cap(h.data) {
		h.resize(cap(h.data) * _sliceGrowthFactor)
	}

	hash := h.clampedHash(h.hash(k))
	for h.data[hash].filled {
		if h.eq(h.data[hash].key, k) {
			h.data[hash].value = v
		}
		hash = h.clampedHash(hash + 1)
	}

	h.data[hash].key = k
	h.data[hash].value = v
	h.data[hash].filled = true
	h.len++
}

func (h *Map[K, V]) resize(newCap int) {
	newHMap := Map[K, V]{
		data: make([]entry[K, V], newCap, newCap),
		len:  0,
		eq:   h.eq,
		hash: h.hash,
	}

	for _, entry := range h.data {
		if entry.filled {
			newHMap.Put(entry.key, entry.value)
		}
	}

	*h = newHMap
}

// Removes the supplied key and associated value from the hash map if it is
// present. If the key is not present in the map then no action will be taken.
func (h *Map[K, V]) Remove(k K) {
	hash := h.clampedHash(h.hash(k))
	for h.data[hash].filled && !h.eq(h.data[hash].key, k) {
		hash = h.clampedHash(hash + 1)
	}
	if !h.data[hash].filled {
		return
	}

	h.data[hash] = entry[K, V]{}
	h.len--

	wrapped := false
	holeHash := hash
	hash = h.clampedHash(hash + 1)
	for h.data[hash].filled {
		requestedHash := h.clampedHash(h.hash(h.data[hash].key))
		moveEntry := !wrapped && requestedHash <= holeHash
		moveEntry = (moveEntry || (wrapped && requestedHash >= holeHash && requestedHash > hash))
		if moveEntry {
			h.data[holeHash] = h.data[hash]
			h.data[hash] = entry[K, V]{}
			holeHash = hash
		}

		wrapped = (hash == uint64(len(h.data)-1))
		hash = h.clampedHash(hash + 1)
	}

	// Original equation:
	// 	len/cap *100 <= _shrinkFactor
	// Except dividing ints is bad, we want more precision. So remove the
	// division and we get this:
	if cap(h.data) > _defaultInitialCap && h.len*100 <= _shrinkFactor*cap(h.data) {
		h.resize(cap(h.data) / 2)
	}
}

// Removes all values from the underlying hash but keeps the maps underlying
// capacity.
func (h *Map[K, V]) Clear() {
	for i, v := range h.data {
		if v.filled {
			h.data[i] = entry[K, V]{}
		}
	}
	h.len = 0
}

// Removes all values from the underlying hash and resets the maps capacity.
func (h *Map[K, V]) Zero() {
	h.data = make([]entry[K, V], _defaultInitialCap, _defaultInitialCap)
	h.len = 0
}

// Creates a copy of the supplied hash map. All values will be copied using
// memcpy, meaning a shallow copy will be made of the values.
func (h *Map[K, V]) Copy() *Map[K, V] {
	newData := make([]entry[K, V], cap(h.data), cap(h.data))
	for i, v := range h.data {
		newData[i] = v
	}

	return &Map[K, V]{
		len:  h.len,
		data: newData,
	}
}

// Iterates over all of the key, value pairs in the map and calls `op` on each
// pair. Any changes to the value will not be propigated back to the hash map.
func (h *Map[K, V]) Each(op func(k K, v V)) {
	for _, v := range h.data {
		if v.filled {
			op(v.key, v.value)
		}
	}
}

// Iterates over all of the key, value pairs in the map and calls `op` on each
// pair. The value may be mutated and the results will be seen by the hash map.
func (h *Map[K, V]) EachPntr(op func(k K, v *V)) {
	for _, v := range h.data {
		if v.filled {
			op(v.key, &v.value)
		}
	}
}

// TODO - std lib iterator?
