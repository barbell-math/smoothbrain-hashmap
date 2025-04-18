//go:build amd64 && sbmap_simd512

#include "textflag.h"

// !!! Note !!!
// These constants must be kept in sync with the constants in the
// AllBuildTargets file
#define Used $1
#define Deleted $2

// func getSlotProbe(
// 		key uint8,
// 		flags [64]uint8,
// 		slotKeys [64]uint8,
// ) (potentialValues uint64, isEmpty uint64)
//
// memory layout of the stack relative to FP
//  +0   				key					argument
//  +1   through +64 	flags				argument
//  +65  through +129	slotKeys			argument
//  +130 through +135	-					alignment padding
//  +136 through +143	potentialValues 	return value
//  +144 through +151	isEmpty 			return value
TEXT ·SlotProbe(SB),NOSPLIT,$0
	MOVB 			Used, R8				// load the used constant into reg
	VPBROADCASTB	R8, Z0					// broadcast used flag
	VMOVDQU8 		flags+1(FP), Z1			// load the flags into z reg

	MOVB 			Deleted, R9				// load the deleted constant into reg
	VPBROADCASTB	R9, Z3					// broadcast deleted flag

	MOVB 			key+0(FP), R10			// load the key constant into reg
	VPBROADCASTB	R10, Z4					// broadcast key
	VMOVDQU8 		slotKeys+65(FP), Z5		// load the slot keys into z reg

	VPANDQ			Z0, Z1, Z2				// used & flags
	VPCMPEQB        Z0, Z2, K1				// (used & flags) == used
	KMOVQ			K1, R8					// get the final used bit flags
	XORQ			R12, R12				// clear register
	ADDQ			R8, R12					// cpy register
	NOTQ			R12						// inverse of (used & flags) == used

	VPANDQ			Z3, Z1, Z2				// deleted & flags
	VPCMPEQB        Z3, Z2, K1				// (deleted & flags) == deleted
	KMOVQ			K1, R9					// get the final deleted bit flags
	NOTQ			R9						// invert the bit mask

	VPCMPEQB        Z4, Z5, K1				// slotKey == key
	KMOVQ			K1, R10					// get the final slot key bit flags

	ANDQ			R8, R9					// ((used & flags) == used) & ((deleted & flags) == deleted)
	ANDQ			R9, R10					// ((used & flags) == used) & ((deleted & flags) == deleted) & (slotKey == key)

	MOVQ			R10, rv+136(FP)			// place the final bit field result in mem
	MOVQ			R12, rv+144(FP)			// place the is empty bit field result in mem
	RET
