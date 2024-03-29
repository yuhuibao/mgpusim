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

	It("should run S_LOAD_DWORD", func() {
		pageTable.EXPECT().
			Find(vm.PID(1), uint64(1040)).
			Return(vm.Page{
				PAddr: uint64(0),
			}, true)
		inst := insts.NewInst()
		inst.FormatType = insts.SMEM
		inst.Opcode = 0
		inst.Base = insts.NewSRegOperand(0, 0, 2)
		inst.Offset = insts.NewSRegOperand(2, 2, 1)
		inst.Data = insts.NewSRegOperand(3, 3, 1)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, uint64(1024))
		wf.WriteReg(insts.SReg(2), 1, 0, uint64(16))
		storage.Write(uint64(1040), insts.Uint32ToBytes(217))

		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(3), 1, 0)
		Expect(results).To(Equal(uint64(217)))
	})

	It("should run S_LOAD_DWORDX2", func() {
		pageTable.EXPECT().
			Find(vm.PID(1), uint64(1040)).
			Return(vm.Page{
				PAddr: uint64(0),
			}, true)
		inst := insts.NewInst()
		inst.FormatType = insts.SMEM
		inst.Opcode = 1
		inst.Base = insts.NewSRegOperand(0, 0, 2)
		inst.Offset = insts.NewIntOperand(0, 16)
		inst.Data = insts.NewSRegOperand(2, 2, 2)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, uint64(1024))
		storage.Write(uint64(1040), insts.Uint32ToBytes(217))
		storage.Write(uint64(1044), insts.Uint32ToBytes(218))

		alu.Run(wf)
		results := wf.ReadReg(insts.SReg(2), 2, 0)
		Expect(results).To(Equal(uint64(218<<32 + 217)))
	})

	It("should run S_LOAD_DWORDX4", func() {
		pageTable.EXPECT().
			Find(vm.PID(1), uint64(1040)).
			Return(vm.Page{
				PAddr: uint64(0),
			}, true)
		inst := insts.NewInst()
		inst.FormatType = insts.SMEM
		inst.Opcode = 2
		inst.Base = insts.NewSRegOperand(0, 0, 2)
		inst.Offset = insts.NewIntOperand(0, 16)
		inst.Data = insts.NewSRegOperand(2, 2, 4)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, uint64(1024))
		storage.Write(uint64(1040), insts.Uint32ToBytes(217))
		storage.Write(uint64(1044), insts.Uint32ToBytes(218))
		storage.Write(uint64(1048), insts.Uint32ToBytes(219))
		storage.Write(uint64(1052), insts.Uint32ToBytes(220))

		alu.Run(wf)
		results := make([]uint32, 4)
		wf.ReadReg4Plus(insts.SReg(2), 4, 0, results)

		Expect(results[0]).To(Equal(uint32(217)))
		Expect(results[1]).To(Equal(uint32(218)))
		Expect(results[2]).To(Equal(uint32(219)))
		Expect(results[3]).To(Equal(uint32(220)))
	})

	It("should run S_LOAD_DWORDX8", func() {
		pageTable.EXPECT().
			Find(vm.PID(1), uint64(1040)).
			Return(vm.Page{
				PAddr: uint64(0),
			}, true)
		inst := insts.NewInst()
		inst.FormatType = insts.SMEM
		inst.Opcode = 3
		inst.Base = insts.NewSRegOperand(0, 0, 2)
		inst.Offset = insts.NewIntOperand(0, 16)
		inst.Data = insts.NewSRegOperand(2, 2, 8)
		wf.inst = inst

		wf.WriteReg(insts.SReg(0), 2, 0, uint64(1024))
		storage.Write(uint64(1040), insts.Uint32ToBytes(217))
		storage.Write(uint64(1044), insts.Uint32ToBytes(218))
		storage.Write(uint64(1048), insts.Uint32ToBytes(219))
		storage.Write(uint64(1052), insts.Uint32ToBytes(220))
		storage.Write(uint64(1056), insts.Uint32ToBytes(221))
		storage.Write(uint64(1060), insts.Uint32ToBytes(222))
		storage.Write(uint64(1064), insts.Uint32ToBytes(223))
		storage.Write(uint64(1068), insts.Uint32ToBytes(224))

		alu.Run(wf)
		results := make([]uint32, 8)
		wf.ReadReg4Plus(insts.SReg(2), 8, 0, results)

		Expect(results[0]).To(Equal(uint32(217)))
		Expect(results[1]).To(Equal(uint32(218)))
		Expect(results[2]).To(Equal(uint32(219)))
		Expect(results[3]).To(Equal(uint32(220)))
		Expect(results[4]).To(Equal(uint32(221)))
		Expect(results[5]).To(Equal(uint32(222)))
		Expect(results[6]).To(Equal(uint32(223)))
		Expect(results[7]).To(Equal(uint32(224)))
	})

	It("should run S_CBRANCH", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 2
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160 + 16*4)))
	})

	It("should run S_CBRANCH, when IMM is negative", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 2
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, int64ToBits(-32))
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 1024)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(1024 - 32*4)))
	})

	It("should run S_CBRANCH_SCC0", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 4
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)
		wf.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160 + 16*4)))
	})

	It("should run S_CBRANCH_SCC0, when IMM is negative", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 4
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, int64ToBits(-32))
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 1024)
		wf.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(1024 - 32*4)))
	})

	It("should skip S_CBRANCH_SCC0, if SCC is 1", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 4
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)
		wf.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160)))
	})

	It("should run S_CBRANCH_SCC1", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 5
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)
		wf.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160 + 16*4)))
	})

	It("should run S_CBRANCH_SCC1, when IMM is negative", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 5
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, int64ToBits(-32))
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 1024)
		wf.WriteReg(insts.Regs[insts.SCC], 1, 0, 1)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(1024 - 32*4)))
	})

	It("should skip S_CBRANCH_SCC1, if SCC is 0", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 5
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)
		wf.WriteReg(insts.Regs[insts.SCC], 1, 0, 0)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160)))
	})

	It("should run S_CBRANCH_VCCZ", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 6
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)
		wf.WriteReg(insts.Regs[insts.VCC], 1, 0, 0)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160 + 16*4)))
	})

	It("should run S_CBRANCH_VCCNZ", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 7
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)
		wf.WriteReg(insts.Regs[insts.VCC], 1, 0, 0xffffffffffffffff)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160 + 16*4)))
	})

	It("should run S_CBRANCH_EXECZ", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 8
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 0)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160 + 16*4)))
	})

	It("should run S_CBRANCH_EXECNZ", func() {
		inst := insts.NewInst()
		inst.FormatType = insts.SOPP
		inst.Opcode = 9
		inst.SImm16 = insts.NewSRegOperand(0, 0, 1)

		wf.inst = inst
		wf.WriteReg(insts.SReg(0), 1, 0, 16)
		wf.WriteReg(insts.Regs[insts.PC], 1, 0, 160)
		wf.WriteReg(insts.Regs[insts.EXEC], 1, 0, 1)

		alu.Run(wf)

		Expect(wf.PC).To(Equal(uint64(160 + 16*4)))
	})

})
