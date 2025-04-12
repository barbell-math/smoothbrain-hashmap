//go:build !sbmap_simd512 && !sbmap_simd256 && !sbmap_simd128

package slotprobes

const (
	GroupSize = 8
)

// This is the slow approach that is Used when no simd is available. It is the
// default operation that can be performed by the CPU in standard registers.
func SlotProbe(
	key uint8,
	flags [GroupSize]uint8,
	slotKeys [GroupSize]uint8,
) (potentialValues uint8, isEmpty uint8) {
	isEmpty = 0
	var usedSplat uint8
	for i := GroupSize - 1; i >= 0; i-- {
		usedSplat = (flags[i] & Used) | (usedSplat << 1)
	}
	isEmpty = ^usedSplat

	var delSplat uint8
	for i := GroupSize - 1; i >= 0; i-- {
		delSplat = ((flags[i] & Deleted) >> 1) | (delSplat << 1)
	}
	delSplat = ^delSplat

	var eqSplat uint8
	for i := GroupSize - 1; i >= 0; i-- {
		var iterV uint8
		if slotKeys[i] == key {
			iterV = 0b1
		}
		eqSplat = iterV | (eqSplat << 1)
	}

	potentialValues = (usedSplat & delSplat & eqSplat)

	return
}
