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

	Completed  bool
	AtBarrier  bool
	inst       *insts.Inst
	scratchpad Scratchpad

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

	wf.scratchpad = make([]byte, 4096)

	return wf
}

// Inst returns the instruction that the wavefront is executing
func (wf *Wavefront) Inst() *insts.Inst {
	return wf.inst
}

// Scratchpad returns the scratchpad that is associated with the wavefront
func (wf *Wavefront) Scratchpad() Scratchpad {
	return wf.scratchpad
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

// ReadOperand returns the operand value as uint64
// use slice buf to handle the case when operand is vgpr in inst X4, X8, X16
func (wf *Wavefront) readOperand(operand *insts.Operand, laneID int, buf []uint32) uint64 {
	switch operand.OperandType {
	case insts.RegOperand:

	}
}

// ReadReg returns the raw register value
//
//nolint:gocyclo
func (wf *Wavefront) ReadReg(reg *insts.Reg, regCount int, laneID int) uint64 {

	// There are some concerns in terms of reading VCC and EXEC (64 or 32? And how to decide?)
	var value uint64
	if reg.IsSReg() {
		if regCount == 1 {
			value = uint64(wf.SRegFile[reg.RegIndex()])
		} else {
			value = uint64(wf.SRegFile[reg.RegIndex()+1]) << 32
			value |= uint64(wf.SRegFile[reg.RegIndex()])
		}
	} else if reg.IsVReg() {
		if regCount == 1 {
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
	} else {
		log.Panicf("Register type %s not supported", reg.Name)
	}

	return value
}

// WriteReg returns the raw register value
//
//nolint:gocyclo
func (wf *Wavefront) WriteReg(
	reg *insts.Reg,
	regCount int,
	laneID int,
	data uint64,
) {
	if reg.IsSReg() {
		if regCount == 1 {
			wf.SRegFile[reg.RegIndex()] = uint32(data)
		} else {
			wf.SRegFile[reg.RegIndex()+1] = uint32(data >> 32)
			wf.SRegFile[reg.RegIndex()] = uint32(data)
		}
	} else if reg.IsVReg() {
		if regCount == 1 {
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
	} else {
		log.Panicf("Register type %s not supported", reg.Name)
	}
}
