package cu

import (
	"github.com/sarchlab/akita/v3/sim"
	"github.com/sarchlab/mgpusim/v3/protocol"
	"github.com/sarchlab/mgpusim/v3/timing/wavefront"
)

// A WfDispatcher initialize wavefronts
type WfDispatcher interface {
	DispatchWf(
		now sim.VTimeInSec,
		wf *wavefront.Wavefront,
		location protocol.WfDispatchLocation,
	)
}

// A WfDispatcherImpl will register the wavefront in wavefront pool and
// initialize all the registers
type WfDispatcherImpl struct {
	cu *ComputeUnit

	Latency int
}

// NewWfDispatcher creates a default WfDispatcher
func NewWfDispatcher(cu *ComputeUnit) *WfDispatcherImpl {
	d := new(WfDispatcherImpl)
	d.cu = cu
	d.Latency = 0
	return d
}

// DispatchWf starts or continues a wavefront dispatching process.
func (d *WfDispatcherImpl) DispatchWf(
	now sim.VTimeInSec,
	wf *wavefront.Wavefront,
	location protocol.WfDispatchLocation,
) {
	d.setWfInfo(wf, location)
	d.initRegisters(wf)
}

func (d *WfDispatcherImpl) setWfInfo(
	wf *wavefront.Wavefront,
	location protocol.WfDispatchLocation,
) {
	wf.SIMDID = location.SIMDID
	wf.SRegOffset = location.SGPROffset
	wf.VRegOffset = location.VGPROffset
	wf.LDSOffset = location.LDSOffset
	wf.PC = wf.Packet.KernelObject + wf.CodeObject.KernelCodeEntryByteOffset
	wf.Exec = wf.InitExecMask
}

//nolint:gocyclo,funlen
func (d *WfDispatcherImpl) initRegisters(wf *wavefront.Wavefront) {
	wf.InitWfRegs()
}
