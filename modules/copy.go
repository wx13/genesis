package modules

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/wx13/genesis"
	"github.com/wx13/genesis/store"
)

type CopyFile struct {
	DestFile string
	SrcFile  string
	Store    *store.Store
}

func (cpf CopyFile) Describe() string {
	return fmt.Sprintf("CopyFile: %s => %s", cpf.SrcFile, cpf.DestFile)
}

func (cpf CopyFile) ID() string {
	return "copyFile" + cpf.DestFile
}

func (cpf CopyFile) Remove() (string, error) {

	cpf.DestFile = genesis.ExpandHome(cpf.DestFile)

	err := cpf.Store.RestoreFile(cpf.DestFile, "")
	if err == nil {
		return "Successfully restored destination file.", nil
	}
	return "Failed to restore destination file.", err
}

func (cpf CopyFile) Install() (string, error) {

	cpf.DestFile = genesis.ExpandHome(cpf.DestFile)

	bytes, err := ioutil.ReadFile(cpf.SrcFile)
	if err != nil {
		return "Could not read source file.", err
	}

	err = cpf.Store.SaveFile(cpf.DestFile, "")
	if err != nil {
		return "Could not save snapshot to file store.", err
	}

	err = ioutil.WriteFile(cpf.DestFile, bytes, 0644)
	if err != nil {
		return "Could not write destination file.", err
	}

	return "Successfully copied file.", nil

}

func (cpf CopyFile) Status() (genesis.Status, string, error) {

	cpf.DestFile = genesis.ExpandHome(cpf.DestFile)

	src, err := ioutil.ReadFile(cpf.SrcFile)
	if err != nil {
		return genesis.StatusFail, "Could not read source file.", err
	}
	dest, err := ioutil.ReadFile(cpf.DestFile)
	if err != nil {
		return genesis.StatusFail, "Could not read destination file.", err
	}

	if string(src) == string(dest) {
		return genesis.StatusPass, "File has been copied.", nil
	}
	return genesis.StatusFail, "File has not been copied.", errors.New("Source and destination files differ.")

}
