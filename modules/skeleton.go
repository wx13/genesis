// Package modules provides some installer modules for the
// genesis package.  But modules can reside anywhere, as long
// as they implement the genesis.Module interface.
package modules

import (
	"github.com/wx13/genesis"
)

// Skeleton can be any type of object.  Commonly it is
// a struct, but it doesn't have to be.
type Skeleton struct {
	Foo string
}

// Describe returns a human-readable description.
func (sk Skeleton) Describe() string {
	return "Skeleton " + sk.Foo
}

// ID returns a deterministic, unique string for the module.
func (sk Skeleton) ID() string {
	return sk.Describe()
}

func (sk Skeleton) Files() []string {
	return []string{}
}

// Install performs an action, and reports on its success.
func (sk Skeleton) Install() (string, error) {
	return "Skeleton did nothing.", nil
}

// Remove reverses the Install action.
func (sk Skeleton) Remove() (string, error) {
	return "Reverse of nothing is nothing.", nil
}

// Status reports on the current state of the system (does
// the install command need to run or not)>
func (sk Skeleton) Status() (genesis.Status, string, error) {
	return genesis.StatusPass, "Skeleton always passes", nil
}
