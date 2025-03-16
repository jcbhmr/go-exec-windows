package unsaferuntime

import (
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/jcbhmr/go-execreplace/internal/unsafereflect"
)

const (
	Gidle = iota
	Grunnable
	Grunning
	Gsyscall
	Gwaiting
	Gmoribund_unused
	Gdead
	Genqueue_unused
	Gcopystack
	Gpreempted
	Gscan          = 0x1000
	Gscanrunnable  = Gscan + Grunnable
	Gscanrunning   = Gscan + Grunning
	Gscansyscall   = Gscan + Gsyscall
	Gscanwaiting   = Gscan + Gwaiting
	Gscanpreempted = Gscan + Gpreempted
)

type G pointerOnly
type M pointerOnly
type P pointerOnly
type Guintptr uintptr
type Muintptr uintptr
type Puintptr uintptr

func (gp Guintptr) Ptr() *G { return (*G)(unsafe.Pointer(gp)) }
func (mp Muintptr) Ptr() *M { return (*M)(unsafe.Pointer(mp)) }
func (pp Puintptr) Ptr() *P { return (*P)(unsafe.Pointer(pp)) }

func (gp *Guintptr) Set(g *G) { *gp = Guintptr(unsafe.Pointer(g)) }
func (mp *Muintptr) Set(m *M) { *mp = Muintptr(unsafe.Pointer(m)) }
func (pp *Puintptr) Set(p *P) { *pp = Puintptr(unsafe.Pointer(p)) }

var gType = sync.OnceValue(func() reflect.Type {
	return unsafereflect.ToType(unsafereflect.TypesByString("*runtime.g")[0]).Elem()
})

var gMOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("m")
	if !ok {
		panic("m not in runtime.g")
	}
	return f.Offset
})

func (g *G) MPtr() **M {
	return (**M)(unsafe.Add(unsafe.Pointer(g), gMOffset()))
}

var gLockedmOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("lockedm")
	if !ok {
		panic("lockedm not in runtime.g")
	}
	return f.Offset
})

func (g *G) LockedmPtr() *Muintptr {
	return (*Muintptr)(unsafe.Add(unsafe.Pointer(g), gLockedmOffset()))
}

var gPreemptStopOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("preemptStop")
	if !ok {
		panic("preemptStop not in runtime.g")
	}
	return f.Offset
})

func (g *G) PreemptStopPtr() *bool {
	return (*bool)(unsafe.Add(unsafe.Pointer(g), gPreemptStopOffset()))
}

var gPaniconfaultOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("paniconfault")
	if !ok {
		panic("paniconfault not in runtime.g")
	}
	return f.Offset
})

func (g *G) PaniconfaultPtr() *bool {
	return (*bool)(unsafe.Add(unsafe.Pointer(g), gPaniconfaultOffset()))
}

var gDeferOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("_defer")
	if !ok {
		panic("_defer not in runtime.g")
	}
	return f.Offset
})

func (g *G) DeferPtr() **Defer {
	return (**Defer)(unsafe.Add(unsafe.Pointer(g), gDeferOffset()))
}

var gPanicOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("_panic")
	if !ok {
		panic("_panic not in runtime.g")
	}
	return f.Offset
})

func (g *G) PanicPtr() **Panic {
	return (**Panic)(unsafe.Add(unsafe.Pointer(g), gPanicOffset()))
}

var gWritebufOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("writebuf")
	if !ok {
		panic("writebuf not in runtime.g")
	}
	return f.Offset
})

func (g *G) WritebufPtr() *[]byte {
	return (*[]byte)(unsafe.Add(unsafe.Pointer(g), gWritebufOffset()))
}

var gWaitreasonOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("waitreason")
	if !ok {
		panic("waitreason not in runtime.g")
	}
	return f.Offset
})

func (g *G) WaitreasonPtr() *WaitReason {
	return (*WaitReason)(unsafe.Add(unsafe.Pointer(g), gWaitreasonOffset()))
}

var gParamOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("param")
	if !ok {
		panic("param not in runtime.g")
	}
	return f.Offset
})

func (g *G) ParamPtr() *unsafe.Pointer {
	return (*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(g), gParamOffset()))
}

var gLabelsOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("labels")
	if !ok {
		panic("labels not in runtime.g")
	}
	return f.Offset
})

func (g *G) LabelsPtr() *unsafe.Pointer {
	return (*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(g), gLabelsOffset()))
}

var gTimerOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("timer")
	if !ok {
		panic("timer not in runtime.g")
	}
	return f.Offset
})

func (g *G) TimerPtr() **Timer {
	return (**Timer)(unsafe.Add(unsafe.Pointer(g), gTimerOffset()))
}

var gSyncGroupOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("syncGroup")
	if !ok {
		panic("syncGroup not in runtime.g")
	}
	return f.Offset
})

func (g *G) SyncGroupPtr() **SynctestGroup {
	return (**SynctestGroup)(unsafe.Add(unsafe.Pointer(g), gSyncGroupOffset()))
}

var gAtomicstatusOffset = sync.OnceValue(func() uintptr {
	gt := gType()
	f, ok := gt.FieldByName("atomicstatus")
	if !ok {
		panic("atomicstatus not in runtime.g")
	}
	return f.Offset
})

func (g *G) MyAtomicstatusPtr() *atomic.Uint32 {
	return (*atomic.Uint32)(unsafe.Add(unsafe.Pointer(g), gAtomicstatusOffset()))
}

//go:linkname Getm runtime.getm
func Getm() uintptr

func MFromUintptr(mp uintptr) *M {
	return (*M)(unsafe.Pointer(mp))
}

var mType = sync.OnceValue(func() reflect.Type {
	return unsafereflect.ToType(unsafereflect.TypesByString("*runtime.m")[0]).Elem()
})

var mCurgOffset = sync.OnceValue(func() uintptr {
	mt := mType()
	f, ok := mt.FieldByName("curg")
	if !ok {
		panic("curg not in runtime.m")
	}
	return f.Offset
})

func (m *M) CurgPtr() **G {
	return (**G)(unsafe.Add(unsafe.Pointer(m), mCurgOffset()))
}

var mPOffset = sync.OnceValue(func() uintptr {
	mt := mType()
	f, ok := mt.FieldByName("p")
	if !ok {
		panic("p not in runtime.m")
	}
	return f.Offset
})

func (m *M) PPtr() *Puintptr {
	return (*Puintptr)(unsafe.Add(unsafe.Pointer(m), mPOffset()))
}

var mPreemptOffOffset = sync.OnceValue(func() uintptr {
	mt := mType()
	f, ok := mt.FieldByName("preemptoff")
	if !ok {
		panic("preemptoff not in runtime.m")
	}
	return f.Offset
})

func (m *M) PreemptOffPtr() *string {
	return (*string)(unsafe.Add(unsafe.Pointer(m), mPreemptOffOffset()))
}

var mLockedgOffset = sync.OnceValue(func() uintptr {
	mt := mType()
	f, ok := mt.FieldByName("lockedg")
	if !ok {
		panic("lockedg not in runtime.m")
	}
	return f.Offset
})

func (m *M) LockedgPtr() *Guintptr {
	return (*Guintptr)(unsafe.Add(unsafe.Pointer(m), mLockedgOffset()))
}

var pType = sync.OnceValue(func() reflect.Type {
	return unsafereflect.ToType(unsafereflect.TypesByString("*runtime.p")[0]).Elem()
})

var pRunqheadOffset = sync.OnceValue(func() uintptr {
	pt := pType()
	f, ok := pt.FieldByName("runqhead")
	if !ok {
		panic("runqhead not in runtime.p")
	}
	return f.Offset
})

func (p *P) RunqheadPtr() *uint32 {
	return (*uint32)(unsafe.Add(unsafe.Pointer(p), pRunqheadOffset()))
}

var pRunqtailOffset = sync.OnceValue(func() uintptr {
	pt := pType()
	f, ok := pt.FieldByName("runqtail")
	if !ok {
		panic("runqtail not in runtime.p")
	}
	return f.Offset
})

func (p *P) RunqtailPtr() *uint32 {
	return (*uint32)(unsafe.Add(unsafe.Pointer(p), pRunqtailOffset()))
}

var pRunqOffset = sync.OnceValue(func() uintptr {
	pt := pType()
	f, ok := pt.FieldByName("runq")
	if !ok {
		panic("runq not in runtime.p")
	}
	return f.Offset
})

func (p *P) RunqPtr() *[256]Guintptr {
	return (*[256]Guintptr)(unsafe.Add(unsafe.Pointer(p), pRunqOffset()))
}

var pRunnextOffset = sync.OnceValue(func() uintptr {
	pt := pType()
	f, ok := pt.FieldByName("runnext")
	if !ok {
		panic("runnext not in runtime.p")
	}
	return f.Offset
})

func (p *P) RunnextPtr() *Guintptr {
	return (*Guintptr)(unsafe.Add(unsafe.Pointer(p), pRunnextOffset()))
}

type WaitReason uint8

const (
	WaitReasonZero WaitReason = iota
)

type Defer pointerOnly
type Panic pointerOnly
