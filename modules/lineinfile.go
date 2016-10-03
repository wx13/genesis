package modules

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/wx13/genesis"
	"github.com/wx13/genesis/store"
)

type LineInFile struct {
	File    string
	Pattern string
	Success string
	Line    string
	Store   *store.Store
	Label   string
	Before  string
	After   string
	Absent  bool
}

func (lif LineInFile) Describe() string {
	return fmt.Sprintf("LineInFile: %s, %s, %s", lif.File, lif.Pattern, lif.Line)
}

func (lif LineInFile) ID() string {
	return "lineInFile" + lif.File + lif.Pattern + lif.Line + lif.Label + lif.Before + lif.After
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

func (lif LineInFile) Remove() (string, error) {
	lif.File = genesis.ExpandHome(lif.File)
	err := lif.Store.ApplyPatch(lif.File, lif.Label)
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

	if lif.Absent {
		for _, line := range lines {
			match, _ := regexp.MatchString(lif.Pattern, line)
			if match {
				return genesis.StatusFail, "Line is in file.", nil
			}
		}
		return genesis.StatusPass, "Line is absent from file.", nil
	}

	isAfter := len(lif.After) == 0
	for _, line := range lines {
		match := lif.Line == line
		if len(lif.Success) > 0 {
			match, _ = regexp.MatchString(lif.Success, line)
		}
		if match && isAfter {
			return genesis.StatusPass, "Line is in file.", nil
		}
		match, _ = regexp.MatchString(lif.After, line)
		if match {
			isAfter = true
		}
		if len(lif.Before) > 0 {
			match, _ = regexp.MatchString(lif.Before, line)
			if match {
				break
			}
		}
	}
	return genesis.StatusFail, "Line not in file.", errors.New("Line not in file.")
}

func (lif LineInFile) Install() (string, error) {

	lif.File = genesis.ExpandHome(lif.File)

	lines, _ := lif.readFile()
	origLines := strings.Join(lines, "\n")

	done := false
	isAfter := len(lif.After) == 0
	for i, line := range lines {
		match, _ := regexp.MatchString(lif.Pattern, line)
		if match && lif.Absent {
			lines = append(lines[:i], lines[i+1:]...)
		} else if match && isAfter {
			lines[i] = lif.Line
			done = true
			break
		}
		match, _ = regexp.MatchString(lif.After, line)
		if match {
			isAfter = true
		}
		if len(lif.Before) > 0 {
			match, _ = regexp.MatchString(lif.Before, line)
			if match {
				lines = append(lines[:i], append([]string{lif.Line}, lines[i:]...)...)
				done = true
				break
			}
		}
	}
	if !done && !lif.Absent {
		lines = append(lines, lif.Line)
	}

	err := lif.writeFile(lines)
	if err != nil {
		return "Unable to write file.", err
	}

	lif.Store.SavePatch(lif.File, origLines, strings.Join(lines, "\n"), lif.Label)

	return "Wrote line to file", nil

}
