package store

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// ApplyPatch applies a patch to a file.
func (store *Store) ApplyPatch(filename, label string) error {

	if store == nil {
		return errors.New("no store")
	}

	// Read file.
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fileStr := string(b)

	// Read patch file.
	patchFile := store.createPath(filename, label)
	b, err = ioutil.ReadFile(patchFile)
	if err != nil {
		return err
	}
	patchStr := string(b)

	// Apply patch.
	dmp := diffmatchpatch.New()
	patches, _ := dmp.PatchFromText(patchStr)
	fileStr2, _ := dmp.PatchApply(patches, fileStr)

	// Write file.
	return store.WriteFile(filename, []byte(fileStr2))

}

// SavePatch computes and stores the patch between to strings.
func (store *Store) SavePatch(filename, origStr, newStr, label string) error {

	if store == nil {
		return errors.New("no store")
	}

	// Create patch
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(newStr, origStr, false)
	diffs = dmp.DiffCleanupSemantic(diffs)
	patches := dmp.PatchMake(diffs)
	strPatch := dmp.PatchToText(patches)

	// Create the destination directory.
	dest := store.createPath(filename, label)
	dir := filepath.Dir(dest)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Write patch to file
	err = ioutil.WriteFile(dest, []byte(strPatch), 0644)

	return err

}
