package modules

import (
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/wx13/genesis"
)

type File struct {
	Path   string
	Mode   os.FileMode
	Owner  string
	Absent bool
	Local  bool // Don't follow links
}

func (file File) Describe() string {
	return fmt.Sprintf("File: %+v", file)
}

func (file File) ID() string {
	return fmt.Sprintf("file %+v", file)
}

func (file File) Status() (genesis.Status, string, error) {
	stat, err := os.Stat(file.Path)
	if file.Local {
		stat, err = os.Lstat(file.Path)
	}
	if file.Absent {
		if err == nil {
			return genesis.StatusFail, "File exists", nil
		}
		return genesis.StatusPass, "File does not exist", nil
	}
	if err != nil {
		return genesis.StatusFail, "Cannot stat file.", err
	}
	if stat.Mode() != file.Mode {
		msg := fmt.Sprintf("File mode should be %o, but is %o", file.Mode, stat.Mode())
		return genesis.StatusFail, msg, fmt.Errorf("Incorrect file permissions")
	}
	return genesis.StatusPass, "File mode is correct.", nil
}

func (file File) Remove() (string, error) {
	return "Cannot undo a file operation (not supported yet).", nil
}

func (file File) Install() (string, error) {
	if file.Absent {
		err := os.Remove(file.Path)
		if err != nil {
			return "Failed to remove file", err
		}
		return "Successfully removed file", nil
	}
	if len(file.Owner) > 0 {
		user, err := user.Lookup(file.Owner)
		if err != nil {
			return "Cannot lookup owner.", err
		}
		uid, _ := strconv.Atoi(user.Uid)
		gid, _ := strconv.Atoi(user.Gid)
		err = os.Chown(file.Path, uid, gid)
		if err != nil {
			return "Cannot change ownership.", err
		}
	}
	err := os.Chmod(file.Path, file.Mode)
	if err != nil {
		return "Cannot change permissions.", err
	}
	return "Successfully changed permissions.", nil
}
