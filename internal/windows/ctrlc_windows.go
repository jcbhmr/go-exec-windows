package windows

import (
	"sync"
	"syscall"
)

var setConsoleCtrlHandler = kernel32.NewProc("SetConsoleCtrlHandler")

var trueHandler = sync.OnceValue(func() uintptr {
	return syscall.NewCallback(func(_ uint32) uintptr {
		return 1
	})
})

func IgnoreCtrlC() error {
	r, _, err := setConsoleCtrlHandler.Call(trueHandler(), 1)
	if r == 0 {
		return err
	}
	return nil
}

func UnignoreCtrlC() error {
	r, _, err := setConsoleCtrlHandler.Call(trueHandler(), 0)
	if r == 0 {
		return err
	}
	return nil
}
