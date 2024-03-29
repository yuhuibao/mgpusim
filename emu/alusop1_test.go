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

	It("should run s_mov_b32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 0
		inst.Src0 = insts.NewSRegOperand(0, 0, 1)
		inst.Dst = insts.NewSRegOperand(1, 1, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 1, 0, uint64(0xffff0000))
		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(1), 1, 0)
		Expect(results).To(Equal(uint64(0xffff0000)))

	})

	It("should run s_mov_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 1
		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(3, 3, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)
		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(3), 2, 0)
		Expect(results).To(Equal(uint64(0x0000ffffffff0000)))
	})

	It("should run s_not_u32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 4
		inst.Src0 = insts.NewSRegOperand(4, 4, 1)
		inst.Dst = insts.NewSRegOperand(5, 5, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(4), 1, 0, uint64(0xff))
		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(5), 1, 0)
		Expect(results).To(Equal(uint64(0xffffff00)))
		Expect(wf.SCC).To(Equal(uint8(0x1)))
	})

	It("should run s_brev_b32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 8

		inst.Src0 = insts.NewSRegOperand(8, 8, 1)
		inst.Dst = insts.NewSRegOperand(9, 9, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(8), 1, 0, uint64(0xffff))

		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(9), 1, 0)
		Expect(results).To(Equal(uint64(0x00000000ffff0000)))
	})

	It("should run s_get_pc_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 28

		inst.Dst = insts.NewSRegOperand(1, 1, 2)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 0xffffffff00000000)
		wf.inst = inst

		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(1), 2, 0)
		Expect(results).To(Equal(uint64(0xffffffff00000004)))
	})

	It("should run s_and_saveexec_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 32

		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.Exec = 0xffffffff00000000
		wf.inst = inst
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)

		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(wf.Exec).To(Equal(uint64(0x0000ffff00000000)))
		Expect(results).To(Equal(uint64(0xffffffff00000000)))
		Expect(wf.SCC).To(Equal(uint8(0x1)))
	})

	It("should run s_or_saveexec_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 33

		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.Exec = 0xffffffff00000000
		wf.inst = inst
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(wf.Exec).To(Equal(uint64(0xffffffffffff0000)))
		Expect(dst).To(Equal(uint64(0xffffffff00000000)))
		Expect(wf.SCC).To(Equal(byte(0x1)))
	})

	It("should run s_xor_saveexec_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 34

		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.Exec = 0xffffffff00000000
		wf.inst = inst
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(wf.Exec).To(Equal(uint64(0xffff0000ffff0000)))
		Expect(dst).To(Equal(uint64(0xffffffff00000000)))
		Expect(wf.SCC).To(Equal(byte(0x1)))
	})

	It("should run s_andn2_saveexec_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 35

		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.Exec = 0xffffffff00000000
		wf.inst = inst
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(wf.Exec).To(Equal(uint64(0x00000000ffff0000)))
		Expect(dst).To(Equal(uint64(0xffffffff00000000)))
		Expect(wf.SCC).To(Equal(byte(0x1)))
	})

	It("should run s_orn2_saveexec_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 36

		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.Exec = 0xffffffff00000000
		wf.inst = inst
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(wf.Exec).To(Equal(uint64(0x0000ffffffffffff)))
		Expect(dst).To(Equal(uint64(0xffffffff00000000)))
		Expect(wf.SCC).To(Equal(byte(0x1)))
	})

	It("should run s_nand_saveexec_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 37

		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.Exec = 0xffffffff00000000
		wf.inst = inst
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(wf.Exec).To(Equal(uint64(0xffff0000ffffffff)))
		Expect(dst).To(Equal(uint64(0xffffffff00000000)))
		Expect(wf.SCC).To(Equal(byte(0x1)))
	})

	It("should run s_nor_saveexec_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 38

		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.Exec = 0xffffffff00000000
		wf.inst = inst
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(wf.Exec).To(Equal(uint64(0x000000000000ffff)))
		Expect(dst).To(Equal(uint64(0xffffffff00000000)))
		Expect(wf.SCC).To(Equal(byte(0x1)))
	})

	It("should run s_nxor_saveexec_b64", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOP1
		inst.Opcode = 39

		inst.Src0 = insts.NewSRegOperand(2, 2, 2)
		inst.Dst = insts.NewSRegOperand(4, 4, 2)
		wf.Exec = 0xffffffff00000000
		wf.inst = inst
		wf.WriteReg(insts.SReg(2), 2, 0, 0x0000ffffffff0000)

		alu.Run(wf)
		dst := wf.ReadReg(insts.SReg(4), 2, 0)
		Expect(wf.Exec).To(Equal(uint64(0x0000ffff0000ffff)))
		Expect(dst).To(Equal(uint64(0xffffffff00000000)))
		Expect(wf.SCC).To(Equal(byte(0x1)))
	})

})
