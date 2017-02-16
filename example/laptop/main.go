// This is a genesis configuration for my own personal laptop.
// It can serve as a simple example for setting up a real-world
// system.

package main

import (
	"fmt"

	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

var inst *installer.Installer

func dotfiles() {

	sect := installer.NewSection("Configure dotfiles")
	defer inst.Add(sect)

	sect.AddTask(modules.LineInFile{
		File:    "~/.bashrc",
		Line:    []string{"source $HOME/.mybashrc"},
		Pattern: []string{`source \$HOME/.mybashrc`},
	})

	sect.AddTask(modules.CopyFile{
		Dest: "~/.mybashrc",
		Src:  "files/mybashrc",
	})

	sect.AddTask(modules.Mkdir{Path: "~/.bash_functions"})
	sect.AddTask(modules.Mkdir{Path: "~/.bash_scripts"})
	sect.AddTask(modules.CopyFile{
		Dest: "~/.bash_functions/battery.sh",
		Src:  "files/bash_functions/battery.sh",
	})
	sect.AddTask(modules.CopyFile{
		Dest: "~/.bash_scripts/battery.sh",
		Src:  "files/bash_scripts/battery.sh",
	})

	sect.AddTask(modules.HttpGet{
		Url:  "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh",
		Dest: "~/.git-prompt.sh",
	})

	sect.AddTask(modules.CopyFile{
		Dest: "~/.gitconfig",
		Src:  "files/gitconfig",
	})

	sect.AddTask(modules.CopyFile{
		Dest: "~/.screenrc",
		Src:  "files/screenrc",
	})

}

func sshConfig() {

	sect := installer.NewSection("SSH configuration")
	defer inst.Add(sect)

	// Ensure SSH directory exists, but don't remove it.
	if inst.Cmd != "remove" {
		sect.AddTask(modules.Mkdir{Path: "~/.ssh"})
	}

	// Enable SSH persistence.
	sect.AddTask(modules.LineInFile{
		File:    "~/.ssh/config",
		Pattern: []string{`^Host \*`, "^ControlPersist"},
		Line: []string{
			"Host *",
			"ControlMaster auto",
			"ControlPath ~/.ssh/master-%r@%h:%p",
			"ControlPersist 30m",
		},
	})

	// Disable host key checking on select local networks.
	ips := []string{"10.0.0.*", "10.0.1.*", "192.168.1.*"}
	for _, ip := range ips {
		sect.AddTask(modules.LineInFile{
			File: "~/.ssh/config",
			Pattern: []string{
				fmt.Sprintf("^Host %s", ip),
				"^UserKnownHostsFile",
			},
			Line: []string{
				fmt.Sprintf("Host %s", ip),
				"StrictHostKeyChecking no",
				"UserKnownHostsFile=/dev/null",
			},
		})
	}
}

func raspbianSetup() {

	sect := installer.NewSection("Configure Raspbian")
	defer inst.Add(sect)

	// Configure monitor
	sect.AddTask(modules.LineInFile{
		File:    "/boot/config.txt",
		Pattern: []string{"hdmi_group"},
		Line:    []string{"hdmi_group=2"},
	})
	sect.AddTask(modules.LineInFile{
		File:    "/boot/config.txt",
		Pattern: []string{"hdmi_mode"},
		Line:    []string{"hdmi_mode=82"},
	})

}

func main() {

	inst = installer.New().Init()
	if inst == nil {
		panic("Unable to create an installer.")
	}
	defer inst.Done()

	if inst.Facts.Distro == "Raspbian" {
		raspbianSetup()
	}
	dotfiles()
	sshConfig()

}
