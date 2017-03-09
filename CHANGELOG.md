Change Log
==========

## [0.3.1] 2017-03-09

- switch argument order in switch/case: now condition comes first
- add "else" to switch/case.  "else" gets run if all other conditions
  are false.
- bugfix: rework internal switch/case storage
- don't show executable name in rerun history -- just show arguments

<hr>

## [0.3.0] 2017-02-23

This release involves many big changes, some which break backwards
compatibility.  Upgrading to this version will require rewriting
your code.

changes:
- Compiled executable can now build its own archive
  + The executable knows which files it needs, so this makes build
    scripts less needed.
- New command-line options and help screen
  + Actions such as 'install', 'remove', etc are now subcommands instead of
    flags.
  + Replace './installer -tags 89a12f -install' with './installer install -tags 89a12f'
  + History lines (rerun) are automatically updated.
- New way to specify 'user' flags to genesis installer.
- LineInFile and BlockInFile have been merged.
- Fixed some bugs related to history file and store directory paths
- New Switch-Case Doer: conditional execution over a set of doers

<hr>

## [0.2.0] 2016-10-20

bugfixes:
- Fixes a few bugs in some modules (blockinfile, file)
- Fixes a bug in 'undo' that was causing it to 'do'
- No longer complains about history file on first run

features:
- Absent and Empty options for mkdir
- Local option for file (don't follow symlinks)
- File globbing in file module
- Readline prompt for history navigation
- Optional timeout in command module


<hr>


## [0.1.1] 2016-10-04

- If a task is done, report "DONE" rather than "PASS" so we know
  what has changed.
- create a pass/fail summary at end of report
- allow user to request that a line be absent from a file
- [bugfix] if status is unknown, don't say the task has failed
- [bugfix] check response error for httpget module
- [bugfix] remove blank lines from history

<hr>

## [0.1.0] 2016-09-11

First official release.  Still in the early stages,
but I am using for real-world systems.

