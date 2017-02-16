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

	// Use a template to create a new file.
	inst.AddTask(modules.Template{
		Src:  "files/file.txt.tmpl",
		Dest: "/tmp/genesis_test/file_from_template.txt",
		Vars: struct{ Name string }{Name: "Bob"},
	})

	// Get a file from the internet.
	inst.AddTask(modules.HttpGet{
		Url:  "https://raw.githubusercontent.com/wx13/genesis/master/README.md",
		Dest: "/tmp/genesis_test/README.md",
	})

}
