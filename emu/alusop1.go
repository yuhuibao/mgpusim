package emu

import (
	"log"
)

//nolint:gocyclo
func (u *ALUImpl) runSOP1(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0:
		u.runSMOVB32(state)
	// case 1:
	// 	u.runSMOVB64(state)
	// case 4:
	// 	u.runSNOTU32(state)
	// case 8:
	// 	u.runSBREVB32(state)
	// case 28:
	// 	u.runSGETPCB64(state)
	// case 32:
	// 	u.runSANDSAVEEXECB64(state)
	// case 33:
	// 	u.runSORSAVEEXECB64(state)
	// case 34:
	// 	u.runSXORSAVEEXECB64(state)
	// case 35:
	// 	u.runSANDN2SAVEEXECB64(state)
	// case 36:
	// 	u.runSORN2SAVEEXECB64(state)
	// case 37:
	// 	u.runSNANDSAVEEXECB64(state)
	// case 38:
	// 	u.runSNORSAVEEXECB64(state)
	// case 39:
	// 	u.runSNXORSAVEEXECB64(state)
	default:
		log.Panicf("Opcode %d for SOP1 format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runSMOVB32(state InstEmuState) {
	inst := state.Inst()

	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	u.WriteOperand(state, inst.Dst, 0, src0, nil)
}

// func (u *ALUImpl) runSMOVB64(state InstEmuState) {
// 	inst := state.Inst()
// 	src0 := state.ReadOperand(inst.Src0, 0, nil)
// 	state.WriteOperand(inst.Dst, 0, src0, nil)
// }

// func (u *ALUImpl) runSNOTU32(state InstEmuState) {
// 	inst := state.Inst()
// 	src0 := state.ReadOperand(inst.Src0, 0, nil)
// 	state.WriteOperand(inst.Dst, 0, ^src0, nil)
// 	if ^src0 != 0 {
// 		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
// 	}
// }

// func (u *ALUImpl) runSBREVB32(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	dst := uint32(0)
// 	for i := 0; i < 32; i++ {
// 		bit := uint32(1 << (31 - i))
// 		bit = uint32(sp.SRC0) & bit
// 		bit = bit >> (31 - i)
// 		bit = bit << i
// 		dst = dst | bit
// 	}
// 	sp.DST = uint64(dst)
// }

// func (u *ALUImpl) runSGETPCB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.PC + 4
// }

// func (u *ALUImpl) runSANDSAVEEXECB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.EXEC
// 	sp.EXEC = sp.SRC0 & sp.EXEC
// 	if sp.EXEC != 0 {
// 		sp.SCC = 1
// 	} else {
// 		sp.SCC = 0
// 	}
// }

// func (u *ALUImpl) runSORSAVEEXECB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.EXEC
// 	sp.EXEC = sp.SRC0 | sp.EXEC
// 	if sp.EXEC != 0 {
// 		sp.SCC = 1
// 	} else {
// 		sp.SCC = 0
// 	}
// }

// func (u *ALUImpl) runSXORSAVEEXECB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.EXEC
// 	sp.EXEC = sp.SRC0 ^ sp.EXEC
// 	if sp.EXEC != 0 {
// 		sp.SCC = 1
// 	} else {
// 		sp.SCC = 0
// 	}
// }

// func (u *ALUImpl) runSANDN2SAVEEXECB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.EXEC
// 	sp.EXEC = sp.SRC0 & (^sp.EXEC)
// 	if sp.EXEC != 0 {
// 		sp.SCC = 1
// 	} else {
// 		sp.SCC = 0
// 	}
// }

// func (u *ALUImpl) runSORN2SAVEEXECB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.EXEC
// 	sp.EXEC = sp.SRC0 | (^sp.EXEC)
// 	if sp.EXEC != 0 {
// 		sp.SCC = 1
// 	} else {
// 		sp.SCC = 0
// 	}
// }

// func (u *ALUImpl) runSNANDSAVEEXECB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.EXEC
// 	sp.EXEC = ^(sp.SRC0 & sp.EXEC)
// 	if sp.EXEC != 0 {
// 		sp.SCC = 1
// 	} else {
// 		sp.SCC = 0
// 	}
// }

// func (u *ALUImpl) runSNORSAVEEXECB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.EXEC
// 	sp.EXEC = ^(sp.SRC0 | sp.EXEC)
// 	if sp.EXEC != 0 {
// 		sp.SCC = 1
// 	} else {
// 		sp.SCC = 0
// 	}
// }

// func (u *ALUImpl) runSNXORSAVEEXECB64(state InstEmuState) {
// 	sp := state.Scratchpad().AsSOP1()
// 	sp.DST = sp.EXEC
// 	sp.EXEC = ^(sp.SRC0 ^ sp.EXEC)
// 	if sp.EXEC != 0 {
// 		sp.SCC = 1
// 	} else {
// 		sp.SCC = 0
// 	}
// }
