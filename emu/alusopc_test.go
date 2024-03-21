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

	It("should run S_CMP_EQ_I32 when input is not equal", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 0

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 2)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_EQ_I32 when input is equal", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 0

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LG_I32 when condition holds", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 1

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 2)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LG_I32 when condition does not hold", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 1

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_GT_I32 when condition holds", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 2

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 2)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_GT_I32 when condition does not hold", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 2

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_GE_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 3

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LT_I32 when condition holds", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 4

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int32ToBits(-2)))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(int32ToBits(-1)))

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LT_I32 when condition does not hold", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 4

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int32ToBits(-1)))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(int32ToBits(-1)))

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_LE_I32 when condition holds", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 5

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int32ToBits(-2)))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(int32ToBits(-1)))

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LE_I32 when condition does not hold", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 5

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int32ToBits(-1)))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(int32ToBits(-2)))

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_EQ_U32 when input is not equal", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 6

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 2)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_EQ_U32 when input is equal", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 6

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LG_U32 when condition holds", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 7

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 2)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LG_U32 when condition does not hold", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 7

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_GT_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 8

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 2)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LT_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPC
		inst.Opcode = 10

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 2)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})
})
