package modules

import (
	"errors"
	"fmt"
	"os"

	"github.com/wx13/genesis"
)

type Mkdir struct {
	Path string
}

func (mkdir Mkdir) Describe() string {
	return fmt.Sprintf("Mkdir: %s", mkdir.Path)
}

func (mkdir Mkdir) ID() string {
	return "mkdir" + mkdir.Path
}

func (mkdir Mkdir) Status() (genesis.Status, string, error) {
	mkdir.Path = genesis.ExpandHome(mkdir.Path)
	s, err := os.Stat(mkdir.Path)
	if err != nil {
		return genesis.StatusFail, "No such file or directory.", err
	}
	if s.IsDir() {
		return genesis.StatusPass, "Directory exists.", nil
	} else {
		return genesis.StatusFail, "Path exists, but is not a directory.", errors.New("not a directory")
	}
}

func (mkdir Mkdir) Remove() (string, error) {
	mkdir.Path = genesis.ExpandHome(mkdir.Path)
	err := os.RemoveAll(mkdir.Path)
	if err != nil {
		return "Could not remove directory", err
	}
	return "Removed directory and all its contents", nil
}

func (mkdir Mkdir) Install() (string, error) {
	mkdir.Path = genesis.ExpandHome(mkdir.Path)
	err := os.MkdirAll(mkdir.Path, 0755)
	if err != nil {
		return "Could not make directory", err
	}
	return "Made directory", nil
}
