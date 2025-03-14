package execreplace

import (
	"errors"
	"os"
	"runtime"
	"sync"
	"syscall"
	_ "unsafe"
)

// https://github.com/golang/go/blob/go1.24.1/src/runtime/proc.go#L7132
//
//go:linkname runtime_procPin runtime.procPin
func runtime_procPin()

// https://github.com/golang/go/blob/go1.24.1/src/runtime/proc.go#L7152
//
//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin()

var prevGODEBUG string
var prevGOMAXPROCS int

func lockSync() {
	prevGOMAXPROCS = runtime.GOMAXPROCS(1)

	runtime.LockOSThread()

	prevGODEBUG = os.Getenv("GODEBUG")
	os.Setenv("GODEBUG", prevGODEBUG+",asyncpreemptoff=1")

	runtime_procPin()
}

func unlockSync() {
	runtime_procUnpin()

	os.Setenv("GODEBUG", prevGODEBUG)
	prevGODEBUG = ""

	runtime.UnlockOSThread()

	runtime.GOMAXPROCS(prevGOMAXPROCS)
	prevGOMAXPROCS = 0
}

type stdio struct {
	Stdin  syscall.Handle
	Stdout syscall.Handle
	Stderr syscall.Handle
}

// Get the standard I/O handles of the current process using the Windows API.
// This doesn't use the Go [os.Stdin], [os.Stdout], and [os.Stderr] variables
// because they may be mutated by user code. [syscall.Stdin], [syscall.Stdout],
// and [syscall.Stderr] are also not used because they are also var variables
// and may be mutated by user code.
//
// [exec(3)]: https://man7.org/linux/man-pages/man3/exec.3.html
func getStdio() (stdio, error) {
	stdin, err1 := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	stdout, err2 := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	stderr, err3 := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)
	return stdio{stdin, stdout, stderr}, errors.Join(err1, err2, err3)
}

// kernel32 is special-cased to always load from the system directory.
var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var setConsoleCtrlHandler = kernel32.NewProc("SetConsoleCtrlHandler")

var trueHandler = sync.OnceValue(func() uintptr {
	return syscall.NewCallback(func(_ uint32) uintptr {
		return 1
	})
})

func ignoreCtrlC() error {
	r, _, err := setConsoleCtrlHandler.Call(trueHandler(), 1)
	if r == 0 {
		return err
	}
	return nil
}

func unignoreCtrlC() error {
	r, _, err := setConsoleCtrlHandler.Call(trueHandler(), 0)
	if r == 0 {
		return err
	}
	return nil
}

// basicProcess is a Windows-specific delibarately simple parallel of the more
// complex [os.Process] type. We only care about:
//
//  1. Starting a process.
//  2. Waiting for the process to exit.
//  3. Getting the exit code of the process.
//
// Everything else (timing information upon completion, sending signals, etc.)
// is not needed. There's not even a Release method because for this use case
// we can rely on the process terminating to clean up resources.
type basicProcess struct {
	pid    int
	handle syscall.Handle
}

// Start a basic process given the provided arguments, environment variables,
// and standard I/O handles.
func startBasicProcess(argv0 string, argv []string, envv []string, stdio stdio) (*basicProcess, error) {
	pid, h, err := syscall.StartProcess(argv0, argv, &syscall.ProcAttr{
		Env:   envv,
		Files: []uintptr{uintptr(stdio.Stdin), uintptr(stdio.Stdout), uintptr(stdio.Stderr)},
	})
	if err != nil {
		return nil, err
	}
	return &basicProcess{pid, syscall.Handle(h)}, nil
}

// Wait for the process to exit and return the exit code.
func (p *basicProcess) Wait() (uint32, error) {
	event, err := syscall.WaitForSingleObject(p.handle, syscall.INFINITE)
	if event == syscall.WAIT_FAILED {
		return 0, err
	}
	if event != syscall.WAIT_OBJECT_0 {
		return 0, errors.New("execreplace: unexpected result from WaitForSingleObject")
	}
	var exitCode uint32
	err = syscall.GetExitCodeProcess(p.handle, &exitCode)
	return exitCode, err
}

// https://github.com/golang/go/blob/go1.24.1/src/runtime/panic.go#L1113
//
//go:linkname runtime_throw runtime.throw
func runtime_throw(s string)

func execReplace(argv0 string, argv []string, envv []string) error {
	lockSync()
	defer unlockSync()

	stdio, err := getStdio()
	if err != nil {
		return err
	}

	err = ignoreCtrlC()
	if err != nil {
		return err
	}
	defer func() {
		err := unignoreCtrlC()
		if err != nil {
			panic(err)
		}
	}()

	bp, err := startBasicProcess(argv0, argv, envv, stdio)
	if err != nil {
		return err
	}

	exitCode, err := bp.Wait()
	if err != nil {
		runtime_throw(err.Error())
	}

	// Skip some (but not all) of Go's cleanup hooks.
	syscall.Exit(int(exitCode))
	panic("unreachable")
}
