package unsaferuntime

import _ "unsafe"

//go:linkname Throw runtime.throw
func Throw(s string)
