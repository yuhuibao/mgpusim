package cu

import (
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sarchlab/mgpusim/v3/emu"
	"github.com/sarchlab/mgpusim/v3/insts"
	"github.com/sarchlab/mgpusim/v3/kernels"
	"github.com/sarchlab/mgpusim/v3/timing/wavefront"
)

var _ = Describe("LDS Unit", func() {

	var (
		mockCtrl *gomock.Controller
		cu       *ComputeUnit
		bu       *LDSUnit
		alu      *MockALU
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		cu = NewComputeUnit("CU", nil)
		alu = NewMockALU(mockCtrl)
		bu = NewLDSUnit(cu, alu)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should allow accepting wavefront", func() {
		bu.toRead = nil
		Expect(bu.CanAcceptWave()).To(BeTrue())
	})

	It("should not allow accepting wavefront is the read stage buffer is occupied", func() {
		wave := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		bu.toRead = wave
		Expect(bu.CanAcceptWave()).To(BeFalse())
	})

	It("should accept wave", func() {
		wave := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		bu.AcceptWave(wave, 10)
		Expect(bu.toRead).To(BeIdenticalTo(wave))
	})

	It("should run", func() {
		wave1 := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		wave2 := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		wave2.WG = wavefront.NewWorkGroup(nil, nil)
		wave2.WG.LDS = make([]byte, 0)
		wave3 := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		inst := wavefront.NewInst(insts.NewInst())
		inst.FormatType = insts.DS
		inst.Opcode = 0
		inst.Addr = insts.NewVRegOperand(0, 0, 1)
		inst.Data = insts.NewVRegOperand(2, 2, 2)
		inst.Data1 = insts.NewVRegOperand(4, 4, 2)
		inst.ByteSize = 4
		wave3.SetDynamicInst(inst)
		wave3.PC = 0x13C
		wave3.InstBuffer = make([]byte, 256)
		wave3.InstBufferStartPC = 0x100

		wave3.State = wavefront.WfRunning

		bu.toRead = wave1
		bu.toExec = wave2
		bu.toWrite = wave3

		alu.EXPECT().SetLDS(wave2.WG.LDS)
		alu.EXPECT().Run(wave2).Times(1)

		bu.Run(10)

		Expect(wave3.State).To(Equal(wavefront.WfReady))
		Expect(wave3.PC).To(Equal(uint64(0x140)))

		Expect(bu.toWrite).To(BeIdenticalTo(wave2))
		Expect(bu.toExec).To(BeIdenticalTo(wave1))
		Expect(bu.toRead).To(BeNil())

		Expect(wave3.InstBuffer).To(HaveLen(192))

	})

	It("should flush the LDS", func() {

		wave1 := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		wave2 := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		wave2.WG = wavefront.NewWorkGroup(nil, nil)
		wave2.WG.LDS = make([]byte, 0)
		wave3 := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		inst := wavefront.NewInst(insts.NewInst())
		inst.FormatType = insts.DS
		inst.Opcode = 0
		inst.Addr = insts.NewVRegOperand(0, 0, 1)
		inst.Data = insts.NewVRegOperand(2, 2, 2)
		inst.Data1 = insts.NewVRegOperand(4, 4, 2)
		inst.ByteSize = 4
		wave3.SetDynamicInst(inst)
		wave3.PC = 0x13C
		wave3.InstBuffer = make([]byte, 256)
		wave3.InstBufferStartPC = 0x100

		wave3.State = wavefront.WfRunning

		bu.toRead = wave1
		bu.toExec = wave2
		bu.toWrite = wave3

		bu.Flush()

		Expect(bu.toRead).To(BeNil())
		Expect(bu.toWrite).To(BeNil())
		Expect(bu.toExec).To(BeNil())

	})
})
