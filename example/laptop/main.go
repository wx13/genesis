// This is a genesis configuration for my own personal laptop.
// It can serve as a simple example for setting up a real-world
// system.

package main

import (
	"fmt"
	"path"

	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

var inst *installer.Installer

func dotfiles() {

	sect := installer.NewSection("Configure dotfiles")
	defer inst.Add(sect)

	sect.AddTask(modules.LineInFile{
		File:    "~/.bashrc",
		Line:    "source $HOME/.mybashrc",
		Pattern: "source $HOME/.mybashrc",
		Store:   inst.Store,
		Label:   "bashrc",
	})

	sect.AddTask(modules.CopyFile{
		Dest:  "~/.mybashrc",
		Src:   path.Join(inst.Tmpdir, "files/mybashrc"),
		Store: inst.Store,
	})

	sect.AddTask(modules.Mkdir{Path: "~/.bash_functions"})
	sect.AddTask(modules.Mkdir{Path: "~/.bash_scripts"})
	sect.AddTask(modules.CopyFile{
		Dest:  "~/.bash_functions/battery.sh",
		Src:   path.Join(inst.Tmpdir, "files/bash_functions/battery.sh"),
		Store: inst.Store,
	})
	sect.AddTask(modules.CopyFile{
		Dest:  "~/.bash_scripts/battery.sh",
		Src:   path.Join(inst.Tmpdir, "files/bash_scripts/battery.sh"),
		Store: inst.Store,
	})

	sect.AddTask(modules.HttpGet{
		Url:   "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh",
		Dest:  "~/.git-prompt.sh",
		Store: inst.Store,
	})

	sect.AddTask(modules.CopyFile{
		Dest:  "~/.gitconfig",
		Src:   path.Join(inst.Tmpdir, "files/gitconfig"),
		Store: inst.Store,
	})

	sect.AddTask(modules.CopyFile{
		Dest:  "~/.screenrc",
		Src:   path.Join(inst.Tmpdir, "files/screenrc"),
		Store: inst.Store,
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
	sect.AddTask(modules.BlockInFile{
		File:     "~/.ssh/config",
		Patterns: []string{`^Host \*`, "^ControlPersist"},
		Lines: []string{
			"Host *",
			"ControlMaster auto",
			"ControlPath ~/.ssh/master-%r@%h:%p",
			"ControlPersist 30m",
		},
		Store: inst.Store,
		Label: "ssh_persistence",
	})

	// Disable host key checking on select local networks.
	ips := []string{"10.0.0.*", "10.0.1.*", "192.168.1.*"}
	for _, ip := range ips {
		sect.AddTask(modules.BlockInFile{
			File: "~/.ssh/config",
			Patterns: []string{
				fmt.Sprintf("^Host %s", ip),
				"^UserKnownHostsFile",
			},
			Lines: []string{
				fmt.Sprintf("Host %s", ip),
				"StrictHostKeyChecking no",
				"UserKnownHostsFile=/dev/null",
			},
			Store: inst.Store,
			Label: "disable_ssh_host_key_checking" + ip,
		})
	}
}

func raspbianSetup() {

	sect := installer.NewSection("Configure Raspbian")
	defer inst.Add(sect)

	// Configure monitor
	sect.AddTask(modules.LineInFile{
		File:    "/boot/config.txt",
		Pattern: "hdmi_group",
		Line:    "hdmi_group=2",
		Store:   inst.Store,
		Label:   "hdmi_group",
	})
	sect.AddTask(modules.LineInFile{
		File:    "/boot/config.txt",
		Pattern: "hdmi_mode",
		Line:    "hdmi_mode=82",
		Store:   inst.Store,
		Label:   "hdmi_mode",
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
