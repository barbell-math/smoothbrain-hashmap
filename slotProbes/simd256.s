//go:build amd64 && simd256

#include "textflag.h"

// These constants must be kept in sync with the vars in the AllBuildTargets file
#define Used $1
#define Deleted $2

// func getSlotProbe(
// 		key int8,
// 		flags [32]int8,
// 		slotKeys [32]int8,
// ) (potentialValues [32]int8, hasPotentialValue bool, hasEmptySlot bool)
//
// memory layout of the stack relative to FP
//  +0   				key
//  +1  through +32 	flags
//  +33 through +64		slotKeys
//  +65 through +71		allignment padding
//  +72 through +103	potentialValues 	return value
//  +104 				hasPotentialValue	return value
//  +105				hasEmptySlot		return value
TEXT ·GetSlotProbe(SB),NOSPLIT,$0
	MOVQ 			Used, R8				// load the used constant into reg
	VPBROADCASTB	R8, Y0					// broadcast used flag
	VMOVDQU8 		flags+1(FP), Y1			// load the flags into y reg
	VPAND			Y0, Y1, Y2				// used & flags

	MOVQ 			$0b11111111, R8 		// load 255 constant into reg
	VPBROADCASTB	R8, Y0 					// broadcast constant
	VPANDN			Y2, Y0, Y4				// invert (used & flag) result
	VPTEST			Y0, Y4					// test that reg is zero
	JZ setHasEmptySlotToTrue				// If ZF==0 jump to false
		MOVB $0, rv2+105(FP)				// Set hasEmptyValue rv to false
		JMP setHasEmptySlot
setHasEmptySlotToTrue:
		MOVB $1, rv2+105(FP)				// Set hasEmptyValue rv to true
setHasEmptySlot:

	MOVQ 			Deleted, R8 			// load the used constant into reg
	VPBROADCASTB	R8, Y0 					// broadcast deleted flag
	VPAND			Y0, Y1, Y3				// deleted & flags
	VPCMPEQB        Y0, Y3, Y3				// (deleted & flags) == deleted
	MOVQ 			$0b11111111, R8 		// load 255 constant into reg
	VPBROADCASTB	R8, Y0 					// broadcast constant
	VPXOR			Y0, Y3, Y3				// ((deleted & flags) == deleted) ^ 0b11111111
	VPAND			Y2, Y3, Y2				// (used & flags) & (((deleted & flags) == deleted) ^ 0b11111111)

	MOVQ 			key+0(FP), R8 			// load key arg into reg
	VPBROADCASTB	R8, Y0 					// broadcast key
	VMOVDQU8 		slotKeys+33(FP), Y3		// load the slotKeys into y reg
	VPCMPEQB        Y0, Y3, Y3				// slotKey == key
	VPAND			Y2, Y3, Y2				// (used & flags) & (((deleted & flags) == deleted) ^ 0b11111111) & (slotKey == key)

	VMOVDQU8 		Y2, rv+72(FP)			// cpy ret val to mem

	MOVQ 			$0b11111111, R8 		// load 255 constant into reg
	VPBROADCASTB	R8, Y0 					// broadcast constant
	VPTEST			Y0, Y2					// test that reg is zero
	JZ setHasPotentialValueToFalse			// If ZF==0 jump to false
		MOVB $1, rv2+104(FP)				// Set bool rv to true
		RET
setHasPotentialValueToFalse:
		MOVB $0, rv2+104(FP)				// Set bool rv to false
		RET
