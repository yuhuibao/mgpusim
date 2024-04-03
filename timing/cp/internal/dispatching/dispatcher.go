package dispatching

import (
	"fmt"
	"log"

	"github.com/sarchlab/akita/v3/monitoring"
	"github.com/sarchlab/akita/v3/sim"
	"github.com/sarchlab/akita/v3/tracing"
	"github.com/sarchlab/mgpusim/v3/kernels"
	"github.com/sarchlab/mgpusim/v3/protocol"
	"github.com/sarchlab/mgpusim/v3/timing/cp/internal/resource"
)

// A Dispatcher is a sub-component of a command processor that can dispatch
// work-groups to compute units.
type Dispatcher interface {
	tracing.NamedHookable
	RegisterCU(cu resource.DispatchableCU)
	IsDispatching() bool
	StartDispatching(req *protocol.LaunchKernelReq)
	Tick(now sim.VTimeInSec) (madeProgress bool)
}

// A DispatcherImpl is a ticking component that can dispatch work-groups.
type DispatcherImpl struct {
	sim.HookableBase

	cp                     tracing.NamedHookable
	name                   string
	respondingPort         sim.Port
	dispatchingPort        sim.Port
	alg                    algorithm
	dispatching            *protocol.LaunchKernelReq
	currWG                 dispatchLocation
	cycleLeft              int
	numDispatchedWGs       int
	numCompletedWGs        int
	inflightWGs            map[string]dispatchLocation
	originalReqs           map[string]*protocol.MapWGReq
	latencyTable           []int
	constantKernelOverhead int

	monitor     *monitoring.Monitor
	progressBar *monitoring.ProgressBar
}

// Name returns the name of the dispatcher
func (d *DispatcherImpl) Name() string {
	return d.name
}

// RegisterCU allows the dispatcher to dispatch work-groups to the CU.
func (d *DispatcherImpl) RegisterCU(cu resource.DispatchableCU) {
	d.alg.RegisterCU(cu)
}

// IsDispatching checks if the dispatcher is dispatching another kernel.
func (d *DispatcherImpl) IsDispatching() bool {
	return d.dispatching != nil
}

// StartDispatching lets the dispatcher to start dispatch another kernel.
func (d *DispatcherImpl) StartDispatching(req *protocol.LaunchKernelReq) {
	d.mustNotBeDispatchingAnotherKernel()

	d.alg.StartNewKernel(kernels.KernelLaunchInfo{
		CodeObject: req.HsaCo,
		Packet:     req.Packet,
		PacketAddr: req.PacketAddress,
		WGFilter:   req.WGFilter,
	})
	d.dispatching = req

	d.numDispatchedWGs = 0
	d.numCompletedWGs = 0

	d.initializeProgressBar(req.ID)
}

func (d *DispatcherImpl) initializeProgressBar(kernelID string) {
	if d.monitor != nil {
		d.progressBar = d.monitor.CreateProgressBar(
			fmt.Sprintf("At %s, Kernel: %s, ", d.Name(), kernelID),
			uint64(d.alg.NumWG()),
		)
	}
}

func (d *DispatcherImpl) mustNotBeDispatchingAnotherKernel() {
	if d.IsDispatching() {
		panic("dispatcher is dispatching another request")
	}
}

// Tick updates the state of the dispatcher.
func (d *DispatcherImpl) Tick(now sim.VTimeInSec) (madeProgress bool) {
	if d.cycleLeft > 0 {
		d.cycleLeft--
		return true
	}

	if d.dispatching != nil {
		if d.kernelCompleted() {
			madeProgress = d.completeKernel(now) || madeProgress
		} else {
			madeProgress = d.dispatchNextWG(now) || madeProgress
		}
	}

	madeProgress = d.processMessagesFromCU(now) || madeProgress

	return madeProgress
}

func (d *DispatcherImpl) processMessagesFromCU(now sim.VTimeInSec) bool {
	msg := d.dispatchingPort.Peek()
	if msg == nil {
		return false
	}

	switch msg := msg.(type) {
	case *protocol.WGCompletionMsg:
		location, ok := d.inflightWGs[msg.RspTo]
		if !ok {
			return false
		}

		d.alg.FreeResources(location)
		delete(d.inflightWGs, msg.RspTo)
		d.numCompletedWGs++
		if d.numCompletedWGs == d.alg.NumWG() {
			d.cycleLeft = d.constantKernelOverhead
		}

		originalReq := d.originalReqs[msg.RspTo]
		delete(d.originalReqs, msg.RspTo)
		tracing.TraceReqFinalize(originalReq, d)

		if d.progressBar != nil {
			d.progressBar.MoveInProgressToFinished(1)
		}

		d.dispatchingPort.Retrieve(now)
		return true
	}

	return false
}

func (d *DispatcherImpl) kernelCompleted() bool {
	if d.currWG.valid {
		return false
	}

	if d.alg.HasNext() {
		return false
	}

	if d.numCompletedWGs < d.numDispatchedWGs {
		return false
	}

	return true
}

func (d *DispatcherImpl) completeKernel(now sim.VTimeInSec) (
	madeProgress bool,
) {
	req := d.dispatching

	rsp := protocol.NewLaunchKernelRsp(now, req.Dst, req.Src, req.ID)

	err := d.respondingPort.Send(rsp)
	if err == nil {
		d.dispatching = nil

		if d.monitor != nil {
			d.monitor.CompleteProgressBar(d.progressBar)
		}

		tracing.TraceReqComplete(req, d.cp)

		return true
	}

	return false
}

func (d *DispatcherImpl) dispatchNextWG(
	now sim.VTimeInSec,
) (madeProgress bool) {
	if !d.currWG.valid {
		if !d.alg.HasNext() {
			return false
		}

		d.currWG = d.alg.Next()
		if !d.currWG.valid {
			return false
		}
	}

	reqBuilder := protocol.MapWGReqBuilder{}.
		WithSrc(d.dispatchingPort).
		WithDst(d.currWG.cu).
		WithSendTime(now).
		WithPID(d.dispatching.PID).
		WithWG(d.currWG.wg)
	for _, l := range d.currWG.locations {
		reqBuilder = reqBuilder.AddWf(l)
	}
	req := reqBuilder.Build()
	err := d.dispatchingPort.Send(req)

	// fmt.Printf("%.10f, %d, %d\n", now, d.currWG.wg.IDX, d.currWG.cuID)

	if err == nil {
		d.currWG.valid = false
		d.numDispatchedWGs++
		d.inflightWGs[req.ID] = d.currWG
		d.originalReqs[req.ID] = req
		d.cycleLeft = d.latencyTable[len(d.currWG.locations)]

		if d.progressBar != nil {
			d.progressBar.IncrementInProgress(1)
		}

		tracing.TraceReqInitiate(req, d,
			tracing.MsgIDAtReceiver(d.dispatching, d.cp))

		return true
	}

	return false
}

// A DispatcherEmu is used in emu mode for dispatching work-groups
type DispatcherEmu struct {
	sim.HookableBase

	cp              tracing.NamedHookable
	name            string
	respondingPort  sim.Port
	dispatchingPort sim.Port
	cuPool          resource.CUResourcePool
	gridBuilder     kernels.GridBuilder
	dispatching     *protocol.LaunchKernelReq
	originalReqs    map[string]*protocol.MapWGReq

	isDoneDispatch bool
	isCompletedExe bool
}

// Name returns the name of the dispatcher
func (d *DispatcherEmu) Name() string {
	return d.name
}

// RegisterCU allows the dispatcher to dispatch work-groups to the CU.
func (d *DispatcherEmu) RegisterCU(cu resource.DispatchableCU) {
	d.cuPool.RegisterCU(cu)
}

// IsDispatching checks if the dispatcher is dispatching another kernel.
func (d *DispatcherEmu) IsDispatching() bool {
	return d.dispatching != nil
}

func (d *DispatcherEmu) mustNotBeDispatchingAnotherKernel() {
	if d.IsDispatching() {
		panic("dispatcher is dispatching another request")
	}
}

// StartDispatching lets the dispatcher to start dispatch another kernel.
func (d *DispatcherEmu) StartDispatching(req *protocol.LaunchKernelReq) {
	d.mustNotBeDispatchingAnotherKernel()

	d.gridBuilder.SetKernel(kernels.KernelLaunchInfo{
		CodeObject: req.HsaCo,
		Packet:     req.Packet,
		PacketAddr: req.PacketAddress,
		WGFilter:   req.WGFilter,
	})
	d.dispatching = req

	d.isDoneDispatch = false
	d.isCompletedExe = false

	// d.initializeProgressBar(req.ID)
}

// Tick updates the state of the dispatcher
func (d *DispatcherEmu) Tick(now sim.VTimeInSec) (madeProgress bool) {
	if d.dispatching != nil {
		if !d.isDoneDispatch {
			madeProgress = d.dispatchALLWG(now) || madeProgress
		}
		if d.isCompletedExe {
			madeProgress = d.completeKernel(now) || madeProgress
		}
	}
	madeProgress = d.processMessagesFromCU(now) || madeProgress

	return madeProgress
}

func (d *DispatcherEmu) processMessagesFromCU(now sim.VTimeInSec) bool {
	msg := d.dispatchingPort.Peek()
	if msg == nil {
		return false
	}

	switch msg := msg.(type) {
	case *protocol.WGCompletionMsg:
		originalReq, ok := d.originalReqs[msg.RspTo]
		if !ok {
			return false
		}
		delete(d.originalReqs, msg.RspTo)
		d.isCompletedExe = true

		tracing.TraceReqFinalize(originalReq, d)

		d.dispatchingPort.Retrieve(now)
		return true
	}
	return false
}

func (d *DispatcherEmu) completeKernel(now sim.VTimeInSec) (
	madeProgress bool,
) {
	req := d.dispatching

	rsp := protocol.NewLaunchKernelRsp(now, req.Dst, req.Src, req.ID)

	err := d.respondingPort.Send(rsp)
	if err == nil {
		d.dispatching = nil

		// if d.monitor != nil {
		// 	d.monitor.CompleteProgressBar(d.progressBar)
		// }

		tracing.TraceReqComplete(req, d.cp)

		return true
	}

	return false
}

func (d *DispatcherEmu) dispatchALLWG(now sim.VTimeInSec) (madeProgress bool) {
	numWGs := d.gridBuilder.NumWG()
	wgReqCollection := make(map[int][]*kernels.WorkGroup)
	cuID := 0
	for i := 0; i < numWGs; i++ {
		currWG := d.gridBuilder.NextWG()
		wgReqCollection[cuID] = append(wgReqCollection[cuID], currWG)
		cuID = (cuID + 1) % d.cuPool.NumCU()
	}

	for cuID, wgs := range wgReqCollection {
		cuPort := d.cuPool.GetCU(cuID).DispatchingPort()
		reqBuilder := protocol.MapWGReqBuilder{}.
			WithSrc(d.dispatchingPort).
			WithDst(cuPort).
			WithSendTime(now).
			WithPID(d.dispatching.PID).
			WithWGs(wgs)

		req := reqBuilder.Build()
		err := d.dispatchingPort.Send(req)

		if err == nil {
			d.originalReqs[req.ID] = req

			tracing.TraceReqInitiate(req, d,
				tracing.MsgIDAtReceiver(d.dispatching, d.cp))
		} else {
			log.Panicf("Fail to send MapWGReq to dispatch WGs to CU#%d in emu\n", cuID)
		}
	}
	d.isDoneDispatch = true

	return true
}
