package sbmap

import (
	"hash/maphash"
	"math/rand"
	"testing"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestEntryFlags(t *testing.T) {
	e := entry[int32, int64]{}
	sbtest.False(t, e.used())
	sbtest.False(t, e.deleted())

	e.flags |= used
	sbtest.True(t, e.used())
	sbtest.False(t, e.deleted())

	e.flags |= deleted
	sbtest.True(t, e.used())
	sbtest.True(t, e.deleted())

	e.flags &= ^deleted
	sbtest.True(t, e.used())
	sbtest.False(t, e.deleted())
}

func TestHashMapPut(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))
}

func TestHashMapGet(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))

	val, ok := h.Get(1)
	sbtest.True(t, ok)
	sbtest.Eq(t, 1, val)
	val, ok = h.Get(2)
	sbtest.True(t, ok)
	sbtest.Eq(t, 2, val)
	val, ok = h.Get(3)
	sbtest.True(t, ok)
	sbtest.Eq(t, 3, val)

	val, ok = h.Get(4)
	sbtest.False(t, ok)
	sbtest.Eq(t, 0, val)
}

func TestHashMapRemove(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))

	h.Remove(1)
	sbtest.Eq(t, 2, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))
	h.Remove(1)
	sbtest.Eq(t, 2, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))

	h.Remove(2)
	sbtest.Eq(t, 1, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))
	h.Remove(3)
	sbtest.Eq(t, 0, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))
}

func TestHashMapClear(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))

	h.Clear()
	sbtest.Eq(t, 0, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))
}

func TestHashMapZero(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))

	h.Zero()
	sbtest.Eq(t, 0, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))
}

func TestGrowAndShrinkFactors(t *testing.T) {
	origGrowFactor := _growFactor
	origShrinkFactor := _shrinkFactor
	origDefaultInitialCap := _defaultInitialCap
	t.Cleanup(func() {
		_growFactor = origGrowFactor
		_shrinkFactor = origShrinkFactor
		_defaultInitialCap = origDefaultInitialCap
	})

	_growFactor = 50
	_shrinkFactor = 50
	_defaultInitialCap = 4
	h := New[int8, int16]()

	h.Put(1, 1)
	sbtest.Eq(t, 1, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))
	h.Put(2, 2)
	sbtest.Eq(t, 2, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))

	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap*2, cap(h.data))

	h.Remove(3)
	sbtest.Eq(t, 2, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.data))

	val, ok := h.Get(1)
	sbtest.True(t, ok)
	sbtest.Eq(t, 1, val)
	val, ok = h.Get(2)
	sbtest.True(t, ok)
	sbtest.Eq(t, 2, val)
	val, ok = h.Get(3)
	sbtest.False(t, ok)
	sbtest.Eq(t, 0, val)
}

func TestLargeishDataset(t *testing.T) {
	op := func() {
		h := New[int32, int64]()

		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < 1000; i++ {
			h.Put(int32(randVals.Int31()), int64(randVals.Int31()))
		}
		randVals = rand.New(rand.NewSource(3))
		for i := 0; i < 1000; i++ {
			randVal := randVals.Int31()
			val, ok := h.Get(int32(randVal))
			sbtest.True(t, ok)
			sbtest.Eq(t, int64(randVals.Int31()), val)
		}

		randVals = rand.New(rand.NewSource(3))
		for i := 0; i < 1000; i++ {
			randVal := randVals.Int31()
			h.Remove(int32(randVal))
			// The next value would have been the value so skip it
			_ = randVals.Int31()

			iterRandVals := rand.New(rand.NewSource(3))
			for j := 0; j < 1000; j++ {
				iterKey := iterRandVals.Int31()
				iterVal := iterRandVals.Int31()

				if j > i {
					val, ok := h.Get(int32(iterKey))
					sbtest.True(t, ok)
					sbtest.Eq(t, int64(iterVal), val)
				}
			}
		}
	}

	for i := 0; i < 20; i++ {
		// Testing with different hash seed values but with the same set of
		// values to produce different map structures
		_comparableSeed = maphash.MakeSeed()
		op()
	}
}
