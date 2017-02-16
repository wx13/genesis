// Package installer is the installer for the genesis package.
// It handles file backup, manages history, reports on progress,
// and invents higher-order tasks (such as if-then tasks)>
package installer

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kardianos/osext"

	"github.com/wx13/genesis"
	"github.com/wx13/genesis/store"
)

// DoTags and SkipTags are global variables, because
// we need to modify them as we descend into sections
// and unmodify them on the way back up.
var DoTags, SkipTags []string

// Installer is a wrapper around modules to provide a nice
// interface for building an installer.
type Installer struct {
	Cmd       string
	Verbose   bool
	Facts     genesis.Facts
	Tasks     []genesis.Doer
	Dir       string
	Gendir    string
	DoTags    string
	SkipTags  string
	UserFlags *flag.FlagSet
	ExecName  string
	BuildDirs []string
}

// New creates a new installer object.
func New() *Installer {

	inst := Installer{}
	inst.Tasks = []genesis.Doer{}
	inst.UserFlags = flag.NewFlagSet("user", flag.ExitOnError)
	inst.UserFlags.Usage = func() {}
	return &inst

}

func (inst *Installer) Init() *Installer {

	inst.ParseFlags()

	// If "rerun" is specified, use the command history to
	// rewrite the command options.
	if inst.Cmd == "rerun" {
		line, err := Rerun(inst.Dir)
		if err == nil {
			os.Args = strings.Fields(line)
			inst.ParseFlags()
		}
	}

	if inst.Cmd == "build" {
		return inst
	}

	if inst.Cmd != "install" && inst.Cmd != "remove" && inst.Cmd != "status" {
		return inst
	}

	SkipTags = strings.Split(inst.SkipTags, ",")
	if len(inst.DoTags) == 0 {
		DoTags = []string{}
	} else {
		DoTags = strings.Split(inst.DoTags, ",")
	}

	var err error
	storedir := filepath.Join(inst.Dir, "store")
	genesis.Store, err = store.New(storedir)
	if err != nil {
		fmt.Println("Cannot access store directory.", err)
		os.Exit(1)
	}

	if inst.Cmd == "install" || inst.Cmd == "remove" {
		err := SaveHistory(inst.Dir, os.Args)
		if err != nil {
			fmt.Println("Error saving command history:", err)
		}
	}

	inst.Facts = genesis.GatherFacts()
	inst.extractFiles()

	return inst

}

func (inst *Installer) extractFiles() error {

	filename, _ := osext.Executable()

	zipRdr, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println("Couldn't extract files.", err, filename)
		return err
	}
	for _, file := range zipRdr.File {
		dest := filepath.Join(genesis.Tmpdir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(dest, file.FileInfo().Mode().Perm())
			continue
		}
		os.MkdirAll(filepath.Dir(dest), 0755)
		perms := file.FileInfo().Mode().Perm()
		out, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR, perms)
		if err != nil {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			continue
		}
		_, err = io.CopyN(out, rc, file.FileInfo().Size())
		if err != nil {
			continue
		}
		rc.Close()
		out.Close()

		mtime := file.FileInfo().ModTime()
		err = os.Chtimes(dest, mtime, mtime)
		if err != nil {
			continue
		}

	}

	return nil
}

// Done finishes up the installer process.
func (inst *Installer) Done() {

	switch inst.Cmd {

	case "remove":
		for k := len(inst.Tasks) - 1; k >= 0; k-- {
			task := inst.Tasks[k]
			task.Undo()
		}

	case "install":
		for _, task := range inst.Tasks {
			task.Do()
		}

	case "status":
		for _, task := range inst.Tasks {
			task.Status()
		}

	case "build":
		inst.Build()
		return

	}

	ReportSummary()
	inst.CleanUp()

}

// CleanUp removes the temporary directory.
func (inst *Installer) CleanUp() {
	fmt.Println("")
	os.RemoveAll(genesis.Tmpdir)
}

func SkipID(id string) string {
	id = genesis.StringHash(id)
	for _, tag := range SkipTags {
		if id == tag {
			return "skip"
		}
	}
	if len(DoTags) == 0 {
		return "do"
	}
	for _, tag := range DoTags {
		if id == tag {
			return "do"
		}
	}
	return "pass"
}

func EmptyDoTags() []string {
	doTags := make([]string, len(DoTags))
	copy(doTags, DoTags)
	DoTags = []string{}
	return doTags
}

func RestoreDoTags(doTags []string) {
	DoTags = doTags
}

func (inst *Installer) AddTask(module genesis.Module) {
	inst.Tasks = append(inst.Tasks, Task{module})
}

func (inst *Installer) Add(task genesis.Doer) {
	inst.Tasks = append(inst.Tasks, task)
}

func (inst *Installer) Files() []string {
	files := []string{}
	for _, task := range inst.Tasks {
		files = append(files, task.Files()...)
	}
	return files
}

func getHistoryFile(dir string) (string, string) {
	if len(dir) == 0 {
		usr, _ := user.Current()
		dir = usr.HomeDir
	}
	filename := filepath.Join(dir, "history.txt")
	return dir, filename
}

func GetHistory(dir string) []string {
	_, filename := getHistoryFile(dir)
	data, err := ioutil.ReadFile(filename)
	lines := []string{}
	if err == nil {
		l := strings.Split(string(data), "\n")
		for _, s := range l {
			if len(s) > 0 {
				lines = append(lines, s)
			}
		}
	}
	if len(lines) > 1000 {
		lines = lines[:1000]
	}
	return lines
}

func SaveHistory(dir string, cmd []string) error {

	lines := GetHistory(dir)
	dir, filename := getHistoryFile(dir)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	line := strings.Join(cmd, " ")
	if len(lines) == 0 || line != lines[0] {
		lines = append([]string{line}, lines...)
	}
	err = ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")+"\n"), 0666)

	return err
}
