package main

import (
	"flag"

	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

// Vars holds data for use in a template.
type Vars struct {
	Name string
}

func main() {

	// Use command line flags to populate data.
	vars := Vars{}
	flag.StringVar(&vars.Name, "name", "You", "Somebody's name")

	// Initialize the installer.
	inst := installer.New()
	if inst == nil {
		panic("Unable to create an installer.")
	}
	// This must be done, or else nothing will happen.
	defer inst.Done()

	//  The simplest thing is to add a "task".
	inst.AddTask(modules.Mkdir{Path: "/tmp/genesis_example"})

	// We can also group tasks into a section.
	section := installer.NewSection("Create some files")
	// If we access files from the zip, we must prefix with inst.Dir.
	section.AddTask(modules.Template{
		DestFile:     "/tmp/genesis_example/file.txt",
		TemplateFile: inst.Dir + "/files/file.txt.tmpl",
		Vars:         vars,
		Store:        inst.Store,
	})
	section.AddTask(modules.CopyFile{
		Dest:  "/tmp/genesis_example/foo",
		Src:   inst.Dir + "/files/file.txt.tmpl",
		Store: inst.Store,
	})
	section.AddTask(modules.File{
		Path: "/tmp/genesis_example/foo",
		Mode: 0755,
	})
	// Don't forget to add the section to the installer.
	inst.Add(section)

	netSection := installer.NewSection("Configure the network")
	netSection.AddTask(modules.LineInFile{
		File:    "/tmp/genesis_example/network.cfg",
		Pattern: "auto eth1",
		Line:    "auto eth1",
		Store:   inst.Store,
		Label:   "auto_eth1",
		Success: "^auto eth1",
	})
	netSection.AddTask(modules.BlockInFile{
		File:     "/tmp/genesis_example/network.cfg",
		Patterns: []string{"iface eth1", ""},
		Lines: []string{
			"iface eth1 inet static",
			"  address 10.1.10.2",
			"  netmask 255.255.255.0",
			"  gateway 10.1.10.1",
			"",
		},
		Label: "iface_eth1",
		Store: inst.Store,
	})
	inst.Add(netSection)

}
