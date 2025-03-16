package unsafereflect

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type pointerOnly struct {
	_ noCopy
	_ [0]func()
}
