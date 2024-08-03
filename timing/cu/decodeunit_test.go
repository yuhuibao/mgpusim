package cu

import (
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sarchlab/akita/v3/sim"
	"github.com/sarchlab/mgpusim/v3/emu"
	"github.com/sarchlab/mgpusim/v3/kernels"
	"github.com/sarchlab/mgpusim/v3/timing/wavefront"
)

var _ = Describe("DecodeUnit", func() {
	var (
		mockCtrl  *gomock.Controller
		cu        *ComputeUnit
		du        *DecodeUnit
		execUnits []*MockSubComponent
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		cu = NewComputeUnit("CU", nil)
		du = NewDecodeUnit(cu)
		execUnits = make([]*MockSubComponent, 4)
		for i := 0; i < 4; i++ {
			// execUnits[i] = new(MockSubComponent)
			execUnits[i] = NewMockSubComponent(mockCtrl)
			// execUnits[i].canAccept = true
			du.AddExecutionUnit(execUnits[i])
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should tell if it cannot accept wave", func() {
		du.toDecode = wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		Expect(du.CanAcceptWave()).To(BeFalse())
	})

	It("should tell if it can accept wave", func() {
		du.toDecode = nil
		Expect(du.CanAcceptWave()).To(BeTrue())
	})

	It("should accept wave", func() {
		wave := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		du.toDecode = nil
		du.AcceptWave(wave, 10)
		Expect(du.toDecode).To(BeIdenticalTo(wave))
	})

	It("should return error if the decoder is busy", func() {
		wave := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		wave2 := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		du.toDecode = wave

		Expect(func() { du.AcceptWave(wave2, 10) }).Should(Panic())
		Expect(du.toDecode).To(BeIdenticalTo(wave))
	})

	It("should deliver the wave to the execution unit", func() {
		wave := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		wave.SIMDID = 1
		du.toDecode = wave

		execUnits[1].EXPECT().CanAcceptWave().Return(true)
		execUnits[1].EXPECT().AcceptWave(wave, sim.VTimeInSec(10)).Times(1)

		du.Run(10)

		// Expect(len(execUnits[0].acceptedWave)).To(Equal(0))
		// Expect(len(execUnits[1].acceptedWave)).To(Equal(1))
		// Expect(len(execUnits[2].acceptedWave)).To(Equal(0))
		// Expect(len(execUnits[3].acceptedWave)).To(Equal(0))
		Expect(du.toDecode).To(BeNil())
	})

	// It("should not deliver to the execution unit, if busy", func() {
	// 	wave := new(wavefront.Wavefront)
	// 	wave.SIMDID = 1
	// 	du.toDecode = wave
	// 	execUnits[1].canAccept = false

	// 	du.Run(10)

	// 	Expect(len(execUnits[0].acceptedWave)).To(Equal(0))
	// 	Expect(len(execUnits[1].acceptedWave)).To(Equal(0))
	// 	Expect(len(execUnits[2].acceptedWave)).To(Equal(0))
	// 	Expect(len(execUnits[3].acceptedWave)).To(Equal(0))
	// })

	It("should flush the decode unit", func() {
		wave := wavefront.NewWavefront(emu.NewWavefront(new(kernels.Wavefront)))
		wave.SIMDID = 1
		du.toDecode = wave

		du.Flush()

		Expect(du.toDecode).To(BeNil())
	})

})
