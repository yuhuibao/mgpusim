package emu

import (
	"log"
	"math"

	"github.com/sarchlab/akita/v3/mem/vm"
	"github.com/sarchlab/mgpusim/v3/insts"
)

// InstEmuState is the interface used by the emulator to track the instruction
// execution status.
type InstEmuState interface {
	PID() vm.PID
	Inst() *insts.Inst
	ReadReg(reg *insts.Reg, regCount int, laneID int) uint64
	WriteReg(reg *insts.Reg, regCount int, laneID int, data uint64)
	ReadReg2Plus(reg *insts.Reg, regCount int, laneID int, buf []uint32)
	WriteReg2Plus(reg *insts.Reg, regCount int, laneID int, buf []uint32)
}

// ALU does its jobs
type ALU interface {
	Run(state InstEmuState)
	ReadOperand(state InstEmuState, operand *insts.Operand, laneID int, buf []uint32) uint64
	WriteOperand(state InstEmuState, operand *insts.Operand, laneID int, data uint64, buf []uint32)

	SetLDS(lds []byte)
	LDS() []byte
}

// ALUImpl is where the instructions get executed.
type ALUImpl struct {
	storageAccessor *storageAccessor
	lds             []byte
}

// NewALU creates a new ALU with a storage as a dependency.
func NewALU(storageAccessor *storageAccessor) *ALUImpl {
	alu := new(ALUImpl)
	alu.storageAccessor = storageAccessor
	return alu
}

// SetLDS assigns the LDS storage to be used in the following instructions.
func (u *ALUImpl) SetLDS(lds []byte) {
	u.lds = lds
}

// LDS returns lds
func (u *ALUImpl) LDS() []byte {
	return u.lds
}

// ReadOperand returns the operand value as uint64
// use slice buf to handle the case when operand is vgpr in inst X4, X8, X16
func (u *ALUImpl) ReadOperand(state InstEmuState, operand *insts.Operand, laneID int, buf []uint32) uint64 {
	var value uint64
	switch operand.OperandType {
	case insts.RegOperand:
		if operand.RegCount <= 2 {
			value = state.ReadReg(operand.Register, operand.RegCount, laneID)
		} else {
			state.ReadReg2Plus(operand.Register, operand.RegCount, laneID, buf)
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
func (u *ALUImpl) WriteOperand(state InstEmuState, operand *insts.Operand, laneID int, data uint64, buf []uint32) {
	if operand.OperandType != insts.RegOperand {
		log.Panic("Can only write into reg operand")
	}

	if operand.RegCount <= 2 {
		state.WriteReg(operand.Register, operand.RegCount, laneID, data)
	} else {
		state.WriteReg2Plus(operand.Register, operand.RegCount, laneID, buf)
	}
}

// Run executes the instruction in the scatchpad of the InstEmuState
//
//nolint:gocyclo
func (u *ALUImpl) Run(state InstEmuState) {
	inst := state.Inst()
	// fmt.Printf("%s\n", inst.String(nil))

	switch inst.FormatType {
	case insts.SOP1:
		u.runSOP1(state)
	case insts.SOP2:
		u.runSOP2(state)
	case insts.SOPC:
		u.runSOPC(state)
	case insts.SMEM:
		u.runSMEM(state)
	case insts.VOP1:
		u.runVOP1(state)
	case insts.VOP2:
		u.runVOP2(state)
	case insts.VOP3a:
		u.runVOP3A(state)
	// case insts.VOP3b:
	// 	u.runVOP3B(state)
	case insts.VOPC:
		u.runVOPC(state)
	case insts.FLAT:
		u.runFlat(state)
	case insts.SOPP:
		u.runSOPP(state)
	case insts.SOPK:
		u.runSOPK(state)
	// case insts.DS:
	// 	u.runDS(state)
	default:
		log.Panicf("Inst format %s is not supported", inst.Format.FormatName)
	}
}
