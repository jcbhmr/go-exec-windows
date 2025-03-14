//go:build unix

package execreplace

import (
	"syscall"
)

var execReplace = syscall.Exec
