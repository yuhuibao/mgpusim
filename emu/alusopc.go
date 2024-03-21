package emu

import (
	"log"

	"github.com/sarchlab/mgpusim/v3/insts"
)

//nolint:gocyclo,funlen
func (u *ALUImpl) runSOPC(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0:
		u.runSCMPEQU32(state)
	case 1:
		u.runSCMPLGU32(state)
	case 2:
		u.runSCMPGTI32(state)
	case 3:
		u.runSCMPGEI32(state)
	case 4:
		u.runSCMPLTI32(state)
	case 5:
		u.runSCMPLEI32(state)
	case 6:
		u.runSCMPEQU32(state)
	case 7:
		u.runSCMPLGU32(state)
	case 8:
		u.runSCMPGTU32(state)
	case 10:
		u.runSCMPLTU32(state)
	default:
		log.Panicf("Opcode %d for SOPC format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runSCMPGTI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := asInt32(uint32(u.ReadOperand(state, inst.Src1, 0, nil)))

	if src0 > src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSCMPLTI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := asInt32(uint32(u.ReadOperand(state, inst.Src1, 0, nil)))

	if src0 < src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSCMPLEI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := asInt32(uint32(u.ReadOperand(state, inst.Src1, 0, nil)))

	if src0 <= src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSCMPGEI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := asInt32(uint32(u.ReadOperand(state, inst.Src1, 0, nil)))

	if src0 >= src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSCMPEQU32(state InstEmuState) {
	inst := state.Inst()
	src0 := uint32(u.ReadOperand(state, inst.Src0, 0, nil))
	src1 := uint32(u.ReadOperand(state, inst.Src1, 0, nil))

	if src0 == src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSCMPLGU32(state InstEmuState) {
	inst := state.Inst()
	src0 := uint32(u.ReadOperand(state, inst.Src0, 0, nil))
	src1 := uint32(u.ReadOperand(state, inst.Src1, 0, nil))

	if src0 != src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSCMPGTU32(state InstEmuState) {
	inst := state.Inst()
	src0 := uint32(u.ReadOperand(state, inst.Src0, 0, nil))
	src1 := uint32(u.ReadOperand(state, inst.Src1, 0, nil))

	if src0 > src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSCMPLTU32(state InstEmuState) {
	inst := state.Inst()
	src0 := uint32(u.ReadOperand(state, inst.Src0, 0, nil))
	src1 := uint32(u.ReadOperand(state, inst.Src1, 0, nil))

	if src0 < src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}
