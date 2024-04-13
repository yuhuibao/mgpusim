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

	It("should run V_ADD_U32 VOP3b", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP3b
		inst.Opcode = 281
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)
		inst.SDst = insts.NewSRegOperand(3, 3, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 3)
		wf.WriteReg(insts.VReg(0), 1, 0, 1)
		wf.WriteReg(insts.VReg(0), 1, 1, 0xffffffff)
		wf.WriteReg(insts.VReg(1), 1, 0, 2)
		wf.WriteReg(insts.VReg(1), 1, 1, 2)

		alu.Run(wf)
		dst_0 := wf.ReadReg(insts.VReg(2), 1, 0)
		dst_1 := wf.ReadReg(insts.VReg(2), 1, 1)
		sdst := wf.ReadReg(insts.SReg(3), 1, 0)
		Expect(dst_0).To(Equal(uint64(3)))
		Expect(dst_1 & 0xffffffff).To(Equal(uint64(1)))
		Expect(sdst).To(Equal(uint64(0x2)))
	})

	// 	It("should run V_SUB_U32 VOP3b", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP3b
	// 		state.inst.Opcode = 282

	// 		sp := state.Scratchpad().AsVOP3B()
	// 		sp.SRC0[0] = 1
	// 		sp.SRC1[0] = 2
	// 		sp.SRC0[1] = 0xffffffff
	// 		sp.SRC1[1] = 2
	// 		sp.EXEC = 3

	// 		alu.Run(state)

	// 		Expect(sp.DST[0] & 0xffffffff).To(Equal(uint64(0xffffffff)))
	// 		Expect(sp.DST[1] & 0xffffffff).To(Equal(uint64(0xfffffffd)))
	// 		Expect(sp.SDST).To(Equal(uint64(0x1)))
	// 	})

	// 	It("should run V_SUBREV_U32 VOP3b", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP3b
	// 		state.inst.Opcode = 283

	// 		sp := state.Scratchpad().AsVOP3B()
	// 		sp.SRC0[0] = uint64(2)
	// 		sp.SRC1[0] = uint64(0xffffffff)
	// 		sp.SRC0[1] = uint64(2)
	// 		sp.SRC1[1] = uint64(0x0)
	// 		sp.EXEC = 3

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(0xfffffffd)))
	// 		Expect(sp.DST[1]).To(Equal(uint64(0xfffffffe)))
	// 		Expect(sp.SDST).To(Equal(uint64(2)))
	// 	})

	It("should run V_ADDC_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP3b
		inst.Opcode = 284
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Src2 = insts.NewVRegOperand(2, 2, 1)
		inst.Dst = insts.NewVRegOperand(3, 3, 1)
		inst.SDst = insts.NewSRegOperand(4, 4, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 3)
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(0xfffffffd))
		wf.WriteReg(insts.VReg(0), 1, 1, uint64(0xfffffffd))
		wf.WriteReg(insts.VReg(1), 1, 0, uint64(2))
		wf.WriteReg(insts.VReg(1), 1, 1, uint64(1))
		wf.WriteReg(insts.VReg(2), 1, 0, uint64(1))
		wf.WriteReg(insts.VReg(2), 1, 1, uint64(1))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x3)

		alu.Run(wf)
		dst_0 := wf.ReadReg(insts.VReg(3), 1, 0)
		dst_1 := wf.ReadReg(insts.VReg(3), 1, 1)
		sdst := wf.ReadReg(insts.SReg(4), 1, 0)

		Expect(dst_0).To(Equal(uint64(0)))
		Expect(dst_1).To(Equal(uint64(0xfffffffe)))
		Expect(sdst).To(Equal(uint64(1)))
	})

	// 	It("should run V_SUBB_U32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP3b
	// 		state.inst.Opcode = 285

	// 		sp := state.scratchpad.AsVOP3B()
	// 		sp.SRC0[0] = uint64(0x1)
	// 		sp.SRC1[0] = uint64(0x2)
	// 		sp.SRC2[0] = uint64(0x1)
	// 		sp.SRC0[1] = uint64(0xfffffffd)
	// 		sp.SRC1[1] = uint64(0x1)
	// 		sp.SRC2[1] = uint64(0x1)
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(0xfffffffe)))
	// 		Expect(sp.DST[1]).To(Equal(uint64(0xfffffffc)))
	// 		Expect(sp.SDST).To(Equal(uint64(1)))
	// 	})

	// 	It("should run V_SUBBREV_U32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP3b
	// 		state.inst.Opcode = 286

	// 		sp := state.scratchpad.AsVOP3B()
	// 		sp.SRC1[0] = uint64(0x1)
	// 		sp.SRC0[0] = uint64(0x2)
	// 		sp.SRC2[0] = uint64(0x1)
	// 		sp.SRC1[1] = uint64(0xfffffffd)
	// 		sp.SRC0[1] = uint64(0x1)
	// 		sp.SRC2[1] = uint64(0x1)
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(0xfffffffe)))
	// 		Expect(sp.DST[1]).To(Equal(uint64(0xfffffffc)))
	// 		Expect(sp.SDST).To(Equal(uint64(1)))
	// 	})

	// 	It("should run V_DIV_SCALE_F64", func() {
	// 		// Need more test case
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP3b
	// 		state.inst.Opcode = 481

	// 		sp := state.scratchpad.AsVOP3B()
	// 		sp.SRC0[0] = uint64(0x3FF0000000000000)
	// 		sp.SRC1[0] = uint64(0x0008A00000000000)
	// 		sp.SRC2[0] = uint64(0x0008A00000000000)
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(math.Float64bits(math.Pow(2.0, 128))))
	// 	})

})
