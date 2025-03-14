//go:build ignore

package main

import (
	"os"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"
	_ "unsafe"
)

type runtime_m pointerOnly

//go:linkname runtime_getm runtime.getm
func runtime_getm() *runtime_m

//go:linkname reflect_typesByString reflect.typesByString
func reflect_typesByString(s string) []unsafe.Pointer

//go:linkname reflect_toType reflect.toType
func reflect_toType(t unsafe.Pointer) reflect.Type

var runtime_mType = sync.OnceValue(func() reflect.Type {
	return reflect_toType(reflect_typesByString("*runtime.m")[0]).Elem()
})

var runtime_m_g0Offset = sync.OnceValue(func() uintptr {
	mType := runtime_mType()
	g0Field, _ := mType.FieldByName("g0")
	return g0Field.Offset
})

func (m *runtime_m) G0() *runtime_g {
	return *(**runtime_g)(unsafe.Add(unsafe.Pointer(m), runtime_m_g0Offset()))
}

var runtime_m_curgOffset = sync.OnceValue(func() uintptr {
	mType := runtime_mType()
	curgField, _ := mType.FieldByName("curg")
	return curgField.Offset
})

func (m *runtime_m) Curg() *runtime_g {
	return *(**runtime_g)(unsafe.Add(unsafe.Pointer(m), runtime_m_curgOffset()))
}

var runtime_m_preemptoffOffset = sync.OnceValue(func() uintptr {
	mType := runtime_mType()
	preemptoffField, _ := mType.FieldByName("preemptoff")
	return preemptoffField.Offset
})

func (m *runtime_m) Preemptoff() string {
	return *(*string)(unsafe.Add(unsafe.Pointer(m), runtime_m_preemptoffOffset()))
}

func (m *runtime_m) SetPreemptoff(v string) {
	*(*string)(unsafe.Add(unsafe.Pointer(m), runtime_m_preemptoffOffset())) = v
}

var runtime_m_pOffset = sync.OnceValue(func() uintptr {
	mType := runtime_mType()
	pField, _ := mType.FieldByName("p")
	return pField.Offset
})

func (m *runtime_m) P() *runtime_p {
	return *(**runtime_p)(unsafe.Add(unsafe.Pointer(m), runtime_m_pOffset()))
}

type runtime_g pointerOnly

var runtime_gType = sync.OnceValue(func() reflect.Type {
	return reflect_toType(reflect_typesByString("*runtime.g")[0]).Elem()
})

func runtime_getg() *runtime_g {
	m := runtime_getm()
	return m.Curg()
}

type runtime_p pointerOnly

var runtime_p_runqheadOffset = sync.OnceValue(func() uintptr {
	pType := reflect_toType(reflect_typesByString("*runtime.p")[0]).Elem()
	runqheadField, _ := pType.FieldByName("runqhead")
	return runqheadField.Offset
})

func (p *runtime_p) Runqhead() uint32 {
	return *(*uint32)(unsafe.Add(unsafe.Pointer(p), runtime_p_runqheadOffset()))
}

func (p *runtime_p) SetRunqhead(v uint32) {
	*(*uint32)(unsafe.Add(unsafe.Pointer(p), runtime_p_runqheadOffset())) = v
}

var runtime_p_runqtailOffset = sync.OnceValue(func() uintptr {
	pType := reflect_toType(reflect_typesByString("*runtime.p")[0]).Elem()
	runqtailField, _ := pType.FieldByName("runqtail")
	return runqtailField.Offset
})

func (p *runtime_p) Runqtail() uint32 {
	return *(*uint32)(unsafe.Add(unsafe.Pointer(p), runtime_p_runqtailOffset()))
}

func (p *runtime_p) SetRunqtail(v uint32) {
	*(*uint32)(unsafe.Add(unsafe.Pointer(p), runtime_p_runqtailOffset())) = v
}

var runtime_p_runqOffset = sync.OnceValue(func() uintptr {
	pType := reflect_toType(reflect_typesByString("*runtime.p")[0]).Elem()
	runqField, _ := pType.FieldByName("runq")
	return runqField.Offset
})

func (p *runtime_p) Runq() [256]*runtime_g {
	return *(*[256]*runtime_g)(unsafe.Add(unsafe.Pointer(p), runtime_p_runqOffset()))
}

func (p *runtime_p) SetRunq(v [256]*runtime_g) {
	*(*[256]*runtime_g)(unsafe.Add(unsafe.Pointer(p), runtime_p_runqOffset())) = v
}

var runtime_p_runnextOffset = sync.OnceValue(func() uintptr {
	pType := reflect_toType(reflect_typesByString("*runtime.p")[0]).Elem()
	runnextField, _ := pType.FieldByName("runnext")
	return runnextField.Offset
})

func (p *runtime_p) Runnext() *runtime_g {
	return *(**runtime_g)(unsafe.Add(unsafe.Pointer(p), runtime_p_runnextOffset()))
}

func (p *runtime_p) SetRunnext(v *runtime_g) {
	*(*uintptr)(unsafe.Add(unsafe.Pointer(p), runtime_p_runnextOffset())) = uintptr(unsafe.Pointer(v))
}

//go:linkname runtime_procPin runtime.procPin
func runtime_procPin()

//go:linkname runtime_throw runtime.throw
func runtime_throw(s string)

func main() {
	go func() {
		time.Sleep(10 * time.Millisecond)
		panic("still running")
	}()

	runtime.LockOSThread()
	runtime.GOMAXPROCS(1)
	runtime_procPin()
	m := runtime_getm()
	m.SetPreemptoff("plsno")
	p := m.P()
	p.SetRunqhead(0)
	p.SetRunqtail(0)
	p.SetRunq([256]*runtime_g{})
	p.SetRunnext(nil)

	_, h, err := syscall.StartProcess("C:\\Windows\\System32\\cmd.exe", []string{"cmd", "/c", "echo", "hello"}, &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{uintptr(syscall.Stdin), uintptr(syscall.Stdout), uintptr(syscall.Stderr)},
	})
	if err != nil {
		runtime_throw(err.Error())
	}

	event, err := syscall.WaitForSingleObject(syscall.Handle(h), syscall.INFINITE)
	if event == syscall.WAIT_FAILED {
		runtime_throw(err.Error())
	}
	if event != syscall.WAIT_OBJECT_0 {
		runtime_throw("unexpected result from WaitForSingleObject")
	}

	println("done!")

	syscall.Exit(0)
}
