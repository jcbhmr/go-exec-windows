package execreplace

// ExecReplace is an abstraction of [exec(3)] that works on Unix and Windows.
//
// On Unix ExecReplace calls [syscall.Exec] with the provided arguments. Nothing
// else is done. This is the behavior that the Windows implementation attempts
// to mimic.
//
// On Windows ExecReplace starts a new child process. After starting it, it
// ignores the [os.Interrupt] signal (Ctrl+C) and waits for the child process to
// exit. When the child process exits, the parent process exits with the same
// exit code.
//
// On Windows the Go runtime does not stop after a successful call to
// ExecReplace. It's up to the caller to ensure that no goroutines are running
// when calling ExecReplace.
//
// On Windows ExecReplace will panic if an error occurs after the process has
// started. At this point the state of the program is unknown and the program
// may be in an inconsistent state. It's not recommended to attempt to recover
// from any panics ExecReplace may cause.
//
// Note that ExecReplace does not attempt to [os.LookPath] the provided argv0
// argument. It's up to the caller to ensure that argv0 is an absolute path to a
// valid executable.
//
// [exec(3)]: https://man7.org/linux/man-pages/man3/exec.3.html
func ExecReplace(argv0 string, argv []string, envv []string) error {
	return execReplace(argv0, argv, envv)
}
