//go:build sbmap_simd512

package slotprobes

const (
	GroupSize = 64
)

func SlotProbe(
	key uint8,
	flags [GroupSize]uint8,
	slotKeys [GroupSize]uint8,
) (potentialValues uint64, isEmpty uint64)
