package emu

import (
	"log"

	"github.com/sarchlab/mgpusim/v3/insts"
)

func (u *ALUImpl) runDS(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 13:
		u.runDSWRITEB32(state)
	case 14:
		u.runDSWRITE2B32(state)
	case 54:
		u.runDSREADB32(state)
	// case 55:
	// 	u.runDSREAD2B32(state)
	case 78:
		u.runDSWRITE2B64(state)
	case 118:
		u.runDSREADB64(state)
	case 119:
		u.runDSREAD2B64(state)
	default:
		log.Panicf("Opcode %d for DS format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runDSWRITEB32(state InstEmuState) {
	inst := state.Inst()
	lds := u.LDS()
	exec := state.ReadReg(insts.Regs[insts.EXEC], 1, 0)

	for i := 0; i < 64; i++ {
		if !laneMasked(exec, uint(i)) {
			continue
		}
		addr := u.ReadOperand(state, inst.Addr, i, nil)

		addr0 := uint32(addr) + inst.Offset0
		data := u.ReadOperand(state, inst.Data, i, nil)

		copy(lds[addr0:addr0+4], insts.Uint32ToBytes(uint32(data)))
	}
}

func (u *ALUImpl) runDSWRITE2B32(state InstEmuState) {
	inst := state.Inst()
	lds := u.LDS()
	exec := state.ReadReg(insts.Regs[insts.EXEC], 1, 0)

	for i := 0; i < 64; i++ {
		if !laneMasked(exec, uint(i)) {
			continue
		}
		addr := u.ReadOperand(state, inst.Addr, i, nil)

		addr0 := uint32(addr) + inst.Offset0*4
		addr1 := uint32(addr) + inst.Offset1*4
		data := u.ReadOperand(state, inst.Data, i, nil)
		data1 := u.ReadOperand(state, inst.Data1, i, nil)

		copy(lds[addr0:addr0+4], insts.Uint32ToBytes(uint32(data)))
		copy(lds[addr1:addr1+4], insts.Uint32ToBytes(uint32(data1)))
	}
}

func (u *ALUImpl) runDSREADB32(state InstEmuState) {
	inst := state.Inst()
	lds := u.LDS()
	exec := state.ReadReg(insts.Regs[insts.EXEC], 1, 0)

	for i := 0; i < 64; i++ {
		if !laneMasked(exec, uint(i)) {
			continue
		}
		addr := u.ReadOperand(state, inst.Addr, i, nil)

		addr0 := uint32(addr) + inst.Offset0

		u.WriteOperand(state, inst.Dst, i, uint64(insts.BytesToUint32(lds[addr0:addr0+4])), nil)
	}
}

// func (u *ALUImpl) runDSREAD2B32(state InstEmuState) {
// 	inst := state.Inst()
// 	sp := state.Scratchpad()
// 	layout := sp.AsDS()
// 	lds := u.LDS()

// 	i := uint(0)
// 	for i = 0; i < 64; i++ {
// 		if !laneMasked(layout.EXEC, i) {
// 			continue
// 		}

// 		addr0 := layout.ADDR[i] + inst.Offset0*4
// 		dstOffset := uint(8 + 64*4 + 256*4*2)
// 		copy(sp[dstOffset+i*16:dstOffset+i*16+4], lds[addr0:addr0+4])

// 		addr1 := layout.ADDR[i] + inst.Offset1*4
// 		copy(sp[dstOffset+i*16+4:dstOffset+i*16+8], lds[addr1:addr1+4])
// 	}
// }

func (u *ALUImpl) runDSWRITE2B64(state InstEmuState) {
	inst := state.Inst()
	lds := u.LDS()
	exec := state.ReadReg(insts.Regs[insts.EXEC], 1, 0)

	for i := 0; i < 64; i++ {
		if !laneMasked(exec, uint(i)) {
			continue
		}

		addr := u.ReadOperand(state, inst.Addr, i, nil)
		addr0 := uint32(addr) + inst.Offset0*8
		addr1 := uint32(addr) + inst.Offset1*8
		data := u.ReadOperand(state, inst.Data, i, nil)
		data1 := u.ReadOperand(state, inst.Data1, i, nil)

		copy(lds[addr0:addr0+8], insts.Uint64ToBytes(data))
		copy(lds[addr1:addr1+8], insts.Uint64ToBytes(data1))

		// addr0 := layout.ADDR[i] + inst.Offset0*8
		// data0Offset := uint(8 + 64*4)
		// copy(lds[addr0:addr0+8], sp[data0Offset+i*16:data0Offset+i*16+8])

		// addr1 := layout.ADDR[i] + inst.Offset1*8
		// data1Offset := uint(8 + 64*4 + 256*4)
		// copy(lds[addr1:addr1+8], sp[data1Offset+i*16:data1Offset+i*16+8])
	}
}

func (u *ALUImpl) runDSREADB64(state InstEmuState) {
	inst := state.Inst()
	lds := u.LDS()
	exec := state.ReadReg(insts.Regs[insts.EXEC], 1, 0)

	for i := 0; i < 64; i++ {
		if !laneMasked(exec, uint(i)) {
			continue
		}
		addr := u.ReadOperand(state, inst.Addr, i, nil)

		u.WriteOperand(state, inst.Dst, i, insts.BytesToUint64(lds[addr:addr+8]), nil)
	}
}

func (u *ALUImpl) runDSREAD2B64(state InstEmuState) {
	inst := state.Inst()
	lds := u.LDS()
	exec := state.ReadReg(insts.Regs[insts.EXEC], 1, 0)

	for i := 0; i < 64; i++ {
		if !laneMasked(exec, uint(i)) {
			continue
		}
		addr := u.ReadOperand(state, inst.Addr, i, nil)
		addr0 := uint32(addr) + inst.Offset0*8
		addr1 := uint32(addr) + inst.Offset1*8

		dst := uint64(insts.BytesToUint64(lds[addr0:addr0+8])) | uint64(insts.BytesToUint64(lds[addr1:addr1+8]))
		u.WriteOperand(state, inst.Dst, i, dst, nil)

		// addr0 := layout.ADDR[i] + inst.Offset0*8
		// dstOffset := uint(8 + 64*4 + 256*4*2)
		// copy(sp[dstOffset+i*16:dstOffset+i*16+8], lds[addr0:addr0+8])

		// addr1 := layout.ADDR[i] + inst.Offset1*8
		// copy(sp[dstOffset+i*16+8:dstOffset+i*16+16], lds[addr1:addr1+8])
	}
}
