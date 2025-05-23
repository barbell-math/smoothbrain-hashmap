//go:build sbmap_simd512

package slotprobes

import (
	"testing"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestGetSlotProbe(t *testing.T) {
	res, isEmpty := SlotProbe(
		3,
		[64]uint8{
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
		},
		[64]uint8{
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
		},
	)
	sbtest.Eq(
		t, res,
		0b0000001000000010000000100000001000000010000000100000001000000010,
	)
	sbtest.Eq(
		t, isEmpty,
		0b1110110111101101111011011110110111101101111011011110110111101101,
	)
	sbtest.True(t, res > 0)
}
