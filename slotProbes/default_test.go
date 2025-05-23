//go:build !sbmap_simd512 && !sbmap_simd256 && !sbmap_simd128

package slotprobes

import (
	"testing"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestSlotProbe(t *testing.T) {
	res, isEmpty := SlotProbe(
		3,
		[8]uint8{0, 1, 2, 0, 1, 2, 0, 0},
		[8]uint8{3, 3, 3, 1, 1, 1, 0, 0},
	)
	sbtest.Eq(t, res, 0b00000010)
	sbtest.Eq(t, isEmpty, 0b11101101)
	sbtest.True(t, res > 0)
}
