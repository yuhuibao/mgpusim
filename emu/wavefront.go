package emu

import (
	"log"
	"math"

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
func (wf *Wavefront) ReadOperand(operand *insts.Operand, laneID int, buf []uint32) uint64 {
	var value uint64
	switch operand.OperandType {
	case insts.RegOperand:
		if operand.RegCount <= 2 {
			value = wf.ReadReg(operand.Register, operand.RegCount, laneID)
		} else {
			wf.ReadRegMore(operand.Register, operand.RegCount, laneID, buf)
			value = uint64(operand.RegCount)
		}
	case insts.IntOperand:
		value = uint64(operand.IntValue)
	case insts.FloatOperand:
		value = uint64(math.Float32bits(float32(operand.FloatValue)))
	case insts.LiteralConstant:
		value = uint64(operand.LiteralConstant)
	default:
		log.Panicf("Operand %s is not supported", operand.String())
	}
	return value
}

// WriteOperand write the operand value either in data of uint64
// or in buf of slice of uint32 to handle the case when operand is vgpr in inst X4, X8, X16
func (wf *Wavefront) WriteOperand(operand *insts.Operand, laneID int, data uint64, buf []uint32) {
	if operand.OperandType != insts.RegOperand {
		log.Panic("Can only write into reg operand")
	}

	if operand.RegCount <= 2 {
		wf.WriteReg(operand.Register, operand.RegCount, laneID, data)
	}

}

// ReadReg returns the raw register value when regCount<=2
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
func (wf *Wavefront) WriteReg(reg *insts.Reg, regCount int, laneID int, data uint64) {
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

// ReadRegMore return the raw register value when regCount > 2
func (wf *Wavefront) ReadRegMore(reg *insts.Reg, regCount int, laneID int, buf []uint32) {
	if reg.IsSReg() {
		copy(buf, wf.SRegFile[reg.RegIndex():reg.RegIndex()+regCount])
	} else if reg.IsVReg() {
		copy(buf, wf.VRegFile[laneID][reg.RegIndex():reg.RegIndex()+regCount])
	} else {
		log.Panicf("Register type %s not supported", reg.Name)
	}
}

// WriteRegMore write the raw register value when regCount > 2
func (wf *Wavefront) WriteRegMore(reg *insts.Reg, regCount int, laneID int, buf []uint32) {
	if reg.IsSReg() {
		copy(wf.SRegFile[reg.RegIndex():reg.RegIndex()+regCount], buf)
	} else if reg.IsVReg() {
		copy(wf.VRegFile[laneID][reg.RegIndex():reg.RegIndex()+regCount], buf)
	} else {
		log.Panicf("Register type %s not supported", reg.Name)
	}
}
