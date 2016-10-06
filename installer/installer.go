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
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kardianos/osext"
	"github.com/robbiev/dilemma"

	"github.com/wx13/genesis"
	"github.com/wx13/genesis/store"
)

var DoTags, SkipTags []string

// Installer is a wrapper around modules to provide a nice
// interface for building an installer.
type Installer struct {
	Status  bool
	Remove  bool
	Install bool
	Verbose bool
	Facts   genesis.Facts
	Dir     string
	Store   *store.Store
	Tasks   []genesis.Doer
}

// New creates a new installer object.
func New() *Installer {

	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("./installer -h")
		fmt.Println("./installer (-status|-install|-remove)")
		fmt.Println("")
		flag.PrintDefaults()
	}

	status := flag.Bool("status", false, "Status.")
	remove := flag.Bool("remove", false, "Remove (uninstall).")
	install := flag.Bool("install", false, "Install.")
	verbose := flag.Bool("verbose", false, "Verbose")
	tmpdir := flag.String("tempdir", "", "Temp directory; empty string == default location")
	storedir := flag.String("store", "", "Storage directory for snapshots. Defaults to user's home directory.")
	dotags := flag.String("tags", "", "Specify comma-separated tags to run.  Defaults to all.")
	skipTags := flag.String("skip-tags", "", "Specify comma-separated tags to skip.  Defaults to none.")
	rerun := flag.Bool("rerun", false, "Rerun a previous command.")
	flag.Parse()

	if *rerun {
		cmds := GetHistory(*storedir)
		if len(cmds) > 10 {
			cmds = cmds[:10]
		}
		menu := dilemma.Config{
			Title:   "Select a command:",
			Help:    "Use the up/down arrow keys, then enter to select.",
			Options: cmds,
		}
		cmd, cancel, err := dilemma.Prompt(menu)
		if err != nil || cancel == dilemma.CtrlC {
			fmt.Println("Cancelled.")
			return nil
		}
		os.Args = strings.Fields(cmd)
		flag.Parse()
	}

	inst := Installer{
		Status:  *status,
		Remove:  *remove,
		Install: *install,
		Verbose: *verbose,
		Tasks:   []genesis.Doer{},
	}

	if !(*install || *remove || *status) {
		return &inst
	}

	SkipTags = strings.Split(*skipTags, ",")
	if len(*dotags) == 0 {
		DoTags = []string{}
	} else {
		DoTags = strings.Split(*dotags, ",")
	}

	inst.Store = store.New(*storedir)
	if inst.Store == nil {
		return nil
	}

	if inst.Install || inst.Remove {
		err := SaveHistory(*storedir, os.Args)
		if err != nil {
			fmt.Println("Error saving command history:", err)
		}
	}

	inst.GatherFacts()
	inst.extractFiles(*tmpdir)

	return &inst

}

func (inst *Installer) extractFiles(tmpdir string) error {

	dir, err := ioutil.TempDir(tmpdir, "installer")
	if err != nil {
		return err
	}
	inst.Dir = dir

	filename, _ := osext.Executable()

	zipRdr, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println("Couldn't extract files.", err, filename)
		return err
	}
	for _, file := range zipRdr.File {
		dest := path.Join(inst.Dir, file.Name)
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

	switch {

	case inst.Remove:
		for k := len(inst.Tasks) - 1; k >= 0; k-- {
			task := inst.Tasks[k]
			task.Undo()
		}

	case inst.Install:
		for _, task := range inst.Tasks {
			task.Do()
		}

	case inst.Status:
		for _, task := range inst.Tasks {
			task.Status()
		}

	}

	ReportSummary()
	inst.CleanUp()

}

// CleanUp removes the temporary directory.
func (inst *Installer) CleanUp() {
	fmt.Println("")
	os.RemoveAll(inst.Dir)
}

// GatherFacts learns stuff about the target system.
func (inst *Installer) GatherFacts() {

	inst.Facts = genesis.Facts{}

	inst.Facts.ArchType = runtime.GOARCH
	inst.Facts.OS = runtime.GOOS
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err == nil {
		inst.Facts.Arch = strings.TrimSpace(string(output))
	}

	inst.Facts.Hostname, _ = os.Hostname()

	u, err := user.Current()
	if err != nil {
		inst.Facts.Username = u.Username
	}

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

func getHistoryFile(dir string) (string, string) {
	if len(dir) == 0 {
		usr, _ := user.Current()
		dir = usr.HomeDir
	}
	dir = path.Join(dir, ".genesis")
	filename := path.Join(dir, "history.txt")
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
