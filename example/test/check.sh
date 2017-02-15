#!/bin/bash

command() {
	out=$($@)
	if [ $? -eq 0 ]
	then
		pass $out
	else
		fail $out
	fi
}

heading() {
	printf "%-50s " "$@..."
}

pass() {
	echo -e "\e[32m PASS \e[0m"
}

fail() {
	echo -e "\e[31m FAIL \e[0m"
	echo $@
	exit
}

test.build() {
	heading "Build the executable"
	command go build
	heading "Build the archive"
	command ./test build
}

cleanup() {
	heading "Remove the target directory"
	command rm -rf /tmp/genesis_test
}

test.install() {
	heading "Run the install on a fresh instance"
	command ./test.x install
	heading "Check that installer worked as intended"
	if [ ! -d /tmp/genesis_test ]
	then
		fail "installer did not create directory"
	fi
	d=$(diff files/file.txt /tmp/genesis_test/file.txt)
	if [ ! -z "$d" ]
	then
		fail "installer did not copy file correctly." $d
	fi
	pass
}

test.tag_remove() {
	heading "Uninstall only one step"
	command ./test.x remove -tags d48e1a
	heading "Check that file was removed, but not directory"
	if [ -e /tmp/genesis_test/file.txt ]
	then
		fail "file was not removed"
	fi
	if [ ! -d /tmp/genesis_test ]
	then
		fail "directory was removed, but it shouldn't have been"
	fi
	pass
}

test.full_remove() {
	heading "Uninstall the rest"
	command ./test.x remove
	heading "Check that removal worked."
	if [ -d /tmp/genesis_test ]
	then
		fail "did not remove directory"
	fi
	pass
}



test.build
test.install
test.tag_remove
test.full_remove




