package main

import (
	"fmt"

	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

var inst *installer.Installer

func dotfiles() {

	inst.AddTask(modules.LineInFile{
		File:    "~/.bashrc",
		Line:    "source $HOME/.mybashrc",
		Pattern: "source $HOME/.mybashrc",
		Store:   inst.Store,
		Label:   "bashrc",
	})

	inst.AddTask(modules.CopyFile{
		DestFile: "~/.mybashrc",
		SrcFile:  "files/mybashrc",
		Store:    inst.Store,
	})

	inst.AddTask(modules.HttpGet{
		Url:   "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh",
		Dest:  "~/.git-prompt.sh",
		Store: inst.Store,
	})

	inst.AddTask(modules.CopyFile{
		DestFile: "~/.gitconfig",
		SrcFile:  "files/gitconfig",
		Store:    inst.Store,
	})

	inst.AddTask(modules.CopyFile{
		DestFile: "~/.screenrc",
		SrcFile:  "files/screenrc",
		Store:    inst.Store,
	})

	sshConfig()

}

func sshConfig() {

	// Ensure SSH directory exists, but don't remove it.
	if !inst.Remove {
		inst.AddTask(modules.Mkdir{Path: "~/.ssh"})
	}

	// Enable SSH persistence.
	inst.AddTask(modules.BlockInFile{
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
		inst.AddTask(modules.BlockInFile{
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

func main() {

	inst = installer.New()
	if inst == nil {
		panic("Unable to create an installer.")
	}
	defer inst.Done()

	dotfiles()

}
