package modules

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/wx13/genesis"
	"github.com/wx13/genesis/store"
)

// LineInFile lets the user insert lines of text into a file.
type LineInFile struct {

	// Required
	File string   // path to the file
	Line []string // line(s) to insert

	// Optional
	Pattern []string     // line(s) to replace
	Success []string     // pattern to check for success (defaults to Line)
	Store   *store.Store // for storing changes
	Before  []string     // insert line before this pattern
	After   []string     // insert line after this pattern
	Absent  bool         // ensure line is absent from file

}

func (lif LineInFile) ID() string {
	short := fmt.Sprintf("LineInFile: file=%s, line=%s, pattern=%s", lif.File, lif.Line, lif.Pattern)
	long := fmt.Sprintf("before=%s, after=%s, success=%s absent=%s", lif.Before, lif.After, lif.Success, lif.Absent)
	return short + "\n" + long
}

func (lif LineInFile) Files() []string {
	return []string{lif.File}
}

func (lif LineInFile) Remove() (string, error) {
	lif.File = genesis.ExpandHome(lif.File)
	err := lif.Store.ApplyPatch(lif.File, lif.ID())
	if err != nil {
		return "Could not apply patch.", err
	}
	return "Patch applied", nil
}

func (lif LineInFile) Status() (genesis.Status, string, error) {

	lif.File = genesis.ExpandHome(lif.File)

	lines, err := lif.readFile()
	if err != nil {
		return genesis.StatusFail, "Could not read file.", err
	}

	_, lines, _ = lif.split(lines, lif.After, lif.Before)

	present, _, _ := lif.find(lines)
	if present {
		if lif.Absent {
			return genesis.StatusFail, "Line is in file.", nil
		} else {
			return genesis.StatusPass, "Line is in file.", nil
		}
	} else {
		if lif.Absent {
			return genesis.StatusPass, "Line is absent from file.", nil
		} else {
			return genesis.StatusFail, "Line is absent from file.", nil
		}
	}

}

func (lif LineInFile) Install() (string, error) {

	lif.File = genesis.ExpandHome(lif.File)

	lines, _ := lif.readFile()
	origLines := strings.Join(lines, "\n")

	beg, mid, end := lif.split(lines, lif.After, lif.Before)

	mid = lif.replace(mid)
	lines = append(beg, append(mid, end...)...)

	err := lif.writeFile(lines)
	if err != nil {
		return "Unable to write file.", err
	}

	lif.Store.SavePatch(lif.File, origLines, strings.Join(lines, "\n"), lif.ID())

	return "Wrote line to file", nil

}

// replace either replaces pattern line with line, or inserts
// the line at the end.
func (lif *LineInFile) replace(lines []string) []string {
	present, start, stop := lif.find(lines)
	if !present {
		return append(lines, lif.Line...)
	}
	if stop == len(lines)-1 {
		return append(lines[:start], lif.Line...)
	}
	return append(lines[:start], append(lif.Line, lines[stop+1:]...)...)
}

func (lif LineInFile) readFile() ([]string, error) {
	content, err := ioutil.ReadFile(lif.File)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(content), "\n")
	return lines, nil
}

func (lif LineInFile) writeFile(lines []string) error {
	content := strings.Join(lines, "\n") + "\n"
	err := ioutil.WriteFile(lif.File, []byte(content), 0644)
	return err
}

// findPattern looks for a line that matches the lif.Pattern (or Success) regex.
// Returns true if it finds it, and the slice index.
func (lif LineInFile) find(lines []string) (bool, int, int) {
	var pattern []string
	if len(lif.Success) == 0 {
		pattern = lif.Pattern
	} else {
		pattern = lif.Success
	}
	return lif.findPattern(lines, pattern)
}

func (lif LineInFile) findPattern(lines []string, pattern []string) (bool, int, int) {
	if len(lines) == 0 || len(pattern) == 0 {
		return false, -1, -1
	}
	idx := 0
	start := -1
	for k, line := range lines {
		match, _ := regexp.MatchString(pattern[idx], line)
		if match {
			if start < 0 {
				start = k
			}
			idx++
			if idx >= len(pattern) {
				return true, start, k
			}
		}
	}
	return false, -1, -1
}

// Grab the lines between start and end (exclusive).
func (lif LineInFile) split(lines []string, sPtrn, ePtrn []string) (beg, mid, end []string) {

	stop := -1
	found := false
	start := 0

	// Assign to 'beg' everything up through the start pattern match (inclusive).
	if len(sPtrn) > 0 {
		found, _, stop = lif.findPattern(lines, sPtrn)
		if found {
			beg = lines[:stop+1]
		}
	}

	// If there are no lines left, we are done.
	if stop >= len(lines)-1 {
		return beg, mid, end
	}

	// If there is no end pattern, we are done.
	if len(ePtrn) == 0 {
		return beg, lines[stop+1:], end
	}

	// Find the end pattern in the remaining text.
	mid = lines[stop+1:]
	found, start, _ = lif.findPattern(mid, ePtrn)
	if found {
		end = mid[start:]
		mid = mid[:start]
	}
	return beg, mid, end
}
