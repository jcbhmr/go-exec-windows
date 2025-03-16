package unsaferuntime

import (
	_ "unsafe"
)

//go:linkname ProcPin runtime.procPin
func ProcPin() int

//go:linkname ProcUnpin runtime.procUnpin
func ProcUnpin()

// func Casgstatus(gp *G, oldval, newval uint32) {
// 	if (oldval&Gscan != 0) || (newval&Gscan != 0) || oldval == newval {
// 		fmt.Printf("unsaferuntime: casgstatus: oldval=%x, newval=%x\n", oldval, newval)
// 		Throw("casgstatus: bad incoming values")
// 	}
// 	// TODO
// }

func MyGdestroy(gp *G) {
	mp := MFromUintptr(Getm())
	pp := mp.PPtr().Ptr()

	gp.MyAtomicstatusPtr().Store(Gdead)

	*gp.MPtr() = nil
	locked := *gp.LockedmPtr() != 0
	*gp.LockedmPtr() = 0
	*mp.LockedgPtr() = 0
	*gp.PreemptStopPtr() = false
	*gp.PaniconfaultPtr() = false
	*gp.DeferPtr() = nil
	*gp.PanicPtr() = nil
	*gp.WritebufPtr() = nil
	*gp.WaitreasonPtr() = WaitReasonZero
	*gp.ParamPtr() = nil
	*gp.LabelsPtr() = nil
	*gp.TimerPtr() = nil
	*gp.SyncGroupPtr() = nil

	// dropg?

	_ = pp
	_ = locked
}
