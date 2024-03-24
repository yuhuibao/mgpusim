package emu

import (
	"log"

	"github.com/sarchlab/mgpusim/v3/insts"
)

func (u *ALUImpl) runSOPK(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0:
		u.runSMOVKI32(state)
	case 3:
		u.runSCMPKLGI32(state)
	case 15:
		u.runSMULKI32(state)
	default:
		log.Panicf("Opcode %d for SOPK format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runSMOVKI32(state InstEmuState) {
	inst := state.Inst()
	imm := int16(u.ReadOperand(state, inst.SImm16, 0, nil))
	u.WriteOperand(state, inst.Dst, 0, uint64(imm), nil)
}

func (u *ALUImpl) runSCMPKLGI32(state InstEmuState) {
	inst := state.Inst()
	imm := int16(u.ReadOperand(state, inst.SImm16, 0, nil))
	dst := int16(u.ReadOperand(state, inst.Dst, 0, nil))
	if dst != imm {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	}
}

func (u *ALUImpl) runSMULKI32(state InstEmuState) {
	inst := state.Inst()
	imm := int16(u.ReadOperand(state, inst.SImm16, 0, nil))
	dst := int32(u.ReadOperand(state, inst.Dst, 0, nil))
	u.WriteOperand(state, inst.Dst, 0, uint64(int64(int32(imm)*dst)), nil)
}
