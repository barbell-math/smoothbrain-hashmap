//go:build simd128

package slotprobes

const (
	GroupSize = 16
)

func SlotProbe(
	key uint8,
	flags [GroupSize]uint8,
	slotKeys [GroupSize]uint8,
) (potentialValues uint16, isEmpty uint16)
