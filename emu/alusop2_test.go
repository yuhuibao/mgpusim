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

	It("should run S_ADD_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 0

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(1<<31-1))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(1<<31+15))

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(14)))
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_SUB_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 1

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(10))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(5))

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(5)))
	})

	It("should run S_SUB_U32 with carry out", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 1

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(5))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(10))

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(^uint64(0) - 4))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ADD_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 2

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0xffffffff)
		wf.WriteReg(insts.SReg(1), 1, 0, 3)

		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(2), 1, 0)
		Expect(results).To(Equal(uint64(2)))
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_SUB_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 3

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 10)
		wf.WriteReg(insts.SReg(1), 1, 0, 6)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(uint64(4)))
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_SUB_I32, when input is negative", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 3

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64((int32ToBits(-6))))
		wf.WriteReg(insts.SReg(1), 1, 0, 15)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(asInt32(uint32(dst))).To(Equal(int32(-21)))
		Expect(wf.SCC).To(Equal(byte(0)))
	})

	It("should run S_SUB_I32, when overflow and src1 is positive", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 3

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0x7ffffffe)
		wf.WriteReg(insts.SReg(1), 1, 0, 0xfffffffc)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_SUB_I32, when overflow and src1 is negtive", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 3

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0x80000001)
		wf.WriteReg(insts.SReg(1), 1, 0, 10)

		alu.Run(wf)
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_ADDC_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 4

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(1<<31-1))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(1<<31))
		wf.SCC = 1

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(uint32(dst)).To(Equal(uint32(0)))
		Expect(wf.SCC).To(Equal(byte(1)))
	})

	It("should run S_SUBB_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 5

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 10)
		wf.WriteReg(insts.SReg(1), 1, 0, 5)
		wf.SCC = 1

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(uint64(4)))
		Expect(wf.SCC).To(Equal(uint8(0)))
	})

	It("should run S_SUBB_U32 with carry out", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 5

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 5)
		wf.WriteReg(insts.SReg(1), 1, 0, 10)
		wf.SCC = 1

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(^uint64(0) - 5))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_SUBB_U32 with carry out", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 5

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0)
		wf.WriteReg(insts.SReg(1), 1, 0, 0)
		wf.SCC = 1

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(^uint64(0)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_MIN_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 6

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int32ToBits(-1)))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(int32ToBits(5)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(asInt32(uint32(dst))).To(Equal(int32(-1)))
	})

	It("should run S_MIN_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 7

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 1)
		wf.WriteReg(insts.SReg(1), 1, 0, 2)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(uint64(1)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_MIN_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 7

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 2)
		wf.WriteReg(insts.SReg(1), 1, 0, 1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(uint64(1)))
		Expect(wf.SCC).To(Equal(uint8(0)))
	})

	It("should run S_MAX_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 8

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(int32ToBits(-1)))
		wf.WriteReg(insts.SReg(1), 1, 0, uint64(int32ToBits(-5)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(asInt32(uint32(dst))).To(Equal(int32(-1)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_MAX_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 9

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0xff)
		wf.WriteReg(insts.SReg(1), 1, 0, 0xffff)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(uint64(0xffff)))
		Expect(wf.SCC).To(Equal(uint8(0)))
	})

	It("should run S_MAX_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 9

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0xffff)
		wf.WriteReg(insts.SReg(1), 1, 0, 0xff)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(uint64(0xffff)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_CSELECT_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 10

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0xffff)
		wf.WriteReg(insts.SReg(1), 1, 0, 0xff)
		wf.SCC = 1

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(dst).To(Equal(uint64(0xffff)))
	})

	It("should run S_AND_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 12

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0xff)
		wf.WriteReg(insts.SReg(1), 1, 0, 0xffff)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(0xff)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_AND_B32 with carry out", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 12

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0x0)
		wf.WriteReg(insts.SReg(1), 1, 0, 0xffff)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(0)))
		Expect(wf.SCC).To(Equal(uint8(0)))
	})

	It("should run S_AND_B64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 13

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, 0xff)
		wf.WriteReg(insts.SReg(2), 2, 0, 0xffff)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(0xff)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_OR_B64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 15

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, 0xf0)
		wf.WriteReg(insts.SReg(2), 2, 0, 0xff)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(0xff)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_XOR_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 16

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0xf0)
		wf.WriteReg(insts.SReg(1), 1, 0, 0xff)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(0x0f)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_XOR_B64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 17

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, 0xf0)
		wf.WriteReg(insts.SReg(2), 2, 0, 0xff)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(0x0f)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ANDN2_B64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 19

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, 0xab)
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0f)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(0xa0)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_LSHL_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 28

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(1, 1, 1)
		inst.Dst = insts.NewSRegOperand(2, 2, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 128)
		wf.WriteReg(insts.SReg(1), 1, 0, 2)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(512)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_LSHL_B64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 29

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, 128)
		wf.WriteReg(insts.SReg(2), 2, 0, 2)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(512)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_LSHL_B64 (To zero)", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 29

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, 0x8000000000000000)
		wf.WriteReg(insts.SReg(2), 2, 0, 1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(0)))
		Expect(wf.SCC).To(Equal(uint8(0)))
	})

	It("should run S_LSHR_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 30

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(2, 2, 1)
		inst.Dst = insts.NewSRegOperand(4, 4, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0x20)
		wf.WriteReg(insts.SReg(2), 1, 0, 0x64)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 1, 0)
		Expect(dst).To(Equal(uint64(0x02)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_LSHR_B64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 31

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, 0x20)
		wf.WriteReg(insts.SReg(2), 2, 0, 0x44)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(0x02)))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ASHR_I32 (Negative)", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 32

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, int64ToBits(-128))
		wf.WriteReg(insts.SReg(2), 2, 0, 2)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(int32ToBits(-32))))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ASHR_I32 (Positive)", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 32

		inst.Src0 = insts.NewSRegOperand(0, 0, 2)
		inst.Src1 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, int64ToBits(128))
		wf.WriteReg(insts.SReg(2), 2, 0, 2)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(dst).To(Equal(uint64(int32ToBits(32))))
		Expect(wf.SCC).To(Equal(uint8(1)))
	})

	It("should run S_BFM_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 34

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(2, 2, 1)
		inst.Dst = insts.NewSRegOperand(4, 4, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0x24)
		wf.WriteReg(insts.SReg(2), 1, 0, 0x64)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 1, 0)
		Expect(dst).To(Equal(uint64(int32ToBits(240))))
	})

	It("should run S_MUL_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 36

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(2, 2, 1)
		inst.Dst = insts.NewSRegOperand(4, 4, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 5)
		wf.WriteReg(insts.SReg(2), 1, 0, 7)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 1, 0)
		Expect(dst).To(Equal(uint64(35)))
	})

	It("should run S_BFE_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP2
		inst.Opcode = 38

		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Src1 = insts.NewSRegOperand(2, 2, 1)
		inst.Dst = insts.NewSRegOperand(4, 4, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, 0b1111_0100)
		wf.WriteReg(insts.SReg(2), 1, 0, 0b000000000_0000001_00000000000_00010)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 1, 0)
		Expect(dst).To(Equal(uint64(1)))
		Expect(wf.SCC).To(Equal(byte(1)))
	})

})
