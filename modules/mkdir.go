package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wx13/genesis"
)

type Mkdir struct {
	Path   string
	Absent bool
	Empty  bool
}

func (mkdir Mkdir) ID() string {
	return fmt.Sprintf("Mkdir: %+v", mkdir)
}

func (mkdir Mkdir) Files() []string {
	return []string{}
}

func (mkdir Mkdir) isExist() bool {
	_, err := os.Stat(mkdir.Path)
	if err == nil {
		return true
	}
	return false
}

func (mkdir Mkdir) isDir() bool {
	s, err := os.Stat(mkdir.Path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func (mkdir Mkdir) isEmpty() bool {
	f, err := os.Open(mkdir.Path)
	if err != nil {
		return true
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err != nil {
		return true
	}
	return false
}

func (mkdir Mkdir) Status() (genesis.Status, string, error) {

	mkdir.Path = genesis.ExpandHome(mkdir.Path)

	// No matter what, if path points to a file (not a dir),
	// then we are in failure.
	if mkdir.isExist() && !mkdir.isDir() {
		return genesis.StatusFail, "This is a file, not a directory.", nil
	}

	if mkdir.Absent {
		if mkdir.isDir() {
			return genesis.StatusFail, "Directory exists.", nil
		}
		return genesis.StatusPass, "Directory does not exist.", nil
	}

	// Whether empty or not, directory must exist.
	if !mkdir.isDir() {
		return genesis.StatusFail, "Directory does not exist.", nil
	}

	if mkdir.Empty {
		if mkdir.isEmpty() {
			return genesis.StatusPass, "Directory is empty.", nil
		}
		return genesis.StatusFail, "Directory is not empty.", nil
	}

	return genesis.StatusPass, "Directory exists.", nil
}

func (mkdir Mkdir) Remove() (string, error) {
	if mkdir.Absent {
		return "Can't un-remove a directory", nil
	}
	mkdir.Path = genesis.ExpandHome(mkdir.Path)
	err := os.RemoveAll(mkdir.Path)
	if err != nil {
		return "Could not remove directory", err
	}
	return "Removed directory and all its contents", nil
}

func (mkdir Mkdir) Install() (string, error) {

	mkdir.Path = genesis.ExpandHome(mkdir.Path)

	if mkdir.Absent {
		err := os.RemoveAll(mkdir.Path)
		if err != nil {
			return "Could not remove directory", err
		}
		return "Removed directory and all its contents", nil
	}

	err := os.MkdirAll(mkdir.Path, 0755)
	if err != nil {
		return "Could not make directory", err
	}

	if mkdir.Empty {
		err := mkdir.makeEmpty()
		if err != nil {
			return "Could not remove directory contents.", err
		}
		return "Created empty directory", nil
	}

	return "Made directory", nil
}

func (mkdir Mkdir) makeEmpty() error {
	d, err := os.Open(mkdir.Path)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(mkdir.Path, name))
		if err != nil {
			return err
		}
	}
	return nil
}
