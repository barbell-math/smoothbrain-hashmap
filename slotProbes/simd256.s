//go:build amd64 && simd256

#include "textflag.h"

// !!! Note !!!
// These constants must be kept in sync with the constants in the
// AllBuildTargets file
#define Used $1
#define Deleted $2

// func getSlotProbe(
// 		key uint8,
// 		flags [32]uint8,
// 		slotKeys [32]uint8,
// ) (potentialValues uint32, isEmpty uint32)
//
// memory layout of the stack relative to FP
//  +0   				key					argument
//  +1  through +32 	flags				argument
//  +33 through +64		slotKeys			argument
//  +65 through +71		-					allignment padding
//  +72 through +75		potentialValues 	return value
//  +76 through +79		isEmpty 			return value
TEXT ·SlotProbe(SB),NOSPLIT,$0
	MOVB 			Used, R8				// load the used constant into reg
	VPBROADCASTB	R8, Y0					// broadcast used flag
	VMOVDQU8 		flags+1(FP), Y1			// load the flags into y reg
	MOVB 			Deleted, R9				// load the deleted constant into reg
	VPBROADCASTB	R9, Y3					// broadcast deleted flag

	VPAND			Y0, Y1, Y2				// used & flags
	VPCMPEQB        Y0, Y2, Y2				// (used & flags) == used
	VEXTRACTI128	$1, Y2, X0				// Get the upper half of (used & flags) == used
	PMOVMSKB		X0, R8					// store the sign flags in reg
	SHLQ			$16, R8					// shift bits to left half of reg
	VEXTRACTI128	$2, Y2, X0				// Get the lower half of (used & flags) == used
	PMOVMSKB		X0, R9					// store the sign flags in reg
	ORQ				R9, R8					// get the final used bit flags
	XORQ			R12, R12				// clear register
	ADDQ			R8, R12					// cpy register
	NOTQ			R12						// inverse of (used & flags) == used

	VPAND			Y3, Y1, Y2				// deleted & flags
	VPCMPEQB        Y3, Y2, Y2				// (deleted & flags) == deleted
	VEXTRACTI128	$1, Y2, X0				// Get the upper half of (deleted & flags) == deleted
	PMOVMSKB		X0, R9					// store the sign flags in reg
	SHLQ			$16, R9					// shift bits to left half of reg
	VEXTRACTI128	$1, Y2, X0				// Get the lower half of (used & flags) == used
	PMOVMSKB		X2, R10					// store the sign flags in reg
	ORQ				R10, R9					// get the final deleted bit flags
	NOTQ			R9						// invert the bit mask

	MOVQ 			key+0(FP), R10			// load the key constant into reg
	VPBROADCASTB	R10, Y0					// broadcast key
	VMOVDQU8 		slotKeys+33(FP), Y1		// load the slot keys into y reg
	VPCMPEQB        Y0, Y1, Y2				// slotKey == key
	VEXTRACTI128	$1, Y2, X0				// Get the upper half of slotKey == key
	PMOVMSKB		X0, R10					// store the sign flags in reg
	SHLQ			$16, R10				// shift bits to left half of reg
	VEXTRACTI128	$2, Y2, X0				// Get the lower half of slotKey == key
	PMOVMSKB		X0, R11					// store the sign flags in reg
	ORQ				R11, R10				// get the final slot key bit flags

	ANDQ			R8, R9					// ((used & flags) == used) & ((deleted & flags) == deleted)
	ANDQ			R9, R10					// ((used & flags) == used) & ((deleted & flags) == deleted) & (slotKey == key)
	// TODO - why overwriting next value???
	MOVD			R10, rv+72(FP)			// place the final bit field result in mem
	MOVD			R12, rv+76(FP)			// place the is empty bit field result in mem
	RET
