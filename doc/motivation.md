The Motivation Behind Genesis
=============================

This document explains the reasons Genesis was created,
and helps illuminate the problems Genesis tries to solve.

## Configuration Management

Some time ago, I.T. departments configured computers by hand.
This is a time-consuming and error prone process.
But, most computers had lifetimes of 3-5 years or more, so a day
or two spent configuring and debugging was hardly a problem.

As software development moved to cloud services (AWS et al.), the
pace of software deployment became more rapid.  Software tools were
created to automate much of the configuration process.  Such tools
now include Chef, Puppet, Ansible, and many others.

These tools provide several common benefits.  They significanly cut
the time it takes to configure a computer (from hours-to-days to
seconds-to-minutes).  They also remove human error from the configuration
process (computers don't "forget" to do a step), which makes for much
more consistent platforms.  They also streamline software testing, by
allowing developers to build their own test systems in an automated way.
There are many other benefits, but I think these are the most important
ones.

## Embedded Systems

While configuration management is very popular in web development,
most embedded systems development does not use it.  Before discussing
challenges and solutions for embedded systems, we must define what
we mean by "embedded."

Not too long ago, a typical embedded system featured a very low power
processor, very limited RAM, and tiny flash storage.  Often times,
these systems didn't even run an proper OS.  While such systems still
abound (wrist watches, microwave ovens, etc.), this is not the focus
of my attention.

A new class of embedded devices has arisen.  Currently, it is possible
to buy a fairly capable computer with minimal size and power requirements
for very little money.  The Raspberry Pi Zero and C.H.I.P. are two examples.
For a few dollars, you can have a 1 GHz quad-core processor, 500 MB of RAM
and a few gigabytes of storage.  This has led to the development of many
devices with high-powered embedded chips.

As part of my day job, I develop sensors for the military.  All of our
sensors contain an embedded board (or two) for processing data,
managing subsystms, communications, and even data analytics.  The newest
systems are running embedded boards with (at least) a few hundred MB of
RAM and storage capacities of 2-32 GB.  It is this class of embedded
systems that I will talk about here.

## Configuration Management for Embedded Systems.

So if these new embedded systems have ample resources and run a full-blown
operating system, why can't we just use existing configuration managment
software?  The short answer is that you can ... to a point.  Until now,
I used Ansible to configure all the devices.  But there are problems with
such an approach.

In order to use traditional configuration managemnt, there are things which
must be true about a system before we can even start.  In the case of Ansible,
the network must be configured and connected, sshd must be installed and enabled,
and python must be installed.  This is trivial on some boards (like
the raspberry pi). But it is not at all trivial on other boards which come with
stripped-down operating systems.  So now we need two configuration management
systems.  One to prepare the system for configuration management, and one to
do the configuration management!

There are other problems with classic configuration management tools, which
are designed around cloud-based servers.  Typically, there is no good way
to roll back changes made to a OS.  This is not a problem in the cloud,
because you can just destroy your cloud VMs and start over.  But you can't
destroy a physical device and start over!  Along the same lines, many
CMSs don't have a good way to examine the current state of the system.

Finally, all CMSs that I am aware of operate over a network connection.
This obviously won't work for offline systems.

More to come...
