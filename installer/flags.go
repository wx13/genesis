// Package installer is the installer for the genesis package.
// It handles file backup, manages history, reports on progress,
// and invents higher-order tasks (such as if-then tasks)>
package installer

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wx13/genesis"
)

type FlagMerger struct {
	vars map[string]*string
}

func NewFlagMerger() *FlagMerger {
	return &FlagMerger{make(map[string]*string)}
}

func (fm *FlagMerger) IsDefined(name string) bool {
	_, ok := fm.vars[name]
	return ok
}

func (fm *FlagMerger) merge(mainFlag, otherFlag *flag.FlagSet) {
	otherFlag.VisitAll(func(f *flag.Flag) {
		s := mainFlag.String(f.Name, f.DefValue, f.Usage)
		fm.vars[f.Name] = s
	})
}

func (fm *FlagMerger) unmerge(otherFlag *flag.FlagSet) {
	for name, ptr := range fm.vars {
		f := otherFlag.Lookup(name)
		if f == nil {
			continue
		}
		f.Value.Set(*ptr)
	}
}

func errln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func errf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
}

// ParseFlags does all the flag parsing for the installer.
func (inst *Installer) ParseFlags() {

	// Grab the executable name for usage printout.
	execName := filepath.Base(os.Args[0])

	// Main help screen.
	flag.Usage = func() {
		errln("")
		errln("Usage:")
		errln("")
		errf("  %s -h\n", execName)
		errf("  %s (status|install|remove) [-verbose] [-tmpdir] [-dir] [-tags] [-skip-tags]\n", execName)
		errf("  %s build [-x file] [dir...]\n", execName)
		errf("  %s rerun\n", execName)
		errln("")
		errln("Commands:")
		errln("")
		errln("  status    Show the current installation.")
		errln("  install   Run the installer.")
		errln("  remove    Reverse the installation process.")
		errln("  rerun     Start a command prompt to search/view/edit/run previous commands.")
		errln("  build     Add file resources to executable to build a stand-alone installer.")
		errln("")
		errln("For details on individual command options, run './installer <cmd> -h'.")
		errln("")
		flag.PrintDefaults()
	}

	// Options for the "run" commands: install, remove, status.
	runFlag := flag.NewFlagSet("run", flag.ExitOnError)
	flagMerger := NewFlagMerger()
	for _, f := range inst.UserFlags {
		flagMerger.merge(runFlag, f)
	}
	verbose := runFlag.Bool("verbose", false, "Verbose")
	tmpdir := runFlag.String("tmpdir", "", "Temp directory for unpacked files; empty string == default location")
	dir := runFlag.String("dir", "~/.genesis", "Storage directory for data. Defaults to ~/.genesis")
	doTags := runFlag.String("tags", "", "Specify comma-separated tags to run.  Defaults to all.")
	skipTags := runFlag.String("skip-tags", "", "Specify comma-separated tags to skip.  Defaults to none.")
	runFlag.Usage = func() {
		errln("")
		errln("Usage:")
		errln("")
		errf("  %s (status|install|remove) [-verbose] [-tmpdir] [-storedir] [-tags] [-skip-tags]\n", execName)
		errln("")
		errln("Genesis options:")
		errln("")
		runFlag.VisitAll(func(f *flag.Flag) {
			if flagMerger.IsDefined(f.Name) {
				return
			}
			errf("  -%-14s", f.Name)
			errf("%s (default: %s)\n", f.Usage, f.DefValue)
		})
		errln("")
		errln("Other options")
		for _, userFlag := range inst.UserFlags {
			errln("")
			userFlag.VisitAll(func(f *flag.Flag) {
				errf("  -%-14s", f.Name)
				errf("%s (default:%s)\n", f.Usage, f.DefValue)
			})
		}
		errln("")
	}

	buildFlag := flag.NewFlagSet("build", flag.ExitOnError)
	xName := buildFlag.String("x", "", "Specify the executable to append zip file to.  Useful for cross compiling.")
	buildFlag.Usage = func() {
		errln("")
		errln("Builds the self-extracting file from the executable. Packages up")
		errln("needed files (and only needed files) as a zip archive and appends")
		errln("to binary executable.  Optionally specify:")
		errln("  - the name of the binary (useful for cross compiling)")
		errln("  - a list of directories to collect files from")
		errln("")
		errln("Usage:")
		errln("")
		errf("  %s build [-x file] [list of directories]\n", execName)
		errln("")
	}

	rerunFlag := flag.NewFlagSet("rerun", flag.ExitOnError)
	rerunFlag.Usage = func() {
		errln("")
		errln("Access command history.")
		errln("")
		errln("Because shell history is unreliable, genesis stores its own")
		errln("command history.  Running with the 'rerun' command will greet")
		errln("you with a command prompt.  This prompt behaves a like a normal")
		errln("readline prompt. Use arrows, ctrl-r, etc. to navigate and edit.")
		errln("")
		errln("Usage:")
		errln("")
		errf("  %s rerun\n", execName)
		errln("")
	}

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
		inst.BuildDirs = buildFlag.Args()
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
	inst.ExecName = *xName

	genesis.Tmpdir, _ = ioutil.TempDir(*tmpdir, "genesis")
	inst.Dir = genesis.ExpandHome(*dir)

	// Put user flags back where they came from.
	for _, f := range inst.UserFlags {
		flagMerger.unmerge(f)
	}

}
