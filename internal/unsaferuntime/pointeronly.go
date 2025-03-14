//go:build gc && go1.24

package unsaferuntime

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type pointerOnly struct {
	_ [0]func()
	_ noCopy
}
