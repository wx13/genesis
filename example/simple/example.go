package main

import (
	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

func main() {

	inst := installer.New()
	defer inst.Done()

	inst.AddTask(modules.Mkdir{Path: "/tmp/genesis_example"})
	inst.AddTask(modules.Mkdir{Path: "/tmp/not_here", Absent: true})
	inst.AddTask(modules.Mkdir{Path: "/tmp/empty_dir", Empty: true})

	aptSection := installer.NewSection("Install some debian packages.")
	pkgs := []string{"git", "gitk", "tig", "screen"}
	for _, pkg := range pkgs {
		aptSection.AddTask(modules.Apt{Name: pkg})
	}
	inst.Add(aptSection)

}
