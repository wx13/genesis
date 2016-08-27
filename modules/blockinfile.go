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

type BlockInFile struct {
	File     string
	Patterns []string
	Lines    []string
	Success  []string
	Store    *store.Store
	Label    string
}

func (bif BlockInFile) Describe() string {
	return fmt.Sprintf("BlockInFile: %s, %s => %s, %s...", bif.File, bif.Patterns[0], bif.Patterns[1], bif.Lines[0])
}

func (bif BlockInFile) ID() string {
	return "blockInFile" + bif.File + strings.Join(bif.Patterns, "") + strings.Join(bif.Lines, "") + bif.Label
}

func (bif BlockInFile) readFile() ([]string, error) {
	content, err := ioutil.ReadFile(bif.File)
	if err != nil {
		return []string{}, err
	}
	lines := strings.Split(string(content), "\n")
	return lines, nil
}

func (bif BlockInFile) writeFile(lines []string) error {
	content := strings.Join(lines, "\n") + "\n"
	err := ioutil.WriteFile(bif.File, []byte(content), 0644)
	return err
}

func (bif BlockInFile) findLine(pattern string, lines []string) int {
	for i, line := range lines {
		match, _ := regexp.MatchString(pattern, line)
		if match {
			return i
		}
	}
	return -1
}

func (bif BlockInFile) Remove() (string, error) {
	bif.File = genesis.ExpandHome(bif.File)
	err := bif.Store.ApplyPatch(bif.File, bif.Label)
	if err != nil {
		return "Could not apply patch.", err
	}
	return "Patch applied", nil
}

func (bif BlockInFile) Status() (genesis.Status, string, error) {

	bif.File = genesis.ExpandHome(bif.File)

	fileLines, err := bif.readFile()
	if err != nil {
		return genesis.StatusFail, "Could not read file.", err
	}
	lines := bif.Lines
	if len(bif.Success) > 0 {
		lines = bif.Success
	} else {
		for k, line := range lines {
			lines[k] = regexp.QuoteMeta(line)
		}
	}
OUTER:
	for fileIdx := range fileLines {
		match, _ := regexp.MatchString(lines[0], fileLines[fileIdx])
		if match {
			for blockIdx, blockLine := range lines {
				k := fileIdx + blockIdx
				if k > len(fileLines) {
					continue OUTER
				}
				match, _ := regexp.MatchString(blockLine, fileLines[k])
				if !match {
					continue OUTER
				}
			}
			return genesis.StatusPass, "Block found in file.", nil
		}
	}
	return genesis.StatusFail, "Block not in file.", errors.New("block not in file")
}

func (bif BlockInFile) Install() (string, error) {

	bif.File = genesis.ExpandHome(bif.File)

	lines, _ := bif.readFile()
	origLines := strings.Join(lines, "\n")

	start := bif.findLine(bif.Patterns[0], lines)
	if start < 0 {
		lines = append(lines, bif.Lines...)
		bif.writeFile(lines)
	} else {
		end := bif.findLine(bif.Patterns[1], lines[start:])
		if end < 0 {
			lines = append(lines, bif.Lines...)
			bif.writeFile(lines)
		} else {
			lines = append(lines[:(start-1)], append(bif.Lines, lines[(end+1):]...)...)
		}
	}

	bif.Store.SavePatch(bif.File, origLines, strings.Join(lines, "\n"), bif.Label)

	return "Wrote block to file", nil

}
