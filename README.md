Genesis
=======

Genesis is a Go library for building stand-alone installers.
It is intended as a configuration management utility for
embedded systems.


Status
------

This is experimental software, with high API instability.


Motivation
----------

Traditional configuration management systems (chef, ansible, etc.)
don't work well for embedded systems because they tend to assume:

- the target is accessible by network and the network has been configured
- the target is running an ssh server
- supporting software has been installed (python, chef, etc).

Genesis is designed to configure a system *from scratch*.  See the `doc`
directory for more information.


Example
-------

See the `example` directory for some examples.


Build
-----

To build the installer, first zip up the supporting files:

    zip -r files.zip files

Now build the executable:

    go build my_installer.go && cat files.zip >> my_installer && zip -A my_installer

This packages the zip file into the installer binary, so that
it is completely standalone.  Run the binary with the `-install`, `-status`,
or `-remove` flags to install / check status / remove.


