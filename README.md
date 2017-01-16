# Genesis: A tool for building installers

Genesis is a Go package for building self-contained, self-extracting
installers. It is intended as a configuration management utility for
offline systems and devices, but it can be used for pretty much any
system.  Genesis can install, roll back, and report status.  You can
construct groups of tasks and define conditional dependencies.

Genesis is designed to configure a system from scratch. No initial setup
is required. Go's ability to cross-compile static binaries makes it easy
to generate installers for other platforms.


## Motivation

Most configuration management systems (chef, ansible, etc.) don't
work well for devices because they make some assumptions that don't
hold for these systems.  These assumptions are:

1. disposability
2. preconfiguration

These assumptions hold true for cloud servers and virtual machines, which
is the typical target of configuration management.  For hardware devices
and similar, these assumptions often fail.


## Design Goals

Genesis is built with a few key design goals in mind.

### 1. Minimize assumptions about the target system.

Genesis tries to make as few assumptions about the target system as we can.
This allows it to target more systems.  It also means less manual a priori setup.
Here are some assumptions Genesis tries to avoid:

- Network connectivity
- Specific OS
- Preinstalled software (python, ruby, bash, etc)
- Availability of command line utilities (grep, awk, etc).

### 2. Provide status information.

A configuration management system should do more than configure a system.
It should also describe the current status of the system and show how
the system's current state differs from the desired state.  This allows
the user to make informed decisions about what to do next, or how to modify
the installer.

### 3. Be able to roll back changes.

The Genesis software should be able to undo changes it makes to the system.
Obviously there are limits to this, and Genesis will never be perfect at this.
However, it should make reasonable attempts to restore the system to its
previous state.

### 4. Able to be run by a non-expert.

Hardware systems are often delivered to a customer, and the customer gains
complete control of the system.  It is critical that non-expert users be
able to update a system.


## Usage

Gensis consists of two parts: a set of modules, and an installer package.
Each module performs a specific configuration task, such as: creating a directory,
inserting a line in a file, or creating a user.  The official Genesis modules
are contained in the `modules` directory.  However, a module can live anywhere.
To create a custom module, it must only satisfy the `Module` interface.

The second component is the installer.  The installer handles all the complicated
stuff of figuring out which tasks should run, and storing information about previous
states.

### Creating an Installer

Here is a very simple example showing how to build an installer.

    package main

    import (
        "github.com/wx13/genesis/installer"
        "github.com/wx13/genesis/modules"
    )

    func main() {

        inst := installer.New()
        defer inst.Done()

        inst.AddTask(modules.Mkdir{Path: "/tmp/genesis_example"})

        aptSection := installer.NewSection("Install some debian packages.")
        pkgs := []string{"git", "gitk", "tig", "screen"}
        for _, pkg := range pkgs {
            aptSection.AddTask(modules.Apt{Name: pkg})
        }
        inst.Add(aptSection)

    }

which produces this:

![genesis screenshot](doc/genesis.png)

See the `example` directory for more examples.

There are a few pieces here. First off, we import the installer package
and the modules package, and create a new installer instance. Notice that
we defer the `inst.Done()` command. Here's why: when we add tasks to the
installer, those tasks don't get run right away. Which order they run in
could depend on circumstances, so they get stored up by the installer.
The `inst.Done()` command actually runs the installation.

To add a task, we can simply run `AddTask` with its argument being a module instance.
As you can see above, we can also create named sections, and add tasks to those sections.


### Building an Installer

There are two ways to build the installer: manually or with genesis's
assistance.  Both begin by building the exectable with:

    go build [FILE]

or

    GOOS=linux GOARCH=arm go build [FILE]

if you are cross-compiling.


#### Manual build

The installer extracts zip data from the end of itself. So you can
create the full installer by appending the zip data.

    zip -r files.zip files
    cat files.zip >> installer
    zip -A installer

The last command fixes the zip file indexing to account for the executable
prepended to the zip data.


### Automated build

There are a couple of issues with building manually.  First off, you must
make sure you zip up all the correct files with correct relative paths.
If you have different versions of your installer (e.g. installer versus updater),
you have to manually manage which files to zip for each.  Finally, you have to
remember to correct the zip file index or else your installer will fail.

Thankfully, genesis has a solution to this.  To build the self-extracting installer
from a binary, just run:

    ./installer build [list of dirs]

This will figure out which files are needed, look for them in the current directory,
and create a self-contained installer at 'installer.x'.  You can optionally specify
a list of directories to look for files in (instead of the current directory).

This will fail for cross-compiled binaries, because you won't be able to execute
the binary on the build system.  Instead run:

    go run installer.go -x installer build [list of dirs]


Running the installer
---------------------

To run the installer, place it on the target system and execute it with
one of the standard commands: "status", "install", or "remove".  Use the "-h"
flag to see the help screen.  Each of the above commands has its own help screen
as well.

