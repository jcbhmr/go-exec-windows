package execreplace_test

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/jcbhmr/go-execreplace"
)

func ExampleExecReplace() {
	command := "go"
	args := []string{"version"}

	argv0, err := exec.LookPath(command)
	if err != nil {
		log.Fatal(err)
	}
	argv := append([]string{command}, args...)
	log.Fatal(execreplace.ExecReplace(argv0, argv, os.Environ()))
}

func ExampleExecReplace_jumper() {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		log.Fatal(err)
	}
	var exeExt string
	if runtime.GOOS == "windows" {
		exeExt = ".exe"
	}
	path := filepath.Join(filepath.Dir(exe), "real"+exeExt)
	log.Fatal(execreplace.ExecReplace(path, os.Args, os.Environ()))
}
