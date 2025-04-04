//go:build !simd512 && !simd256

package slotprobes

import (
	"testing"

	sbtest "github.com/barbell-math/smoothbrain-test"
)

func TestGetSlotProbe(t *testing.T) {
	res, hasPotentialValue, hasEmptyValue := GetSlotProbe(
		3,
		[8]int8{0, 1, 2, 0, 1, 2, 0, 0},
		[8]int8{3, 3, 3, 1, 1, 1, 0, 0},
	)
	sbtest.Eq(t, res, [8]int8{0, 1, 0, 0, 0, 0, 0, 0})
	sbtest.True(t, hasPotentialValue)
	sbtest.True(t, hasEmptyValue)
}
