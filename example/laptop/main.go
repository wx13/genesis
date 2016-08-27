package main

import (
	"fmt"

	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

var inst *installer.Installer

func aptGet() {

	sect := installer.NewSection("Apt-Get install some software")
	defer inst.Add(sect)

	pkgs := []string{"git", "gitk", "tig", "screen", "w3m"}
	for _, pkg := range pkgs {
		sect.AddTask(modules.Apt{Name: pkg})
	}

}

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
		DestFile: "~/.mybashrc",
		SrcFile:  "files/mybashrc",
		Store:    inst.Store,
	})

	sect.AddTask(modules.HttpGet{
		Url:   "https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh",
		Dest:  "~/.git-prompt.sh",
		Store: inst.Store,
	})

	sect.AddTask(modules.CopyFile{
		DestFile: "~/.gitconfig",
		SrcFile:  "files/gitconfig",
		Store:    inst.Store,
	})

	sect.AddTask(modules.CopyFile{
		DestFile: "~/.screenrc",
		SrcFile:  "files/screenrc",
		Store:    inst.Store,
	})

}

func sshConfig() {

	sect := installer.NewSection("SSH configuration")
	defer inst.Add(sect)

	// Ensure SSH directory exists, but don't remove it.
	if !inst.Remove {
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

func main() {

	inst = installer.New()
	if inst == nil {
		panic("Unable to create an installer.")
	}
	defer inst.Done()

	dotfiles()
	sshConfig()
	aptGet()

}
