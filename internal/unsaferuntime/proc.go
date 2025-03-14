//go:build gc && go1.24

package unsaferuntime

type G pointerOnly

func Gdestroy(gp *G) {
	mp := Getg().m
	pp := mp.p.ptr()

	casgstatus(gp, _Grunning, _Gdead)
	gcController.addScannableStack(pp, -int64(gp.stack.hi-gp.stack.lo))
	if isSystemGoroutine(gp, false) {
		sched.ngsys.Add(-1)
	}
	gp.m = nil
	locked := gp.lockedm != 0
	gp.lockedm = 0
	mp.lockedg = 0
	gp.preemptStop = false
	gp.paniconfault = false
	gp._defer = nil // should be true already but just in case.
	gp._panic = nil // non-nil for Goexit during panic. points at stack-allocated data.
	gp.writebuf = nil
	gp.waitreason = waitReasonZero
	gp.param = nil
	gp.labels = nil
	gp.timer = nil
	gp.syncGroup = nil

	if gcBlackenEnabled != 0 && gp.gcAssistBytes > 0 {
		// Flush assist credit to the global pool. This gives
		// better information to pacing if the application is
		// rapidly creating an exiting goroutines.
		assistWorkPerByte := gcController.assistWorkPerByte.Load()
		scanCredit := int64(assistWorkPerByte * float64(gp.gcAssistBytes))
		gcController.bgScanCredit.Add(scanCredit)
		gp.gcAssistBytes = 0
	}

	dropg()

	if GOARCH == "wasm" { // no threads yet on wasm
		gfput(pp, gp)
		return
	}

	if locked && mp.lockedInt != 0 {
		print("runtime: mp.lockedInt = ", mp.lockedInt, "\n")
		if mp.isextra {
			throw("runtime.Goexit called in a thread that was not created by the Go runtime")
		}
		throw("exited a goroutine internally locked to the OS thread")
	}
	gfput(pp, gp)
	if locked {
		// The goroutine may have locked this thread because
		// it put it in an unusual kernel state. Kill it
		// rather than returning it to the thread pool.

		// Return to mstart, which will release the P and exit
		// the thread.
		if GOOS != "plan9" { // See golang.org/issue/22227.
			gogo(&mp.g0.sched)
		} else {
			// Clear lockedExt on plan9 since we may end up re-using
			// this thread.
			mp.lockedExt = 0
		}
	}
}
