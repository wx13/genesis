package store_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/wx13/genesis/store"
)

func TestFile(t *testing.T) {

	dir, err := ioutil.TempDir("", "genesis_test")
	if err != nil {
		t.Error("Could not create temp dir for testing purposes")
	}
	defer os.RemoveAll(dir)

	storeDir := path.Join(dir, "store")
	os.Mkdir(storeDir, 0755)
	s := store.New(storeDir)
	if s == nil {
		t.Error("Could not create store")
	}

	filename := path.Join(dir, "myfile.txt")
	text := "This is line 1,\nand this is line two.\n\nNow line four.\n"
	ioutil.WriteFile(filename, []byte(text), 0644)

	s.SaveFile(filename, "")
	os.Remove(filename)
	s.RestoreFile(filename, "")

	data, err := ioutil.ReadFile(filename)
	if string(data) != text {
		t.Error("Restored file is not equal to original file.")
	}

}
