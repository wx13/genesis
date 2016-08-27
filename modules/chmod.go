package modules

import (
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/wx13/genesis"
)

type Chmod struct {
	Path  string
	Mode  os.FileMode
	Owner string
}

func (chmod Chmod) Describe() string {
	return fmt.Sprintf("Chmod: file=%s mode=%o owner=%s", chmod.Path, chmod.Mode, chmod.Owner)
}

func (chmod Chmod) ID() string {
	return fmt.Sprintf("chmod %s %#v %s", chmod.Path, chmod.Mode, chmod.Owner)
}

func (chmod Chmod) Status() (genesis.Status, string, error) {
	stat, err := os.Stat(chmod.Path)
	if err != nil {
		return genesis.StatusFail, "Cannot stat file.", err
	}
	if stat.Mode() != chmod.Mode {
		msg := fmt.Sprintf("File mode should be %o, but is %o", chmod.Mode, stat.Mode())
		return genesis.StatusFail, msg, fmt.Errorf("Incorrect file permissions")
	}
	return genesis.StatusPass, "File mode is correct.", nil
}

func (chmod Chmod) Remove() (string, error) {
	return "Cannot undo a chmod.", nil
}

func (chmod Chmod) Install() (string, error) {
	if len(chmod.Owner) > 0 {
		user, err := user.Lookup(chmod.Owner)
		if err != nil {
			return "Cannot lookup owner.", err
		}
		uid, _ := strconv.Atoi(user.Uid)
		gid, _ := strconv.Atoi(user.Gid)
		err = os.Chown(chmod.Path, uid, gid)
		if err != nil {
			return "Cannot change ownership.", err
		}
	}
	err := os.Chmod(chmod.Path, chmod.Mode)
	if err != nil {
		return "Cannot change permissions.", err
	}
	return "Successfully changed permissions.", nil
}
