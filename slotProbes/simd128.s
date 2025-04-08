//go:build amd64 && simd128

#include "textflag.h"

// !!! Note !!!
// These constants must be kept in sync with the constants in the
// AllBuildTargets file
#define Used $1
#define Deleted $2

// func getSlotProbe(
// 		key uint8,
// 		flags [16]uint8,
// 		slotKeys [16]uint8,
// ) (potentialValues uint16, isEmpty uint16)
//
// memory layout of the stack relative to FP
//  +0   				key					argument
//  +1  through +16 	flags				argument
//  +17 through +33		slotKeys			argument
//  +34 through +39		-					allignment padding
//  +40 through +41		potentialValues 	return value
//  +42 through +43		isEmpty 			return value
TEXT ·SlotProbe(SB),NOSPLIT,$0
	MOVB 			Used, R8				// load the used constant into reg
	VPBROADCASTB	R8, X0					// broadcast used flag
	VMOVDQU8 		flags+1(FP), X1			// load the flags into y reg
	MOVB 			Deleted, R9				// load the deleted constant into reg
	VPBROADCASTB	R9, X3					// broadcast deleted flag

	VPAND			X0, X1, X2				// used & flags
	VPCMPEQB        X0, X2, X2				// (used & flags) == used
	PMOVMSKB		X2, R8					// store the sign flags in reg
	XORQ			R12, R12				// clear register
	ADDQ			R8, R12					// cpy register
	NOTQ			R12						// inverse of (used & flags) == used

	VPAND			X3, X1, X2				// deleted & flags
	VPCMPEQB        X3, X2, X2				// (deleted & flags) == deleted
	PMOVMSKB		X2, R9					// store the sign flags in reg
	NOTW			R9						// invert the bit mask

	MOVQ 			key+0(FP), R10			// load the key constant into reg
	VPBROADCASTB	R10, X0					// broadcast key
	VMOVDQU8 		slotKeys+17(FP), X1		// load the slot keys into y reg
	VPCMPEQB        X0, X1, X2				// slotKey == key
	PMOVMSKB		X2, R10					// store the sign flags in reg

	ANDW			R8, R9					// ((used & flags) == used) & ((deleted & flags) == deleted)
	ANDW			R9, R10					// ((used & flags) == used) & ((deleted & flags) == deleted) & (slotKey == key)
	// TODO - why overwriting next value???
	MOVW			R10, rv+40(FP)			// place the final bit field result in mem
	MOVW			R12, rv+42(FP)			// place the is empty bit field result in mem
	RET
