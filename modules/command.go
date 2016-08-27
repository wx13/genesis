package modules

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/wx13/genesis"
)

type Command struct {
	Cmd  string
	Opts []string
}

func (cmd Command) Describe() string {
	return fmt.Sprintf("Command: %s %s", cmd.Cmd, strings.Join(cmd.Opts, " "))
}

func (cmd Command) ID() string {
	return "command" + cmd.Cmd + strings.Join(cmd.Opts, "")
}

func (cmd Command) Status() (genesis.Status, string, error) {
	return genesis.StatusUnknown, "Cannot discern whether a command has been run or not.", nil
}

func (cmd Command) Remove() (string, error) {
	return "Cannot undo a command.", nil
}

func (cmd Command) Install() (string, error) {
	out, err := exec.Command(cmd.Cmd, cmd.Opts...).CombinedOutput()
	return string(out), err
}
