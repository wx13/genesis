package modules_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/wx13/genesis"
	"github.com/wx13/genesis/modules"
)

func TestEmptyFile(t *testing.T) {

	// Create an empty file.
	file, _ := ioutil.TempFile("", "genesis")
	file.Close()
	defer os.Remove(file.Name())

	// Put a line in the file.
	lif := modules.LineInFile{
		File:    file.Name(),
		Line:    "This is the line",
		Pattern: "^This is",
	}

	// Check status.
	status, _, err := lif.Status()
	if status != genesis.StatusFail {
		t.Error("Status should be", genesis.StatusFail, "but got", status)
	}

	// Run install.
	msg, err := lif.Install()
	if err != nil {
		t.Error(err, msg)
	}

	// Check that it worked.
	b, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Error("Could not read file.")
	}
	line := strings.Split(string(b), "\n")[1]
	if line != "This is the line" {
		t.Errorf("Line should be:\n%s\nBut is:\n%s\n", "This is the line", line)
	}

	// Check status.
	status, _, err = lif.Status()
	if err != nil {
		t.Error("Error getting status:", err)
	}
	if status != genesis.StatusPass {
		t.Error("Status should be", genesis.StatusPass, "but got", status)
	}

	// Do it again, and nothing should happen.
	msg, err = lif.Install()
	if err != nil {
		t.Error(err, msg)
	}

	// Undo. Nothing should happen, because we have no store.
	msg, err = lif.Remove()

}

func TestInsertLine(t *testing.T) {

	// Create an empty file.
	file, _ := ioutil.TempFile("", "genesis")
	file.Write([]byte("Line one\nLine two\n"))
	file.Close()
	defer os.Remove(file.Name())

	// Put a line in the file.
	lif := modules.LineInFile{
		File:    file.Name(),
		Line:    "This is the line",
		Pattern: "^This is",
		After:   "Line one",
		Before:  "Line two",
	}
	msg, err := lif.Install()
	if err != nil {
		t.Error(err, msg)
	}

	// Check that it worked.
	b, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Error("Could not read file.")
	}
	lines := strings.Split(string(b), "\n")
	if lines[0] != "Line one" || lines[1] != "This is the line" || lines[2] != "Line two" {
		t.Errorf("file should be:\n%s\nBut is:\n%s\n", "Line one\nThis is the line\nLine two\n", string(b))
	}

	// Do it again, and nothing should happen.
	msg, err = lif.Install()
	if err != nil {
		t.Error(err, msg)
	}

	// Undo. Nothing should happen, because we have no store.
	msg, err = lif.Remove()

}

func TestModifyLine(t *testing.T) {

	// Create an empty file.
	file, _ := ioutil.TempFile("", "genesis")
	file.Write([]byte("Foo false\nBar false\nBuz false\n"))
	file.Close()
	defer os.Remove(file.Name())

	// Put a line in the file.
	lif := modules.LineInFile{
		File:    file.Name(),
		Line:    "Bar true",
		Pattern: "^Bar",
	}
	msg, err := lif.Install()
	if err != nil {
		t.Error(err, msg)
	}

	// Check that it worked.
	b, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Error("Could not read file.")
	}
	lines := strings.Split(string(b), "\n")
	if lines[0] != "Foo false" || lines[1] != "Bar true" || lines[2] != "Buz false" {
		t.Errorf("file should be:\n%s\nBut is:\n%s\n", "Line one\nThis is the line\nLine two\n", string(b))
	}

	// Do it again, and nothing should happen.
	msg, err = lif.Install()
	if err != nil {
		t.Error(err, msg)
	}

	// Undo. Nothing should happen, because we have no store.
	msg, err = lif.Remove()

}
