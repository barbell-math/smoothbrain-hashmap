// A very simple library that implements a generic, linear probing map.
package sbmap

import (
	"hash/maphash"
	"iter"
	"math/bits"
	"reflect"

	slotprobes "github.com/barbell-math/smoothbrain-hashmap/slotProbes"
)

type (
	slot[K any, V any] struct {
		key   K
		value V
	}

	group[K any, V any] struct {
		flags    [slotprobes.GroupSize]uint8
		slotKeys [slotprobes.GroupSize]uint8
		slots    [slotprobes.GroupSize]slot[K, V]
	}

	Map[K any, V any] struct {
		groups []group[K, V]
		len    int
		del    int
		eq     func(l K, r K) bool
		hash   func(l K) uint64
	}
)

var (
	_comparableSeed = maphash.MakeSeed()
	// The default initial capacity that will be used if no capacity or zero
	// capacity is supplied
	_defaultInitialCap = max(1, 64/slotprobes.GroupSize)
	// A value between 0 and 100 that determines how full the map can get before
	// the hash map doubles the underlying slice.
	_growFactor = 75
	// A value between 0 and 100 that determines how empty the map can get
	// before the hash map halves the underlying slice.
	_shrinkFactor = 25
	// The power of two to use when increasing the backing slices capacity.
	_sliceGrowthFactor = 1
)

// An equality function that can be passed to [NewCustom] when using a
// comparable type. If the key type is comparable then you can simply use [New]
// instead of [NewCustom] and this function will be Used by default.
func ComparableEqual[T comparable](l T, r T) bool {
	return l == r
}

// A hash function that can be passed to [NewCustom] when using a comparable
// type. If the key type is comparable then you can simply use [New] instead of
// [NewCustom] and this function will be Used by default.
func ComparableHash[T comparable]() func(v T) uint64 {
	// For speed if the underlying type is an int then just return the int
	// value as the hash. This might be less evenly distributed but is much
	// faster than the [maphash.Comparable] function.
	switch reflect.TypeFor[T]().Kind() {
	case reflect.Int:
		return func(v T) uint64 {
			return uint64(any(v).(int))
		}
	case reflect.Int8:
		return func(v T) uint64 {
			return uint64(any(v).(int8))
		}
	case reflect.Int16:
		return func(v T) uint64 {
			return uint64(any(v).(int16))
		}
	case reflect.Int32:
		return func(v T) uint64 {
			return uint64(any(v).(int32))
		}
	case reflect.Int64:
		return func(v T) uint64 {
			return uint64(any(v).(int64))
		}
	case reflect.Uint:
		return func(v T) uint64 {
			return uint64(any(v).(uint))
		}
	case reflect.Uint8:
		return func(v T) uint64 {
			return uint64(any(v).(uint8))
		}
	case reflect.Uint16:
		return func(v T) uint64 {
			return uint64(any(v).(uint16))
		}
	case reflect.Uint32:
		return func(v T) uint64 {
			return uint64(any(v).(uint32))
		}
	case reflect.Uint64:
		return func(v T) uint64 {
			return uint64(any(v).(uint64))
		}
	default:
		return func(v T) uint64 {
			return maphash.Comparable(_comparableSeed, v)
		}
	}
}

// Creates a Map where K is the key type and V is the value type.
// [ComparableEqual] and [ComparableHash] functions will be Used by the returned
// Map. For creating a Map with non-comparable types or custom hash and equality
// functions refer to [NewCustom].
func New[K comparable, V comparable]() Map[K, V] {
	return Map[K, V]{
		groups: make([]group[K, V], _defaultInitialCap, _defaultInitialCap),
		len:    0,
		eq:     ComparableEqual[K],
		hash:   ComparableHash[K](),
	}
}

// Creates a Map where K is the key type and V is the value type with a capacity
// of `_cap`. [ComparableEqual] and [ComparableHash] functions will be Used by
// the returned Map. For creating a Map with non-comparable types or custom hash
// and equality functions refer to [NewCustom].
func NewCap[K comparable, V comparable](_cap int) Map[K, V] {
	return Map[K, V]{
		groups: make([]group[K, V], _cap, _cap),
		len:    0,
		eq:     ComparableEqual[K],
		hash:   ComparableHash[K](),
	}
}

// Creates a Map where K is the key type and V is the value type with a capacity
// of `_cap`. The supplied `eq` and `hash` functions will be Used by the Map. If
// two values are equal the `hash` function hash function should return the same
// hash for both values.
func NewCustom[K any, V any](
	_cap int,
	eq func(l K, r K) bool,
	hash func(v K) uint64,
) Map[K, V] {
	return Map[K, V]{
		groups: make([]group[K, V], _cap, _cap),
		len:    0,
		eq:     eq,
		hash:   hash,
	}
}

// Returns the number of elements in the hash map. This is different than the
// maps capacity.
func (m *Map[K, V]) Len() int {
	return m.len - m.del
}

// The group hash is the upper 57 bits of the original hash.
// The slot hash is the lower 7 bits of the original hash.
func (m *Map[K, V]) splitHash(hash uint64) (uint64, uint8) {
	return hash >> 7, uint8(hash & 0b1111111)
}

// The double hash function that the hash map will use when a collision occurs
// to perform probing of the underlying slice.
func (m *Map[K, V]) doubleHash(hash uint64) uint64 {
	// TODO - TEST
	// return maphash.Comparable(_comparableSeed, hash) | 0b1
	return 3
}

// Clamps the hash to always be within the groups slice length.
func (m *Map[K, V]) clampedGroupHash(hash uint64) uint64 {
	return hash & (uint64(cap(m.groups)) - 1)
}

// Gets the value that is related to the supplied key. If the key is found the
// boolean return value will be true and the value will be returned. If the key
// is not found the boolean return value will be false and a zero-initialized
// value of type V will be returned.
func (m *Map[K, V]) Get(k K) (V, bool) {
	groupHash, slotHash := m.splitHash(m.hash(k))
	groupHash = m.clampedGroupHash(groupHash)
	// All probing is performed on the group level
	doubleHash := m.doubleHash(groupHash)

	for i := uint64(1); ; i++ {

		potentialMatches, emptySlots := slotprobes.SlotProbe(
			slotHash,
			m.groups[groupHash].flags,
			m.groups[groupHash].slotKeys,
		)

		for j := 0; potentialMatches > 0; {
			tz := bits.TrailingZeros(uint(potentialMatches))
			potentialMatches >>= tz
			emptySlots >>= tz
			j += tz

			if m.eq(m.groups[groupHash].slots[j].key, k) {
				return m.groups[groupHash].slots[j].value, true
			}
			potentialMatches = potentialMatches >> 1
			emptySlots = emptySlots >> 1
			j++
		}
		// There should never be a potential match after an empty slot
		// Meaning, if there is any remaining empty slots, the value was not found
		if emptySlots > 0 {
			var tmp V
			return tmp, false
		}

		groupHash = m.clampedGroupHash(groupHash + i*doubleHash)
	}
}

// Places the supplied key, value pair in the hash map. If the key was already
// present in the map the old value will be overwritten. The map will rehash as
// necessary.
func (m *Map[K, V]) Put(k K, v V) {
	// Original equation:
	// 	len/cap *100 >= _growFactor
	// Except dividing ints is bad, we want more precision. So remove the
	// division and we get this:
	if m.len*100 >= _growFactor*len(m.groups)*slotprobes.GroupSize {
		m.rehash(cap(m.groups) << _sliceGrowthFactor)
	}

	groupHash, slotHash := m.splitHash(m.hash(k))
	groupHash = m.clampedGroupHash(groupHash)
	// All probing is performed on the group level
	doubleHash := m.doubleHash(groupHash)

	for i := uint64(1); ; i++ {

		potentialMatches, emptySlots := slotprobes.SlotProbe(
			slotHash,
			m.groups[groupHash].flags,
			m.groups[groupHash].slotKeys,
		)

		// Potential matches bit field will be 0 when no slots matched
		// Empty slots bit field will be 0 when no slots are empty
		for j := 0; potentialMatches > 0 || emptySlots > 0; {
			tz := min(
				bits.TrailingZeros(uint(potentialMatches)),
				bits.TrailingZeros(uint(emptySlots)),
			)
			potentialMatches >>= tz
			emptySlots >>= tz
			j += tz

			if emptySlots&0b1 == 1 {
				m.groups[groupHash].slots[j] = slot[K, V]{key: k, value: v}
				m.groups[groupHash].slotKeys[j] = slotHash
				m.groups[groupHash].flags[j] |= slotprobes.Used
				m.len++
				return
			}
			if potentialMatches&0b1 == 1 && m.eq(m.groups[groupHash].slots[j].key, k) {
				m.groups[groupHash].slots[j].value = v
				m.del -= int((m.groups[groupHash].flags[j] & slotprobes.Deleted) >> 1)
				return
			}
			potentialMatches = potentialMatches >> 1
			emptySlots = emptySlots >> 1
			j++
		}

		groupHash = m.clampedGroupHash(groupHash + i*doubleHash)
	}
}

func (m *Map[K, V]) rehash(newCap int) {
	newHMap := Map[K, V]{
		groups: make([]group[K, V], newCap, newCap),
		len:    0,
		eq:     m.eq,
		hash:   m.hash,
	}

	for i := range m.groups {
		for j := range m.groups[i].slots {
			if m.groups[i].flags[j]&(slotprobes.Used|slotprobes.Deleted) == 0b1 {
				newHMap.Put(m.groups[i].slots[j].key, m.groups[i].slots[j].value)
			}
		}
	}

	*m = newHMap
}

// Removes the supplied key and associated value from the hash map if it is
// present. If the key is not present in the map then no action will be taken.
func (m *Map[K, V]) Remove(k K) {
	groupHash, slotHash := m.splitHash(m.hash(k))
	groupHash = m.clampedGroupHash(groupHash)
	// All probing is performed on the group level
	doubleHash := m.doubleHash(groupHash)

	for i := uint64(1); ; i++ {
		potentialMatches, emptySlots := slotprobes.SlotProbe(
			slotHash,
			m.groups[groupHash].flags,
			m.groups[groupHash].slotKeys,
		)

		// Potential matches bit field will be 0 when no slots matched
		// Empty slots bit field will be 0 when no slots are empty
		for j := 0; potentialMatches > 0 || emptySlots > 0; {
			tz := min(
				bits.TrailingZeros(uint(potentialMatches)),
				bits.TrailingZeros(uint(emptySlots)),
			)
			potentialMatches >>= tz
			emptySlots >>= tz
			j += tz

			if emptySlots&0b1 == 1 {
				goto end
			}
			if potentialMatches&0b1 == 1 && m.eq(m.groups[groupHash].slots[j].key, k) {
				m.del += int(((^m.groups[groupHash].flags[j]) & slotprobes.Deleted) >> 1)
				m.groups[groupHash].flags[j] |= slotprobes.Deleted
				goto end
			}
			potentialMatches = potentialMatches >> 1
			emptySlots = emptySlots >> 1
			j++
		}

		groupHash = m.clampedGroupHash(groupHash + i*doubleHash)
	}

end:
	// Original equation:
	// 	len/cap *100 <= _shrinkFactor
	// Except dividing ints is bad, we want more precision. So remove the
	// division and we get this:
	if cap(m.groups) > _defaultInitialCap && m.Len()*100 <= _shrinkFactor*cap(m.groups)*slotprobes.GroupSize {
		m.rehash(cap(m.groups) >> _sliceGrowthFactor)
	}
}

// Removes all values from the underlying hash but keeps the maps underlying
// capacity.
func (m *Map[K, V]) Clear() {
	for i := range m.groups {
		for j := range m.groups[i].slots {
			if m.groups[i].flags[j]&(slotprobes.Used|slotprobes.Deleted) == 0b1 {
				m.groups[i].slots[j] = slot[K, V]{}
			}
		}
	}
	m.len = 0
}

// Removes all values from the underlying hash and resets the maps capacity to
// the default initial capacity.
func (m *Map[K, V]) Zero() {
	m.groups = make([]group[K, V], _defaultInitialCap, _defaultInitialCap)
	m.len = 0
}

// Creates a copy of the supplied hash map. All values will be copied using
// memcpy, meaning a shallow copy will be made of the values.
func (m *Map[K, V]) Copy() *Map[K, V] {
	newData := make([]group[K, V], cap(m.groups), cap(m.groups))
	for i := range m.groups {
		for j := range m.groups[i].slots {
			newData[i].slots[j] = m.groups[i].slots[j]
		}
		newData[i].flags = m.groups[i].flags
		newData[i].slotKeys = m.groups[i].slotKeys
	}

	return &Map[K, V]{
		groups: newData,
		len:    m.len,
		eq:     m.eq,
		hash:   m.hash,
	}
}

// Iterates over all of the keys in the map. Uses the stdlib `iter` package so
// this function can be Used in a standard `for` loop.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(k K) bool) {
		for i := range m.groups {
			for j := range m.groups[i].slots {
				if m.groups[i].flags[j]&(slotprobes.Used|slotprobes.Deleted) == 0b1 &&
					!yield(m.groups[i].slots[j].key) {
					return
				}
			}
		}
	}
}

// Iterates over all of the values in the map. Uses the stdlib `iter` package so
// this function can be Used in a standard `for` loop.
func (m *Map[K, V]) Vals() iter.Seq[V] {
	return func(yield func(v V) bool) {
		for i := range m.groups {
			for j := range m.groups[i].slots {
				if m.groups[i].flags[j]&(slotprobes.Used|slotprobes.Deleted) == 0b1 &&
					!yield(m.groups[i].slots[j].value) {
					return
				}
			}
		}
	}
}

// Iterates over all of the values in the map. Uses the stdlib `iter` package so
// this function can be Used in a standard `for` loop. The value may be mutated
// and the results will be seen by the hash map.
func (m *Map[K, V]) PntrVals() iter.Seq[*V] {
	return func(yield func(v *V) bool) {
		for i := range m.groups {
			for j := range m.groups[i].slots {
				if m.groups[i].flags[j]&(slotprobes.Used|slotprobes.Deleted) == 0b1 &&
					!yield(&m.groups[i].slots[j].value) {
					return
				}
			}
		}
	}
}
