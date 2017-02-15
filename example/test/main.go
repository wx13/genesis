package main

import (
	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

func main() {

	inst := installer.New().Init()
	defer inst.Done()

	// Ensure a directory exists.
	inst.AddTask(modules.Mkdir{Path: "/tmp/genesis_test"})

	// Copy a file from the tempdir to the system.
	inst.AddTask(modules.CopyFile{
		Src:  "files/file.txt",
		Dest: "/tmp/genesis_test/file.txt",
	})

}
