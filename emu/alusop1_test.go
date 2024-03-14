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

	// It("should run s_mov_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 1

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)
	// 	Expect(sp.DST).To(Equal(uint64(0x0000ffffffff0000)))
	// })

	// It("should run s_not_u32", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 4

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.SRC0 = 0xff

	// 	alu.Run(state)

	// 	Expect(sp.DST).To(Equal(uint64(0xffffffffffffff00)))
	// 	Expect(sp.SCC).To(Equal(uint8(0x1)))
	// })

	// It("should run s_brev_b32", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 8

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.SRC0 = 0xffff

	// 	alu.Run(state)

	// 	Expect(sp.DST).To(Equal(uint64(0x00000000ffff0000)))
	// })

	// It("should run s_get_pc_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 28

	// 	sp := state.Scratchpad().AsSOP1()

	// 	sp.PC = 0xffffffff00000000

	// 	alu.Run(state)

	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000004)))
	// })

	// It("should run s_and_saveexec_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 32

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.EXEC = 0xffffffff00000000
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)

	// 	Expect(sp.EXEC).To(Equal(uint64(0x0000ffff00000000)))
	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000000)))
	// 	Expect(sp.SCC).To(Equal(byte(0x1)))
	// })

	// It("should run s_or_saveexec_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 33

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.EXEC = 0xffffffff00000000
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)

	// 	Expect(sp.EXEC).To(Equal(uint64(0xffffffffffff0000)))
	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000000)))
	// 	Expect(sp.SCC).To(Equal(byte(0x1)))
	// })

	// It("should run s_xor_saveexec_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 34

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.EXEC = 0xffffffff00000000
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)

	// 	Expect(sp.EXEC).To(Equal(uint64(0xffff0000ffff0000)))
	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000000)))
	// 	Expect(sp.SCC).To(Equal(byte(0x1)))
	// })

	// It("should run s_andn2_saveexec_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 35

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.EXEC = 0xffffffff00000000
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)

	// 	Expect(sp.EXEC).To(Equal(uint64(0x00000000ffff0000)))
	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000000)))
	// 	Expect(sp.SCC).To(Equal(byte(0x1)))
	// })

	// It("should run s_orn2_saveexec_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 36

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.EXEC = 0xffffffff00000000
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)

	// 	Expect(sp.EXEC).To(Equal(uint64(0x0000ffffffffffff)))
	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000000)))
	// 	Expect(sp.SCC).To(Equal(byte(0x1)))
	// })

	// It("should run s_nand_saveexec_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 37

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.EXEC = 0xffffffff00000000
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)

	// 	Expect(sp.EXEC).To(Equal(uint64(0xffff0000ffffffff)))
	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000000)))
	// 	Expect(sp.SCC).To(Equal(byte(0x1)))
	// })

	// It("should run s_nor_saveexec_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 38

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.EXEC = 0xffffffff00000000
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)

	// 	Expect(sp.EXEC).To(Equal(uint64(0x000000000000ffff)))
	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000000)))
	// 	Expect(sp.SCC).To(Equal(byte(0x1)))
	// })

	// It("should run s_nxor_saveexec_b64", func() {
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.SOP1
	// 	state.inst.Opcode = 39

	// 	sp := state.Scratchpad().AsSOP1()
	// 	sp.EXEC = 0xffffffff00000000
	// 	sp.SRC0 = 0x0000ffffffff0000

	// 	alu.Run(state)

	// 	Expect(sp.EXEC).To(Equal(uint64(0x0000ffff0000ffff)))
	// 	Expect(sp.DST).To(Equal(uint64(0xffffffff00000000)))
	// 	Expect(sp.SCC).To(Equal(byte(0x1)))
	// })

})
