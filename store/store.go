// Package store provides support for keeping track of changes
// to files.  It can keep a copy of a file, or a patch to reverse
// changes to a file.
package store

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Store is for storing change information for a set of files.
type Store struct {
	Dir string
}

// New generates a new Store object.
func New(dir string) (*Store, error) {
	store := Store{Dir: dir}
	err := os.MkdirAll(store.Dir, 0755)
	if err != nil {
		return &store, err
	}
	return &store, nil
}

func (store *Store) createPath(filename, label string) string {
	if len(label) > 0 {
		label = "." + label
	}
	filename = strings.Replace(filename, ":", "", 1)
	return filepath.Join(store.Dir, filename+label)
}

// Hash computes a hash of a set of strings.
func Hash(things ...string) string {
	data := []byte(strings.Join(things, "_"))
	hash := fmt.Sprintf("%x", md5.Sum(data))
	if len(hash) < 5 {
		return hash
	}
	return hash[:5]
}

// WriteFile writes data to a file, handling permissions.
func (store *Store) WriteFile(filename string, bytes []byte) error {

	// If destination exists, keep the same permissions.
	info, err := os.Stat(filename)
	if err == nil {
		mode := info.Mode()
		err := ioutil.WriteFile(filename, bytes, mode)
		return err
	}

	// Default permissions.
	return ioutil.WriteFile(filename, bytes, 0644)

}
