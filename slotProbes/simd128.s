//go:build amd64 && sbmap_simd128

#include "textflag.h"

// !!! Note !!!
// These constants must be kept in sync with the constants in the
// AllBuildTargets file
#define Used $1
#define Deleted $2

// func SlotProbe(
// 		key uint8,
// 		flags [16]uint8,
// 		slotKeys [16]uint8,
// ) (potentialValues uint16, isEmpty uint16)
//
// memory layout of the stack relative to FP
//  +0   				key					argument
//  +1  through +16 	flags				argument
//  +17 through +33		slotKeys			argument
//  +34 through +39		-					alignment padding
//  +40 through +41		potentialValues 	return value
//  +42 through +43		isEmpty 			return value
TEXT ·SlotProbe(SB),NOSPLIT,$0
	MOVB 			Used, R8				// load the used constant into reg
	VPBROADCASTB	R8, X0					// broadcast used flag
	VMOVDQU8 		flags+1(FP), X1			// load the flags into x reg

	MOVB 			Deleted, R9				// load the deleted constant into reg
	VPBROADCASTB	R9, X3					// broadcast deleted flag

	MOVB 			key+0(FP), R10			// load the key constant into reg
	VPBROADCASTB	R10, X4					// broadcast key
	VMOVDQU8 		slotKeys+17(FP), X5		// load the slot keys into x reg

	VPAND			X0, X1, X2				// used & flags
	VPCMPEQB        X0, X2, K1				// (used & flags) == used
	KMOVW			K1, R8					// get the final used bit flags
	XORQ			R12, R12				// clear register
	ADDQ			R8, R12					// cpy register
	NOTQ			R12						// inverse of (used & flags) == used

	VPAND			X3, X1, X2				// deleted & flags
	VPCMPEQB        X3, X2, K1				// (deleted & flags) == deleted
	KMOVW			K1, R9					// get the final deleted bit flags
	NOTW			R9						// invert the bit mask

	VPCMPEQB        X4, X5, K1				// slotKey == key
	KMOVW			K1, R10					// get the final slot key bit flags

	ANDW			R8, R9					// ((used & flags) == used) & ((deleted & flags) == deleted)
	ANDW			R9, R10					// ((used & flags) == used) & ((deleted & flags) == deleted) & (slotKey == key)

	MOVW			R10, rv+40(FP)			// place the final bit field result in mem
	MOVW			R12, rv+42(FP)			// place the is empty bit field result in mem
	RET
