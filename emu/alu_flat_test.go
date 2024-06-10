package emu

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sarchlab/akita/v3/mem/mem"
	"github.com/sarchlab/akita/v3/mem/vm"
	"github.com/sarchlab/mgpusim/v3/insts"
	"github.com/sarchlab/mgpusim/v3/kernels"
)

var _ = Describe("ALU", func() {

	var (
		mockCtrl  *gomock.Controller
		pageTable *MockPageTable

		alu           *ALUImpl
		wf            *Wavefront
		storage       *mem.Storage
		addrConverter *mem.InterleavingConverter
		sAccessor     *storageAccessor
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		pageTable = NewMockPageTable(mockCtrl)

		storage = mem.NewStorage(1 * mem.GB)
		addrConverter = &mem.InterleavingConverter{
			InterleavingSize:    1 * mem.GB,
			TotalNumOfElements:  1,
			CurrentElementIndex: 0,
			Offset:              0,
		}
		sAccessor = newStorageAccessor(storage, pageTable, 12, addrConverter)
		alu = NewALU(sAccessor)
		rawWf := kernels.NewWavefront()
		wf = NewWavefront(rawWf)
		wf.pid = vm.PID(1)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should run FLAT_LOAD_UBYTE", func() {
		for i := 0; i < 64; i++ {
			pageTable.EXPECT().Find(vm.PID(1), uint64(i*4)).
				Return(vm.Page{
					PAddr: uint64(0),
				}, true)
		}
		inst := insts.NewInst()
		inst.FormatType = insts.FLAT
		inst.Opcode = 16
		inst.Addr = insts.NewVRegOperand(0, 0, 2)
		inst.Dst = insts.NewVRegOperand(2, 2, 1)
		wf.inst = inst

		for i := 0; i < 64; i++ {
			wf.WriteReg(insts.VReg(0), 2, i, uint64(i*4))
			storage.Write(uint64(i*4), insts.Uint32ToBytes(uint32(i)))
		}
		wf.WriteReg(insts.Regs[insts.EXEC], 2, 0, 0xffffffffffffffff)

		alu.Run(wf)

		for i := 0; i < 64; i++ {
			results := wf.ReadReg(insts.VReg(2), 1, i)
			buf := insts.Uint32ToBytes(uint32(results))
			Expect(buf[0]).To(Equal(byte(i)))

			Expect(buf[1]).To(Equal(byte(0)))
			Expect(buf[2]).To(Equal(byte(0)))
			Expect(buf[3]).To(Equal(byte(0)))
		}
	})

	// It("should run FLAT_LOAD_USHORT", func() {
	// 	for i := 0; i < 64; i++ {
	// 		pageTable.EXPECT().
	// 			Find(vm.PID(1), uint64(i*4)).
	// 			Return(vm.Page{
	// 				PAddr: uint64(0),
	// 			}, true)
	// 	}
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.FLAT
	// 	state.inst.Opcode = 18

	// 	layout := state.Scratchpad().AsFlat()
	// 	for i := 0; i < 64; i++ {
	// 		layout.ADDR[i] = uint64(i * 4)
	// 		storage.Write(uint64(i*4), insts.Uint32ToBytes(uint32(i)))
	// 	}
	// 	layout.EXEC = 0xffffffffffffffff

	// 	alu.Run(state)

	// 	for i := 0; i < 64; i++ {
	// 		Expect(layout.DST[i*4]).To(Equal(uint32(i)))

	// 	}
	// })

	It("should run FLAT_LOAD_DWORD", func() {
		for i := 0; i < 64; i++ {
			pageTable.EXPECT().
				Find(vm.PID(1), uint64(i*4)).
				Return(vm.Page{
					PAddr: uint64(0),
				}, true)
		}
		inst := insts.NewInst()
		inst.FormatType = insts.FLAT
		inst.Opcode = 20
		inst.Addr = insts.NewVRegOperand(0, 0, 2)
		inst.Dst = insts.NewVRegOperand(2, 2, 2)
		wf.inst = inst

		// state.inst = insts.NewInst()
		// state.inst.FormatType = insts.FLAT
		// state.inst.Opcode = 20

		// layout := state.Scratchpad().AsFlat()
		for i := 0; i < 64; i++ {
			wf.WriteReg(insts.VReg(0), 2, i, uint64(i*4))
			// layout.ADDR[i] = uint64(i * 4)
			storage.Write(uint64(i*4), insts.Uint32ToBytes(uint32(i)))
		}
		// layout.EXEC = 0xffffffffffffffff
		wf.WriteReg(insts.Regs[insts.EXEC], 2, 0, 0xffffffffffffffff)

		alu.Run(wf)

		for i := 0; i < 64; i++ {
			results := wf.ReadReg(insts.VReg(2), 2, i)
			Expect(results).To(Equal((uint64(i))))
			// Expect(layout.DST[i*4]).To(Equal(uint32(i)))
		}
	})

	It("should run FLAT_LOAD_DWORDX2", func() {
		for i := 0; i < 64; i++ {
			pageTable.EXPECT().
				Find(vm.PID(1), uint64(i*8)).
				Return(vm.Page{
					PAddr: uint64(0),
				}, true)
		}
		inst := insts.NewInst()
		inst.FormatType = insts.FLAT
		inst.Opcode = 21
		inst.Addr = insts.NewVRegOperand(0, 0, 2)
		inst.Dst = insts.NewVRegOperand(2, 2, 2)
		wf.inst = inst

		// state.inst = insts.NewInst()
		// state.inst.FormatType = insts.FLAT
		// state.inst.Opcode = 21

		// layout := state.Scratchpad().AsFlat()
		for i := 0; i < 64; i++ {
			wf.WriteReg(insts.VReg(0), 2, i, uint64(i*8))
			// layout.ADDR[i] = uint64(i * 8)
			storage.Write(uint64(i*8), insts.Uint32ToBytes(uint32(i)))
			storage.Write(uint64(i*8+4), insts.Uint32ToBytes(uint32(i)))
		}
		// layout.EXEC = 0xffffffffffffffff
		wf.WriteReg(insts.Regs[insts.EXEC], 2, 0, 0xffffffffffffffff)

		alu.Run(wf)

		for i := 0; i < 64; i++ {
			results := wf.ReadReg(insts.VReg(2), 2, i)
			buf_0 := results & 0x00000000ffffffff
			buf_1 := (results & 0xffffffff00000000) >> 32
			// buf := insts.Uint32ToBytes(uint32(results))
			Expect(buf_0).To(Equal(uint64(i)))
			Expect(buf_1).To(Equal(uint64(i)))
			// Expect(layout.DST[i*4]).To(Equal(uint32(i)))
			// Expect(layout.DST[i*4+1]).To(Equal(uint32(i)))
		}
	})

	It("should run FLAT_LOAD_DWORDX4", func() {
		for i := 0; i < 64; i++ {
			pageTable.EXPECT().
				Find(vm.PID(1), uint64(i*16)).
				Return(vm.Page{
					PAddr: uint64(0),
				}, true)
		}
		inst := insts.NewInst()
		inst.FormatType = insts.FLAT
		inst.Opcode = 23
		inst.Addr = insts.NewVRegOperand(0, 0, 4)
		inst.Dst = insts.NewVRegOperand(4, 4, 4)
		wf.inst = inst

		for i := 0; i < 64; i++ {
			wf.WriteReg(insts.VReg(0), 4, i, uint64(i*16))
			// layout.ADDR[i] = uint64(i * 16)
			storage.Write(uint64(i*16), insts.Uint32ToBytes(uint32(i)))
			storage.Write(uint64(i*16+4), insts.Uint32ToBytes(uint32(i)))
			storage.Write(uint64(i*16+8), insts.Uint32ToBytes(uint32(i)))
			storage.Write(uint64(i*16+12), insts.Uint32ToBytes(uint32(i)))
		}
		wf.WriteReg(insts.Regs[insts.EXEC], 2, 0, 0xffffffffffffffff)

		alu.Run(wf)

		for i := 0; i < 64; i++ {
			results := wf.ReadReg(insts.VReg(4), 4, i)
			buf_0 := results & 0x000000000000ffff
			buf_1 := (results & 0x00000000ffff0000) >> 16
			buf_2 := (results & 0x0000ffff00000000) >> 32
			buf_3 := (results & 0xffff000000000000) >> 48

			Expect(buf_0).To(Equal(uint64(i)))
			Expect(buf_1).To(Equal(uint64(i)))
			Expect(buf_2).To(Equal(uint64(i)))
			Expect(buf_3).To(Equal(uint64(i)))
		}
	})

	// It("should run FLAT_STORE_DWORD", func() {
	// 	for i := 0; i < 64; i++ {
	// 		pageTable.EXPECT().
	// 			Find(vm.PID(1), uint64(i*4)).
	// 			Return(vm.Page{
	// 				PAddr: uint64(0),
	// 			}, true)
	// 	}
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.FLAT
	// 	state.inst.Opcode = 28

	// 	layout := state.Scratchpad().AsFlat()
	// 	for i := 0; i < 64; i++ {
	// 		layout.ADDR[i] = uint64(i * 4)
	// 		layout.DATA[i*4] = uint32(i)
	// 	}
	// 	layout.EXEC = 0xffffffffffffffff

	// 	alu.Run(state)

	// 	for i := 0; i < 64; i++ {
	// 		buf, err := storage.Read(uint64(i*4), uint64(4))
	// 		Expect(err).To(BeNil())
	// 		Expect(insts.BytesToUint32(buf)).To(Equal(uint32(i)))
	// 	}
	// })

	// It("should run FLAT_STORE_DWORDX2", func() {
	// 	for i := 0; i < 64; i++ {
	// 		pageTable.EXPECT().
	// 			Find(vm.PID(1), uint64(i*16)).
	// 			Return(vm.Page{
	// 				PAddr: uint64(0),
	// 			}, true)
	// 	}
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.FLAT
	// 	state.inst.Opcode = 29

	// 	layout := state.Scratchpad().AsFlat()
	// 	for i := 0; i < 64; i++ {
	// 		layout.ADDR[i] = uint64(i * 16)
	// 		layout.DATA[i*4] = uint32(i)
	// 		layout.DATA[(i*4)+1] = uint32(i)
	// 	}
	// 	layout.EXEC = 0xffffffffffffffff

	// 	alu.Run(state)

	// 	for i := 0; i < 64; i++ {
	// 		buf, err := storage.Read(uint64(i*16), uint64(16))
	// 		Expect(err).To(BeNil())
	// 		Expect(insts.BytesToUint32(buf[0:4])).To(Equal(uint32(i)))
	// 	}
	// })

	// It("should run FLAT_STORE_DWORDX3", func() {
	// 	for i := 0; i < 64; i++ {
	// 		pageTable.EXPECT().
	// 			Find(vm.PID(1), uint64(i*16)).
	// 			Return(vm.Page{
	// 				PAddr: uint64(0),
	// 			}, true)
	// 	}
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.FLAT
	// 	state.inst.Opcode = 30

	// 	layout := state.Scratchpad().AsFlat()
	// 	for i := 0; i < 64; i++ {
	// 		layout.ADDR[i] = uint64(i * 16)
	// 		layout.DATA[i*4] = uint32(i)
	// 		layout.DATA[(i*4)+1] = uint32(i)
	// 		layout.DATA[(i*4)+2] = uint32(i)
	// 	}
	// 	layout.EXEC = 0xffffffffffffffff

	// 	alu.Run(state)

	// 	for i := 0; i < 64; i++ {
	// 		buf, err := storage.Read(uint64(i*16), uint64(16))
	// 		Expect(err).To(BeNil())
	// 		Expect(insts.BytesToUint32(buf[0:4])).To(Equal(uint32(i)))
	// 		Expect(insts.BytesToUint32(buf[4:8])).To(Equal(uint32(i)))
	// 		Expect(insts.BytesToUint32(buf[8:12])).To(Equal(uint32(i)))
	// 	}
	// })

	// It("should run FLAT_STORE_DWORDX4", func() {
	// 	for i := 0; i < 64; i++ {
	// 		pageTable.EXPECT().
	// 			Find(vm.PID(1), uint64(i*16)).
	// 			Return(vm.Page{
	// 				PAddr: uint64(0),
	// 			}, true)
	// 	}
	// 	state.inst = insts.NewInst()
	// 	state.inst.FormatType = insts.FLAT
	// 	state.inst.Opcode = 31

	// 	layout := state.Scratchpad().AsFlat()
	// 	for i := 0; i < 64; i++ {
	// 		layout.ADDR[i] = uint64(i * 16)
	// 		layout.DATA[i*4] = uint32(i)
	// 		layout.DATA[(i*4)+1] = uint32(i)
	// 		layout.DATA[(i*4)+2] = uint32(i)
	// 		layout.DATA[(i*4)+3] = uint32(i)
	// 	}
	// 	layout.EXEC = 0xffffffffffffffff

	// 	alu.Run(state)

	// 	for i := 0; i < 64; i++ {
	// 		buf, err := storage.Read(uint64(i*16), uint64(16))
	// 		Expect(err).To(BeNil())
	// 		Expect(insts.BytesToUint32(buf[0:4])).To(Equal(uint32(i)))
	// 		Expect(insts.BytesToUint32(buf[4:8])).To(Equal(uint32(i)))
	// 		Expect(insts.BytesToUint32(buf[8:12])).To(Equal(uint32(i)))
	// 		Expect(insts.BytesToUint32(buf[12:16])).To(Equal(uint32(i)))
	// 	}
	// })
})
