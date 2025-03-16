package execreplace_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestGowrapper(t *testing.T) {
	cmd := exec.Command("go", "run", "testdata/gowrapper.go", "version")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	t.Logf("$ %v", cmd)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGoroutine(t *testing.T) {
	cmd := exec.Command("go", "run", "testdata/goroutine.go")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	t.Logf("$ %v", cmd)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}
