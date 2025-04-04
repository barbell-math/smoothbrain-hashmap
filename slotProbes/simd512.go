//go:build simd512

package slotprobes

const (
	GroupSize = 64
)

func GetSlotProbe(
	key int8,
	flags [GroupSize]int8,
	slotKeys [GroupSize]int8,
) [GroupSize]int8
