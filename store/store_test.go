package store_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/wx13/genesis/store"
)

func TestStore(t *testing.T) {

	dir, err := ioutil.TempDir("", "genesis_test")
	if err != nil {
		t.Error("Could not create temp dir for testing purposes")
	}
	defer os.RemoveAll(dir)

	storeDir := path.Join(dir, "store")
	os.Mkdir(storeDir, 0755)
	s, err := store.New(storeDir)
	if err != nil {
		t.Error("Could not create store:", err)
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

	text2 := "This is line 27,\nand this is line owt.\n3\nNow line four.\n"
	s.SavePatch(filename, text, text2, "foo")
	ioutil.WriteFile(filename, []byte(text2), 0644)
	s.ApplyPatch(filename, "foo")
	data, err = ioutil.ReadFile(filename)
	if string(data) != text {
		t.Error("Patched file is not equal to original file.")
	}

}
