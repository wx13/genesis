package installer

import (
	"flag"
	"fmt"
	"os/exec"
)

type Installer struct {
	Pretend bool
	Remove  bool
}

func New() *Installer {
	pretend := flag.Bool("p", false, "Prentend; don't actually install.")
	remove := flag.Bool("r", false, "Remove (uninstall).")
	flag.Parse()
	inst := Installer{
		Pretend: *pretend,
		Remove:  *remove,
	}
	return &inst
}

func (inst *Installer) Dpkg(path string) error {
	cmd := exec.Command("dpkg-deb", "-W", "--showformat", "'${Package}'", path)
	output, err := cmd.Output()
	pkgName := string(output)
	if err != nil {
		fmt.Println("Error: Could not get deb package name.", err)
		return err
	}
	if inst.Remove {
		fmt.Println("Removing debian package:", pkgName)
		if inst.Pretend {
			return nil
		}
		cmd = exec.Command("dpkg", "-r", pkgName)
	} else {
		fmt.Println("Installing debian package", pkgName, "from file", path)
		if inst.Pretend {
			return nil
		}
		cmd = exec.Command("dpkg", "-r", pkgName)
	}
	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Error:", output, err)
		return err
	}
	return nil
}
