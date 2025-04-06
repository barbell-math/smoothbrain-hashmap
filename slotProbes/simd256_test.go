//go:build simd256

package slotprobes

import (
	"testing"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestSlotProbe(t *testing.T) {
	res, isEmpty := SlotProbe(
		3,
		[32]uint8{
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
		},
		[32]uint8{
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
		},
	)
	sbtest.Eq(t, res, 0b00000010000000100000001000000010)
	sbtest.Eq(t, isEmpty, 0b11101101111011011110110111101101)
	sbtest.True(t, res > 0)

	// res, hasPotentialValue, hasEmptySlot = GetSlotProbe(
	// 	2,
	// 	[32]uint8{
	// 		0, 1, 2, 0, 1, 2, 0, 0,
	// 		0, 1, 2, 0, 1, 2, 0, 0,
	// 		0, 1, 2, 0, 1, 2, 0, 0,
	// 		0, 1, 2, 0, 1, 2, 0, 0,
	// 	},
	// 	[32]uint8{
	// 		3, 3, 3, 1, 1, 1, 0, 0,
	// 		3, 3, 3, 1, 1, 1, 0, 0,
	// 		3, 3, 3, 1, 1, 1, 0, 0,
	// 		3, 3, 3, 1, 1, 1, 0, 0,
	// 	},
	// )
	// sbtest.Eq(t, res, [32]uint8{
	// 	0, 0, 0, 0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0, 0, 0, 0,
	// 	0, 0, 0, 0, 0, 0, 0, 0,
	// })
	// sbtest.False(t, hasPotentialValue)
	// sbtest.True(t, hasEmptySlot)
}
