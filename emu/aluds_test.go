package emu

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sarchlab/mgpusim/v3/insts"
	"github.com/sarchlab/mgpusim/v3/kernels"
)

var _ = Describe("ALU", func() {

	var (
		alu *ALUImpl
		wf  *Wavefront
	)

	BeforeEach(func() {
		alu = NewALU(nil)
		alu.lds = make([]byte, 4096)

		rawWf := kernels.NewWavefront()
		wf = NewWavefront(rawWf)
	})

	It("should run DS_WRITE_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.DS
		inst.Opcode = 13
		inst.Offset0 = 0
		inst.Addr = insts.NewVRegOperand(0, 0, 1)
		inst.Data = insts.NewVRegOperand(1, 1, 1)
		wf.inst = inst

		wf.Exec = 0x01
		wf.WriteReg(insts.VReg(0), 1, 0, 100)
		wf.WriteReg(insts.VReg(1), 1, 0, 1)

		alu.Run(wf)

		lds := alu.LDS()
		Expect(insts.BytesToUint32(lds[100:])).To(Equal(uint32(1)))
	})

	It("should run DS_WRITE2_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.DS
		inst.Opcode = 14
		inst.Offset0 = 0
		inst.Offset1 = 4
		inst.Addr = insts.NewVRegOperand(0, 0, 1)
		inst.Data = insts.NewVRegOperand(1, 1, 1)
		inst.Data1 = insts.NewVRegOperand(2, 2, 1)
		wf.inst = inst

		wf.Exec = 0x01
		wf.WriteReg(insts.VReg(0), 1, 0, 100)
		wf.WriteReg(insts.VReg(1), 1, 0, 1)
		wf.WriteReg(insts.VReg(2), 1, 0, 2)

		alu.Run(wf)

		lds := alu.LDS()
		Expect(insts.BytesToUint32(lds[100:])).To(Equal(uint32(1)))
		Expect(insts.BytesToUint32(lds[116:])).To(Equal(uint32(2)))
	})

	It("should run DS_READ_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.DS
		inst.Opcode = 54
		inst.Offset0 = 0
		inst.Addr = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(1, 1, 1)
		wf.inst = inst

		wf.Exec = 0x1
		wf.WriteReg(insts.VReg(0), 1, 0, 100)
		lds := alu.LDS()
		copy(lds[100:], insts.Uint32ToBytes(12))

		alu.Run(wf)

		result := wf.ReadReg(insts.VReg(1), 1, 0)
		Expect(uint32(result)).To(Equal(uint32(12)))
	})

	// 	It("should run DS_READ2_B32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.DS
	// 		state.inst.Opcode = 55
	// 		state.inst.Offset0 = 0
	// 		state.inst.Offset1 = 4

	// 		sp := state.scratchpad.AsDS()
	// 		sp.EXEC = 0x1
	// 		sp.ADDR[0] = 100

	// 		lds := alu.LDS()
	// 		copy(lds[100:], insts.Uint32ToBytes(1))
	// 		copy(lds[116:], insts.Uint32ToBytes(2))

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint32(1)))
	// 		Expect(sp.DST[1]).To(Equal(uint32(2)))
	// 	})

	// 	It("should run DS_WRITE2_B64", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.DS
	// 		state.inst.Opcode = 78
	// 		state.inst.Offset0 = 1
	// 		state.inst.Offset1 = 3

	// 		sp := state.scratchpad.AsDS()
	// 		sp.EXEC = 0x1
	// 		sp.ADDR[0] = 100
	// 		sp.DATA[0] = 1
	// 		sp.DATA[1] = 2
	// 		sp.DATA1[0] = 3
	// 		sp.DATA1[1] = 4

	// 		alu.Run(state)

	// 		lds := alu.LDS()
	// 		Expect(insts.BytesToUint32(lds[108:])).To(Equal(uint32(1)))
	// 		Expect(insts.BytesToUint32(lds[112:])).To(Equal(uint32(2)))
	// 		Expect(insts.BytesToUint32(lds[124:])).To(Equal(uint32(3)))
	// 		Expect(insts.BytesToUint32(lds[128:])).To(Equal(uint32(4)))
	// 	})

	It("should run DS_READ_B64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.DS
		inst.Opcode = 118
		inst.Addr = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(1, 1, 1)
		wf.inst = inst

		wf.Exec = 0x1
		wf.WriteReg(insts.VReg(0), 1, 0, 100)
		lds := alu.LDS()
		copy(lds[100:], insts.Uint64ToBytes(12))

		alu.Run(wf)

		result := wf.ReadReg(insts.VReg(1), 1, 0)
		Expect(result).To(Equal(uint64(12)))
	})

	// 	It("should run DS_READ2_B64", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.DS
	// 		state.inst.Opcode = 119
	// 		state.inst.Offset0 = 1
	// 		state.inst.Offset1 = 3

	// 		sp := state.scratchpad.AsDS()
	// 		sp.EXEC = 0x1
	// 		sp.ADDR[0] = 100

	// 		lds := alu.LDS()
	// 		copy(lds[108:], insts.Uint32ToBytes(12))
	// 		copy(lds[124:], insts.Uint32ToBytes(156))

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint32(12)))
	// 		Expect(sp.DST[2]).To(Equal(uint32(156)))
	// 	})

})
