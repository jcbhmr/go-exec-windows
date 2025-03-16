package execreplace

// ExecReplace is an abstraction of [exec(3)] that works on Unix and Windows.
//
// Note that ExecReplace does not attempt to [os.LookPath] the provided argv0
// argument. It's up to the caller to ensure that argv0 is an absolute path to a
// valid executable.
//
// # Unix
//
// On Unix ExecReplace calls [syscall.Exec] with the provided arguments. Nothing
// else is done. This is the behavior that the Windows implementation attempts
// to mimic.
//
// # Windows
//
// On Windows ExecReplace starts a new child process. After starting it, it
// ignores the [os.Interrupt] signal (Ctrl+C) and waits for the child process to
// exit. When the child process exits, the parent process exits with the same
// exit code.
//
// On Windows ExecReplace will abort if an error occurs after the process has
// started.
//
// [exec(3)]: https://man7.org/linux/man-pages/man3/exec.3.html
func ExecReplace(argv0 string, argv []string, envv []string) error {
	return execReplace(argv0, argv, envv)
}
