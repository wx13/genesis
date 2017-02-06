package modules

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/wx13/genesis"
)

type Apt struct {
	Name   string
	Absent bool
}

func (apt Apt) ID() string {
	if apt.Absent {
		return "Apt remove " + apt.Name
	}
	return "Apt install " + apt.Name
}

func (apt Apt) Files() []string {
	return []string{}
}

func (apt Apt) Install() (string, error) {
	var cmd *exec.Cmd
	if apt.Absent {
		cmd = exec.Command("apt-get", "--yes", "--force-yes", "remove", apt.Name)
	} else {
		cmd = exec.Command("apt-get", "--yes", "--force-yes", "install", apt.Name)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(output)), err
	}
	return "Install was successful", nil
}

func (apt Apt) Remove() (string, error) {
	var cmd *exec.Cmd
	if apt.Absent {
		cmd = exec.Command("apt-get", "--yes", "--force-yes", "install", apt.Name)
	} else {
		cmd = exec.Command("apt-get", "--yes", "--force-yes", "remove", apt.Name)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(output)), err
	}
	return "Removal was successful", nil
}

func (apt Apt) Status() (genesis.Status, string, error) {

	var err error

	cmd := exec.Command("dpkg-query", "-W", "-f", "${Status}", apt.Name)
	output, err := cmd.CombinedOutput()
	resp := strings.TrimSpace(string(output))
	if err != nil {
		return genesis.StatusFail, resp, err
	}

	words := strings.Split(string(output), " ")
	if len(words) >= 3 && words[2] == "installed" {
		if apt.Absent {
			return genesis.StatusFail, "Package is installed", nil
		}
		return genesis.StatusPass, "Package is installed", nil
	}

	if apt.Absent {
		return genesis.StatusPass, resp, nil
	}
	return genesis.StatusFail, resp, errors.New("Package not installed")

}
