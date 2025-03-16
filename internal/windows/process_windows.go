package windows

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

func WaitForExitCodeProcess(p syscall.Handle) (uint32, error) {
	event, err := syscall.WaitForSingleObject(p, syscall.INFINITE)
	if event == syscall.WAIT_FAILED {
		return 0, err
	}
	if event != syscall.WAIT_OBJECT_0 {
		return 0, errors.New("windows: unexpected result from WaitForSingleObject")
	}
	var exitCode uint32
	err = syscall.GetExitCodeProcess(p, &exitCode)
	return exitCode, err
}

func ChildProcessEntries() ([]syscall.ProcessEntry32, error) {
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(snapshot)

	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = syscall.Process32First(snapshot, &entry)
	if err != nil {
		return nil, err
	}

	entries := []syscall.ProcessEntry32{}

	pid := os.Getpid()
	for {
		if entry.ParentProcessID == uint32(pid) {
			entries = append(entries, entry)
		}
		err = syscall.Process32Next(snapshot, &entry)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return nil, err
		}
	}

	return entries, nil
}

func OpenChildProcesses(desiredAccess uint32) ([]syscall.Handle, error) {
	entries, err := ChildProcessEntries()
	if err != nil {
		return nil, err
	}

	processes := make([]syscall.Handle, 0, len(entries))
	for _, entry := range entries {
		p, err := syscall.OpenProcess(desiredAccess, false, entry.ProcessID)
		if err != nil {
			for _, p := range processes {
				err := syscall.CloseHandle(p)
				if err != nil {
					panic(err)
				}
			}
			return nil, err
		}
		processes = append(processes, p)
	}

	return processes, nil
}