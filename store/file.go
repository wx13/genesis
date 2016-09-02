package store

import (
	"io/ioutil"
	"os"
	"path"
)

// RestoreFile restores a file from backup.
func (store *Store) RestoreFile(filename, label string) error {

	if store == nil {
		return nil
	}

	src := store.createPath(filename, label)

	// If the backup file doesn't exist, then we remove the original.
	_, err := os.Stat(src)
	if err != nil {
		os.Remove(filename)
	}

	// If we can't read the backup file, then there is
	// nothing to revert.
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		return nil
	}

	// Write file.
	return store.WriteFile(filename, bytes)

}

// SaveFile makes a backup of a file.
func (store *Store) SaveFile(filename, label string) error {

	if store == nil {
		return nil
	}

	// If we can't read the source file, then we can't back it up.
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	dest := store.createPath(filename, label)

	// Create the destination directory.
	dir := path.Dir(dest)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Write to the backup file, but only if it doesn't exist already.
	info, _ := os.Stat(filename)
	mode := info.Mode()
	f, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	_, err = f.Write(bytes)

	return err

}
