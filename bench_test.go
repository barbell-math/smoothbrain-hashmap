package sbmap

import (
	"math/rand"
	"testing"

	slotprobes "github.com/barbell-math/smoothbrain-hashmap/slotProbes"
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
		_, _ = slotprobes.GetSlotProbe(3, flags, slotKeys)
	}
}

func BenchmarkDifferentGrowthFactors(b *testing.B) {
	op := func(growthFactor int) {
		origGrowthFactor := _growFactor
		_growFactor = growthFactor

		h := New[int32, int64]()
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < 1000; i++ {
			h.Put(int32(randVals.Int31()), int64(randVals.Int31()))
		}
		randVals = rand.New(rand.NewSource(3))
		for i := 0; i < 1000; i++ {
			randVal := randVals.Int31()
			val, ok := h.Get(int32(randVal))
			if !ok || val != int64(randVals.Int31()) {
				b.Fatal("Map malfunctioned!!")
			}
		}

		_growFactor = origGrowthFactor
	}

	b.Run("85% Full Growth Factor", func(b *testing.B) {
		for b.Loop() {
			op(85)
		}
	})
	b.Run("70% Full Growth Factor", func(b *testing.B) {
		for b.Loop() {
			op(70)
		}
	})
	b.Run("65% Full Growth Factor", func(b *testing.B) {
		for b.Loop() {
			op(65)
		}
	})
	b.Run("60% Full Growth Factor", func(b *testing.B) {
		for b.Loop() {
			op(60)
		}
	})
	b.Run("55% Full Growth Factor", func(b *testing.B) {
		for b.Loop() {
			op(55)
		}
	})
	b.Run("50% Full Growth Factor", func(b *testing.B) {
		for b.Loop() {
			op(50)
		}
	})
}

func BenchmarkAgainstMapInsertOnly(b *testing.B) {
	customMapOp := func(size int) {
		h := New[int32, int64]()
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < size; i++ {
			h.Put(int32(randVals.Int31()), int64(randVals.Int31()))
		}
	}
	builtinMapOp := func(size int) {
		h := map[int32]int64{}
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < size; i++ {
			h[int32(randVals.Int31())] = int64(randVals.Int31())
		}
	}

	_growFactor = 50
	b.Run("1e2 elements", func(b *testing.B) {
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(1e2)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(1e2)
			}
		})
	})
	b.Run("1e3 elements", func(b *testing.B) {
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(1e3)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(1e3)
			}
		})
	})
	b.Run("1e4 elements", func(b *testing.B) {
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(1e4)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(1e4)
			}
		})
	})
	b.Run("1e6 elements", func(b *testing.B) {
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(1e6)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(1e6)
			}
		})
	})
}

func BenchmarkAgainstMapGetOnly(b *testing.B) {
	buildCustomHashMap := func(size int) Map[int32, int64] {
		customHashMap := New[int32, int64]()
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < size; i++ {
			customHashMap.Put(int32(randVals.Int31()), int64(randVals.Int31()))
		}
		return customHashMap
	}
	buildBuiltinMap := func(size int) map[int32]int64 {
		h := map[int32]int64{}
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < size; i++ {
			h[int32(randVals.Int31())] = int64(randVals.Int31())
		}
		return h
	}

	customMapOp := func(m *Map[int32, int64]) {
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < m.Len(); i++ {
			randVal := randVals.Int31()
			val, ok := m.Get(int32(randVal))
			_ = val
			_ = ok
		}
	}
	builtinMapOp := func(m map[int32]int64) {
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < len(m); i++ {
			randVal := randVals.Int31()
			val, ok := m[int32(randVal)]
			_ = val
			_ = ok
		}
	}

	_growFactor = 50
	b.Run("1e2 elements", func(b *testing.B) {
		customHashMap := buildCustomHashMap(1e2)
		builtinMap := buildBuiltinMap(1e2)
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(&customHashMap)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(builtinMap)
			}
		})
	})
	b.Run("1e3 elements", func(b *testing.B) {
		customHashMap := buildCustomHashMap(1e3)
		builtinMap := buildBuiltinMap(1e3)
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(&customHashMap)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(builtinMap)
			}
		})
	})
	b.Run("1e4 elements", func(b *testing.B) {
		customHashMap := buildCustomHashMap(1e4)
		builtinMap := buildBuiltinMap(1e4)
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(&customHashMap)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(builtinMap)
			}
		})
	})
	b.Run("1e6 elements", func(b *testing.B) {
		customHashMap := buildCustomHashMap(1e6)
		builtinMap := buildBuiltinMap(1e6)
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(&customHashMap)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(builtinMap)
			}
		})
	})
}

func BenchmarkAgainstMapRemoveOnly(b *testing.B) {
	buildCustomHashMap := func(size int) Map[int32, int64] {
		customHashMap := New[int32, int64]()
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < size; i++ {
			customHashMap.Put(int32(randVals.Int31()), int64(randVals.Int31()))
		}
		return customHashMap
	}
	buildBuiltinMap := func(size int) map[int32]int64 {
		h := map[int32]int64{}
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < size; i++ {
			h[int32(randVals.Int31())] = int64(randVals.Int31())
		}
		return h
	}

	customMapOp := func(m *Map[int32, int64]) {
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < m.Len(); i++ {
			randVal := randVals.Int31()
			m.Remove(randVal)
		}
		// fmt.Println("Max collision: ", m.maxChain)
	}
	builtinMapOp := func(m map[int32]int64) {
		randVals := rand.New(rand.NewSource(3))
		for i := 0; i < len(m); i++ {
			randVal := randVals.Int31()
			delete(m, randVal)
		}
	}

	_growFactor = 50
	b.Run("1e2 elements", func(b *testing.B) {
		customHashMap := buildCustomHashMap(1e2)
		builtinMap := buildBuiltinMap(1e2)
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(&customHashMap)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(builtinMap)
			}
		})
	})
	b.Run("1e3 elements", func(b *testing.B) {
		customHashMap := buildCustomHashMap(1e3)
		builtinMap := buildBuiltinMap(1e3)
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(&customHashMap)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(builtinMap)
			}
		})
	})
	b.Run("1e4 elements", func(b *testing.B) {
		customHashMap := buildCustomHashMap(1e4)
		builtinMap := buildBuiltinMap(1e4)
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(&customHashMap)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(builtinMap)
			}
		})
	})
	b.Run("1e6 elements", func(b *testing.B) {
		customHashMap := buildCustomHashMap(1e6)
		builtinMap := buildBuiltinMap(1e6)
		b.Run("Custom Map", func(b *testing.B) {
			for b.Loop() {
				customMapOp(&customHashMap)
			}
		})
		b.Run("Builtin Map", func(b *testing.B) {
			for b.Loop() {
				builtinMapOp(builtinMap)
			}
		})
	})
}
