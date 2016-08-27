package main

import (
	"github.com/wx13/genesis/installer"
	"github.com/wx13/genesis/modules"
)

func main() {

	inst := installer.New()
	if inst == nil {
		panic("Unable to create an installer.")
	}
	defer inst.Done()

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

}
