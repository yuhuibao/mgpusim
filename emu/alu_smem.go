package emu

import (
	"log"

	"github.com/sarchlab/mgpusim/v3/insts"
)

func (u *ALUImpl) runSMEM(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0:
		u.runSLOADDWORD(state)
	case 1:
		u.runSLOADDWORDX2(state)
	case 2:
		u.runSLOADDWORDX4(state)
	case 3:
		u.runSLOADDWORDX8(state)
	default:
		log.Panicf("Opcode %d for SMEM format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runSLOADDWORD(state InstEmuState) {
	inst := state.Inst()
	offset := u.ReadOperand(state, inst.Offset, 0, nil)
	base := u.ReadOperand(state, inst.Base, 0, nil)
	pid := state.PID()

	buf := u.storageAccessor.Read(pid, base+offset, 4)

	u.WriteOperand(state, inst.Data, 0, uint64(insts.BytesToUint32(buf)), nil)
}

func (u *ALUImpl) runSLOADDWORDX2(state InstEmuState) {
	inst := state.Inst()
	offset := u.ReadOperand(state, inst.Offset, 0, nil)
	base := u.ReadOperand(state, inst.Base, 0, nil)
	pid := state.PID()

	buf := u.storageAccessor.Read(pid, base+offset, 8)
	u.WriteOperand(state, inst.Data, 0, insts.BytesToUint64(buf), nil)
}

func (u *ALUImpl) runSLOADDWORDX4(state InstEmuState) {
	inst := state.Inst()
	offset := u.ReadOperand(state, inst.Offset, 0, nil)
	base := u.ReadOperand(state, inst.Base, 0, nil)
	pid := state.PID()

	buf := u.storageAccessor.Read(pid, base+offset, 16)
	var buffer []uint32
	for i := 0; i < 16; i += 4 {
		num := uint32(buf[i]) | uint32(buf[i+1])<<8 | uint32(buf[i+2])<<16 | uint32(buf[i+3])<<24
		buffer = append(buffer, num)
	}
	u.WriteOperand(state, inst.Data, 0, 0, buffer)
}

func (u *ALUImpl) runSLOADDWORDX8(state InstEmuState) {
	inst := state.Inst()
	offset := u.ReadOperand(state, inst.Offset, 0, nil)
	base := u.ReadOperand(state, inst.Base, 0, nil)
	pid := state.PID()

	buf := u.storageAccessor.Read(pid, base+offset, 32)
	var buffer []uint32
	for i := 0; i < 32; i += 4 {
		num := uint32(buf[i]) | uint32(buf[i+1])<<8 | uint32(buf[i+2])<<16 | uint32(buf[i+3])<<24
		buffer = append(buffer, num)
	}
	u.WriteOperand(state, inst.Data, 0, 0, buffer)
}

// //nolint:gocyclo
func (u *ALUImpl) runSOPP(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0: // S_NOP
	// Do nothing
	case 2: // S_CBRANCH
		u.runSCBRANCH(state)
	case 4: // S_CBRANCH_SCC0
		u.runSCBRANCHSCC0(state)
	case 5: // S_CBRANCH_SCC1
		u.runSCBRANCHSCC1(state)
	case 6: // S_CBRANCH_VCCZ
		u.runSCBRANCHVCCZ(state)
	case 7: // S_CBRANCH_VCCNZ
		u.runSCBRANCHVCCNZ(state)
	case 8: // S_CBRANCH_EXECZ
		u.runSCBRANCHEXECZ(state)
	case 9: // S_CBRANCH_EXECNZ
		u.runSCBRANCHEXECNZ(state)
	case 12: // S_WAITCNT
	// Do nothing
	default:
		log.Panicf("Opcode %d for SOPP format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runSCBRANCH(state InstEmuState) {
	inst := state.Inst()
	imm := int16(uint16(u.ReadOperand(state, inst.SImm16, 0, nil) & 0xffff))
	pc := uint64(int64(state.ReadReg(insts.Regs[insts.PC], 1, 0)) + int64(imm)*4)
	state.WriteReg(insts.Regs[insts.PC], 1, 0, pc)
}

func (u *ALUImpl) runSCBRANCHSCC0(state InstEmuState) {
	inst := state.Inst()
	imm := int16(uint16(u.ReadOperand(state, inst.SImm16, 0, nil) & 0xffff))

	scc := state.ReadReg(insts.Regs[insts.SCC], 1, 0)
	if scc == 0 {
		pc := uint64(int64(state.ReadReg(insts.Regs[insts.PC], 1, 0)) + int64(imm)*4)
		state.WriteReg(insts.Regs[insts.PC], 1, 0, pc)
	}
}

func (u *ALUImpl) runSCBRANCHSCC1(state InstEmuState) {
	inst := state.Inst()
	imm := int16(uint16(u.ReadOperand(state, inst.SImm16, 0, nil) & 0xffff))

	scc := state.ReadReg(insts.Regs[insts.SCC], 1, 0)
	if scc == 1 {
		pc := uint64(int64(state.ReadReg(insts.Regs[insts.PC], 1, 0)) + int64(imm)*4)
		state.WriteReg(insts.Regs[insts.PC], 1, 0, pc)
	}
}

func (u *ALUImpl) runSCBRANCHVCCZ(state InstEmuState) {
	inst := state.Inst()
	imm := int16(uint16(u.ReadOperand(state, inst.SImm16, 0, nil) & 0xffff))

	vcc := state.ReadReg(insts.Regs[insts.VCC], 1, 0)
	if vcc == 0 {
		pc := uint64(int64(state.ReadReg(insts.Regs[insts.PC], 1, 0)) + int64(imm)*4)
		state.WriteReg(insts.Regs[insts.PC], 1, 0, pc)
	}
}

func (u *ALUImpl) runSCBRANCHVCCNZ(state InstEmuState) {
	inst := state.Inst()
	imm := int16(uint16(u.ReadOperand(state, inst.SImm16, 0, nil) & 0xffff))

	vcc := state.ReadReg(insts.Regs[insts.VCC], 1, 0)
	if vcc != 0 {
		pc := uint64(int64(state.ReadReg(insts.Regs[insts.PC], 1, 0)) + int64(imm)*4)
		state.WriteReg(insts.Regs[insts.PC], 1, 0, pc)
	}
}

func (u *ALUImpl) runSCBRANCHEXECZ(state InstEmuState) {
	inst := state.Inst()
	imm := int16(uint16(u.ReadOperand(state, inst.SImm16, 0, nil) & 0xffff))

	exec := state.ReadReg(insts.Regs[insts.EXEC], 1, 0)
	if exec == 0 {
		pc := uint64(int64(state.ReadReg(insts.Regs[insts.PC], 1, 0)) + int64(imm)*4)
		state.WriteReg(insts.Regs[insts.PC], 1, 0, pc)
	}
}

func (u *ALUImpl) runSCBRANCHEXECNZ(state InstEmuState) {
	inst := state.Inst()
	imm := int16(uint16(u.ReadOperand(state, inst.SImm16, 0, nil) & 0xffff))

	exec := state.ReadReg(insts.Regs[insts.EXEC], 1, 0)
	if exec != 0 {
		pc := uint64(int64(state.ReadReg(insts.Regs[insts.PC], 1, 0)) + int64(imm)*4)
		state.WriteReg(insts.Regs[insts.PC], 1, 0, pc)
	}
}

func laneMasked(Exec uint64, laneID uint) bool {
	return Exec&(1<<laneID) > 0
}

// func (u *ALUImpl) sdwaSrcSelect(src uint32, sel insts.SDWASelect) uint32 {
// 	switch sel {
// 	case insts.SDWASelectByte0:
// 		return src & 0x000000ff
// 	case insts.SDWASelectByte1:
// 		return (src & 0x0000ff00) >> 8
// 	case insts.SDWASelectByte2:
// 		return (src & 0x00ff0000) >> 16
// 	case insts.SDWASelectByte3:
// 		return (src & 0xff000000) >> 24
// 	case insts.SDWASelectWord0:
// 		return src & 0x0000ffff
// 	case insts.SDWASelectWord1:
// 		return (src & 0xffff0000) >> 16
// 	case insts.SDWASelectDWord:
// 		return src
// 	}
// 	return src
// }

// func (u *ALUImpl) sdwaDstSelect(
// 	dstOld uint32,
// 	dstNew uint32,
// 	sel insts.SDWASelect,
// 	unused insts.SDWAUnused,
// ) uint32 {
// 	value := dstNew
// 	switch sel {
// 	case insts.SDWASelectByte0:
// 		value = value & 0x000000ff
// 	case insts.SDWASelectByte1:
// 		value = (value << 8) & 0x0000ff00
// 	case insts.SDWASelectByte2:
// 		value = (value << 16) & 0x00ff0000
// 	case insts.SDWASelectByte3:
// 		value = (value << 24) & 0xff000000
// 	case insts.SDWASelectWord0:
// 		value = value & 0x0000ffff
// 	case insts.SDWASelectWord1:
// 		value = (value << 16) & 0xffff0000
// 	}

// 	return value
// }

// func (u *ALUImpl) dumpScratchpadAsSop2(state InstEmuState, byteCount int) string {
// 	scratchpad := state.Scratchpad()
// 	layout := new(SOP2Layout)

// 	binary.Read(bytes.NewBuffer(scratchpad), binary.LittleEndian, layout)

// 	output := fmt.Sprintf(
// 		`
// 			SRC0: 0x%[1]x(%[1]d),
// 			SRC1: 0x%[2]x(%[2]d),
// 			SCC: 0x%[3]x(%[3]d),
// 			DST: 0x%[4]x(%[4]d)\n",
// 		`,
// 		layout.SRC0, layout.SRC1, layout.SCC, layout.DST)

// 	return output
// }
