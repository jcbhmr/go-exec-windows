package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/jcbhmr/go-execreplace"
)

func main() {
	path, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(execreplace.ExecReplace(path, os.Args, os.Environ()))
}
