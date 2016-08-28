Genesis Installer
=================

This the manual for the Genesis installer.

## Building an installer

Building an installer is as simple as:

1. Compile the code like usual, with e.g. `GOOS=linux GOARCH=arm go build`.
2. Zip up supporting files, with e.g. `zip -r files.zip files/`.
3. Append the zip file to the binary: `cat files.zip >> my_installer`.
4. Fix the zip file indexing: `zip -A my_installer`.

## Programming an installer

The Genesis installer is just a Go library.  Here is a simple example
of an installer:

	package main

	import (
		"github.com/wx13/genesis/installer"
		"github.com/wx13/genesis/modules"
	)

	func main() {
		inst = installer.New()
		defer inst.Done()
		pkgs := []string{"git", "gitk", "tig", "screen", "w3m"}
		for _, pkg := range pkgs {
			inst.AddTask(modules.Apt{Name: pkg})
		}
	}

This is a simple program which installs several packages using debian's apt.
The apt modeule lives inside the Genesis repository, but modules can be
imported from anywhere.

This code may not look like much, but it performs three roles.  On `install`,
it installs the listed packages.  On `remove` it removes the packages. And on
`status` it reports the install state of the packages.

