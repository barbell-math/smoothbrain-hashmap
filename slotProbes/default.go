//go:build !simd512 && !simd256

package slotprobes

const (
	GroupSize = 8
)

// This is the slow approach that is Used when no simd is available. It is the
// default operation that can be performed by the CPU in standard registers.
func GetSlotProbe(
	key int8,
	flags [GroupSize]int8,
	slotKeys [GroupSize]int8,
) (potentialValues [GroupSize]int8, hasPotentialValue bool, hasEmptySlot bool) {
	hasEmptySlot = false
	UsedSplat := [GroupSize]int8{}
	for i := 0; i < GroupSize; i++ {
		UsedSplat[i] = flags[i] & Used
		hasEmptySlot = hasEmptySlot || UsedSplat[i] == 0
	}

	delSplat := [GroupSize]int8{}
	for i := 0; i < GroupSize; i++ {
		delSplat[i] = ((flags[i] & Deleted) >> 1) ^ 0b1
	}

	eqSplat := [GroupSize]int8{}
	for i := 0; i < GroupSize; i++ {
		if slotKeys[i] == key {
			eqSplat[i] = 1
		}
	}

	hasPotentialValue = false
	for i := 0; i < GroupSize; i++ {
		potentialValues[i] = UsedSplat[i] & delSplat[i] & eqSplat[i]
		hasPotentialValue = hasPotentialValue || potentialValues[i] == 0b1
	}

	return
}
