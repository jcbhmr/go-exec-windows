package windows

import (
	"syscall"
)

func Stdio() (syscall.Handle, syscall.Handle, syscall.Handle) {
	stdin, err := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err != nil {
		panic("could not get stdin handle")
	}
	stdout, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		panic("could not get stdout handle")
	}
	stderr, err := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)
	if err != nil {
		panic("could not get stderr handle")
	}
	return stdin, stdout, stderr
}
