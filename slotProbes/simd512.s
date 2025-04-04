//go:build amd64 && simd512

#include "textflag.h"

// func getSlotProbe(
// 		key int8,
// 		flags [64]int8,
// 		slotKeys [64]int8,
// ) [64]int8
//
// memory layout of the stack relative to FP
//  +0   				key
//  +1  through +64 	flags
//  +65 through +130	slotKeys
TEXT ·GetSlotProbe(SB),NOSPLIT,$0
	VMOVDQU8 	flags+1(FP), Z0
	VMOVDQU8 	slotKeys+65(FP), Z1
	VPADDB 		Z0, Z1, Z2
	VMOVDQU8 	Z2, flags+1(FP)
	// MOVB 		flags+1(FP), BX
	RET
