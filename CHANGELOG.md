Change Log
==========

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

