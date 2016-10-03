package modules

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/wx13/genesis"
	"github.com/wx13/genesis/store"
)

type HttpGet struct {
	Dest  string
	Url   string
	Store *store.Store
}

func (get HttpGet) Describe() string {
	return fmt.Sprintf("HttpGet: %s => %s", get.Url, get.Dest)
}

func (get HttpGet) ID() string {
	return get.Describe()
}

func (get HttpGet) Remove() (string, error) {

	get.Dest = genesis.ExpandHome(get.Dest)

	err := get.Store.RestoreFile(get.Dest, "")
	if err == nil {
		return "Successfully restored destination file.", nil
	}
	return "Failed to restore destination file.", err
}

func (get HttpGet) Install() (string, error) {

	get.Dest = genesis.ExpandHome(get.Dest)

	err := get.Store.SaveFile(get.Dest, "")
	if err != nil {
		return "Could not save snapshot to file store.", err
	}

	out, err := os.Create(get.Dest)
	if err != nil {
		return "Could not create destination file.", err
	}
	defer out.Close()

	resp, err := http.Get(get.Url)
	if err != nil {
		return "Could not fetch file.", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "Could not write to destination file.", err
	}

	return "Successfully copied file.", nil

}

func (get HttpGet) Status() (genesis.Status, string, error) {

	get.Dest = genesis.ExpandHome(get.Dest)

	dest, err := ioutil.ReadFile(get.Dest)
	if err != nil {
		return genesis.StatusFail, "Could not read destination file.", err
	}

	resp, err := http.Get(get.Url)
	if err != nil {
		return genesis.StatusFail, "Could not fetch file.", err
	}
	defer resp.Body.Close()
	src, err := ioutil.ReadAll(resp.Body)

	if string(src) == string(dest) {
		return genesis.StatusPass, "File has been downloaded.", nil
	}

	return genesis.StatusFail, "File has not been downloaeded.", errors.New("Source and destination files differ.")

}
