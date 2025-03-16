package unsaferuntime

import _ "unsafe"

//go:linkname Procyield runtime.procyield
func Procyield(cycles uint32)
