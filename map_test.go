package sbmap

import (
	"hash/maphash"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"slices"
	"testing"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestSplitHash(t *testing.T) {
	m := New[uint32, uint64]()
	group, slot := m.splitHash(0b11111111)
	sbtest.Eq(t, group, 0b1)
	sbtest.Eq(t, slot, 0b1111111)

	group, slot = m.splitHash(0b1011111111)
	sbtest.Eq(t, group, 0b101)
	sbtest.Eq(t, slot, 0b1111111)

	group, slot = m.splitHash(0b1011111110)
	sbtest.Eq(t, group, 0b101)
	sbtest.Eq(t, slot, 0b1111110)
}

func TestHashMapPut(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))
}

func TestHashMapGet(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

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
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

	h.Remove(1)
	sbtest.Eq(t, 2, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))
	h.Remove(1)
	sbtest.Eq(t, 2, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

	h.Remove(2)
	sbtest.Eq(t, 1, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))
	h.Remove(3)
	sbtest.Eq(t, 0, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))
}

func TestHashMapClear(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

	h.Clear()
	sbtest.Eq(t, 0, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))
}

func TestHashMapZero(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

	h.Zero()
	sbtest.Eq(t, 0, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))
}

func TestCopy(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

	h2 := h.Copy()
	sbtest.Eq(t, 3, h2.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h2.groups))

	val, ok := h2.Get(1)
	sbtest.True(t, ok)
	sbtest.Eq(t, 1, val)
	val, ok = h2.Get(2)
	sbtest.True(t, ok)
	sbtest.Eq(t, 2, val)
	val, ok = h2.Get(3)
	sbtest.True(t, ok)
	sbtest.Eq(t, 3, val)
}

func TestKeys(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

	keys := slices.Collect(h.Keys())
	sbtest.SlicesMatchUnordered(t, []int8{1, 2, 3}, keys)
}

func TestValues(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

	vals := slices.Collect(h.Vals())
	sbtest.SlicesMatchUnordered(t, []int16{1, 2, 3}, vals)
}

func TestValuesPntrs(t *testing.T) {
	h := New[int8, int16]()

	h.Put(1, 1)
	h.Put(2, 2)
	h.Put(3, 3)
	sbtest.Eq(t, 3, h.Len())
	sbtest.Eq(t, _defaultInitialCap, cap(h.groups))

	vals := []int16{0, 0, 0}
	i := 0
	for v := range h.PntrVals() {
		vals[i] = *v
		i++
	}
	sbtest.SlicesMatchUnordered(t, []int16{1, 2, 3}, vals)
}

func TestLargeishDataset(t *testing.T) {
	f, err := os.Create("./bs/testProf.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	op := func() {
		h := New[int32, int64]()

		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < 10000; i++ {
			h.Put(int32(randVals.Int31()), int64(randVals.Int31()))
		}
		randVals = rand.New(rand.NewSource(3))
		for i := 0; i < 10000; i++ {
			randVal := randVals.Int31()
			val, ok := h.Get(int32(randVal))
			sbtest.True(t, ok)
			sbtest.Eq(t, int64(randVals.Int31()), val)
		}

		randVals = rand.New(rand.NewSource(3))
		for i := 0; i < 10000; i++ {
			randVal := randVals.Int31()
			h.Remove(int32(randVal))
			// The next value would have been the value so skip it
			_ = randVals.Int31()

			iterRandVals := rand.New(rand.NewSource(3))
			for j := 0; j < 10000; j++ {
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

	for i := 0; i < 10; i++ {
		// Testing with different hash seed values but with the same set of
		// values to produce different map structures
		_comparableSeed = maphash.MakeSeed()
		op()
	}
}
