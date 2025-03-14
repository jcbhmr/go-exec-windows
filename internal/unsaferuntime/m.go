//go:build gc && go1.24

package unsaferuntime

import "github.com/jcbhmr/go-execreplace/internal/goexperiment"

type M struct {g0      *g     // goroutine with scheduling stack
	morebuf gobuf  // gobuf arg to morestack
	divmod  uint32 // div/mod denominator for arm - known to liblink (cmd/internal/obj/arm/obj5.go)

	// Fields not known to debuggers.
	procid          uint64            // for debuggers, but offset not hard-coded
	gsignal         *g                // signal-handling g
	goSigStack      gsignalStack      // Go-allocated signal handling stack
	sigmask         sigset            // storage for saved signal mask
	tls             [tlsSlots]uintptr // thread-local storage (for x86 extern register)
	mstartfn        func()
	curg            *g       // current running goroutine
	caughtsig       guintptr // goroutine running during fatal signal
	p               puintptr // attached p for executing go code (nil if not executing go code)
	nextp           puintptr
	oldp            puintptr // the p that was attached before executing a syscall
	id              int64
	mallocing       int32
	throwing        throwType
	preemptoff      string // if != "", keep curg running on this m
	locks           int32
	dying           int32
	profilehz       int32
	spinning        bool // m is out of work and is actively looking for work
	blocked         bool // m is blocked on a note
	newSigstack     bool // minit on C thread called sigaltstack
	printlock       int8
	incgo           bool          // m is executing a cgo call
	isextra         bool          // m is an extra m
	isExtraInC      bool          // m is an extra m that is not executing Go code
	isExtraInSig    bool          // m is an extra m in a signal handler
	freeWait        atomic.Uint32 // Whether it is safe to free g0 and delete m (one of freeMRef, freeMStack, freeMWait)
	needextram      bool
	g0StackAccurate bool // whether the g0 stack has accurate bounds
	traceback       uint8
	ncgocall        uint64        // number of cgo calls in total
	ncgo            int32         // number of cgo calls currently in progress
	cgoCallersUse   atomic.Uint32 // if non-zero, cgoCallers in use temporarily
	cgoCallers      *cgoCallers   // cgo traceback if crashing in cgo call
	park            note
	alllink         *m // on allm
	schedlink       muintptr
	lockedg         guintptr
	createstack     [32]uintptr // stack that created this thread, it's used for StackRecord.Stack0, so it must align with it.
	lockedExt       uint32      // tracking for external LockOSThread
	lockedInt       uint32      // tracking for internal lockOSThread
	mWaitList       mWaitList   // list of runtime lock waiters

	mLockProfile mLockProfile // fields relating to runtime.lock contention
	profStack    []uintptr    // used for memory/block/mutex stack traces

	// wait* are used to carry arguments from gopark into park_m, because
	// there's no stack to put them on. That is their sole purpose.
	waitunlockf          func(*g, unsafe.Pointer) bool
	waitlock             unsafe.Pointer
	waitTraceSkip        int
	waitTraceBlockReason traceBlockReason

	syscalltick uint32
	freelink    *m // on sched.freem
	trace       mTraceState

	// these are here because they are too large to be on the stack
	// of low-level NOSPLIT functions.
	libcall    libcall
	libcallpc  uintptr // for cpu profiler
	libcallsp  uintptr
	libcallg   guintptr
	winsyscall winlibcall // stores syscall parameters on windows

	vdsoSP uintptr // SP for traceback while in VDSO call (0 if not in call)
	vdsoPC uintptr // PC for traceback while in VDSO call

	// preemptGen counts the number of completed preemption
	// signals. This is used to detect when a preemption is
	// requested, but fails.
	preemptGen atomic.Uint32

	// Whether this is a pending preemption signal on this M.
	signalPending atomic.Uint32

	// pcvalue lookup cache
	pcvalueCache pcvalueCache

	dlogPerM

	mOS

	chacha8   chacha8rand.State
	cheaprand uint64

	// Up to 10 locks held by this m, maintained by the lock ranking code.
	locksHeldLen int
	locksHeld    [10]heldLockInfo

	// Size the runtime.m structure so it fits in the 2048-byte size class, and
	// not in the next-smallest (1792-byte) size class. That leaves the 11 low
	// bits of muintptr values available for flags, as required for
	// GOEXPERIMENT=spinbitmutex.
	_ [goexperiment.SpinbitMutexInt * 64 * goarch.PtrSize / 8]byte
	_ [goexperiment.SpinbitMutexInt * 700 * (2 - goarch.PtrSize / 4)]byte
}

//go:linkname Getm runtime.getm
func Getm() *M
