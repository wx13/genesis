package store_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/wx13/genesis/store"
)

func TestHash(t *testing.T) {
	hash1 := store.Hash("hello", "bye", "yo")
	hash2 := store.Hash("hello", "bye", "yo")
	hash3 := store.Hash("hello", "bYe", "yo")
	if hash1 != hash2 || hash2 == hash3 {
		t.Error("Hashing is broken")
	}
}

func TestStore(t *testing.T) {

	// Create a temp directory for testing.
	dir, err := ioutil.TempDir("", "genesis_test")
	if err != nil {
		t.Error("Could not create temp dir for testing purposes")
	}
	defer os.RemoveAll(dir)

	// Generate a file for testing.
	filename := filepath.Join(dir, "myfile.txt")
	text := "This is line 1,\nand this is line two.\n\nNow line four.\n"
	ioutil.WriteFile(filename, []byte(text), 0644)

	// Create a new store.
	storeDir := filepath.Join(dir, "store")
	os.Mkdir(storeDir, 0755)
	s, err := store.New(storeDir)
	if err != nil {
		t.Error("Could not create store:", err)
	}

	// Save a snapshot, remove the file, and then restore it.
	err = s.SaveFile(filename, "")
	if err != nil {
		t.Error("Could not save file to store.", err)
	}
	os.Remove(filename)
	err = s.RestoreFile(filename, "")
	if err != nil {
		t.Error("Could not restore file from store.", err)
	}

	// Test if restored file matches original.
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("Cannot read file:", err)
	}
	if string(data) != text {
		t.Error("Restored file is not equal to original file.", err)
	}

	// Modify the file to test patching.
	text2 := "This is line 27,\nand this is line owt.\n3\nNow line four.\n"
	err = s.SavePatch(filename, text, text2, "foo")
	if err != nil {
		t.Error("Error saving file patch.", err)
	}
	ioutil.WriteFile(filename, []byte(text2), 0644)
	err = s.ApplyPatch(filename, "foo")
	if err != nil {
		t.Error("Error patching file.", err)
	}
	data, err = ioutil.ReadFile(filename)
	if string(data) != text {
		t.Error("Patched file is not equal to original file.")
	}

}
