package sbmap

import (
	"fmt"
	"math/rand"
	"testing"

	slotprobes "github.com/barbell-math/smoothbrain-hashmap/slotProbes"
)

type (
	benchOps[T any] struct {
		SetupOp  func(size int) T
		PutOp    func(input T, size int)
		GetOp    func(input T, size int)
		RemoveOp func(input T, size int)
		MixedOp  func(input T, size int)
	}
)

func BenchmarkGetSlotProbe(b *testing.B) {
	flags := [slotprobes.GroupSize]uint8{}
	slotKeys := [slotprobes.GroupSize]uint8{}
	flagsRow := [8]uint8{0, 1, 2, 0, 1, 2, 0, 0}
	slotKeysRow := [8]uint8{3, 3, 3, 1, 1, 1, 0, 0}
	for i := 0; i < slotprobes.GroupSize; i += 8 {
		copy(flags[i:], flagsRow[:])
		copy(slotKeys[i:], slotKeysRow[:])
	}

	for b.Loop() {
		_, _ = slotprobes.SlotProbe(3, flags, slotKeys)
	}
}

func BenchmarkBuiltinMap(b *testing.B) {
	benchmarkOps(b, benchOps[map[int32]int64]{
		SetupOp:  builtinMapSetup,
		PutOp:    builtinMapPut,
		GetOp:    builtinMapGet,
		RemoveOp: builtinMapRemove,
		MixedOp:  builtinMapMixedUsage,
	})
}

func BenchmarkCustomMap(b *testing.B) {
	benchmarkDifferentGrowthFactors(b, benchOps[*Map[int32, int64]]{
		SetupOp:  customMapSetup,
		PutOp:    customMapPut,
		GetOp:    customMapGet,
		RemoveOp: customMapRemove,
		MixedOp:  customMapMixedUsage,
	})
}

func benchmarkDifferentGrowthFactors[T any](b *testing.B, ops benchOps[T]) {
	subTests := func(b *testing.B, growFactor int) {
		origGrowthFactor := _growFactor
		_growFactor = growFactor
		b.Cleanup(func() {
			_growFactor = origGrowthFactor
		})
		benchmarkOps(b, ops)
	}

	for i := 50; i < 95; i += 5 {
		b.Run(
			fmt.Sprintf("%d%% GrowthFactor", i),
			func(b *testing.B) { subTests(b, i) },
		)
	}
}

func benchmarkOps[T any](b *testing.B, ops benchOps[T]) {
	b.Run("Put", func(b *testing.B) {
		benchmarkSizeRange(b, 1e2, 1e8, 10, ops.SetupOp, ops.PutOp)
	})
	b.Run("Get", func(b *testing.B) {
		benchmarkSizeRange(b, 1e2, 1e8, 10, ops.SetupOp, ops.GetOp)
	})
	b.Run("PutAndRemove", func(b *testing.B) {
		benchmarkSizeRange(b, 1e2, 1e8, 10, ops.SetupOp, ops.RemoveOp)
	})
	b.Run("Mixed", func(b *testing.B) {
		benchmarkSizeRange(b, 1e2, 1e5, 10, ops.SetupOp, ops.MixedOp)
	})
}

func benchmarkSizeRange[T any](
	b *testing.B,
	start int,
	end int,
	step int,
	setupOp func(size int) T,
	benchOp func(input T, size int),
) {
	for i := start; i < end; i *= step {
		b.Run(
			fmt.Sprintf("%d Elements", i),
			func(b *testing.B) {
				setupVal := setupOp(i)
				for b.Loop() {
					benchOp(setupVal, i)
				}
			},
		)
	}
}

func customMapSetup(size int) *Map[int32, int64] {
	rv := New[int32, int64]()
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		rv.Put(int32(randVals.Int31()), int64(randVals.Int31()))
	}
	return &rv
}

func builtinMapSetup(size int) map[int32]int64 {
	rv := map[int32]int64{}
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		rv[int32(randVals.Int31())] = int64(randVals.Int31())
	}
	return rv
}

func customMapPut(input *Map[int32, int64], size int) {
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		input.Put(int32(randVals.Int31()), int64(randVals.Int31()))
	}
}

func builtinMapPut(input map[int32]int64, size int) {
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		input[int32(randVals.Int31())] = int64(randVals.Int31())
	}
}

func customMapGet(input *Map[int32, int64], size int) {
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		randVal := randVals.Int31()
		val, ok := input.Get(int32(randVal))
		_ = val
		_ = ok
	}
}

func builtinMapGet(input map[int32]int64, size int) {
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		randVal := randVals.Int31()
		val, ok := input[int32(randVal)]
		_ = val
		_ = ok
	}
}

func customMapRemove(m *Map[int32, int64], size int) {
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		randVal := randVals.Int31()
		m.Remove(randVal)
	}
}

func builtinMapRemove(input map[int32]int64, size int) {
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		randVal := randVals.Int31()
		delete(input, randVal)
	}
}

func customMapMixedUsage(input *Map[int32, int64], size int) {
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		input.Put(int32(randVals.Int31()), int64(randVals.Int31()))
	}
	randVals = rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		randVal := randVals.Int31()
		val, ok := input.Get(int32(randVal))
		_ = val
		_ = ok
	}

	randVals = rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		randVal := randVals.Int31()
		input.Remove(int32(randVal))
		// The next value would have been the value so skip it
		_ = randVals.Int31()

		iterRandVals := rand.New(rand.NewSource(3))
		for j := 0; j < size; j++ {
			iterKey := iterRandVals.Int31()
			_ = iterRandVals.Int31()

			if j > i {
				val, ok := input.Get(int32(iterKey))
				_ = val
				_ = ok
			}
		}
	}
}

func builtinMapMixedUsage(input map[int32]int64, size int) {
	randVals := rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		input[randVals.Int31()] = int64(randVals.Int31())
	}
	randVals = rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		randVal := randVals.Int31()
		val, ok := input[int32(randVal)]
		_ = val
		_ = ok
	}

	randVals = rand.New(rand.NewSource(3))
	for i := 0; i < size; i++ {
		randVal := randVals.Int31()
		delete(input, int32(randVal))
		// The next value would have been the value so skip it
		_ = randVals.Int31()

		iterRandVals := rand.New(rand.NewSource(3))
		for j := 0; j < size; j++ {
			iterKey := iterRandVals.Int31()
			_ = iterRandVals.Int31()

			if j > i {
				val, ok := input[int32(iterKey)]
				_ = val
				_ = ok
			}
		}
	}
}
