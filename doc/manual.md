Genesis Installer
=================

This the manual for the Genesis installer.

## Quick Start

This sections shows a quick introduction into how to
use the Genesis installer.

### Building an installer

Building an installer is as simple as:

1. Compile the code like usual, with e.g. `GOOS=linux GOARCH=arm go build`.
2. Zip up supporting files, with e.g. `zip -r files.zip files/`.
3. Append the zip file to the binary: `cat files.zip >> my_installer`.
4. Fix the zip file indexing: `zip -A my_installer`.

### Programming an installer

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

### Running an installer

Once you have build an installer, it behaves like any ordinary executable.
Run it with `-h` to see all the options.  Most of the time, you just use
`-install`, `-status`, or `-remove`.  Other options are covered later in
this manual.


## Types of Doers

In the above code, we created instances of the Apt module and
added them to the installer.  Truth be told, we didn't actually
add the module instance directly to the installer.  By using the
`AddTask` method, we actually created a 'Task' containing a module
instance.  A Task is just a thin wrapper around a module instance.
We do this because a Task is a type of `Doer` (whereas a module is not).

Other types of Doers include Sections, Groups, Customs, and IfThens.
A Group is a list of Doers, and a Section is a list of Doers with a title.
An IfThen is a pair of Doers, where the second Doer only is run if the
first Doer changes state.  Finally, a Custom is a Doer with mutable
methods. Customs are very useful for specifying a custom Status method.

Notices that all of the Doers (except Tasks) are collections of Doers.
So hierarchies of Doers can be created.  In this way, you can add
IfThens to Sections, and then have Sections coupled with IfThens!
