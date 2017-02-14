package modules

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/wx13/genesis"
)

type CopyFile struct {
	Dest string
	Src  string
}

func (cpf CopyFile) src() string {
	match, _ := regexp.MatchString("^[.]?/", cpf.Src)
	if match {
		return cpf.Src
	}
	return filepath.Join(genesis.Tmpdir, cpf.Src)
}

func (cpf CopyFile) ID() string {
	return fmt.Sprintf("CopyFile: %s => %s", cpf.Src, cpf.Dest)
}

func (cpf CopyFile) Files() []string {
	return []string{cpf.src()}
}

func (cpf CopyFile) Remove() (string, error) {

	cpf.Dest = genesis.ExpandHome(cpf.Dest)

	err := genesis.Store.RestoreFile(cpf.Dest, "")
	if err == nil {
		return "Successfully restored destination file.", nil
	}
	return "Failed to restore destination file.", err
}

func (cpf CopyFile) Install() (string, error) {

	cpf.Dest = genesis.ExpandHome(cpf.Dest)

	bytes, err := ioutil.ReadFile(cpf.src())
	if err != nil {
		return "Could not read source file.", err
	}

	err = genesis.Store.SaveFile(cpf.Dest, "")

	err = ioutil.WriteFile(cpf.Dest, bytes, 0644)
	if err != nil {
		return "Could not write destination file.", err
	}

	return "Successfully copied file.", nil

}

func (cpf CopyFile) Status() (genesis.Status, string, error) {

	cpf.Dest = genesis.ExpandHome(cpf.Dest)

	src, err := ioutil.ReadFile(cpf.src())
	if err != nil {
		return genesis.StatusFail, "Could not read source file.", err
	}
	dest, err := ioutil.ReadFile(cpf.Dest)
	if err != nil {
		return genesis.StatusFail, "Could not read destination file.", err
	}

	if string(src) == string(dest) {
		return genesis.StatusPass, "File has been copied.", nil
	}
	return genesis.StatusFail, "File has not been copied.", errors.New("Source and destination files differ.")

}
