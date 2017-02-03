package modules

import (
	"strings"
	"testing"
)

func TestFindPattern(t *testing.T) {

	lif := LineInFile{}
	var found bool
	var start, end int

	// Empty pattern in empty file should return false.
	found, _, _ = lif.findPattern([]string{}, []string{})
	if found == true {
		t.Errorf("Empty pattern should not be found in empty file.")
	}

	// Non-empty pattern in empty file should return false.
	found, _, _ = lif.findPattern([]string{}, []string{"hello"})
	if found == true {
		t.Errorf("Non-empty pattern should not be found in empty file.")
	}

	// Empty pattern in non-empty file should return false.
	found, _, _ = lif.findPattern([]string{"hello", "bob"}, []string{})
	if found == true {
		t.Errorf("Empty pattern should not be found in non-empty file.")
	}

	// Single line file with matching single-line pattern
	found, start, end = lif.findPattern([]string{"hello"}, []string{"^h"})
	if found == false || start != 0 || end != 0 {
		t.Error("Expected false, 0, 0; but got:", found, start, end)
	}

	// Multi-line pattern in a multi-line file
	found, start, end = lif.findPattern([]string{"hello", "bye", "yo", "hey"}, []string{"^h", "^y"})
	if found == false || start != 0 || end != 2 {
		t.Error("Expected false, 0, 2; but got:", found, start, end)
	}

}

func TestSplit(t *testing.T) {

	lif := LineInFile{}
	var beg, mid, end []string

	beg, mid, end = lif.split([]string{}, []string{}, []string{})
	if len(beg) != 0 || len(mid) != 0 || len(end) != 0 {
		t.Error("Empty file split with empty patterns -- all fields should be empty.")
	}

	beg, mid, end = lif.split([]string{"hello", "bye"}, []string{}, []string{})
	if len(beg) != 0 || len(mid) != 2 || len(end) != 0 {
		t.Error("Non-empty file split with empty patterns should all be empty except 'mid'.")
	}

	beg, mid, end = lif.split([]string{"hello", "bye"}, []string{"hello"}, []string{"bye"})
	if len(beg) != 1 || len(mid) != 0 || len(end) != 1 {
		t.Error("2 lines. start matches first.  end matches second.")
	}

	beg, mid, end = lif.split([]string{"hello", "bye"}, []string{"hello"}, []string{"byeeee"})
	if len(beg) != 1 || len(mid) != 1 || len(end) != 0 {
		t.Error("2 lines. start matches first.  end matches none.")
	}

	beg, mid, end = lif.split([]string{"hello", "bye", "yo", "hey", "sup"}, []string{"hello", "^b"}, []string{"hey"})
	if len(beg) != 2 || len(mid) != 1 || len(end) != 2 {
		t.Error("Start and end matching, with multi-line start")
	}

}

func TestReplace(t *testing.T) {

	lif := LineInFile{}
	var lines []string

	// Empty file, empty pattern, empty replace => still empty file.
	lines = lif.replace([]string{})
	if len(lines) > 0 {
		t.Error("Empty file, empty pattern, empty replace ==>", lines)
	}

	// Empty file, empty pattern, non-empty replace.
	lif.Line = []string{"hello"}
	lines = lif.replace([]string{})
	if lines[0] != "hello" {
		t.Error("Empty file, empty pattern, non-empty replace ==>", lines)
	}

	// Non-empty file, empty pattern, non-empty replace.
	lif.Line = []string{"hello"}
	lines = lif.replace([]string{"yo", "hey"})
	if strings.Join(lines, ":") != "yo:hey:hello" {
		t.Error("Non-empty file, empty pattern, non-empty replace ==>", lines)
	}

	// Non-empty file, non-empty pattern, non-empty replace.
	lif.Line = []string{"hello"}
	lif.Pattern = []string{"hey"}
	lines = lif.replace([]string{"yo", "hey"})
	if strings.Join(lines, ":") != "yo:hello" {
		t.Error("Non-empty file, non-empty pattern, non-empty replace ==>", lines)
	}

}
