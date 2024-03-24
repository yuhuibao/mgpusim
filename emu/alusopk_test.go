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
		rawWf := kernels.NewWavefront()
		wf = NewWavefront(rawWf)
	})

	It("should run s_movk_i32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPK
		inst.Opcode = 0

		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)
		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int16ToBits(-12)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 1, 0)
		Expect(asInt16(uint16(dst))).To(Equal(int16(-12)))
	})

	It("should run s_cmpk_lg_i32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPK
		inst.Opcode = 3

		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int16ToBits(100)))
		wf.WriteReg(insts.SReg(2), 1, 0, 200)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run s_mulk_i32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPK
		inst.Opcode = 15

		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int16ToBits(100)))
		wf.WriteReg(insts.SReg(2), 2, 0, 200)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(uint64(20000)))
	})

	It("should run s_mulk_i32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPK
		inst.Opcode = 15

		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int16ToBits(-100)))
		wf.WriteReg(insts.SReg(2), 2, 0, 200)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(asInt64(dst)).To(Equal(int64(-20000)))
	})

})
