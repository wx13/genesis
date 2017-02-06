package modules

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/wx13/genesis"
	"github.com/wx13/genesis/store"
)

type CopyFile struct {
	Dest  string
	Src   string
	Store *store.Store
}

func (cpf CopyFile) ID() string {
	return fmt.Sprintf("CopyFile: %s => %s", cpf.Src, cpf.Dest)
}

func (cpf CopyFile) Files() []string {
	return []string{cpf.Src}
}

func (cpf CopyFile) Remove() (string, error) {

	cpf.Dest = genesis.ExpandHome(cpf.Dest)

	err := cpf.Store.RestoreFile(cpf.Dest, "")
	if err == nil {
		return "Successfully restored destination file.", nil
	}
	return "Failed to restore destination file.", err
}

func (cpf CopyFile) Install() (string, error) {

	cpf.Dest = genesis.ExpandHome(cpf.Dest)

	bytes, err := ioutil.ReadFile(cpf.Src)
	if err != nil {
		return "Could not read source file.", err
	}

	err = cpf.Store.SaveFile(cpf.Dest, "")
	if err != nil {
		return "Could not save snapshot to file store.", err
	}

	err = ioutil.WriteFile(cpf.Dest, bytes, 0644)
	if err != nil {
		return "Could not write destination file.", err
	}

	return "Successfully copied file.", nil

}

func (cpf CopyFile) Status() (genesis.Status, string, error) {

	cpf.Dest = genesis.ExpandHome(cpf.Dest)

	src, err := ioutil.ReadFile(cpf.Src)
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
