package modules

import (
	"errors"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wx13/genesis"
)

type Dpkg struct {
	Path   string
	Name   string
	Force  bool
	Absent bool
}

func (dpkg Dpkg) path() string {
	if dpkg.Path == "" {
		return ""
	}
	match, _ := regexp.MatchString("^[.]?/", dpkg.Path)
	if match {
		return dpkg.Path
	}
	return filepath.Join(genesis.Tmpdir, dpkg.Path)
}

func (dpkg Dpkg) ID() string {
	if dpkg.Absent {
		return "Dpkg remove " + dpkg.Name
	}
	name, _ := dpkg.packageName()
	return "Dpkg install " + name + " from " + dpkg.Path
}

func (dpkg Dpkg) Files() []string {
	if dpkg.Path == "" {
		return []string{}
	}
	return []string{dpkg.path()}
}

func (dpkg *Dpkg) packageName() (string, error) {
	cmd := exec.Command("dpkg-deb", "-W", "--showformat", "${Package}", dpkg.path())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	pkgName := string(output)
	return pkgName, nil
}

func (dpkg Dpkg) Install() (string, error) {
	var cmd *exec.Cmd
	if dpkg.Absent {
		cmd = exec.Command("dpkg", "-r", dpkg.Name)
	} else {
		if dpkg.Force {
			cmd = exec.Command("dpkg", "--force-depends", "-i", dpkg.path())
		} else {
			cmd = exec.Command("dpkg", "-i", dpkg.path())
		}
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(output)), err
	}
	return "Install was successful", nil
}

func (dpkg Dpkg) Remove() (string, error) {
	if dpkg.Absent {
		return "Not gonna unremove package", nil
	}
	pkgName, err := dpkg.packageName()
	if err != nil {
		return "Couldn't get package name.", err
	}
	cmd := exec.Command("dpkg", "-r", pkgName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(output)), err
	}
	return "Removal was successful", nil
}

func (dpkg Dpkg) Status() (genesis.Status, string, error) {

	var pkgName string
	var err error

	if dpkg.Absent {
		pkgName = dpkg.Name
	} else {
		pkgName, err = dpkg.packageName()
		if err != nil {
			return genesis.StatusFail, "Couldn't get package name.", err
		}
	}

	cmd := exec.Command("dpkg-query", "-W", "-f", "${Status}", pkgName)
	output, err := cmd.CombinedOutput()
	resp := strings.TrimSpace(string(output))
	if err != nil {
		if dpkg.Absent {
			return genesis.StatusPass, "Package is not installed", nil
		}
		return genesis.StatusFail, resp, err
	}

	words := strings.Split(string(output), " ")
	if len(words) >= 3 && words[2] == "installed" {
		if dpkg.Absent {
			return genesis.StatusFail, "Package is installed", nil
		}
		return genesis.StatusPass, "Package is installed", nil
	}

	if dpkg.Absent {
		return genesis.StatusPass, resp, nil
	}
	return genesis.StatusFail, resp, errors.New("Package not installed")

}
