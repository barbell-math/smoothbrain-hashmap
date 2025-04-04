//go:build simd256

package slotprobes

import (
	"testing"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestGetSlotProbe(t *testing.T) {
	res, hasPotentialValue, hasEmptySlot := GetSlotProbe(
		3,
		[32]int8{
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
		},
		[32]int8{
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
		},
	)
	sbtest.Eq(t, res, [32]int8{
		0, 1, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0, 0,
	})
	sbtest.True(t, hasPotentialValue)
	sbtest.True(t, hasEmptySlot)

	res, hasPotentialValue, hasEmptySlot = GetSlotProbe(
		2,
		[32]int8{
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
			0, 1, 2, 0, 1, 2, 0, 0,
		},
		[32]int8{
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
			3, 3, 3, 1, 1, 1, 0, 0,
		},
	)
	sbtest.Eq(t, res, [32]int8{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})
	sbtest.False(t, hasPotentialValue)
	sbtest.True(t, hasEmptySlot)
}
