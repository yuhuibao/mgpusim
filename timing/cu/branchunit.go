package cu

import (
	"github.com/sarchlab/akita/v3/sim"
	"github.com/sarchlab/mgpusim/v3/emu"
	"github.com/sarchlab/mgpusim/v3/timing/wavefront"
)

// A BranchUnit performs branch operations
type BranchUnit struct {
	cu *ComputeUnit

	alu emu.ALU

	toRead  *wavefront.Wavefront
	toExec  *wavefront.Wavefront
	toWrite *wavefront.Wavefront

	isIdle bool
}

// NewBranchUnit creates a new branch unit, injecting the dependency of
// the compute unit.
func NewBranchUnit(
	cu *ComputeUnit,
	alu emu.ALU,
) *BranchUnit {
	u := new(BranchUnit)
	u.cu = cu
	u.alu = alu
	return u
}

// CanAcceptWave checks if the buffer of the read stage is occupied or not
func (u *BranchUnit) CanAcceptWave() bool {
	return u.toRead == nil
}

// IsIdle checks idleness
func (u *BranchUnit) IsIdle() bool {
	u.isIdle = (u.toRead == nil) && (u.toWrite == nil) && (u.toExec == nil)
	return u.isIdle
}

// AcceptWave moves one wavefront into the read buffer of the branch unit
func (u *BranchUnit) AcceptWave(
	wave *wavefront.Wavefront,
	now sim.VTimeInSec,
) {
	u.toRead = wave
}

// Run executes three pipeline stages that are controlled by the BranchUnit
func (u *BranchUnit) Run(now sim.VTimeInSec) bool {
	madeProgress := false
	madeProgress = u.runWriteStage(now) || madeProgress
	madeProgress = u.runExecStage(now) || madeProgress
	madeProgress = u.runReadStage(now) || madeProgress
	return madeProgress
}

func (u *BranchUnit) runReadStage(now sim.VTimeInSec) bool {
	if u.toRead == nil {
		return false
	}

	if u.toExec == nil {

		u.toExec = u.toRead
		u.toRead = nil

		return true
	}
	return false
}

func (u *BranchUnit) runExecStage(now sim.VTimeInSec) bool {
	if u.toExec == nil {
		return false
	}

	if u.toWrite == nil {
		u.alu.Run(u.toExec)

		u.toWrite = u.toExec
		u.toExec = nil
		return true
	}
	return false
}

func (u *BranchUnit) runWriteStage(now sim.VTimeInSec) bool {
	if u.toWrite == nil {
		return false
	}

	u.cu.logInstTask(now, u.toWrite, u.toWrite.DynamicInst(), true)

	u.toWrite.InstBuffer = nil
	u.cu.UpdatePCAndSetReady(u.toWrite)
	u.toWrite.InstBufferStartPC = u.toWrite.PC & 0xffffffffffffffc0
	u.toWrite = nil
	u.isIdle = false
	return true
}

// Flush clear the unit
func (u *BranchUnit) Flush() {
	u.toRead = nil
	u.toWrite = nil
	u.toExec = nil
}
