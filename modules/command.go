package modules

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/wx13/genesis"
)

type Command struct {
	Cmd          string
	Opts         []string
	PSPattern    string
	IgnoreErrors bool
	Timeout      time.Duration
}

func MakeCommand(cmd string, opts ...string) Command {
	return Command{Cmd: cmd, Opts: opts}
}

func (cmd Command) ID() string {
	return fmt.Sprintf("Command: %s %s", cmd.Cmd, strings.Join(cmd.Opts, " "))
}

func (cmd Command) Files() []string {
	return []string{}
}

func (cmd Command) Status() (genesis.Status, string, error) {
	if len(cmd.PSPattern) == 0 {
		return genesis.StatusUnknown, "Cannot discern whether a command has been run or not.", nil
	}
	isRunning, err := genesis.IsRunning(cmd.PSPattern)
	if err != nil {
		return genesis.StatusUnknown, "Cannot tell if the process is running.", err
	}
	if isRunning {
		return genesis.StatusPass, cmd.PSPattern + " is running", nil
	}
	return genesis.StatusFail, cmd.PSPattern + " is not running", err
}

func (cmd Command) Remove() (string, error) {
	return "Cannot undo a command.", nil
}

func (cmd Command) Install() (string, error) {

	// Create an output channel, so we can implement a timeout.
	type output struct {
		msg string
		err error
	}
	done := make(chan output)

	// Go run the command and write to channel when done.
	go func() {
		out, err := exec.Command(cmd.Cmd, cmd.Opts...).CombinedOutput()
		if cmd.IgnoreErrors {
			err = nil
		}
		msg := strings.Replace(string(out), "\n", "; ", -1)
		done <- output{msg, err}
	}()

	// If timeout is not set, set it to a really long time.
	if cmd.Timeout <= 0 {
		cmd.Timeout = time.Hour * 100000
	}

	// Wait for output, but with a timeout.
	select {
	case <-time.After(cmd.Timeout):
		if cmd.IgnoreErrors {
			return "Command timed out", nil
		}
		return "Command timed out", fmt.Errorf("command timed out")
	case out := <-done:
		return out.msg, out.err
	}

}
