//go:build simd256

package slotprobes

const (
	GroupSize = 32
)

func SlotProbe(
	key uint8,
	flags [GroupSize]uint8,
	slotKeys [GroupSize]uint8,
) (potentialValues uint32, isEmpty uint32)
