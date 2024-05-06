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

	It("should run V_MOV_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 1

		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		for i := 0; i < 32; i++ {
			wf.WriteReg(insts.VReg(0), 1, i, uint64(i))
		}
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x00000000ffffffff)

		alu.Run(wf)

		for i := 0; i < 32; i++ {
			src0 := wf.ReadReg(insts.VReg(0), 1, i)
			dst := wf.ReadReg(insts.VReg(2), 1, i)
			Expect(src0).To(Equal(dst))
		}

		for i := 32; i < 64; i++ {
			src0 := wf.ReadReg(insts.VReg(0), 1, i)
			dst := wf.ReadReg(insts.VReg(2), 1, i)
			Expect(src0).To(Equal(dst))
		}
	})

	It("should run V_READFIRSTLANE_B32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 2
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 8, 1)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x0000000000000100)

		alu.Run(wf)

		for i := 0; i < 64; i++ {
			src0 := wf.ReadReg(insts.VReg(0), 1, 8)
			dst := wf.ReadReg(insts.VReg(2), 1, i)
			Expect(src0).To(Equal(dst))
		}

	})

	// 	It("should run V_CVT_F64_I32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 4

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(1)
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(math.Float64frombits(sp.DST[0])).To(Equal(float64(1.0)))
	// 	})

	// 	It("should run V_CVT_F32_I32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 5

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(int32ToBits(-1))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(-1.0)))
	// 	})

	It("should run V_CVT_F32_U32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 6
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, 1)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(math.Float32frombits(uint32(dst))).To(Equal(float32(1.0)))
	})

	It("should run V_CVT_U32_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 7
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(1.0)))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(1)))
	})

	It("should run V_CVT_U32_F32, when input is nan", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 7
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(float32(math.NaN()))))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(0)))
	})

	It("should run V_CVT_U32_F32, when the input is negative", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 7
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(-1.0)))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(0)))
	})

	It("should run V_CVT_U32_F32, when the input is very large", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 7
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(float32(math.MaxUint32+1))))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x1)

		alu.Run(wf)
		dst := wf.ReadReg(insts.VReg(2), 1, 0)
		Expect(dst).To(Equal(uint64(math.MaxUint32)))
	})

	// 	It("should run V_CVT_I32_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 8

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(1.5))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(1)))
	// 	})

	// 	It("should run V_CVT_I32_F32, when input is nan", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 8

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(float32(0 - math.NaN())))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(0)))
	// 	})

	// 	It("should run V_CVT_I32_F32, when the input is negative", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 8

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(-1.5))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(int32ToBits(-1))))
	// 	})

	// 	It("should run V_CVT_I32_F32, when the input is very large", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 8

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(0 - float32(math.MaxInt32) - 1))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(sp.DST[0]).To(Equal(uint64(int32ToBits(0 - math.MaxInt32))))
	// 	})

	// 	It("should run V_TRUNC_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 28

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(1.1))
	// 		sp.SRC0[1] = uint64(math.Float32bits(-2.2))
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(1.0)))
	// 		Expect(math.Float32frombits(uint32(sp.DST[1]))).To(Equal(float32(-2.0)))
	// 	})

	// 	It("should run V_RNDNE_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 30

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(1.1))
	// 		sp.SRC0[1] = uint64(math.Float32bits(-2.6))
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(1.0)))
	// 		Expect(math.Float32frombits(uint32(sp.DST[1]))).To(Equal(float32(-3.0)))
	// 	})

	// 	It("should run V_EXP_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 32

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(1.1))
	// 		sp.SRC0[1] = uint64(math.Float32bits(-2.6))
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).
	// 			To(BeNumerically("~", float32(2.1436), 1e-3))
	// 		Expect(math.Float32frombits(uint32(sp.DST[1]))).
	// 			To(BeNumerically("~", float32(0.1649), 1e-3))
	// 	})

	// 	It("should run V_LOG_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 33

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(1.1))
	// 		sp.SRC0[1] = uint64(math.Float32bits(-2.6))
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).
	// 			To(BeNumerically("~", float32(0.1375), 1e-3))
	// 		Expect(math.IsNaN(float64(math.Float32frombits(uint32(sp.DST[1]))))).
	// 			To(BeTrue())
	// 	})

	It("should run V_RCP_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 34
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(1.0)))
		wf.WriteReg(insts.VReg(0), 1, 1, uint64(math.Float32bits(2.0)))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x3)

		alu.Run(wf)
		dst_0 := wf.ReadReg(insts.VReg(2), 1, 0)
		dst_1 := wf.ReadReg(insts.VReg(2), 1, 1)
		Expect(math.Float32frombits(uint32(dst_0))).To(Equal(float32(1.0)))
		Expect(math.Float32frombits(uint32(dst_1))).To(Equal(float32(0.5)))
	})

	It("should run V_RCP_IFLAG_F32", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.VOP1
		inst.Opcode = 35
		inst.Src0 = insts.NewVRegOperand(0, 0, 1)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)

		wf.inst = inst
		wf.WriteReg(insts.VReg(0), 1, 0, uint64(math.Float32bits(1.0)))
		wf.WriteReg(insts.VReg(0), 1, 1, uint64(math.Float32bits(2.0)))
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0x3)

		alu.Run(wf)
		dst_0 := wf.ReadReg(insts.VReg(2), 1, 0)
		dst_1 := wf.ReadReg(insts.VReg(2), 1, 1)
		Expect(math.Float32frombits(uint32(dst_0))).To(Equal(float32(1.0)))
		Expect(math.Float32frombits(uint32(dst_1))).To(Equal(float32(0.5)))
	})

	// 	It("should run V_RSQ_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 36

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(4.0))
	// 		sp.SRC0[1] = uint64(math.Float32bits(625.0))
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(0.5)))
	// 		Expect(math.Float32frombits(uint32(sp.DST[1]))).To(Equal(float32(0.04)))
	// 	})

	// 	It("should run V_SQRT_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 39

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(4.0))
	// 		sp.SRC0[1] = uint64(math.Float32bits(625.0))
	// 		sp.EXEC = 0x3

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(2.0)))
	// 		Expect(math.Float32frombits(uint32(sp.DST[1]))).To(Equal(float32(25.0)))
	// 	})

	// 	It("should run V_CVT_F32_UBYTE0", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 17

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(256.0))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(0)))
	// 	})

	// 	It("should run V_CVT_F64_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 16

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(-1.0))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)
	// 		Expect(sp.DST[0]).To(Equal(math.Float64bits(float64(-1.0))))
	// 	})

	// 	It("should run V_RCP_F64", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 37

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = math.Float64bits(25.0)
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(math.Float64frombits(sp.DST[0])).To(Equal(float64(0.04)))
	// 	})

	// 	It("should run V_CVT_F32_F64", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 15

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = math.Float64bits(25.0)
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)

	// 		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(25.0)))
	// 	})

	// 	It("should run V_CVT_F16_F32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 10

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(math.Float32bits(8.0))
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)
	// 		// value 8.0 => half - precision : 0x4800
	// 		Expect(uint16(sp.DST[0])).To(Equal(uint16(0x4800)))
	// 	})

	// 	It("should run V_BREV_B32", func() {
	// 		state.inst = insts.NewInst()
	// 		state.inst.FormatType = insts.VOP1
	// 		state.inst.Opcode = 44

	// 		sp := state.Scratchpad().AsVOP1()
	// 		sp.SRC0[0] = uint64(0xffff)
	// 		sp.EXEC = 0x1

	// 		alu.Run(state)
	// 		Expect(uint32(sp.DST[0])).To(Equal(uint32(0xffff0000)))
	// 	})

})
