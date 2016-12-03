// Package installer is the installer for the genesis package.
// It handles file backup, manages history, reports on progress,
// and invents higher-order tasks (such as if-then tasks)>
package installer

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
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
	Cmd      string
	Verbose  bool
	Facts    genesis.Facts
	Store    *store.Store
	Tasks    []genesis.Doer
	Tmpdir   string
	Dir      string
	Gendir   string
	DoTags   string
	SkipTags string
}

// ParseFlags does all the flag parsing for the installer.
func (inst *Installer) ParseFlags() {

	// Grab the executable name for usage printout.
	execName := path.Base(os.Args[0])

	// Main help screen.
	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Printf("  %s -h\n", execName)
		fmt.Printf("  %s (status|install|remove) [-verbose] [-tmpdir] [-dir] [-tags] [-skip-tags]\n", execName)
		fmt.Printf("  %s build [dir...]\n", execName)
		fmt.Printf("  %s rerun\n", execName)
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("")
		fmt.Println("  status    Show the current installation.")
		fmt.Println("  install   Run the installer.")
		fmt.Println("  remove    Reverse the installation process.")
		fmt.Println("  rerun     Start a command prompt to search/view/edit/run previous commands.")
		fmt.Println("  build     Add file resources to executable to build a stand-alone installer.")
		fmt.Println("")
		fmt.Println("For details on individual command options, run './installer <cmd> -h'.")
		fmt.Println("")
		flag.PrintDefaults()
	}

	// Options for the "run" commands: install, remove, status.
	runFlag := flag.NewFlagSet("run", flag.ExitOnError)
	verbose := runFlag.Bool("verbose", false, "Verbose")
	tmpdir := runFlag.String("tmpdir", "", "Temp directory for unpacked files; empty string == default location")
	dir := runFlag.String("dir", "~/.genesis", "Storage directory for data. Defaults to ~/.genesis")
	doTags := runFlag.String("tags", "", "Specify comma-separated tags to run.  Defaults to all.")
	skipTags := runFlag.String("skip-tags", "", "Specify comma-separated tags to skip.  Defaults to none.")
	runFlag.Usage = func() {
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Printf("  %s (status|install|remove) [-verbose] [-tmpdir] [-storedir] [-tags] [-skip-tags]\n", execName)
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("")
		runFlag.PrintDefaults()
		fmt.Println("")
	}

	buildFlag := flag.NewFlagSet("build", flag.ExitOnError)
	buildFlag.Usage = func() {
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Printf("  %s build [list of directories]\n", execName)
		fmt.Println("")
	}
	rerunFlag := flag.NewFlagSet("rerun", flag.ExitOnError)

	// Print help screen if no arguments are given.
	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Get the subcommand.
	cmd := os.Args[1]

	// Parse the subcommand options.
	switch cmd {
	case "install", "remove", "status":
		runFlag.Parse(os.Args[2:])
	case "build":
		buildFlag.Parse(os.Args[2:])
	case "rerun":
		rerunFlag.Parse(os.Args[2:])
	default:
		flag.Usage()
		os.Exit(1)
	}

	inst.Cmd = cmd
	inst.Verbose = *verbose
	inst.DoTags = *doTags
	inst.SkipTags = *skipTags

	inst.Tmpdir, _ = ioutil.TempDir(*tmpdir, "genesis")
	inst.Dir = genesis.ExpandHome(*dir)

}

// New creates a new installer object.
func New() *Installer {

	inst := Installer{}
	inst.Tasks = []genesis.Doer{}
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
		return &inst
	}

	if inst.Cmd != "install" && inst.Cmd != "remove" && inst.Cmd != "status" {
		return &inst
	}

	SkipTags = strings.Split(inst.SkipTags, ",")
	if len(inst.DoTags) == 0 {
		DoTags = []string{}
	} else {
		DoTags = strings.Split(inst.DoTags, ",")
	}

	var err error
	inst.Store, err = store.New(inst.Dir)
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

	return &inst

}

func (inst *Installer) extractFiles() error {

	filename, _ := osext.Executable()

	zipRdr, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println("Couldn't extract files.", err, filename)
		return err
	}
	for _, file := range zipRdr.File {
		dest := path.Join(inst.Tmpdir, file.Name)
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
		inst.Build(os.Args[2:])
		return

	}

	ReportSummary()
	inst.CleanUp()

}

// CleanUp removes the temporary directory.
func (inst *Installer) CleanUp() {
	fmt.Println("")
	os.RemoveAll(inst.Tmpdir)
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

func (inst *Installer) Build(dirs []string) {

	fmt.Println("Building the self-contained executable...")

	// Create list of files to archive.
	files := []string{}
	for _, file := range inst.Files() {
		if strings.HasPrefix(file, inst.Tmpdir) {
			p, err := filepath.Rel(inst.Tmpdir, file)
			if err == nil {
				files = append(files, p)
			}
		}
	}

	execname, _ := osext.Executable()
	execbody, err := ioutil.ReadFile(execname)
	if err != nil {
		fmt.Println("Cannot read executable (self):", execname, err)
		return
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)
	w.SetOffset(int64(len(execbody)))

	// Add files to the archive.
	fmt.Println("Adding files to archive:")
	for _, file := range files {
		fmt.Println("   ", file)
		f, err := w.Create(file)
		if err != nil {
			fmt.Println("Cannot add file to archive:", file, err)
			continue
		}
		body, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Cannot read file:", file, err)
			continue
		}
		_, err = f.Write(body)
		if err != nil {
			fmt.Println("Cannot write file contents to archive:", file, err)
			continue
		}
	}

	err = w.Close()
	if err != nil {
		fmt.Println("Cannot close archive:", err)
	}

	execbody = append(execbody, buf.Bytes()...)

	err = ioutil.WriteFile(execname+".x", execbody, 0755)
	if err != nil {
		fmt.Println("Error writing to zip file:", err)
		return
	}

	fmt.Println("Done building archive.")

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
