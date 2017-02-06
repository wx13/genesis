package main

import (
	"path/filepath"

	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

func main() {

	inst := installer.New().Init()
	defer inst.Done()

	inst.AddTask(modules.Mkdir{Path: "/tmp/genesis_test"})

	inst.AddTask(modules.CopyFile{
		Src:  filepath.Join(inst.Tmpdir, "files/file.txt"),
		Dest: "/tmp/genesis_test/file.txt",
	})

}
