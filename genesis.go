package genesis

import (
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
)

// Facts stores discovered information about the target system.
type Facts struct {
	Arch     string
	ArchType string
	OS       string
	Hostname string
	Username string
}

// Status represents a Pass/Fail/Unknown.
type Status int

const (
	StatusPass Status = iota
	StatusFail
	StatusUnknown
)

// Module is an interface for all the modules.
type Module interface {
	Install() (string, error)
	Remove() (string, error)
	Status() (Status, string, error)
	Describe() string
	ID() string
}

// Doer can do and undo things.
type Doer interface {
	Do() (bool, error)
	Undo() (bool, error)
	Status() (Status, error)
	ID() string
}

func DoerHash(doer Doer) string {
	id := doer.ID()
	return StringHash(id)
}

func StringHash(id string) string {
	id = fmt.Sprintf("%x", md5.Sum([]byte(id)))
	return id[:6]
}

// FileExists is a helper function to check if a file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

// IsRunning checks to see if a process is running.
func IsRunning(pattern string) (bool, error) {
	out, err := exec.Command("pgrep", pattern).CombinedOutput()
	if err != nil {
		return false, err
	}
	if len(out) > 0 {
		return true, nil
	}
	return false, nil
}

// ExpandHome expands a leading tilde to the user's home directory.
func ExpandHome(name string) string {
	if name[0] != '~' {
		return name
	}
	user, err := user.Current()
	if err != nil {
		return name
	}
	return path.Join(user.HomeDir, name[1:])
}
