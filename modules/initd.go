package modules

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/wx13/genesis"
)

type Initd struct {
	Name string
}

func (initd Initd) Describe() string {
	return fmt.Sprintf("Initd: %s", initd.Name)
}

func (initd Initd) ID() string {
	return "initd" + initd.Name
}

func (initd Initd) Status() (genesis.Status, string, error) {
	out, err := exec.Command("service", initd.Name, "status").CombinedOutput()
	if err != nil {
		return genesis.StatusFail, "Error checking service status.", err
	}
	text := strings.Split(string(out), "\n")
	if len(text) < 3 {
		return genesis.StatusFail, "Service is not running.", nil
	}
	fields := strings.Fields(text[2])
	if len(fields) < 2 {
		return genesis.StatusFail, "Service is not running.", nil
	}
	if fields[1] == "active" {
		return genesis.StatusPass, "Service is running.", nil
	}
	return genesis.StatusFail, "Service is not running.", nil
}

func (initd Initd) Remove() (string, error) {
	exec.Command("service", initd.Name, "stop").Run()
	out, err := exec.Command("update-rc.d", initd.Name, "disable").CombinedOutput()
	if err != nil {
		return "Error removing init.d service: " + string(out), err
	}
	return "Successfully removed service", nil
}

func (initd Initd) Install() (string, error) {
	file := path.Join("/etc/init.d", initd.Name)
	out, err := exec.Command("chmod", "+x", file).CombinedOutput()
	if err != nil {
		return "Unable to set init.d script to executable. " + string(out), err
	}
	out, err = exec.Command("update-rc.d", initd.Name, "defaults").CombinedOutput()
	if err != nil {
		return "Error running update-rc.d. " + string(out), err
	}
	out, err = exec.Command("service", initd.Name, "restart").CombinedOutput()
	if err != nil {
		return "Error restarting service. " + string(out), err
	}
	return "Successfully installed service.", nil
}
