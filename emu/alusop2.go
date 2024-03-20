package emu

import (
	"log"
	"math"

	"github.com/sarchlab/mgpusim/v3/bitops"
	"github.com/sarchlab/mgpusim/v3/insts"
)

//nolint:gocyclo,funlen
func (u *ALUImpl) runSOP2(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0:
		u.runSADDU32(state)
	case 1:
		u.runSSUBU32(state)
	case 2:
		u.runSADDI32(state)
	case 3:
		u.runSSUBI32(state)
	case 4:
		u.runSADDCU32(state)
	case 5:
		u.runSSUBBU32(state)
	case 6:
		u.runSMINI32(state)
	case 7:
		u.runSMINU32(state)
	case 8:
		u.runSMAXI32(state)
	case 9:
		u.runSMAXU32(state)
	case 10:
		u.runSCSELECTB32(state)
	case 12:
		u.runSANDB32(state)
	case 13:
		u.runSANDB64(state)
	case 15:
		u.runSORB64(state)
	case 16, 17:
		u.runSXORB64(state)
	case 19:
		u.runSANDN2B64(state)
	case 28:
		u.runSLSHLB32(state)
	case 29:
		u.runSLSHLB64(state)
	case 30:
		u.runSLSHRB32(state)
	case 31:
		u.runSLSHRB64(state)
	case 32:
		u.runSASHRI32(state)
	case 34:
		u.runSBFMB32(state)
	case 36:
		u.runSMULI32(state)
	case 38:
		u.runSBFEI32(state)
	default:
		log.Panicf("Opcode %d for SOP2 format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runSADDU32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)

	dst := src0 + src1
	if src0 > math.MaxUint32-src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}

	u.WriteOperand(state, inst.Dst, 0, dst, nil)
}

func (u *ALUImpl) runSSUBU32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)

	if src0 < src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	}
	u.WriteOperand(state, inst.Dst, 0, src0-src1, nil)
}

func (u *ALUImpl) runSADDI32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)

	if src0 > math.MaxUint32-src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
	u.WriteOperand(state, inst.Dst, 0, src0+src1, nil)
}

func (u *ALUImpl) runSSUBI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := asInt32(uint32(u.ReadOperand(state, inst.Src1, 0, nil)))

	dst := src0 - src1

	if src1 > 0 && dst > src0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else if src1 < 0 && dst < src0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}

	u.WriteOperand(state, inst.Dst, 0, uint64(int32ToBits(dst)), nil)
}

func (u *ALUImpl) runSADDCU32(state InstEmuState) {
	inst := state.Inst()
	src0 := uint32(u.ReadOperand(state, inst.Src0, 0, nil))
	src1 := uint32(u.ReadOperand(state, inst.Src1, 0, nil))

	scc := uint32(state.ReadReg(insts.Regs[insts.SCC], 1, 0))

	u.WriteOperand(state, inst.Dst, 0, uint64(src0+src1+scc), nil)
	if src0 < math.MaxUint32-uint32(scc)-src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	}
}

func (u *ALUImpl) runSSUBBU32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	scc := state.ReadReg(insts.Regs[insts.SCC], 1, 0)

	u.WriteOperand(state, inst.Dst, 0, src0-src1-scc, nil)

	if src0 < src1+scc {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSMINI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := asInt32(uint32(u.ReadOperand(state, inst.Src1, 0, nil)))

	if src0 < src1 {
		u.WriteOperand(state, inst.Dst, 0, uint64(src0), nil)
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		u.WriteOperand(state, inst.Dst, 0, uint64(src1), nil)
	}
}

func (u *ALUImpl) runSMINU32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)

	if src0 < src1 {
		u.WriteOperand(state, inst.Dst, 0, src0, nil)
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		u.WriteOperand(state, inst.Dst, 0, src1, nil)
	}
}

func (u *ALUImpl) runSMAXI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := asInt32(uint32(u.ReadOperand(state, inst.Src1, 0, nil)))

	if src0 > src1 {
		u.WriteOperand(state, inst.Dst, 0, uint64(src0), nil)
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		u.WriteOperand(state, inst.Dst, 0, uint64(src1), nil)
	}
}

func (u *ALUImpl) runSMAXU32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)

	if src0 > src1 {
		u.WriteOperand(state, inst.Dst, 0, src0, nil)
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		u.WriteOperand(state, inst.Dst, 0, src1, nil)
	}
}

func (u *ALUImpl) runSCSELECTB32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	scc := state.ReadReg(insts.Regs[insts.SCC], 1, 0)

	if scc == 1 {
		u.WriteOperand(state, inst.Dst, 0, src0, nil)
	} else {
		u.WriteOperand(state, inst.Dst, 0, src1, nil)
	}
}

func (u *ALUImpl) runSANDB32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	dst := src0 & src1
	u.WriteOperand(state, inst.Dst, 0, dst, nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSANDB64(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	dst := src0 & src1
	u.WriteOperand(state, inst.Dst, 0, dst, nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSORB64(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	dst := src0 | src1
	u.WriteOperand(state, inst.Dst, 0, dst, nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSXORB64(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	dst := src0 ^ src1
	u.WriteOperand(state, inst.Dst, 0, dst, nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSANDN2B64(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	dst := src0 &^ src1
	u.WriteOperand(state, inst.Dst, 0, dst, nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSLSHLB32(state InstEmuState) {
	inst := state.Inst()
	src0 := uint32(u.ReadOperand(state, inst.Src0, 0, nil))
	src1 := uint8(u.ReadOperand(state, inst.Src1, 0, nil))
	dst := src0 << (src1 & 0x1f)
	u.WriteOperand(state, inst.Dst, 0, uint64(dst), nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSLSHLB64(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := uint8(u.ReadOperand(state, inst.Src1, 0, nil))
	dst := src0 << (src1 & 0x3f)
	u.WriteOperand(state, inst.Dst, 0, dst, nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSLSHRB32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	dst := src0 >> (src1 & 0x1f)
	u.WriteOperand(state, inst.Dst, 0, dst, nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSLSHRB64(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	dst := src0 >> (src1 & 0x3f)
	u.WriteOperand(state, inst.Dst, 0, dst, nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSASHRI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := uint8(u.ReadOperand(state, inst.Src1, 0, nil))
	dst := src0 >> src1
	u.WriteOperand(state, inst.Dst, 0, uint64(int32ToBits(dst)), nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}

func (u *ALUImpl) runSBFMB32(state InstEmuState) {
	inst := state.Inst()
	src0 := u.ReadOperand(state, inst.Src0, 0, nil)
	src1 := u.ReadOperand(state, inst.Src1, 0, nil)
	dst := ((1 << (src0 & 0x1f)) - 1) << (src1 & 0x1f)
	u.WriteOperand(state, inst.Dst, 0, uint64(dst), nil)
}

func (u *ALUImpl) runSMULI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := asInt32(uint32(u.ReadOperand(state, inst.Src1, 0, nil)))
	dst := src0 * src1
	u.WriteOperand(state, inst.Dst, 0, uint64(int32ToBits(dst)), nil)

	if src0 != 0 && dst/src0 != src1 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	}
}

func (u *ALUImpl) runSBFEI32(state InstEmuState) {
	inst := state.Inst()
	src0 := asInt32(uint32(u.ReadOperand(state, inst.Src0, 0, nil)))
	src1 := uint32(u.ReadOperand(state, inst.Src1, 0, nil))
	offset := bitops.ExtractBitsFromU32(src1, 0, 4)
	width := bitops.ExtractBitsFromU32(src1, 16, 22)
	dst := (src0 >> offset) & ((1 << width) - 1)
	u.WriteOperand(state, inst.Dst, 0, uint64(int32ToBits(dst)), nil)

	if dst != 0 {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)
	} else {
		state.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)
	}
}
