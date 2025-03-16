package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/jcbhmr/go-execreplace"
)

func main() {
	path, err := exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		log.Fatal("still runs goroutines")
	}()
	log.Fatal(execreplace.ExecReplace(path, []string{"go", "mod", "graph"}, os.Environ()))
}
