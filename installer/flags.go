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

// ParseFlags does all the flag parsing for the installer.
func (inst *Installer) ParseFlags() {

	// Grab the executable name for usage printout.
	execName := filepath.Base(os.Args[0])

	// Main help screen.
	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Printf("  %s -h\n", execName)
		fmt.Printf("  %s (status|install|remove) [-verbose] [-tmpdir] [-dir] [-tags] [-skip-tags]\n", execName)
		fmt.Printf("  %s build [-x file] [dir...]\n", execName)
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
	runFlag := flag.NewFlagSet("run", flag.ContinueOnError)
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
		fmt.Println("User options:")
		fmt.Println("")
		inst.UserFlags.PrintDefaults()
		fmt.Println("")
	}

	buildFlag := flag.NewFlagSet("build", flag.ExitOnError)
	xName := buildFlag.String("x", "", "Specify the executable to append zip file to.  Useful for cross compiling.")
	buildFlag.Usage = func() {
		fmt.Println("")
		fmt.Println("Builds the self-extracting file from the executable. Packages up")
		fmt.Println("needed files (and only needed files) as a zip archive and appends")
		fmt.Println("to binary executable.  Optionally specify:")
		fmt.Println("  - the name of the binary (useful for cross compiling)")
		fmt.Println("  - a list of directories to collect files from")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Printf("  %s build [-x file] [list of directories]\n", execName)
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
		inst.UserFlags.Parse(os.Args[2:])
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

	inst.Tmpdir, _ = ioutil.TempDir(*tmpdir, "genesis")
	inst.Dir = genesis.ExpandHome(*dir)

}
