package windows

import (
	"os"
	"structs"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var thread32First = kernel32.NewProc("Thread32First")
var thread32Next = kernel32.NewProc("Thread32Next")
var terminateThread = kernel32.NewProc("TerminateThread")
var getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
var openThread = kernel32.NewProc("OpenThread")

const THREAD_TERMINATE = 0x0001

type ThreadEntry32 struct {
	_              structs.HostLayout
	Size           uint32
	_              uint32
	ThreadID       uint32
	OwnerProcessID uint32
	BasePriority   uint32
	_              uint32
	_              uint32
}

func ProcessThreadEntries() ([]ThreadEntry32, error) {
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(snapshot)

	var entry ThreadEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	ok, _, err := thread32First.Call(uintptr(snapshot), uintptr(unsafe.Pointer(&entry)))
	if ok == 0 {
		return nil, err
	}

	entries := []ThreadEntry32{}

	pid := os.Getpid()
	for {
		if entry.OwnerProcessID == uint32(pid) {
			entries = append(entries, entry)
		}
		ok, _, err = thread32Next.Call(uintptr(snapshot), uintptr(unsafe.Pointer(&entry)))
		if ok == 0 {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return nil, err
		}
	}

	return entries, nil
}

func OpenOtherProcessThreads(desiredAccess uint32) ([]syscall.Handle, error) {
	entries, err := ProcessThreadEntries()
	if err != nil {
		return nil, err
	}

	threads := make([]syscall.Handle, 0, len(entries)-1)
	currentThreadID := GetCurrentThreadID()
	for _, entry := range entries {
		if entry.ThreadID == currentThreadID {
			continue
		}
		t, err := OpenThread(desiredAccess, false, entry.ThreadID)
		if err != nil {
			for _, t := range threads {
				err := syscall.CloseHandle(t)
				if err != nil {
					panic(err)
				}
			}
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, nil
}

func OpenThread(desiredAccess uint32, inheritHandle bool, id uint32) (syscall.Handle, error) {
	var inheritHandleU32 uint32
	if inheritHandle {
		inheritHandleU32 = 1
	}
	h, _, err := openThread.Call(uintptr(desiredAccess), uintptr(inheritHandleU32), uintptr(id))
	if h == 0 {
		return syscall.InvalidHandle, err
	}
	return syscall.Handle(h), nil
}

func GetCurrentThreadID() uint32 {
	id, _, _ := getCurrentThreadId.Call()
	return uint32(id)
}

func TerminateThread(t syscall.Handle, exitCode uint32) error {
	ok, _, err := terminateThread.Call(uintptr(t), uintptr(exitCode))
	if ok == 0 {
		return err
	}
	return nil
}
