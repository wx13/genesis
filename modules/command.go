package modules

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/wx13/genesis"
)

type Command struct {
	Cmd       string
	Opts      []string
	PSPattern string
}

func MakeCommand(cmd string, opts ...string) Command {
	return Command{Cmd: cmd, Opts: opts}
}

func (cmd Command) Describe() string {
	return fmt.Sprintf("Command: %s %s", cmd.Cmd, strings.Join(cmd.Opts, " "))
}

func (cmd Command) ID() string {
	return "command" + cmd.Cmd + strings.Join(cmd.Opts, "")
}

func (cmd Command) Status() (genesis.Status, string, error) {
	if len(cmd.PSPattern) == 0 {
		return genesis.StatusUnknown, "Cannot discern whether a command has been run or not.", nil
	}
	out, err := exec.Command("pgrep", cmd.PSPattern).CombinedOutput()
	if err != nil {
		return genesis.StatusUnknown, "Cannot discern whether a command has been run or not.", err
	}
	if len(out) > 0 {
		return genesis.StatusPass, cmd.PSPattern + " is running", nil
	}
	return genesis.StatusFail, cmd.PSPattern + " is not running", err
}

func (cmd Command) Remove() (string, error) {
	return "Cannot undo a command.", nil
}

func (cmd Command) Install() (string, error) {
	out, err := exec.Command(cmd.Cmd, cmd.Opts...).CombinedOutput()
	return strings.Replace(string(out), "\n", "; ", -1), err
}
