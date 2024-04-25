package emu

import (
	"log"

	"github.com/sarchlab/akita/v3/mem/vm"
	"github.com/sarchlab/mgpusim/v3/insts"
	"github.com/sarchlab/mgpusim/v3/kernels"
)

// A Wavefront in the emu package is a wrapper for the kernels.Wavefront
type Wavefront struct {
	*kernels.Wavefront

	pid vm.PID

	Completed bool
	AtBarrier bool
	inst      *insts.Inst

	PC       uint64
	Exec     uint64
	SCC      byte
	VCC      uint64
	M0       uint32
	SRegFile []uint32
	VRegFile [][]uint32
	LDS      []byte
}

// NewWavefront returns the Wavefront that wraps the nativeWf
func NewWavefront(nativeWf *kernels.Wavefront) *Wavefront {
	wf := new(Wavefront)
	wf.Wavefront = nativeWf

	wf.SRegFile = make([]uint32, 102)
	wf.VRegFile = make([][]uint32, 64)
	for i := 0; i < 64; i++ {
		wf.VRegFile[i] = make([]uint32, 256)
	}

	return wf
}

// Inst returns the instruction that the wavefront is executing
func (wf *Wavefront) Inst() *insts.Inst {
	return wf.inst
}

// PID returns pid
func (wf *Wavefront) PID() vm.PID {
	return wf.pid
}

// SRegValue returns s(i)'s value
func (wf *Wavefront) SRegValue(i int) uint32 {
	return wf.SRegFile[i]
}

// VRegValue returns the value of v(i) of a certain lain
func (wf *Wavefront) VRegValue(lane int, i int) uint32 {
	return wf.VRegFile[lane][i]
}

// ReadReg returns the raw register value when regCount<=2
//
//nolint:gocyclo
func (wf *Wavefront) ReadReg(reg *insts.Reg, regCount int, laneID int) uint64 {
	// There are some concerns in terms of reading VCC and EXEC (64 or 32? And how to decide?)
	var value uint64
	if reg.IsSReg() {
		if regCount <= 1 {
			value = uint64(wf.SRegFile[reg.RegIndex()])
		} else {
			value = uint64(wf.SRegFile[reg.RegIndex()+1]) << 32
			value |= uint64(wf.SRegFile[reg.RegIndex()])
		}
	} else if reg.IsVReg() {
		if regCount <= 1 {
			value = uint64(wf.VRegFile[laneID][reg.RegIndex()])
		} else {
			value = uint64(wf.VRegFile[laneID][reg.RegIndex()+1]) << 32
			value |= uint64(wf.VRegFile[laneID][reg.RegIndex()])
		}
	} else if reg.RegType == insts.SCC {
		value = uint64(wf.SCC)
	} else if reg.RegType == insts.VCC {
		value = wf.VCC
	} else if reg.RegType == insts.VCCLO && regCount == 1 {
		value = wf.VCC & 0x00000000ffffffff
	} else if reg.RegType == insts.VCCHI && regCount == 1 {
		value = wf.VCC >> 32
	} else if reg.RegType == insts.VCCLO && regCount == 2 {
		value = wf.VCC
	} else if reg.RegType == insts.EXEC {
		value = wf.Exec
	} else if reg.RegType == insts.EXECLO && regCount == 2 {
		value = wf.Exec
	} else if reg.RegType == insts.M0 {
		value = uint64(wf.M0)
	} else if reg.RegType == insts.PC {
		value = wf.PC
	} else {
		log.Panicf("Register type %s not supported", reg.Name)
	}

	return value
}

// WriteReg returns the raw register value
//
//nolint:gocyclo
func (wf *Wavefront) WriteReg(reg *insts.Reg, regCount int, laneID int, data uint64) {
	if reg.IsSReg() {
		if regCount <= 1 {
			wf.SRegFile[reg.RegIndex()] = uint32(data)
		} else {
			wf.SRegFile[reg.RegIndex()+1] = uint32(data >> 32)
			wf.SRegFile[reg.RegIndex()] = uint32(data)
		}
	} else if reg.IsVReg() {
		if regCount <= 1 {
			wf.VRegFile[laneID][reg.RegIndex()] = uint32(data)
		} else {
			wf.VRegFile[laneID][reg.RegIndex()+1] = uint32(data >> 32)
			wf.VRegFile[laneID][reg.RegIndex()] = uint32(data)
		}
	} else if reg.RegType == insts.SCC {
		wf.SCC = byte(data)
	} else if reg.RegType == insts.VCC {
		wf.VCC = data
	} else if reg.RegType == insts.VCCLO && regCount == 2 {
		wf.VCC = data
	} else if reg.RegType == insts.VCCLO && regCount == 1 {
		wf.VCC &= uint64(0xffffffff00000000)
		wf.VCC |= data
	} else if reg.RegType == insts.VCCHI && regCount == 1 {
		wf.VCC &= uint64(0x00000000ffffffff)
		wf.VCC |= data << 32
	} else if reg.RegType == insts.EXEC {
		wf.Exec = data
	} else if reg.RegType == insts.EXECLO && regCount == 2 {
		wf.Exec = data
	} else if reg.RegType == insts.M0 {
		wf.M0 = uint32(data)
	} else if reg.RegType == insts.PC {
		wf.PC = data
	} else {
		log.Panicf("Register type %s not supported", reg.Name)
	}
}

// ReadReg2Plus return the raw register value when regCount > 2
func (wf *Wavefront) ReadReg2Plus(reg *insts.Reg, regCount int, laneID int, buf []uint32) {
	if reg.IsSReg() {
		copy(buf, wf.SRegFile[reg.RegIndex():reg.RegIndex()+regCount])
	} else if reg.IsVReg() {
		copy(buf, wf.VRegFile[laneID][reg.RegIndex():reg.RegIndex()+regCount])
	} else {
		log.Panicf("Register type %s not supported", reg.Name)
	}
}

// WriteReg2Plus write the raw register value when regCount > 2
func (wf *Wavefront) WriteReg2Plus(reg *insts.Reg, regCount int, laneID int, buf []uint32) {
	if reg.IsSReg() {
		copy(wf.SRegFile[reg.RegIndex():reg.RegIndex()+regCount], buf)
	} else if reg.IsVReg() {
		copy(wf.VRegFile[laneID][reg.RegIndex():reg.RegIndex()+regCount], buf)
	} else {
		log.Panicf("Register type %s not supported", reg.Name)
	}
}

//nolint:funlen,gocyclo
func (wf *Wavefront) InitWfRegs() {
	co := wf.CodeObject
	pkt := wf.Packet

	wf.PC = pkt.KernelObject + co.KernelCodeEntryByteOffset
	wf.Exec = wf.InitExecMask

	SGPRPtr := 0
	if co.EnableSgprPrivateSegmentBuffer() {
		// log.Printf("EnableSgprPrivateSegmentBuffer is not supported")
		//fmt.Printf("s%d SGPRPrivateSegmentBuffer\n", SGPRPtr/4)
		SGPRPtr += 4
	}

	if co.EnableSgprDispatchPtr() {
		wf.SRegFile[SGPRPtr+1] = uint32(wf.PacketAddress >> 32)
		wf.SRegFile[SGPRPtr] = uint32(wf.PacketAddress)
		//fmt.Printf("s%d SGPRDispatchPtr\n", SGPRPtr/4)
		SGPRPtr += 2
	}

	if co.EnableSgprQueuePtr() {
		log.Printf("EnableSgprQueuePtr is not supported")
		//fmt.Printf("s%d SGPRQueuePtr\n", SGPRPtr/4)
		SGPRPtr += 2
	}

	if co.EnableSgprKernelArgSegmentPtr() {
		wf.SRegFile[SGPRPtr+1] = uint32(pkt.KernargAddress >> 32)
		wf.SRegFile[SGPRPtr] = uint32(pkt.KernargAddress)
		//fmt.Printf("s%d SGPRKernelArgSegmentPtr\n", SGPRPtr/4)
		SGPRPtr += 2
	}

	if co.EnableSgprDispatchID() {
		log.Printf("EnableSgprDispatchID is not supported")
		//fmt.Printf("s%d SGPRDispatchID\n", SGPRPtr/4)
		SGPRPtr += 2
	}

	if co.EnableSgprFlatScratchInit() {
		log.Printf("EnableSgprFlatScratchInit is not supported")
		//fmt.Printf("s%d SGPRFlatScratchInit\n", SGPRPtr/4)
		SGPRPtr += 2
	}

	if co.EnableSgprPrivateSegementSize() {
		log.Printf("EnableSgprPrivateSegmentSize is not supported")
		//fmt.Printf("s%d SGPRPrivateSegmentSize\n", SGPRPtr/4)
		SGPRPtr += 1
	}

	if co.EnableSgprGridWorkGroupCountX() {
		wf.SRegFile[SGPRPtr] =
			(pkt.GridSizeX + uint32(pkt.WorkgroupSizeX) - 1) / uint32(pkt.WorkgroupSizeX)
		//fmt.Printf("s%d WorkGroupCountX\n", SGPRPtr/4)
		SGPRPtr += 1
	}

	if co.EnableSgprGridWorkGroupCountY() {
		wf.SRegFile[SGPRPtr] =
			(pkt.GridSizeY + uint32(pkt.WorkgroupSizeY) - 1) / uint32(pkt.WorkgroupSizeY)
		//fmt.Printf("s%d WorkGroupCountY\n", SGPRPtr/4)
		SGPRPtr += 1
	}

	if co.EnableSgprGridWorkGroupCountZ() {
		wf.SRegFile[SGPRPtr] =
			(pkt.GridSizeZ + uint32(pkt.WorkgroupSizeZ) - 1) / uint32(pkt.WorkgroupSizeZ)
		//fmt.Printf("s%d WorkGroupCountZ\n", SGPRPtr/4)
		SGPRPtr += 1
	}

	if co.EnableSgprWorkGroupIDX() {
		wf.SRegFile[SGPRPtr] =
			uint32(wf.WG.IDX)
		//fmt.Printf("s%d WorkGroupIdX\n", SGPRPtr/4)
		SGPRPtr += 1
	}

	if co.EnableSgprWorkGroupIDY() {
		wf.SRegFile[SGPRPtr] =
			uint32(wf.WG.IDY)
		//fmt.Printf("s%d WorkGroupIdY\n", SGPRPtr/4)
		SGPRPtr += 1
	}

	if co.EnableSgprWorkGroupIDZ() {
		wf.SRegFile[SGPRPtr] =
			uint32(wf.WG.IDZ)
		//fmt.Printf("s%d WorkGroupIdZ\n", SGPRPtr/4)
		SGPRPtr += 1
	}

	if co.EnableSgprWorkGroupInfo() {
		log.Printf("EnableSgprPrivateSegmentSize is not supported")
		SGPRPtr += 1
	}

	if co.EnableSgprPrivateSegmentWaveByteOffset() {
		log.Printf("EnableSgprPrivateSegentWaveByteOffset is not supported")
		SGPRPtr += 4
	}

	var x, y, z int
	for i := wf.FirstWiFlatID; i < wf.FirstWiFlatID+64; i++ {
		z = i / (wf.WG.SizeX * wf.WG.SizeY)
		y = i % (wf.WG.SizeX * wf.WG.SizeY) / wf.WG.SizeX
		x = i % (wf.WG.SizeX * wf.WG.SizeY) % wf.WG.SizeX
		laneID := i - wf.FirstWiFlatID

		wf.WriteReg(insts.VReg(0), 1, laneID, uint64(x))

		if co.EnableVgprWorkItemID() > 0 {
			wf.WriteReg(insts.VReg(1), 1, laneID, uint64(y))
		}

		if co.EnableVgprWorkItemID() > 1 {
			wf.WriteReg(insts.VReg(2), 1, laneID, uint64(z))
		}
	}
}
