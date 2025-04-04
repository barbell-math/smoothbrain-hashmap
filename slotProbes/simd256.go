//go:build simd256

package slotprobes

const (
	GroupSize = 32
)

func GetSlotProbe(
	key int8,
	flags [GroupSize]int8,
	slotKeys [GroupSize]int8,
) (potentialValues [GroupSize]int8, hasPotentialValue bool, hasEmptySlot bool)
