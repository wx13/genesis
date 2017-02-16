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
	heading "Remove store"
	command rm -rf ~/.genesis/store
}

test.install() {

	heading "Run the install on a fresh instance"
	command ./test.x install

	heading "Check that directory was created"
	if [ ! -d /tmp/genesis_test ]
	then
		fail "installer did not create directory"
	else
		pass
	fi

	heading "Check that file was copied"
	d=$(diff files/file.txt /tmp/genesis_test/file.txt)
	if [ ! -z "$d" ]
	then
		fail "installer did not copy file correctly." $d
	else
		pass
	fi

	heading "Check that template was run"
	if [ -e /tmp/genesis_test/file_from_template.txt ]
	then
		pass
	else
		fail "Template was not run"
	fi

	heading "Check that the README file was retrieved"
	title=$(head -1 /tmp/genesis_test/README.md)
	if [ "$title" == "Genesis" ]
	then
		pass
	else
		fail "Failed to get file from internet"
	fi

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

test.tag_remove_overwrite() {

	mkdir /tmp/genesis_test
	echo "foo" > /tmp/genesis_test/file.txt
	test.install

	heading "Uninstall text file which had overwritten a file"

	command ./test.x remove -tags d48e1a
	heading "Check that the old file was restored"
	if [ ! -e /tmp/genesis_test/file.txt ]
	then
		fail "Old file was not restored"
	fi
	content=$(cat /tmp/genesis_test/file.txt)
	if [ ! "$content" == "foo" ]
	then
		fail "Old file was restored with wrong content"
	fi

	pass
}

test.full_remove() {

	heading "Uninstall all"
	command ./test.x remove

	heading "Check that removal worked."
	if [ -d /tmp/genesis_test ]
	then
		fail "did not remove directory"
	fi
	pass

}



cleanup
test.build
test.install
test.tag_remove
test.full_remove

cleanup
test.tag_remove_overwrite



