package emu

import (
	"math"

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

	It("should run V_CNDMASK_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 0
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 1)
		wf.WriteReg(insts.VReg(0), 1, 1, 2)
		wf.WriteReg(insts.VReg(1), 1, 0, 3)
		wf.WriteReg(insts.VReg(1), 1, 1, 4)
		wf.WriteReg(insts.Regs[insts.VCC], 1, 0, 1)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 3)

		alu.Run(wf)
		Expect(wf.ReadReg(insts.VReg(2), 1, 0)).To(Equal(uint64(3)))
		Expect(wf.ReadReg(insts.VReg(2), 1, 1)).To(Equal(uint64(2)))
	})

	It("should run V_ADD_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 1

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(2.0)))
		wf.WriteReg(insts.VReg(1), 1, 0, uint64(math.Float32bits(3.1)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(math.Float32bits(float32(5.1)))))
	})

	It("should run V_SUB_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 2
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(2.0)))
		wf.WriteReg(insts.VReg(1), 1, 0, uint64(math.Float32bits(3.1)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(math.Float32frombits(uint32(dst))).To(BeNumerically("~", -1.1, 1e-4))
	})

	It("should run V_SUBREV_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 3
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(2.0)))
		wf.WriteReg(insts.VReg(1), 1, 0, uint64(math.Float32bits(3.1)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(math.Float32frombits(uint32(dst))).To(
			BeNumerically("~", 1.1, 1e-4))
	})

	It("should run V_MUL_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 5

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(2.0)))
		wf.WriteReg(insts.VReg(1), 1, 0, uint64(math.Float32bits(3.1)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(math.Float32bits(float32(6.2)))))
	})

	// 	It("should run V_MUL_I32_I24", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 6

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = uint64(int32ToBits(-10))
	// 		sp.SRC1[0] = uint64(int32ToBits(20))
	// 		sp.EXEC = 1

	// 		alu.Run(state)

	// 		Expect(int32(sp.DST[0] & 0xffffffff)).To(Equal(int32(-200)))
	// 	})

	// 	It("should run V_MUL_U32_U24", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 8

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = 2
	// 		sp.SRC1[0] = 0x1000001
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(2)))
	// 	})

	// 	It("should run V_MIN_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 10

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = uint64(math.Float32bits(2.0))
	// 		sp.SRC1[0] = uint64(math.Float32bits(3.1))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(math.Float32bits(float32(2.0)))))
	// 	})

	It("should run V_MAX_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 11

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(2.0)))
		wf.WriteReg(insts.VReg(1), 1, 0, uint64(math.Float32bits(3.1)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(math.Float32bits(float32(3.1)))))
	})

	It("should run V_MIN_U32, with src0 > src1", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 14
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0x64)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x20)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x20)))
	})

	It("should run V_MIN_U32, with src0 = src1", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 14
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0x64)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x64)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x64)))
	})

	It("should run V_MIN_U32, with src0 < src1", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 14
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0x20)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x23)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x20)))
	})

	It("should run V_MAX_U32, with src0 > src1", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 15
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0x64)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x20)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x64)))
	})

	It("should run V_MAX_U32, with src0 = src1", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 15
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0x64)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x64)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x64)))
	})

	It("should run V_MAX_U32, with src0 < src1", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 15
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0x20)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x23)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x23)))
	})

	It("should run V_LSHRREV_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 16
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0x64)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x20)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x02)))
	})

	It("should run V_ASHRREV_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 17
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)
		wf.WriteReg(insts.VReg(0), 1, 0, 97)
		wf.WriteReg(insts.VReg(1), 1, 0, uint64(int32ToBits(-64)))

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(asInt32(uint32(dst))).To(Equal(int32(-32)))

	})

	It("should run V_LSHLREV_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 18
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0x64)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x02)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x20)))
	})

	It("should run V_AND_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 19
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 2) //10
		wf.WriteReg(insts.VReg(1), 1, 0, 3) //11
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(2)))
	})

	It("should run V_AND_B32 SDWA", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 19
		inst.IsSdwa = true
		inst.Src0Sel = insts.SDWASelectByte0
		inst.Src1Sel = insts.SDWASelectByte3
		inst.DstSel = insts.SDWASelectWord1
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 0xfedcba98)
		wf.WriteReg(insts.VReg(1), 1, 0, 0x12345678)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)

		alu.Run(wf)

		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0x00100000)))
	})

	// 	It("should run V_OR_B32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 20

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = 2 // 10
	// 		sp.SRC1[0] = 3 // 11
	// 		sp.EXEC = 1

	// 		alu.Run(state)

	// 		Expect(uint32(sp.DST[0])).To(Equal(uint32(3)))
	// 	})

	// 	It("should run V_OR_B32 SDWA", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 20
	// 		state.inst.IsSdwa = true
	// 		state.inst.Src0Sel = insts.SDWASelectByte0
	// 		state.inst.Src1Sel = insts.SDWASelectByte3
	// 		state.inst.DstSel = insts.SDWASelectWord1

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = 0xfedcba98
	// 		sp.SRC1[0] = 0x12345678
	// 		sp.EXEC = 1

	// 		alu.Run(state)

	// 		Expect(uint32(sp.DST[0])).To(Equal(uint32(0x009a0000)))
	// 	})

	// 	It("should run V_XOR_B32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 21

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = 2 // 10
	// 		sp.SRC1[0] = 3 // 11
	// 		sp.EXEC = 1

	// 		alu.Run(state)

	// 		Expect(uint32(sp.DST[0])).To(Equal(uint32(1)))
	// 	})
	// 	It("should run V_OR_B32 SDWA", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 21
	// 		state.inst.IsSdwa = true
	// 		state.inst.Src0Sel = insts.SDWASelectByte0
	// 		state.inst.Src1Sel = insts.SDWASelectByte3
	// 		state.inst.DstSel = insts.SDWASelectWord1

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = 0xfedcba98
	// 		sp.SRC1[0] = 0x12345678
	// 		sp.EXEC = 1

	// 		alu.Run(state)

	// 		Expect(uint32(sp.DST[0])).To(Equal(uint32(0x008a0000)))
	// 	})

	It("should run V_MAC_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 22
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(float32ToBits(4)))
		wf.WriteReg(insts.VReg(1), 1, 0, uint64(float32ToBits(16)))
		wf.WriteReg(insts.VReg(2), 1, 0, uint64(float32ToBits(1024)))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)

		alu.Run(wf)
		result := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(asFloat32(uint32(result))).To(Equal(float32(1024.0 + 16.0*4.0)))
	})

	// 	It("should run V_MADAK_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 24

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = uint64(float32ToBits(4))
	// 		sp.SRC1[0] = uint64(float32ToBits(16))
	// 		sp.LiteralConstant = uint64(float32ToBits(1024))
	// 		sp.EXEC = 1

	// 		alu.Run(state)

	// 		Expect(asFloat32(uint32(sp.DST[0]))).To(Equal(float32(1024.0 + 16.0*4.0)))
	// 	})

	It("should run V_ADD_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 25

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)
		wf.inst = inst

		for i := 0; i < 64; i++ {
			wf.WriteReg(insts.VReg(0), 1, i, uint64(int32ToBits(-100)))
			wf.WriteReg(insts.VReg(1), 1, i, uint64(int32ToBits(10)))
		}
		wf.WriteReg(insts.Regs[insts.EXEC], 2, 0, 0xffffffffffffffff)

		alu.Run(wf)

		for i := 0; i < 64; i++ {
			dst := wf.ReadReg(insts.VReg(2), 1, i)
			Expect(asInt32(uint32(dst))).To(Equal(int32(-90)))
		}
	})

	It("should run V_ADD_I32_SDWA", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 25
		inst.IsSdwa = true
		inst.Src0Sel = insts.SDWASelectByte0
		inst.Src1Sel = insts.SDWASelectByte0
		inst.DstSel = insts.SDWASelectDWord
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)
		wf.inst = inst

		for i := 0; i < 64; i++ {
			wf.WriteReg(insts.VReg(0), 1, i, uint64(int32ToBits(-100)))
			wf.WriteReg(insts.VReg(1), 1, i, uint64(int32ToBits(10)))
		}
		wf.WriteReg(insts.Regs[insts.EXEC], 2, 0, 0xffffffffffffffff)

		alu.Run(wf)

		for i := 0; i < 64; i++ {
			dst := wf.ReadReg(insts.VReg(2), 1, i)
			Expect(asInt32(uint32(dst))).To(Equal(int32(166)))
		}
	})

	// 	// It("should run V_ADD_I32, with positive overflow", func() {
	// 	// 	state.inst = insts.NewInst()
	// 	// 	state.inst.FormatType = insts.VOP2
	// 	// 	state.inst.Opcode = 25

	// 	// 	sp := state.Scratchpad().AsVOP2()
	// 	// 	for i := 0; i < 64; i++ {
	// 	// 		sp.SRC0[i] = uint64(int32ToBits(math.MaxInt32 - 10))
	// 	// 		sp.SRC1[i] = uint64(int32ToBits(12))
	// 	// 	}
	// 	// 	sp.EXEC = 0xffffffffffffffff

	// 	// 	alu.Run(state)

	// 	// 	for i := 0; i < 64; i++ {
	// 	// 		Expect(asInt32(uint32(sp.DST[0]))).To(
	// 	// 			Equal(int32(math.MinInt32 + 1)))
	// 	// 	}
	// 	// 	Expect(sp.VCC).To(Equal(uint64(0xffffffffffffffff)))
	// 	// })

	// 	// It("should run V_ADD_I32, with negative overflow", func() {
	// 	// 	state.inst = insts.NewInst()
	// 	// 	state.inst.FormatType = insts.VOP2
	// 	// 	state.inst.Opcode = 25

	// 	// 	sp := state.Scratchpad().AsVOP2()
	// 	// 	for i := 0; i < 64; i++ {
	// 	// 		sp.SRC0[i] = uint64(int32ToBits(math.MinInt32 + 10))
	// 	// 		sp.SRC1[i] = uint64(int32ToBits(-12))
	// 	// 	}
	// 	// 	sp.EXEC = 0xffffffffffffffff

	// 	// 	alu.Run(state)

	// 	// 	for i := 0; i < 64; i++ {
	// 	// 		Expect(asInt32(uint32(sp.DST[0]))).To(
	// 	// 			Equal(int32(math.MaxInt32 - 1)))
	// 	// 	}
	// 	// 	Expect(sp.VCC).To(Equal(uint64(0xffffffffffffffff)))
	// 	// })

	It("should run V_SUB_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 26

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)
		wf.WriteReg(insts.VReg(0), 1, 0, 10)
		wf.WriteReg(insts.VReg(1), 1, 0, 4)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		vcc := wf.ReadReg(insts.Regs[insts.VCC], 1, 0)
		Expect(dst).To(Equal(uint64(6)))
		Expect(vcc).To(Equal(uint64(0)))
	})

	It("should run V_SUB_I32, when underflow", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 26

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)
		wf.WriteReg(insts.VReg(0), 1, 0, 4)
		wf.WriteReg(insts.VReg(1), 1, 0, 10)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		vcc := wf.ReadReg(insts.Regs[insts.VCC], 1, 0)
		Expect(dst).To(Equal(uint64(0xfffffffa)))
		Expect(vcc).To(Equal(uint64(1)))
	})

	It("should run V_SUBREV_I32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 27

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)
		wf.WriteReg(insts.VReg(0), 1, 0, 4)
		wf.WriteReg(insts.VReg(1), 1, 0, 10)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		vcc := wf.ReadReg(insts.Regs[insts.VCC], 1, 0)
		Expect(dst).To(Equal(uint64(6)))
		Expect(vcc).To(Equal(uint64(0)))
	})

	It("should run V_SUBREV_I32, when underflow", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 27

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)
		wf.WriteReg(insts.VReg(0), 1, 0, 10)
		wf.WriteReg(insts.VReg(1), 1, 0, 4)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		vcc := wf.ReadReg(insts.Regs[insts.VCC], 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0xfffffffa)))
		Expect(vcc).To(Equal(uint64(1)))
	})

	It("should run V_ADDC_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP2
		inst.Opcode = 28

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Src1 = insts.NewVRegOperand(1, 1, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, math.MaxUint32-10)
		wf.WriteReg(insts.VReg(1), 1, 0, 10)
		wf.WriteReg(insts.Regs[insts.VCC], 1, 0, uint64(1))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)

		alu.Run(wf)

		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		vcc := wf.ReadReg(insts.Regs[insts.VCC], 1, 0)
		Expect(uint32(dst)).To(Equal(uint32(0)))
		Expect(vcc).To(Equal(uint64(1)))
	})

	// 	It("should run V_SUBB_U32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 29

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = 10
	// 		sp.SRC1[0] = 5
	// 		sp.SRC0[1] = 5
	// 		sp.SRC1[1] = 10
	// 		sp.VCC = uint64(3)
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(uint32(sp.DST[0])).To(Equal(uint32(4)))
	// 		Expect(uint32(sp.DST[1])).To(Equal(^uint32(0) - 5))
	// 		Expect(sp.VCC).To(Equal(uint64(2)))
	// 	})

	// 	It("should run V_SUBBREV_U32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 30

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = 10
	// 		sp.SRC1[0] = 11
	// 		sp.VCC = uint64(0)
	// 		sp.EXEC = 1

	// 		alu.Run(state)

	// 		Expect(uint32(sp.DST[0])).To(Equal(uint32(1)))
	// 		Expect(sp.VCC).To(Equal(uint64(0)))
	// 	})

	// 	It("should run V_SUBBREV_U32, when underflow", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP2
	// 		state.inst.Opcode = 30

	// 		sp := state.Scratchpad().AsVOP2()
	// 		sp.SRC0[0] = 10
	// 		sp.SRC1[0] = 4
	// 		sp.VCC = uint64(1)
	// 		sp.EXEC = 1

	// 		alu.Run(state)

	// 		Expect(uint32(sp.DST[0])).To(Equal(uint32(0xfffffff9)))
	// 		Expect(sp.VCC).To(Equal(uint64(1)))
	// 	})

})
