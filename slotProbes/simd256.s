//go:build amd64 && sbmap_simd256

#include "textflag.h"

// !!! Note !!!
// These constants must be kept in sync with the constants in the
// AllBuildTargets file
#define Used $1
#define Deleted $2

// func SlotProbe(
// 		key uint8,
// 		flags [32]uint8,
// 		slotKeys [32]uint8,
// ) (potentialValues uint32, isEmpty uint32)
//
// memory layout of the stack relative to FP
//  +0   				key					argument
//  +1  through +32 	flags				argument
//  +33 through +64		slotKeys			argument
//  +65 through +71		-					alignment padding
//  +72 through +75		potentialValues 	return value
//  +76 through +79		isEmpty 			return value
TEXT ·SlotProbe(SB),NOSPLIT,$0
	MOVB 			Used, R8				// load the used constant into reg
	VPBROADCASTB	R8, Y0					// broadcast used flag
	VMOVDQU8 		flags+1(FP), Y1			// load the flags into y reg

	MOVB 			Deleted, R9				// load the deleted constant into reg
	VPBROADCASTB	R9, Y3					// broadcast deleted flag

	MOVB 			key+0(FP), R10			// load the key constant into reg
	VPBROADCASTB	R10, Y4					// broadcast key
	VMOVDQU8 		slotKeys+33(FP), Y5		// load the slot keys into y reg

	VPAND			Y0, Y1, Y2				// used & flags
	VPCMPEQB        Y0, Y2, K1				// (used & flags) == used
	KMOVD			K1, R8					// get the final used bit flags
	XORQ			R12, R12				// clear register
	ADDQ			R8, R12					// cpy register
	NOTQ			R12						// inverse of (used & flags) == used

	VPAND			Y3, Y1, Y2				// deleted & flags
	VPCMPEQB        Y3, Y2, K1				// (deleted & flags) == deleted
	KMOVD			K1, R9					// get the final deleted bit flags
	NOTQ			R9						// invert the bit mask

	VPCMPEQB        Y4, Y5, K1				// slotKey == key
	KMOVD			K1, R10					// get the final slot key bit flags

	ANDQ			R8, R9					// ((used & flags) == used) & ((deleted & flags) == deleted)
	ANDQ			R9, R10					// ((used & flags) == used) & ((deleted & flags) == deleted) & (slotKey == key)

	MOVD			R10, rv+72(FP)			// place the final bit field result in mem
	MOVD			R12, rv+76(FP)			// place the is empty bit field result in mem
	RET
