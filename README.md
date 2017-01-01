Genesis
=======

Genesis is a Go library for building stand-alone installers. It is
intended as a configuration management utility for embedded systems, but
it can be used for pretty much any system.

Genesis can install, uninstall (reverse all installation steps), and
report status. You can construct groups of tasks, define sections, and
have tasks/groups/sections run conditionally based on the actions of
other tasks/groups/sections. You can tell it to re-run certain steps or
sections, and/or to skip some.



Motivation
----------

Traditional configuration management systems (chef, ansible, etc.) don't
work well for embedded systems because they tend to assume:

- the target is accessible by network and the network has been configured
- the target is running an ssh server
- supporting software has been installed (python, chef client, etc).

Genesis is designed to configure a system *from scratch*. See the `doc`
directory for more information.


Example
-------

Here is a very simple example:

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


Build
-----

There are two ways to build the installer: manually or with genesis's
assistance.  Both begin by building the exectable with:

    go build [FILE]

or

    GOOS=linux GOARCH=arm go build [FILE]

if you are cross compiling.


### Manual build

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

