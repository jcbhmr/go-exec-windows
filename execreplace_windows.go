package execreplace

import (
	"errors"
	"runtime"
	"runtime/debug"
	"syscall"

	"github.com/jcbhmr/go-execreplace/internal/unsaferuntime"
	"github.com/jcbhmr/go-execreplace/internal/windows"
)

func execReplace(argv0 string, argv []string, envv []string) error {
	prevGOMAXPROCS := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prevGOMAXPROCS)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	stdin, stdout, stderr := windows.Stdio()

	threads, err := windows.OpenOtherProcessThreads(windows.THREAD_TERMINATE)
	if err != nil {
		return err
	}
	defer func() {
		errs := []error{}
		for _, thread := range threads {
			err := syscall.CloseHandle(thread)
			if err != nil {
				errs = append(errs, err)
			}
		}
		err := errors.Join(errs...)
		if err != nil {
			panic(err)
		}
	}()

	processes, err := windows.OpenChildProcesses(syscall.PROCESS_TERMINATE)
	if err != nil {
		return err
	}
	defer func() {
		errs := []error{}
		for _, process := range processes {
			err := syscall.CloseHandle(process)
			if err != nil {
				errs = append(errs, err)
			}
		}
		err := errors.Join(errs...)
		if err != nil {
			panic(err)
		}
	}()

	err = windows.IgnoreCtrlC()
	if err != nil {
		return err
	}
	defer func() {
		err := windows.UnignoreCtrlC()
		if err != nil {
			panic(err)
		}
	}()

	_, hu, err := syscall.StartProcess(argv0, argv, &syscall.ProcAttr{
		Env:   envv,
		Files: []uintptr{uintptr(stdin), uintptr(stdout), uintptr(stderr)},
		Sys: &syscall.SysProcAttr{
			HideWindow: true,
		},
	})
	if err != nil {
		return err
	}
	h := syscall.Handle(hu)
	defer unsaferuntime.Throw("cannot return from execReplace")

	unsaferuntime.ProcPin()
	defer unsaferuntime.ProcUnpin()

	mp := unsaferuntime.MFromUintptr(unsaferuntime.Getm())
	gp := *mp.CurgPtr()
	pp := mp.PPtr().Ptr()
	*mp.PreemptOffPtr() = "execReplace"

	pp.RunnextPtr().Set(nil)

	i := 0
	for {
	}

	*pp.RunqPtr() = [256]unsaferuntime.Guintptr{}
	*pp.RunqheadPtr() = 0
	*pp.RunqtailPtr() = 0

	runtime.SetBlockProfileRate(0)
	runtime.SetCPUProfileRate(0)
	runtime.SetMutexProfileFraction(0)

	for _, thread := range threads {
		err := windows.TerminateThread(thread, 1)
		if err != nil {
			unsaferuntime.Throw(err.Error())
		}
		err = syscall.CloseHandle(thread)
		if err != nil {
			unsaferuntime.Throw(err.Error())
		}
	}

	for _, process := range processes {
		err := syscall.TerminateProcess(process, 1)
		if err != nil {
			unsaferuntime.Throw(err.Error())
		}
		err = syscall.CloseHandle(process)
		if err != nil {
			unsaferuntime.Throw(err.Error())
		}
	}

	debug.SetGCPercent(-1)

	exitCode, err := windows.WaitForExitCodeProcess(h)
	if err != nil {
		unsaferuntime.Throw(err.Error())
	}

	syscall.Exit(int(exitCode))
	panic("unreachable")
}
